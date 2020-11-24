package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) ValidateMsgCreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
	params := k.GetParams(ctx)
	var poolType types.LiquidityPoolType

	// check poolType exist, get poolType from param
	if len(params.LiquidityPoolTypes) >= int(msg.PoolTypeIndex) {
		poolType = params.LiquidityPoolTypes[msg.PoolTypeIndex-1]
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

	denomA, denomB := types.AlphabeticalDenomPair(msg.ReserveCoinDenoms[0], msg.ReserveCoinDenoms[1])
	if denomA != msg.ReserveCoinDenoms[0] || denomB != msg.ReserveCoinDenoms[1] {
		return types.ErrBadOrderingReserveCoin
	}

	poolKey := types.GetPoolKey(msg.ReserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	_, found := k.GetLiquidityPoolByReserveAccIndex(ctx, reserveAcc)
	if found {
		return types.ErrPoolAlreadyExists
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
		//PoolId: k.GetNextLiquidityPoolIdWithUpdate(ctx),
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
		panic(err)
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, poolCreator, mintPoolCoin); err != nil {
		panic(err)
	}

	liquidityPool = k.SetLiquidityPoolAtomic(ctx, liquidityPool)
	batch := types.NewLiquidityPoolBatch(liquidityPool.PoolId, 1)
	k.SetLiquidityPoolBatch(ctx, batch)
	// TODO: params.LiquidityPoolCreationFee logic
	// TODO: refactoring, LiquidityPoolCreationFee, check event on handler

	// TODO: remove result state check, debugging
	reserveCoins := k.GetReserveCoins(ctx, liquidityPool)
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("createPool", msg, "pool", liquidityPool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	fmt.Println("createPool", msg, "pool", liquidityPool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return nil
}

// Get reserve Coin from the liquidity pool
func (k Keeper) GetReserveCoins(ctx sdk.Context, pool types.LiquidityPool) (reserveCoins sdk.Coins) {
	for _, denom := range pool.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, pool.GetReserveAccount(), denom))
	}
	// TODO: if reserveCoins.Empty(), return zero coin
	return
}

func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, pool types.LiquidityPool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(pool.PoolCoinDenom)
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
	k.SetLiquidityPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)

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

	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	depositableAmount := coinB.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
	depositableAmountA := coinA.Amount
	depositableAmountB := coinB.Amount
	var inputs []banktypes.Input
	var outputs []banktypes.Output

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	reserveAcc := pool.GetReserveAccount()
	depositor := msg.Msg.GetDepositor()

	if coinA.Amount.LT(depositableAmount) {
		depositableAmountB = coinA.Amount.ToDec().Quo(lastReserveRatio).TruncateInt()
		refundAmtB := coinB.Amount.Sub(depositableAmountB)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableAmountB))))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableAmountB))))
		// refund
		if refundAmtB.IsPositive() {
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
			outputs = append(outputs, banktypes.NewOutput(depositor, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
		}
	} else if coinA.Amount.GT(depositableAmount) {
		depositableAmountA = coinB.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
		refundAmtA := coinA.Amount.Sub(depositableAmountA)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableAmountA))))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableAmountA))))
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
	// TODO: verify only use depositableAmount?
	//reserveCoins = k.GetReserveCoins(ctx, pool)
	//reserveCoins.Sort()
	//reserveCoinA := reserveCoins[0]
	poolCoinAmt := k.GetPoolCoinTotalSupply(ctx, pool).Mul(depositableAmountA).Quo(reserveCoins[0].Amount) // TODO: coinA after executed ?
	poolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt))

	// mint pool token to Depositor
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, poolCoin); err != nil {
		panic(err)
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, poolCoin); err != nil {
		panic(err)
	}
	msg.Succeed = true
	msg.ToDelete = true
	k.SetLiquidityPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)
	// TODO: add events for batch result, each err cases

	// TODO: remove result state check, debugging
	reserveCoins = k.GetReserveCoins(ctx, pool)
	lastReserveRatio = sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("deposit", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	fmt.Println("deposit", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
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
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg)

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
		panic(err)
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins); err != nil {
		panic(err)
	}
	msg.Succeed = true
	msg.ToDelete = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg)
	// TODO: add events for batch result, each err cases

	// TODO: remove result state check, debugging
	reserveCoins = k.GetReserveCoins(ctx, pool)
	if !reserveCoins.Empty() {
		lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
		logger := k.Logger(ctx)
		logger.Info("withdraw", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
		fmt.Println("withdraw", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	}
	return nil
}

// TODO: testcodes
func (k Keeper) RefundDepositLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolDepositMsg) error {
	batchMsg, _ = k.GetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
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
	k.SetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, msg)
	k.DeleteLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

func (k Keeper) RefundWithdrawLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolWithdrawMsg) error {
	batchMsg, _ = k.GetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
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
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, msg)
	k.DeleteLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	return err
}

// TODO: WIP
//func (k Keeper) RefundSwapLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolSwapMsg) error {
//	batchMsg, _ = k.GetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
//	if !batchMsg.Executed || batchMsg.Succeed {
//		panic("can't refund not executed or succeed msg")
//	}
//	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetSwapRequester(), sdk.NewCoins(batchMsg.Msg.OfferCoin))
//	if err != nil {
//		panic(err)
//	}
//	msg, found := k.GetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
//	if !found {
//		panic(err)
//	}
//	msg.ToDelete = true
//	k.SetLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, msg)
//	k.DeleteLiquidityPoolBatchSwapMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
//	return err
//}

//func (k Keeper) FractionalRefundSwapLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolSwapMsg) error {
//	if !batchMsg.Executed {
//		panic("can't refund not executed msg")
//	}
//}

func (k Keeper) GetLiquidityPoolMetaData(ctx sdk.Context, pool types.LiquidityPool) *types.LiquidityPoolMetaData {
	totalSupply := sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
	reserveCoin := k.GetReserveCoins(ctx, pool)
	return &types.LiquidityPoolMetaData{PoolId:pool.PoolId, PoolCoinTotalSupply: totalSupply, ReserveCoins:reserveCoin}
}