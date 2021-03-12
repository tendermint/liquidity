package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) ValidateMsgCreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	params := k.GetParams(ctx)
	var poolType types.PoolType

	// check poolType exist, get poolType from param
	if len(params.PoolTypes) >= int(msg.PoolTypeIndex) {
		poolType = params.PoolTypes[msg.PoolTypeIndex-1]
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

	poolKey := types.PoolName(reserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	_, found := k.GetPoolByReserveAccIndex(ctx, reserveAcc)
	if found {
		return types.ErrPoolAlreadyExists
	}
	return nil
}

// Validate logic for Liquidity pool after set or before export
func (k Keeper) ValidatePool(ctx sdk.Context, pool *types.Pool) error {
	params := k.GetParams(ctx)
	var poolType types.PoolType

	// check poolType exist, get poolType from param
	if len(params.PoolTypes) >= int(pool.PoolTypeIndex) {
		poolType = params.PoolTypes[pool.PoolTypeIndex-1]
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

	poolKey := types.PoolName(pool.ReserveCoinDenoms, pool.PoolTypeIndex)
	poolCoin := k.GetPoolCoinTotal(ctx, *pool)
	if poolCoin.Denom != types.GetPoolCoinDenom(poolKey) {
		return types.ErrBadPoolCoinDenom
	}

	_, found := k.GetPoolBatch(ctx, pool.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}

	return nil
}

func (k Keeper) CreatePool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) (types.Pool, error) {
	if err := k.ValidateMsgCreateLiquidityPool(ctx, msg); err != nil {
		return types.Pool{}, err
	}
	params := k.GetParams(ctx)

	denom1, denom2 := types.AlphabeticalDenomPair(msg.DepositCoins[0].Denom, msg.DepositCoins[1].Denom)
	reserveCoinDenoms := []string{denom1, denom2}

	poolKey := types.PoolName(reserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)

	poolCreator := msg.GetPoolCreator()
	accPoolCreator := k.accountKeeper.GetAccount(ctx, poolCreator)
	poolCreatorBalances := k.bankKeeper.GetAllBalances(ctx, accPoolCreator.GetAddress())

	if !poolCreatorBalances.IsAllGTE(msg.DepositCoins) {
		return types.Pool{}, types.ErrInsufficientBalance
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.LT(params.MinInitDepositToPool) {
			return types.Pool{}, types.ErrLessThanMinInitDeposit
		}
	}

	if !poolCreatorBalances.IsAllGTE(params.LiquidityPoolCreationFee.Add(msg.DepositCoins...)) {
		return types.Pool{}, types.ErrInsufficientPoolCreationFee
	}

	PoolCoinDenom := types.GetPoolCoinDenom(poolKey)

	pool := types.Pool{
		//PoolId: will set on SetPoolAtomic
		PoolTypeIndex:         msg.PoolTypeIndex,
		ReserveCoinDenoms:     reserveCoinDenoms,
		ReserveAccountAddress: reserveAcc.String(),
		PoolCoinDenom:         PoolCoinDenom,
	}

	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, params.InitPoolCoinMintAmount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		return types.Pool{}, err
	}

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	inputs = append(inputs, banktypes.NewInput(poolCreator, msg.DepositCoins))
	outputs = append(outputs, banktypes.NewOutput(reserveAcc, msg.DepositCoins))

	inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, mintPoolCoin))
	outputs = append(outputs, banktypes.NewOutput(poolCreator, mintPoolCoin))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return types.Pool{}, err
	}

	// pool creation fees are collected in community pool
	if err := k.distrKeeper.FundCommunityPool(ctx, params.LiquidityPoolCreationFee, poolCreator); err != nil {
		return types.Pool{}, err
	}

	pool = k.SetPoolAtomic(ctx, pool)
	batch := types.NewLiquidityPoolBatch(pool.PoolId, 1)

	k.SetPoolBatch(ctx, batch)

	// TODO: remove result state check, debugging
	reserveCoins := k.GetReserveCoins(ctx, pool)
	lastReserveRatio := sdk.NewDecFromInt(reserveCoins[0].Amount).Quo(sdk.NewDecFromInt(reserveCoins[1].Amount))
	logger := k.Logger(ctx)
	logger.Info("createPool", msg, "pool", pool, "reserveCoins", reserveCoins, "lastReserveRatio", lastReserveRatio)
	return pool, nil
}

// Get reserve Coin from the liquidity pool
func (k Keeper) GetReserveCoins(ctx sdk.Context, pool types.Pool) (reserveCoins sdk.Coins) {
	reserveAcc := pool.GetReserveAccount()
	for _, denom := range pool.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, reserveAcc, denom))
	}
	return
}

// Get total supply of pool coin of the pool as sdk.Int
func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, pool types.Pool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(pool.PoolCoinDenom)
}

// Get total supply of pool coin of the pool as sdk.Coin
func (k Keeper) GetPoolCoinTotal(ctx sdk.Context, pool types.Pool) sdk.Coin {
	return sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
}

// Get meta data of the pool, containing pool coin total supply, Reserved Coins, It used for result of queries
func (k Keeper) GetPoolMetaData(ctx sdk.Context, pool types.Pool) types.PoolMetadata {
	return types.PoolMetadata{
		PoolId:              pool.PoolId,
		PoolCoinTotalSupply: k.GetPoolCoinTotal(ctx, pool),
		ReserveCoins:        k.GetReserveCoins(ctx, pool),
	}
}

// Get Liquidity Pool Record
func (k Keeper) GetPoolRecord(ctx sdk.Context, pool types.Pool) (*types.PoolRecord, bool) {
	batch, found := k.GetPoolBatch(ctx, pool.PoolId)
	if !found {
		return nil, found
	}
	return &types.PoolRecord{
		Pool:              pool,
		PoolMetadata:      k.GetPoolMetaData(ctx, pool),
		PoolBatch:         batch,
		DepositMsgStates:  k.GetAllPoolBatchDepositMsgs(ctx, batch),
		WithdrawMsgStates: k.GetAllPoolBatchWithdrawMsgStates(ctx, batch),
		SwapMsgStates:     k.GetAllPoolBatchSwapMsgStates(ctx, batch),
	}, true
}

func (k Keeper) SetPoolRecord(ctx sdk.Context, record *types.PoolRecord) {
	k.SetPoolAtomic(ctx, record.Pool)
	//k.SetPool(ctx, record.Pool)
	//k.SetPoolByReserveAccIndex(ctx, record.Pool)
	k.GetNextPoolBatchIndexWithUpdate(ctx, record.Pool.PoolId)
	record.PoolBatch.BeginHeight = ctx.BlockHeight()
	k.SetPoolBatch(ctx, record.PoolBatch)
	k.SetPoolBatchDepositMsgStates(ctx, record.Pool.PoolId, record.DepositMsgStates)
	k.SetPoolBatchWithdrawMsgStates(ctx, record.Pool.PoolId, record.WithdrawMsgStates)
	k.SetPoolBatchSwapMsgStates(ctx, record.Pool.PoolId, record.SwapMsgStates)
}

func (k Keeper) ValidateMsgDepositLiquidityPool(ctx sdk.Context, msg types.MsgDepositToLiquidityPool) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	pool, found := k.GetPool(ctx, msg.PoolId)
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

func (k Keeper) DepositLiquidityPool(ctx sdk.Context, msg types.DepositMsgState, batch types.PoolBatch) error {
	msg.Executed = true
	k.SetPoolBatchDepositMsgState(ctx, msg.Msg.PoolId, msg)

	if err := k.ValidateMsgDepositLiquidityPool(ctx, *msg.Msg); err != nil {
		return err
	}

	depositCoins := msg.Msg.DepositCoins.Sort()

	pool, found := k.GetPool(ctx, msg.Msg.PoolId)
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
		k.SetPoolBatchDepositMsgState(ctx, msg.Msg.PoolId, msg)

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
	k.SetPoolBatchDepositMsgState(ctx, msg.Msg.PoolId, msg)

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
	pool, found := k.GetPool(ctx, msg.PoolId)
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
	pool, found := k.GetPool(ctx, msg.PoolId)
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

func (k Keeper) WithdrawLiquidityPool(ctx sdk.Context, msg types.WithdrawMsgState, batch types.PoolBatch) error {
	msg.Executed = true
	k.SetPoolBatchWithdrawMsgState(ctx, msg.Msg.PoolId, msg)

	if err := k.ValidateMsgWithdrawLiquidityPool(ctx, *msg.Msg); err != nil {
		return err
	}
	// TODO: validate reserveCoin balance

	poolCoins := sdk.NewCoins(msg.Msg.PoolCoin)

	pool, found := k.GetPool(ctx, msg.Msg.PoolId)
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
	k.SetPoolBatchWithdrawMsgState(ctx, msg.Msg.PoolId, msg)

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

func (k Keeper) RefundDepositLiquidityPool(ctx sdk.Context, batchMsg types.DepositMsgState, batch types.PoolBatch) error {
	batchMsg, _ = k.GetPoolBatchDepositMsgState(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !batchMsg.Executed || batchMsg.Succeeded {
		panic("can't refund not executed or already succeed msg")
	}
	pool, _ := k.GetPool(ctx, batchMsg.Msg.PoolId)
	err := k.ReleaseEscrow(ctx, batchMsg.Msg.GetDepositor(), batchMsg.Msg.DepositCoins)
	if err != nil {
		panic(err)
	}
	// not delete now, set ToBeDeleted true for delete on next block beginblock
	batchMsg.ToBeDeleted = true
	k.SetPoolBatchDepositMsgState(ctx, batchMsg.Msg.PoolId, batchMsg)
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

func (k Keeper) RefundWithdrawLiquidityPool(ctx sdk.Context, batchMsg types.WithdrawMsgState, batch types.PoolBatch) error {
	batchMsg, _ = k.GetPoolBatchWithdrawMsgState(ctx, batchMsg.Msg.PoolId, batchMsg.MsgIndex)
	if !batchMsg.Executed || batchMsg.Succeeded {
		panic("can't refund not executed or already succeed msg")
	}
	pool, _ := k.GetPool(ctx, batchMsg.Msg.PoolId)
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
	k.SetPoolBatchWithdrawMsgState(ctx, batchMsg.Msg.PoolId, batchMsg)
	return err
}

// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapLiquidityPool
func (k Keeper) TransactAndRefundSwapLiquidityPool(ctx sdk.Context, batchMsgs []*types.SwapMsgState,
	matchResultMap map[uint64]types.MatchResult, pool types.Pool, batchResult types.BatchResult) error {

	var inputs []banktypes.Input
	var outputs []banktypes.Output
	batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	poolReserveAcc := pool.GetReserveAccount()
	batch, found := k.GetPoolBatch(ctx, pool.PoolId)
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
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.TruncateInt()))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.TruncateInt()))))
					inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt).TruncateInt()))))
					outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt).TruncateInt()))))

					// Add swap offer coin fee to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt.TruncateInt()))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt.TruncateInt()))))

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
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.TruncateInt()))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.TruncateInt()))))
					inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt).TruncateInt()))))
					outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
						sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt).TruncateInt()))))

					// Add swap offer coin fee to multisend
					inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt.TruncateInt()))))
					outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
						sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt.TruncateInt()))))

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
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.TruncateInt()))))
				outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.ExchangedOfferCoin.Denom, msgAfter.TransactedCoinAmt.TruncateInt()))))
				inputs = append(inputs, banktypes.NewInput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt).TruncateInt()))))
				outputs = append(outputs, banktypes.NewOutput(batchMsg.Msg.GetSwapRequester(),
					sdk.NewCoins(sdk.NewCoin(batchMsg.Msg.DemandCoinDenom, msgAfter.ExchangedDemandCoinAmt.Sub(msgAfter.ExchangedCoinFeeAmt).TruncateInt()))))

				// Add swap offer coin fee to multisend
				inputs = append(inputs, banktypes.NewInput(batchEscrowAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt.TruncateInt()))))
				outputs = append(outputs, banktypes.NewOutput(poolReserveAcc,
					sdk.NewCoins(sdk.NewCoin(batchMsg.OfferCoinFeeReserve.Denom, msgAfter.OfferCoinFeeAmt.TruncateInt()))))

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
	k.SetPoolBatchSwapMsgStatesByPointer(ctx, pool.PoolId, batchMsgs)
	return nil
}

//func (k Keeper) GetPoolMetaData(ctx sdk.Context, pool types.Pool) *types.PoolMetadata {
//	totalSupply := sdk.NewCoin(pool.PoolCoinDenom, k.GetPoolCoinTotalSupply(ctx, pool))
//	reserveCoin := k.GetReserveCoins(ctx, pool).Sort()
//	return &types.PoolMetadata{PoolId: pool.PoolId, PoolCoinTotalSupply: totalSupply, ReserveCoins: reserveCoin}
//}

func (k Keeper) ValidatePoolMetadata(ctx sdk.Context, pool *types.Pool, metaData *types.PoolMetadata) error {
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
func (k Keeper) ValidatePoolRecord(ctx sdk.Context, record *types.PoolRecord) error {
	params := k.GetParams(ctx)
	if err := params.Validate(); err != nil {
		return err
	}

	// validate liquidity pool
	if err := k.ValidatePool(ctx, &record.Pool); err != nil {
		return err
	}

	// validate metadata
	if err := k.ValidatePoolMetadata(ctx, &record.Pool, &record.PoolMetadata); err != nil {
		return err
	}

	// validate each msgs with batch state

	if len(record.DepositMsgStates) != 0 && record.PoolBatch.DepositMsgIndex != record.DepositMsgStates[len(record.DepositMsgStates)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}
	if len(record.WithdrawMsgStates) != 0 && record.PoolBatch.WithdrawMsgIndex != record.WithdrawMsgStates[len(record.WithdrawMsgStates)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}
	if len(record.SwapMsgStates) != 0 && record.PoolBatch.SwapMsgIndex != record.SwapMsgStates[len(record.SwapMsgStates)-1].MsgIndex+1 {
		return types.ErrBadBatchMsgIndex
	}

	// TODO: add verify of escrow amount and poolcoin amount with compare to remaining msgs
	return nil
}
