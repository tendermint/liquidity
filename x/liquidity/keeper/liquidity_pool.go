package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"strconv"
)

func (k Keeper) ValidateMsgCreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
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

	reserveCoinNum := uint32(msg.DepositCoins.Len())
	if reserveCoinNum > poolType.MaxReserveCoinNum || poolType.MinReserveCoinNum > reserveCoinNum {
		return types.ErrNumOfReserveCoin
	}

	reserveCoinDenoms := make([]string, reserveCoinNum)
	for i := 0; i < int(reserveCoinNum); i++ {
		reserveCoinDenoms[i] = msg.DepositCoins.GetDenomByIndex(i)
	}

	denomA, denomB := types.AlphabeticalDenomPair(reserveCoinDenoms[0], reserveCoinDenoms[1])
	if denomA != msg.DepositCoins[0].Denom || denomB != msg.DepositCoins[1].Denom {
		return types.ErrBadOrderingReserveCoin
	}

	if denomA == denomB {
		return types.ErrEqualDenom
	}

	if err := types.ValidateReserveCoinLimit(params.ReserveCoinLimitAmount, msg.DepositCoins); err != nil {
		return err
	}

	poolKey := types.GetPoolKey(reserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	_, found := k.GetLiquidityPoolByReserveAccIndex(ctx, reserveAcc)
	if found {
		return types.ErrPoolAlreadyExists
	}
	return nil
}

// Validate logic for Liquidity pool after set or before export
func (k Keeper) ValidateLiquidityPool(ctx sdk.Context, pool *types.LiquidityPool) error {
	params := k.GetParams(ctx)
	var poolType types.LiquidityPoolType

	// check poolType exist, get poolType from param
	if len(params.LiquidityPoolTypes) >= int(pool.PoolTypeIndex) {
		poolType = params.LiquidityPoolTypes[pool.PoolTypeIndex-1]
		if poolType.PoolTypeIndex != pool.PoolTypeIndex {
			return types.ErrPoolTypeNotExists
		}
	} else {
		return types.ErrPoolTypeNotExists
	}

	if poolType.MaxReserveCoinNum > types.MaxReserveCoinNum || types.MinReserveCoinNum > poolType.MinReserveCoinNum {
		return types.ErrNumOfReserveCoin
	}

	reserveCoins := k.GetReserveCoins(ctx, *pool)
	if uint32(reserveCoins.Len()) > poolType.MaxReserveCoinNum || poolType.MinReserveCoinNum > uint32(reserveCoins.Len()) {
		return types.ErrNumOfReserveCoin
	}

	if len(pool.ReserveCoinDenoms) != reserveCoins.Len() {
		return types.ErrNumOfReserveCoin
	}
	for i, denom := range pool.ReserveCoinDenoms {
		if denom != reserveCoins[i].Denom {
			return types.ErrInvalidDenom
		}
	}

	denomA, denomB := types.AlphabeticalDenomPair(pool.ReserveCoinDenoms[0], pool.ReserveCoinDenoms[1])
	if denomA != pool.ReserveCoinDenoms[0] || denomB != pool.ReserveCoinDenoms[1] {
		return types.ErrBadOrderingReserveCoin
	}

	poolKey := types.GetPoolKey(pool.ReserveCoinDenoms, pool.PoolTypeIndex)
	poolCoin := k.GetPoolCoinTotal(ctx, *pool)
	if poolCoin.Denom != types.GetPoolCoinDenom(poolKey) {
		return types.ErrBadPoolCoinDenom
	}

	_, found := k.GetLiquidityPoolBatch(ctx, pool.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}

	return nil
}

func (k Keeper) CreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) (types.LiquidityPool, error) {
	if err := k.ValidateMsgCreateLiquidityPool(ctx, msg); err != nil {
		return types.LiquidityPool{}, err
	}
	params := k.GetParams(ctx)

	denom1, denom2 := types.AlphabeticalDenomPair(msg.DepositCoins[0].Denom, msg.DepositCoins[1].Denom)
	reserveCoinDenoms := []string{denom1, denom2}

	poolKey := types.GetPoolKey(reserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)

	poolCreator := msg.GetPoolCreator()
	accPoolCreator := k.accountKeeper.GetAccount(ctx, poolCreator)
	poolCreatorBalances := k.bankKeeper.GetAllBalances(ctx, accPoolCreator.GetAddress())

	if !poolCreatorBalances.IsAllGTE(msg.DepositCoins) {
		return types.LiquidityPool{}, types.ErrInsufficientBalance
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.LT(params.MinInitDepositToPool) {
			return types.LiquidityPool{}, types.ErrLessThanMinInitDeposit
		}
	}

	if !poolCreatorBalances.IsAllGTE(params.LiquidityPoolCreationFee.Add(msg.DepositCoins...)) {
		return types.LiquidityPool{}, types.ErrInsufficientPoolCreationFee
	}

	PoolCoinDenom := types.GetPoolCoinDenom(poolKey)

	pool := types.LiquidityPool{
		//PoolId: will set on SetLiquidityPoolAtomic
		PoolTypeIndex:         msg.PoolTypeIndex,
		ReserveCoinDenoms:     reserveCoinDenoms,
		ReserveAccountAddress: reserveAcc.String(),
		PoolCoinDenom:         PoolCoinDenom,
	}

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	// pool creation fees are collected in community pool
	// TODO: verify distr community pool is consistency safe for collected fee
	poolCreationFeePoolAcc := k.accountKeeper.GetModuleAddress(distrtypes.ModuleName)
	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, params.InitPoolCoinMintAmount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		return types.LiquidityPool{}, err
	}

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	inputs = append(inputs, banktypes.NewInput(poolCreator, params.LiquidityPoolCreationFee))
	outputs = append(outputs, banktypes.NewOutput(poolCreationFeePoolAcc, params.LiquidityPoolCreationFee))

	inputs = append(inputs, banktypes.NewInput(poolCreator, msg.DepositCoins))
	outputs = append(outputs, banktypes.NewOutput(reserveAcc, msg.DepositCoins))

	inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoin))
	outputs = append(outputs, banktypes.NewOutput(poolCreator, mintPoolCoin))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return types.LiquidityPool{}, err
	}

	pool = k.SetLiquidityPoolAtomic(ctx, pool)
	batch := types.NewLiquidityPoolBatch(pool.PoolId, 1)

	k.SetLiquidityPoolBatch(ctx, batch)

	// TODO: remove result state check, debugging
	reserveCoins := k.GetReserveCoins(ctx, pool)
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("createPool", msg, "pool", pool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return pool, nil
}

// Get reserve Coin from the liquidity pool
func (k Keeper) GetReserveCoins(ctx sdk.Context, pool types.LiquidityPool) (reserveCoins sdk.Coins) {
	reserveAcc := pool.GetReserveAccount()
	for _, denom := range pool.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, reserveAcc, denom))
	}
	return
}

// Get total supply of pool coin of the pool as sdk.Int
func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, pool types.LiquidityPool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(pool.PoolCoinDenom)
}

// Get total supply of pool coin of the pool as sdk.Coin
func (k Keeper) GetPoolCoinTotal(ctx sdk.Context, pool types.LiquidityPool) sdk.Coin {
	return sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
}

// Get meta data of the pool, containing pool coin total supply, Reserved Coins, It used for result of queries
func (k Keeper) GetPoolMetaData(ctx sdk.Context, pool types.LiquidityPool) types.LiquidityPoolMetadata {
	return types.LiquidityPoolMetadata{
		PoolId:              pool.PoolId,
		PoolCoinTotalSupply: k.GetPoolCoinTotal(ctx, pool),
		ReserveCoins:        k.GetReserveCoins(ctx, pool),
	}
}

// Get Liquidity Pool Record
func (k Keeper) GetLiquidityPoolRecord(ctx sdk.Context, pool types.LiquidityPool) (*types.LiquidityPoolRecord, bool) {
	batch, found := k.GetLiquidityPoolBatch(ctx, pool.PoolId)
	if !found {
		return nil, found
	}
	return &types.LiquidityPoolRecord{
		LiquidityPool:         pool,
		LiquidityPoolMetadata: k.GetPoolMetaData(ctx, pool),
		LiquidityPoolBatch:    batch,
		BatchPoolDepositMsgs:  k.GetAllLiquidityPoolBatchDepositMsgs(ctx, batch),
		BatchPoolWithdrawMsgs: k.GetAllLiquidityPoolBatchWithdrawMsgs(ctx, batch),
		BatchPoolSwapMsgs:     k.GetAllLiquidityPoolBatchSwapMsgs(ctx, batch),
	}, true
}
func (k Keeper) SetLiquidityPoolRecord(ctx sdk.Context, record *types.LiquidityPoolRecord) {
	k.SetLiquidityPoolAtomic(ctx, record.LiquidityPool)
	//k.SetLiquidityPool(ctx, record.LiquidityPool)
	//k.SetLiquidityPoolByReserveAccIndex(ctx, record.LiquidityPool)
	k.GetNextBatchIndexWithUpdate(ctx, record.LiquidityPool.PoolId)
	record.LiquidityPoolBatch.BeginHeight = ctx.BlockHeight()
	k.SetLiquidityPoolBatch(ctx, record.LiquidityPoolBatch)
	k.SetLiquidityPoolBatchDepositMsgs(ctx, record.LiquidityPool.PoolId, record.BatchPoolDepositMsgs)
	k.SetLiquidityPoolBatchWithdrawMsgs(ctx, record.LiquidityPool.PoolId, record.BatchPoolWithdrawMsgs)
	k.SetLiquidityPoolBatchSwapMsgs(ctx, record.LiquidityPool.PoolId, record.BatchPoolSwapMsgs)
}

func (k Keeper) ValidateMsgDepositLiquidityPool(ctx sdk.Context, msg types.MsgDepositToLiquidityPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	pool, found := k.GetLiquidityPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	if msg.DepositCoins.Len() != len(pool.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	params := k.GetParams(ctx)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	if err := types.ValidateReserveCoinLimit(params.ReserveCoinLimitAmount, reserveCoins.Add(msg.DepositCoins...)); err != nil {
		return err
	}
	// TODO: validate msgIndex

	denomA, denomB := types.AlphabeticalDenomPair(msg.DepositCoins[0].Denom, msg.DepositCoins[1].Denom)
	if denomA != pool.ReserveCoinDenoms[0] || denomB != pool.ReserveCoinDenoms[1] {
		return types.ErrNotMatchedReserveCoin
	}
	return nil
}

func (k Keeper) DepositLiquidityPool(ctx sdk.Context, msg types.BatchPoolDepositMsg, batch types.LiquidityPoolBatch) error {
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

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	reserveAcc := pool.GetReserveAccount()
	depositor := msg.Msg.GetDepositor()
	params := k.GetParams(ctx)

	reserveCoins := k.GetReserveCoins(ctx, pool)

	// case of reserve coins has run out, ReinitializePool
	if reserveCoins.IsZero() {
		mintPoolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, params.InitPoolCoinMintAmount))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
			return err
		}
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, msg.Msg.DepositCoins))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, msg.Msg.DepositCoins))

		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoin))
		outputs = append(outputs, banktypes.NewOutput(depositor, mintPoolCoin))

		// execute multi-send
		if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
			return err
		}

		msg.Succeeded = true
		msg.ToBeDeleted = true
		k.SetLiquidityPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)

		// TODO: remove result state check, debugging
		reserveCoins := k.GetReserveCoins(ctx, pool)
		lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
		logger := k.Logger(ctx)
		logger.Info("ReinitializePool", msg, "pool", pool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
		return nil
	} else if reserveCoins.Len() != msg.Msg.DepositCoins.Len() {
		return types.ErrNumOfReserveCoin
	}

	reserveCoins.Sort()

	coinA := depositCoins[0]
	coinB := depositCoins[1]

	// Decimal Error, divide the Int coin amount by the Decimal Rate and erase the decimal point to deposit a lower value
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	depositableAmount := coinB.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
	depositableAmountA := coinA.Amount
	depositableAmountB := coinB.Amount
	acceptedCoins := sdk.NewCoins()
	refundedCoins := sdk.NewCoins()

	if coinA.Amount.LT(depositableAmount) {
		depositableAmountB = coinA.Amount.ToDec().Quo(lastReserveRatio).TruncateInt()
		refundAmtB := coinB.Amount.Sub(depositableAmountB)
		acceptedCoins = sdk.NewCoins(
			coinA,
			sdk.NewCoin(coinB.Denom, depositableAmountB))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, acceptedCoins))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, acceptedCoins))
		// refund
		if refundAmtB.IsPositive() {
			refundedCoins = sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, refundedCoins))
			outputs = append(outputs, banktypes.NewOutput(depositor, refundedCoins))
		}
	} else if coinA.Amount.GT(depositableAmount) {
		depositableAmountA = coinB.Amount.ToDec().Mul(lastReserveRatio).TruncateInt()
		refundAmtA := coinA.Amount.Sub(depositableAmountA)
		acceptedCoins = sdk.NewCoins(
			coinB,
			sdk.NewCoin(coinA.Denom, depositableAmountA))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, acceptedCoins))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, acceptedCoins))
		// refund
		if refundAmtA.IsPositive() {
			refundedCoins = sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))
			inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, refundedCoins))
			outputs = append(outputs, banktypes.NewOutput(depositor, refundedCoins))
		}
	} else {
		acceptedCoins = sdk.NewCoins(coinA, coinB)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, acceptedCoins))
		outputs = append(outputs, banktypes.NewOutput(reserveAcc, acceptedCoins))
	}

	// calculate pool token mint amount
	poolCoinAmt := k.GetPoolCoinTotalSupply(ctx, pool).Mul(depositableAmountA).Quo(reserveCoins[0].Amount) // TODO: coinA after executed ?
	mintPoolCoin := sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt)
	mintPoolCoins := sdk.NewCoins(mintPoolCoin)

	// mint pool token to Depositor
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoins); err != nil {
		panic(err)
	}
	inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoins))
	outputs = append(outputs, banktypes.NewOutput(depositor, mintPoolCoins))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	msg.Succeeded = true
	msg.ToBeDeleted = true
	k.SetLiquidityPoolBatchDepositMsg(ctx, msg.Msg.PoolId, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositToLiquidityPool,
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(pool.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(batch.BatchIndex, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(msg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueDepositor, depositor.String()),
			sdk.NewAttribute(types.AttributeValueAcceptedCoins, acceptedCoins.String()),
			sdk.NewAttribute(types.AttributeValueRefundedCoins, refundedCoins.String()),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, mintPoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, mintPoolCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueSuccess, types.Success),
		))

	// TODO: remove result state check, debugging
	reserveCoins = k.GetReserveCoins(ctx, pool)
	lastReserveRatio = sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("deposit", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return nil
}

func (k Keeper) ValidateMsgWithdrawLiquidityPool(ctx sdk.Context, msg types.MsgWithdrawFromLiquidityPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	pool, found := k.GetLiquidityPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	if msg.PoolCoin.Denom != pool.PoolCoinDenom {
		return types.ErrBadPoolCoinDenom
	}

	poolCoinTotalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	if msg.PoolCoin.Amount.GT(poolCoinTotalSupply) {
		return types.ErrBadPoolCoinAmount
	}
	return nil
}

func (k Keeper) ValidateMsgSwap(ctx sdk.Context, msg types.MsgSwap) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	pool, found := k.GetLiquidityPool(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	denomA, denomB := types.AlphabeticalDenomPair(msg.OfferCoin.Denom, msg.DemandCoinDenom)
	if denomA != pool.ReserveCoinDenoms[0] || denomB != pool.ReserveCoinDenoms[1] {
		return types.ErrNotMatchedReserveCoin
	}

	params := k.GetParams(ctx)

	// can not exceed max order ratio  of reserve coins that can be ordered at a order
	reserveCoinAmt := k.GetReserveCoins(ctx, pool).AmountOf(msg.OfferCoin.Denom)
	// Decimal Error, Multiply the Int coin amount by the Decimal Rate and erase the decimal point to order a lower value
	maximumOrderableAmt := reserveCoinAmt.ToDec().Mul(params.MaxOrderAmountRatio).TruncateInt()
	if msg.OfferCoin.Amount.GT(maximumOrderableAmt) {
		return types.ErrExceededMaxOrderable
	}
	// TODO: half-half invariant check, need to after msg created
	if msg.OfferCoinFee.Denom != msg.OfferCoin.Denom {
		return types.ErrBadOfferCoinFee
	}
	// TODO: half-half fee refund when over
	if !msg.OfferCoinFee.Equal(types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)) {
		return types.ErrBadOfferCoinFee
	}

	return nil
}

func (k Keeper) WithdrawLiquidityPool(ctx sdk.Context, msg types.BatchPoolWithdrawMsg, batch types.LiquidityPoolBatch) error {
	msg.Executed = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg)

	if err := k.ValidateMsgWithdrawLiquidityPool(ctx, *msg.Msg); err != nil {
		return err
	}
	// TODO: validate reserveCoin balance

	poolCoins := sdk.NewCoins(msg.Msg.PoolCoin)

	pool, found := k.GetLiquidityPool(ctx, msg.Msg.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}

	totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	reserveAcc := pool.GetReserveAccount()
	withdrawer := msg.Msg.GetWithdrawer()

	params := k.GetParams(ctx)
	withdrawProportion := sdk.OneDec().Sub(params.WithdrawFeeRate)
	withdrawCoins := sdk.NewCoins()

	for _, reserveCoin := range reserveCoins {
		// Decimal Error, Multiply the Int coin amount by the Decimal proportion and erase the decimal point to withdraw a conservative value
		withdrawAmt := reserveCoin.Amount.Mul(msg.Msg.PoolCoin.Amount).ToDec().Mul(withdrawProportion).TruncateInt().Quo(totalSupply)
		withdrawCoins = withdrawCoins.Add(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))
	}
	if withdrawCoins.IsValid() {
		inputs = append(inputs, banktypes.NewInput(reserveAcc, withdrawCoins))
		outputs = append(outputs, banktypes.NewOutput(withdrawer, withdrawCoins))
	}

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	// burn the escrowed pool coins
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins); err != nil {
		panic(err)
	}

	msg.Succeeded = true
	msg.ToBeDeleted = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, msg.Msg.PoolId, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawFromLiquidityPool,
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(pool.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(batch.BatchIndex, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(msg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueWithdrawer, withdrawer.String()),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, msg.Msg.PoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, msg.Msg.PoolCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueWithdrawCoins, withdrawCoins.String()),
			sdk.NewAttribute(types.AttributeValueSuccess, types.Success),
		))

	// TODO: remove result state check, debugging
	reserveCoins = k.GetReserveCoins(ctx, pool)

	var lastReserveRatio sdk.Dec
	if reserveCoins.Empty() {
		lastReserveRatio = sdk.ZeroDec()
	} else {
		lastReserveRatio = sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	}
	logger := k.Logger(ctx)
	logger.Info("withdraw", msg, "pool", pool, "inputs", inputs, "outputs", outputs, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return nil
}

func (k Keeper) RefundDepositLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolDepositMsg, batch types.LiquidityPoolBatch) error {
	batchMsg, _ = k.GetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !batchMsg.Executed || batchMsg.Succeeded {
		panic("can't refund not executed or already succeed msg")
	}
	pool, _ := k.GetLiquidityPool(ctx, batchMsg.Msg.PoolId)
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetDepositor(), batchMsg.Msg.DepositCoins)
	if err != nil {
		panic(err)
	}
	// not delete now, set ToBeDeleted true for delete on next block beginblock
	batchMsg.ToBeDeleted = true
	k.SetLiquidityPoolBatchDepositMsg(ctx, batchMsg.Msg.PoolId, batchMsg)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositToLiquidityPool,
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(pool.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(batch.BatchIndex, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueDepositor, batchMsg.Msg.GetDepositor().String()),
			sdk.NewAttribute(types.AttributeValueAcceptedCoins, sdk.NewCoins().String()),
			sdk.NewAttribute(types.AttributeValueRefundedCoins, batchMsg.Msg.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeValueSuccess, types.Failure),
		))
	return err
}

func (k Keeper) RefundWithdrawLiquidityPool(ctx sdk.Context, batchMsg types.BatchPoolWithdrawMsg, batch types.LiquidityPoolBatch) error {
	batchMsg, _ = k.GetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !batchMsg.Executed || batchMsg.Succeeded {
		panic("can't refund not executed or already succeed msg")
	}
	pool, _ := k.GetLiquidityPool(ctx, batchMsg.Msg.PoolId)
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetWithdrawer(), sdk.NewCoins(batchMsg.Msg.PoolCoin))
	if err != nil {
		panic(err)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawFromLiquidityPool,
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(pool.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(batch.BatchIndex, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueWithdrawer, batchMsg.Msg.GetWithdrawer().String()),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, batchMsg.Msg.PoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, batchMsg.Msg.PoolCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueWithdrawCoins, sdk.NewCoins().String()),
			sdk.NewAttribute(types.AttributeValueSuccess, types.Failure),
		))

	// not delete now, set ToBeDeleted true for delete on next block beginblock
	batchMsg.ToBeDeleted = true
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, batchMsg.Msg.PoolId, batchMsg)
	return err
}

// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapLiquidityPool
func (k Keeper) TransactAndRefundSwapLiquidityPool(ctx sdk.Context, batchMsgs []*types.BatchPoolSwapMsg,
	matchResultMap map[uint64]types.MatchResult, pool types.LiquidityPool, batchResult types.BatchResult) error {

	var inputs []banktypes.Input
	var outputs []banktypes.Output
	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	poolReserveAcc := pool.GetReserveAccount()
	batch, found := k.GetLiquidityPoolBatch(ctx, pool.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	for _, batchMsg := range batchMsgs {
		// TODO: make Validate function to batchMsg
		if !batchMsg.Executed && batchMsg.Succeeded {
			panic("can't refund not executed with succeed msg")
		}
		if pool.PoolId != batchMsg.Msg.PoolId {
			panic("broken msg pool consistency")
		}

		// full matched, fractional matched
		if msgAfter, ok := matchResultMap[batchMsg.MsgIndex]; ok {
			if batchMsg.MsgIndex != msgAfter.OrderMsgIndex {
				panic("broken msg consistency")
			}

			if (*msgAfter.BatchMsg) != (*batchMsg) {
				panic("broken msg consistency")
			}

			// TODO: fix invariant for half-half fee
			//if msgAfter.TransactedCoinAmt.Sub(msgAfter.OfferCoinFeeAmt).IsNegative() ||
			//	msgAfter.OfferCoinFeeAmt.GT(msgAfter.TransactedCoinAmt) {
			//	panic("fee over offer")
			//}

			// fractional match, but expired order case
			if batchMsg.RemainingOfferCoin.IsPositive() {
				// not to delete, but expired case
				if !batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight <= ctx.BlockHeight() {
					panic("impossible case")
					// TODO: set to Delete
				} else if !batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight > ctx.BlockHeight() {
					// fractional matched, to be remaining order, not refund, only transact fractional exchange amt
					// Add transacted coins to multisend
					// TODO: coverage
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt))))
					inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt)))))
					outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt)))))

					// Add swap offer coin fee to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt))))

					// Add swap exchanged coin fee to multisend, It cause temporary insufficient funds when InputOutputCoins, skip offsetting input, output
					//inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
					//	sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.ExchangedCoinFeeAmt))))
					//outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					//	sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.ExchangedCoinFeeAmt))))

					batchMsg.Succeeded = true

				} else if batchMsg.ToBeDeleted || batchMsg.OrderExpiryHeight == ctx.BlockHeight() {
					// fractional matched, but expired order, transact with refund remaining offer coin

					// Add transacted coins to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt))))
					inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt)))))
					outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt)))))

					// Add swap offer coin fee to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt))))

					// Add swap exchanged coin fee to multisend, It cause temporary insufficient funds when InputOutputCoins, skip offsetting input, output
					//inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
					//	sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.ExchangedCoinFeeAmt))))
					//outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					//	sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.ExchangedCoinFeeAmt))))

					// refund remaining coins
					if input, output, err := k.ReleaseEscrowForMultiSend(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(batchMsg.RemainingOfferCoin)); err != nil {
						panic(err)
					} else {
						inputs = append(inputs, input)
						outputs = append(outputs, output)
					}
					batchMsg.Succeeded = true
					batchMsg.ToBeDeleted = true
				} else {
					panic("impossible case")
					// TODO: set to Delete
				}
			} else if batchMsg.RemainingOfferCoin.IsZero() {
				// full matched case, Add transacted coins to multisend
				inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt))))
				outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt))))
				inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt)))))
				outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt)))))

				// Add swap offer coin fee to multisend
				inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt))))
				outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt))))

				// Add swap exchanged coin fee to multisend, It cause temporary insufficient funds when InputOutputCoins, skip offsetting input, output
				//inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
				//	sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.ExchangedCoinFeeAmt))))
				//outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
				//	sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.ExchangedCoinFeeAmt))))

				batchMsg.Succeeded = true
				batchMsg.ToBeDeleted = true
			} else {
				panic("impossible case")
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeSwapTransacted,
					sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(pool.PoolId, 10)),
					sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(batch.BatchIndex, 10)),
					sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(msgAfter.BatchMsg.MsgIndex, 10)),
					sdk.NewAttribute(types.AttributeValueSwapRequester, msgAfter.BatchMsg.Msg.GetSwapRequester().String()),
					sdk.NewAttribute(types.AttributeValueSwapType, strconv.FormatUint(uint64(msgAfter.BatchMsg.Msg.SwapType), 10)),
					sdk.NewAttribute(types.AttributeValueOfferCoinDenom, msgAfter.BatchMsg.Msg.OfferCoin.Denom),
					sdk.NewAttribute(types.AttributeValueOfferCoinAmount, msgAfter.BatchMsg.Msg.OfferCoin.Amount.String()),
					sdk.NewAttribute(types.AttributeValueOrderPrice, msgAfter.BatchMsg.Msg.OrderPrice.String()),
					sdk.NewAttribute(types.AttributeValueSwapPrice, batchResult.SwapPrice.String()),
					sdk.NewAttribute(types.AttributeValueTransactedCoinAmount, msgAfter.TransactedCoinAmt.String()),
					sdk.NewAttribute(types.AttributeValueRemainingOfferCoinAmount, msgAfter.BatchMsg.RemainingOfferCoin.Amount.String()),
					sdk.NewAttribute(types.AttributeValueExchangedOfferCoinAmount, msgAfter.BatchMsg.ExchangedOfferCoin.Amount.String()),
					sdk.NewAttribute(types.AttributeValueOfferCoinFeeAmount, msgAfter.OfferCoinFeeAmt.String()),
					sdk.NewAttribute(types.AttributeValueOfferCoinFeeReserveAmount, msgAfter.BatchMsg.OfferCoinFeeReserve.Amount.String()),
					sdk.NewAttribute(types.AttributeValueOrderExpiryHeight, strconv.FormatInt(msgAfter.OrderExpiryHeight, 10)),
					sdk.NewAttribute(types.AttributeValueSuccess, types.Success),
				))
		} else {
			// not matched, remaining
			if !batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight > ctx.BlockHeight() {
				// have fractional matching history, not matched and expired, remaining refund
				// refund remaining coins
				// TODO: coverage
				if input, output, err := k.ReleaseEscrowForMultiSend(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(batchMsg.RemainingOfferCoin.Add(batchMsg.OfferCoinFeeReserve))); err != nil {
					panic(err)
				} else {
					inputs = append(inputs, input)
					outputs = append(outputs, output)
				}

				batchMsg.Succeeded = false
				batchMsg.ToBeDeleted = true

			} else if batchMsg.ToBeDeleted && batchMsg.OrderExpiryHeight == ctx.BlockHeight() {
				// not matched and expired, remaining refund
				// refund remaining coins
				if input, output, err := k.ReleaseEscrowForMultiSend(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(batchMsg.RemainingOfferCoin.Add(batchMsg.OfferCoinFeeReserve))); err != nil {
					panic(err)
				} else {
					inputs = append(inputs, input)
					outputs = append(outputs, output)
				}

				batchMsg.Succeeded = false
				batchMsg.ToBeDeleted = true

			} else {
				panic("impossible case")
			}
		}
	}
	// remove zero coins
	newI := 0
	for _, i := range inputs {
		if !i.Coins.IsValid() {
			i.Coins = sdk.NewCoins(i.Coins...) // for sanitizeCoins, remove zero coin
		}
		if !i.Coins.Empty() {
			inputs[newI] = i
			newI++
		}
	}
	inputs = inputs[:newI]
	newI = 0
	for _, i := range outputs {
		if !i.Coins.IsValid() {
			i.Coins = sdk.NewCoins(i.Coins...) // for sanitizeCoins, remove zero coin
		}
		if !i.Coins.Empty() {
			outputs[newI] = i
			newI++
		}
	}
	outputs = outputs[:newI]
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
	k.SetLiquidityPoolBatchSwapMsgPointers(ctx, pool.PoolId, batchMsgs)
	return nil
}

//func (k Keeper) GetPoolMetaData(ctx sdk.Context, pool types.LiquidityPool) *types.LiquidityPoolMetadata {
//	totalSupply := sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
//	reserveCoin := k.GetReserveCoins(ctx, pool).Sort()
//	return &types.LiquidityPoolMetadata{PoolId: pool.PoolId, PoolCoinTotalSupply: totalSupply, ReserveCoins: reserveCoin}
//}

func (k Keeper) ValidateLiquidityPoolMetadata(ctx sdk.Context, pool *types.LiquidityPool, metaData *types.LiquidityPoolMetadata) error {
	if err := metaData.ReserveCoins.Validate(); err != nil {
		return err
	}
	if !metaData.ReserveCoins.IsEqual(k.GetReserveCoins(ctx, *pool)) {
		return types.ErrNumOfReserveCoin
	}
	if !metaData.PoolCoinTotalSupply.IsEqual(sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, *pool))) {
		return types.ErrBadPoolCoinAmount
	}
	return nil
}

// Validate Liquidity Pool Record after init or after export
func (k Keeper) ValidateLiquidityPoolRecord(ctx sdk.Context, record *types.LiquidityPoolRecord) error {
	params := k.GetParams(ctx)
	if err := params.Validate(); err != nil {
		return err
	}

	// validate liquidity pool
	if err := k.ValidateLiquidityPool(ctx, &record.LiquidityPool); err != nil {
		return err
	}

	// validate metadata
	if err := k.ValidateLiquidityPoolMetadata(ctx, &record.LiquidityPool, &record.LiquidityPoolMetadata); err != nil {
		return err
	}

	// validate each msgs with batch state

	if len(record.BatchPoolDepositMsgs) != 0 && record.LiquidityPoolBatch.DepositMsgIndex != record.BatchPoolDepositMsgs[len(record.BatchPoolDepositMsgs)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}
	if len(record.BatchPoolWithdrawMsgs) != 0 && record.LiquidityPoolBatch.WithdrawMsgIndex != record.BatchPoolWithdrawMsgs[len(record.BatchPoolWithdrawMsgs)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}
	if len(record.BatchPoolSwapMsgs) != 0 && record.LiquidityPoolBatch.SwapMsgIndex != record.BatchPoolSwapMsgs[len(record.BatchPoolSwapMsgs)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}

	// TODO: add verify of escrow amount and poolcoin amount with compare to remaining msgs
	return nil
}
