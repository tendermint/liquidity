package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) ValidateMsgCreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
	params := k.GetParams(ctx)
	var poolType types.LiquidityPoolType

	// check poolType exist, get poolType from param
	if len(params.LiquidityPoolTypes)-1 >= int(msg.PoolTypeIndex) {
		poolType = params.LiquidityPoolTypes[msg.PoolTypeIndex]
		if poolType.PoolTypeIndex != msg.PoolTypeIndex {
			return types.ErrPoolTypeNotExists
		}
	} else {
		return types.ErrPoolTypeNotExists
	}

	if poolType.MaxReserveCoinNum > types.MaxReserveCoinNum || types.MinReserveCoinNum > poolType.MinReserveCoinNum {
		return types.ErrNumOfReserveCoin
	}

	// TODO: duplicated with ValidateBasic
	if len(msg.ReserveCoinDenoms) != msg.DepositCoins.Len() {
		return types.ErrNumOfReserveCoin
	}

	if uint32(msg.DepositCoins.Len()) > poolType.MaxReserveCoinNum && poolType.MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return types.ErrNumOfReserveCoin
	}

	return nil
}

func (k Keeper) CreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
	if err := k.ValidateMsgCreateLiquidityPool(ctx, msg); err != nil {
		return err
	}
	params := k.GetParams(ctx)
	poolKey := types.GetPoolKey(msg.ReserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)

	if _, found := k.GetLiquidityPoolByReserveAccIndex(ctx, reserveAcc); found {
		return types.ErrPoolAlreadyExists
	}

	poolCreator := msg.GetPoolCreator()
	accPoolCreator := k.accountKeeper.GetAccount(ctx, poolCreator)
	poolCreatorBalances := k.bankKeeper.GetAllBalances(ctx, accPoolCreator.GetAddress())
	if !poolCreatorBalances.IsAllGTE(msg.DepositCoins) {
		return types.ErrInsufficientBalance
	}

	for _, coin := range msg.DepositCoins {
		if coin.Amount.LT(params.MinInitDepositToPool) {
			return types.ErrLessThanMinInitDeposit
		}
	}

	denom1, denom2 := types.AlphabeticalDenomPair(msg.ReserveCoinDenoms[0], msg.ReserveCoinDenoms[1])
	reserveCoinDenoms := []string{denom1, denom2}

	PoolCoinDenom := types.GetPoolCoinDenom(reserveAcc)

	liquidityPool := types.LiquidityPool{
		PoolTypeIndex:         msg.PoolTypeIndex,
		ReserveCoinDenoms:     reserveCoinDenoms,
		ReserveAccountAddress: reserveAcc.String(),
		PoolCoinDenom:         PoolCoinDenom,
	}

	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(liquidityPool.PoolCoinDenom, params.InitPoolCoinMintAmount))
	if err := k.bankKeeper.SendCoins(ctx, poolCreator, reserveAcc, msg.DepositCoins); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, poolCreator, mintPoolCoin); err != nil {
		return err
	}

	k.SetLiquidityPoolAtomic(ctx, liquidityPool)
	batch := types.NewLiquidityPoolBatch(liquidityPool.PoolId, 0)
	k.SetLiquidityPoolBatch(ctx, batch)
	// TODO: params.LiquidityPoolCreationFee logic
	// TODO: refactoring, LiquidityPoolCreationFee, check event on handler
	return nil
}

// Get reserve Coin from the liquidity pool
func (k Keeper) GetReserveCoins(ctx sdk.Context, lp types.LiquidityPool) (reserveCoins sdk.Coins) {
	for _, denom := range lp.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, lp.GetReserveAccount(), denom))
	}
	return
}

func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, lp types.LiquidityPool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(lp.PoolCoinDenom)
}

func (k Keeper) ValidateMsgDepositLiquidityPool(ctx sdk.Context, msg types.MsgDepositToLiquidityPool) error {
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}

	pool, found := k.GetLiquidityPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	for _, coin := range msg.DepositCoins {
		if !types.StringInSlice(coin.Denom, pool.ReserveCoinDenoms) {
			return types.ErrInvalidDenom
		}
	}

	if msg.DepositCoins.Len() != len(pool.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	// TODO: duplicated with ValidateBasic
	if uint32(msg.DepositCoins.Len()) > types.MaxReserveCoinNum ||
		types.MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return types.ErrNumOfReserveCoin
	}

	return nil
}

func (k Keeper) DepositLiquidityPool(ctx sdk.Context, msg types.BatchPoolDepositMsg) error {
	msg.Executed = true
	// TODO: verify executed after err
	k.SetLiquidityPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg.MsgIndex, msg)

	if err := k.ValidateMsgDepositLiquidityPool(ctx, *msg.Msg); err != nil {
		return err
	}

	depositCoins := msg.Msg.DepositCoins.Sort()

	pool, found := k.GetLiquidityPool(ctx, msg.Msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	coinA := depositCoins[0]
	coinB := depositCoins[1]

	//lastReserveRatio := coinA.Amount.Quo(coinB.Amount)
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	//lastReserveRatio := reserveCoins[0].Amount.Quo(reserveCoins[1].Amount)
	// TODO: To truncateInt
	depositableCoinA := coinA.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
	//depositableCoinB := coinB.Amount.Mul(lastReserveRatio.Ceil().RoundInt())
	var inputs []banktypes.Input
	var outputs []banktypes.Output

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	reserveAcc := pool.GetReserveAccount()
	depositor := msg.Msg.GetDepositor()

	if coinB.Amount.GT(depositableCoinA) {
		depositableCoinB := depositableCoinA
		refundAmtB := coinA.Amount.Sub(depositableCoinB)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableCoinB))))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableCoinB))))
		// refund
		if refundAmtB.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
			outputs = append(outputs, banktypes.NewOutput(depositor, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
		}
	} else if coinB.Amount.LT(depositableCoinA) {
		// TODO: handle dec to Int error, roundInt, TruncateInt to using ToDec,  MulInt
		depositableCoinA = coinB.Amount.Quo(lastReserveRatio.RoundInt())
		refundAmtA := coinA.Amount.Sub(depositableCoinA)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableCoinA))))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableCoinA))))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinB)))
		// refund
		if refundAmtA.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
			outputs = append(outputs, banktypes.NewOutput(depositor, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
		}
	} else {
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinB)))
	}

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	// calculate pool token mint amount
	poolCoinAmt := k.GetPoolCoinTotalSupply(ctx, pool).Mul(depositableCoinA).Quo(coinA.Amount)
	poolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt))

	// mint pool token to Depositor
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, poolCoin); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, poolCoin); err != nil {
		return err
	}
	msg.Succeed = true
	msg.ToDelete = true
	k.SetLiquidityPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg.MsgIndex, msg)
	// TODO: add events for batch result, each err cases
	return nil
}

func (k Keeper) ValidateMsgWithdrawLiquidityPool(ctx sdk.Context, msg types.MsgWithdrawFromLiquidityPool) error {
	// TODO: add validate logic
	return nil
}

func (k Keeper) ValidateMsgSwap(ctx sdk.Context, msg types.MsgSwap) error {
	// TODO: add validate logic
	return nil
}

func (k Keeper) WithdrawLiquidityPool(ctx sdk.Context, msg types.BatchPoolWithdrawMsg) error {
	msg.Executed = true
	// TODO: verify executed after err
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg.MsgIndex, msg)

	if err := k.ValidateMsgWithdrawLiquidityPool(ctx, *msg.Msg); err != nil {
		return err
	}
	// TODO: validate reserveCoin balance

	poolCoin := msg.Msg.PoolCoin
	poolCoins := sdk.NewCoins(poolCoin)

	pool, found := k.GetLiquidityPool(ctx, msg.Msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	for _, reserveCoin := range reserveCoins {
		withdrawAmt := reserveCoin.Amount.Mul(poolCoin.Amount).Quo(totalSupply)
		// TODO: apply fee, (sdk.NewDec(1).Sub(params.LiquidityPoolFeeRate)
		// TODO: to using k.bankKeeper.SendCoinsFromModuleToAccount() with poolKey
		inputs = append(inputs, banktypes.NewInput(pool.GetReserveAccount(),
			sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
		outputs = append(outputs, banktypes.NewOutput(msg.Msg.GetWithdrawer(),
			sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
	}

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName),
		types.ModuleName, poolCoins); err != nil {
		return err
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins); err != nil {
		return err
	}
	msg.Succeed = true
	msg.ToDelete = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg.MsgIndex, msg)
	// TODO: add events for batch result, each err cases
	return nil
}

// TODO: testcodes
func (k Keeper) RefundDepositLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolDepositMsg) error {
	if !batchMsg.Executed || batchMsg.Succeed {
		// TODO: make Err type
		panic("can't refund not executed or succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetDepositor(), batchMsg.Msg.DepositCoins)
	if err != nil {
		panic(err)
	}
	msg, found := k.GetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !found {
		panic(err)
	}
	msg.ToDelete = true
	k.SetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex, msg)
	k.DeleteLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

func (k Keeper) RefundWithdrawLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolWithdrawMsg) error {
	if !batchMsg.Executed || batchMsg.Succeed {
		panic("can't refund not executed or succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetWithdrawer(), sdk.NewCoins(batchMsg.Msg.PoolCoin))
	if err != nil {
		panic(err)
	}
	msg, found := k.GetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !found {
		panic(err)
	}
	msg.ToDelete = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex, msg)
	k.DeleteLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

// TODO: WIP
func (k Keeper) RefundSwapLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolSwapMsg) error {
	if !batchMsg.Executed || batchMsg.Succeed {
		panic("can't refund not executed or succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetSwapRequester(), sdk.NewCoins(batchMsg.Msg.OfferCoin))
	if err != nil {
		panic(err)
	}
	msg, found := k.GetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !found {
		panic(err)
	}
	msg.ToDelete = true
	k.SetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex, msg)
	k.DeleteLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

//func (k Keeper) FractionalRefundSwapLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolSwapMsg) error {
//	if !batchMsg.Executed {
//		panic("can't refund not executed msg")
//	}
//}
