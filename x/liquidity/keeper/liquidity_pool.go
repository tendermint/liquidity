package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) CreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
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

	if poolType.MinReserveCoinNum != 2 && poolType.MaxReserveCoinNum != 2 {
		return types.ErrNotImplementedYet
	}

	if len(msg.ReserveCoinDenoms) > int(poolType.MaxReserveCoinNum) && int(poolType.MinReserveCoinNum) > len(msg.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	poolKey := types.GetPoolKey(msg.ReserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)

	if _, found := k.GetLiquidityPoolByReserveAccIndex(ctx, reserveAcc); found {
		return types.ErrPoolAlreadyExists
	}

	accPoolCreator := k.accountKeeper.GetAccount(ctx, msg.PoolCreator)
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
		PoolTypeIndex:     msg.PoolTypeIndex,
		ReserveCoinDenoms: reserveCoinDenoms,
		ReserveAccount:    reserveAcc,
		PoolCoinDenom:     PoolCoinDenom,
	}

	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(liquidityPool.PoolCoinDenom, params.InitPoolCoinMintAmount))
	if err := k.bankKeeper.SendCoins(ctx, msg.PoolCreator, liquidityPool.ReserveAccount, msg.DepositCoins); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.PoolCreator, mintPoolCoin); err != nil {
		return err
	}

	k.SetLiquidityPoolAtomic(ctx, liquidityPool)
	batch := types.NewLiquidityPoolBatch(liquidityPool.PoolId, 0)
	k.SetLiquidityPoolBatch(ctx, batch)
	// TODO: params.LiquidityPoolCreationFee logic
	// TODO: refactoring, LiquidityPoolCreationFee, check event on handler
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Get reserve Coin from the liquidity pool
func (k Keeper) GetReserveCoins(ctx sdk.Context, lp types.LiquidityPool) (reserveCoins sdk.Coins) {
	for _, denom := range lp.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, lp.ReserveAccount, denom))
	}
	return
}

func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, lp types.LiquidityPool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(lp.PoolCoinDenom)
}

func (k Keeper) DepositLiquidityPool(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	depositCoins := (sdk.Coins)(msg.DepositCoins)
	depositCoins = depositCoins.Sort()
	if err := depositCoins.Validate(); err != nil {
		return err
	}

	pool, found := k.GetLiquidityPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}
	for _, coin := range depositCoins {
		if !stringInSlice(coin.Denom, pool.ReserveCoinDenoms) {
			return types.ErrInvalidDenom
		}
	}

	if depositCoins.Len() != 2 || len(pool.ReserveCoinDenoms) != 2 {
		return types.ErrNotImplementedYet
	}

	if depositCoins.Len() != len(pool.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	coinA := depositCoins[0]
	coinB := depositCoins[1]

	//lastReserveRatio := coinA.Amount.Quo(coinB.Amount)
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	//lastReserveRatio := reserveCoins[0].Amount.Quo(reserveCoins[1].Amount)
	// TODO: handle dec to Int error, roundInt, TruncateInt to using ToDec,  MulInt
	depositableCoinA := coinA.Amount.Mul(lastReserveRatio.Ceil().RoundInt())
	//depositableCoinB := coinB.Amount.Mul(lastReserveRatio.Ceil().RoundInt())
	var inputs []banktypes.Input
	var outputs []banktypes.Output

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)

	if coinB.Amount.GT(depositableCoinA) {
		depositableCoinB := depositableCoinA
		refundAmtB := coinA.Amount.Sub(depositableCoinB)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableCoinB))))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableCoinB))))
		// refund
		if refundAmtB.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
			outputs = append(outputs, banktypes.NewOutput(msg.Depositor, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
		}
	} else if coinB.Amount.LT(depositableCoinA) {
		// TODO: handle dec to Int error, roundInt, TruncateInt to using ToDec,  MulInt
		depositableCoinA = coinB.Amount.Quo(lastReserveRatio.RoundInt())
		refundAmtA := coinA.Amount.Sub(depositableCoinA)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableCoinA))))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableCoinA))))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinB)))
		// refund
		if refundAmtA.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
			outputs = append(outputs, banktypes.NewOutput(msg.Depositor, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
		}
	} else {
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinB)))
	}

	// execute multi-send
	// TODO: check events
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
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Depositor, poolCoin); err != nil {
		return err
	}
	// TODO: add events for batch result, each err cases
	return nil
}

func (k Keeper) WithdrawLiquidityPool(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	poolCoins := (sdk.Coins)(msg.PoolCoin)
	poolCoins = poolCoins.Sort()
	if err := poolCoins.Validate(); err != nil {
		return err
	}
	if poolCoins.Len() != 1 {
		return types.ErrNumOfReserveCoin
	}
	poolCoin := poolCoins[0]

	pool, found := k.GetLiquidityPool(ctx, msg.PoolId)
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
		inputs = append(inputs, banktypes.NewInput(pool.ReserveAccount,
			sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
		outputs = append(outputs, banktypes.NewOutput(msg.Withdrawer,
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
	// TODO: add events for batch result, each err cases
	return nil
}

// TODO: testcodes
func (k Keeper) RefundDepositLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolDepositMsg) error {
	if !batchMsg.Executed || batchMsg.Succeed {
		// TODO: make Err type
		panic("can't refund not executed or succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.Depositor, batchMsg.Msg.DepositCoins)
	if err != nil {
		panic(err)
	}
	msg, found := k.GetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !found {
		panic(err)
	}
	msg.ReadyToDelete = true
	k.SetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex, msg)
	k.DeleteLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

func (k Keeper) RefundWithdrawLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolWithdrawMsg) error {
	if !batchMsg.Executed || batchMsg.Succeed {
		panic("can't refund not executed or succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.Withdrawer, batchMsg.Msg.PoolCoin)
	if err != nil {
		panic(err)
	}
	msg, found := k.GetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !found {
		panic(err)
	}
	msg.ReadyToDelete = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex, msg)
	k.DeleteLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

// TODO: WIP
func (k Keeper) RefundSwapLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolSwapMsg) error {
	if !batchMsg.Executed || batchMsg.Succeed {
		panic("can't refund not executed or succeed msg")
	}
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.SwapRequester, sdk.NewCoins(batchMsg.Msg.OfferCoin))
	if err != nil {
		panic(err)
	}
	msg, found := k.GetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !found {
		panic(err)
	}
	msg.ReadyToDelete = true
	k.SetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex, msg)
	k.DeleteLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

//func (k Keeper) FractionalRefundSwapLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolSwapMsg) error {
//	if !batchMsg.Executed {
//		panic("can't refund not executed msg")
//	}
//}