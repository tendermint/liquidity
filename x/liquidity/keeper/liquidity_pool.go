package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// return a specific liquidityPool
func (k Keeper) GetLiquidityPool(ctx sdk.Context, poolID uint64) (liquidityPool types.LiquidityPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolKey(poolID)

	value := store.Get(key)
	if value == nil {
		return liquidityPool, false
	}

	liquidityPool = types.MustUnmarshalLiquidityPool(k.cdc, value)

	return liquidityPool, true
}

func (k Keeper) SetLiquidityPool(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPool(k.cdc, liquidityPool)
	store.Set(types.GetLiquidityPoolKey(liquidityPool.PoolID), b)
}

// IterateAllLiquidityPools iterate through all of the liquidityPools
func (k Keeper) IterateAllLiquidityPools(ctx sdk.Context, cb func(liquidityPool types.LiquidityPool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.LiquidityPoolKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		liquidityPool := types.MustUnmarshalLiquidityPool(k.cdc, iterator.Value())
		if cb(liquidityPool) {
			break
		}
	}
}

// GetAllLiquidityPools returns all liquidityPools used during genesis dump
func (k Keeper) GetAllLiquidityPools(ctx sdk.Context) (liquidityPools []types.LiquidityPool) {
	k.IterateAllLiquidityPools(ctx, func(liquidityPool types.LiquidityPool) bool {
		liquidityPools = append(liquidityPools, liquidityPool)
		return false
	})

	return liquidityPools
}

// GetNextLiquidityID returns and increments the global Liquidity Pool ID counter.
// If the global account number is not set, it initializes it with value 0.
func (k Keeper) GetNextLiquidityPoolID(ctx sdk.Context) uint64 {
	var poolID uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalLiquidityPoolIDKey)
	if bz == nil {
		// initialize the LiquidityPoolID
		poolID = 0
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.UnmarshalBinaryBare(bz, &val)
		if err != nil {
			panic(err)
		}

		poolID = val.GetValue()
	}

	bz = k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: poolID + 1})
	store.Set(types.GlobalLiquidityPoolIDKey, bz)

	return poolID
}

func (k Keeper) GetLiquidityPoolByReserveAccIndex(ctx sdk.Context, reserveAcc sdk.AccAddress) (liquidityPool types.LiquidityPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolByReserveAccIndexKey(reserveAcc)

	value := store.Get(key)
	if value == nil {
		return liquidityPool, false
	}

	liquidityPool = types.MustUnmarshalLiquidityPool(k.cdc, value)

	return liquidityPool, true
}

func (k Keeper) SetLiquidityPoolByReserveAccIndex(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPool(k.cdc, liquidityPool)
	store.Set(types.GetLiquidityPoolByReserveAccIndexKey(liquidityPool.ReserveAccount), b)
}

func (k Keeper) SetLiquidityPoolAtomic(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	liquidityPool.PoolID = k.GetNextLiquidityPoolID(ctx)
	k.SetLiquidityPool(ctx, liquidityPool)
	k.SetLiquidityPoolByReserveAccIndex(ctx, liquidityPool)
}

// return a specific GetLiquidityPoolBatchIndex
func (k Keeper) GetLiquidityPoolBatchIndex(ctx sdk.Context, poolID uint64) (liquidityPoolBatchIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolBatchIndex(poolID)

	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	liquidityPoolBatchIndex = sdk.BigEndianToUint64(bz)
	return liquidityPoolBatchIndex
}

func (k Keeper) SetLiquidityPoolBatchIndex(ctx sdk.Context, poolID, batchIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	b := sdk.Uint64ToBigEndian(batchIndex)
	store.Set(types.GetLiquidityPoolBatchIndex(poolID), b)
}

// return a specific liquidityPoolBatch
func (k Keeper) GetLiquidityPoolBatch(ctx sdk.Context, poolID uint64) (liquidityPoolBatch types.LiquidityPoolBatch, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolKey(poolID)

	value := store.Get(key)
	if value == nil {
		return liquidityPoolBatch, false
	}

	liquidityPoolBatch = types.MustUnmarshalLiquidityPoolBatch(k.cdc, value)

	return liquidityPoolBatch, true
}

func (k Keeper) GetNextBatchIndex(ctx sdk.Context, poolID uint64) (batchIndex uint64) {
	return k.GetLiquidityPoolBatchIndex(ctx, poolID) + 1
}

func (k Keeper) GetAllLiquidityPoolBatches(ctx sdk.Context) (liquidityPoolBatches []types.LiquidityPoolBatch) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		liquidityPoolBatches = append(liquidityPoolBatches, liquidityPoolBatch)
		return false
	})

	return liquidityPoolBatches
}

// IterateAllLiquidityPoolBatches iterate through all of the liquidityPoolBatches
func (k Keeper) IterateAllLiquidityPoolBatches(ctx sdk.Context, cb func(liquidityPoolBatch types.LiquidityPoolBatch) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.LiquidityPoolBatchKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		liquidityPoolBatch := types.MustUnmarshalLiquidityPoolBatch(k.cdc, iterator.Value())
		if cb(liquidityPoolBatch) {
			break
		}
	}
}

func (k Keeper) DeleteLiquidityPoolBatch(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetLiquidityPoolBatchKey(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex)
	store.Delete(batchKey)
}

func (k Keeper) SetLiquidityPoolBatch(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPoolBatch(k.cdc, liquidityPoolBatch)
	store.Set(types.GetLiquidityPoolBatchKey(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex), b)
}

//
func (k Keeper) SetLiquidityPoolBatchDepositMsg(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, msgIndex uint64, msg types.BatchPoolDepositMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolDepositMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchDepositMsgIndex(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex, msgIndex), b)
}

// IterateAllLiquidityPoolBatchDepositMsgs iterate through all of the LiquidityPoolBatchDepositMsgs
func (k Keeper) IterateAllLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, cb func(msg types.BatchPoolDepositMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetLiquidityPoolBatchDepositMsgsPrefix(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllLiquidityPoolBatchDepositMsgs returns all BatchDepositMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolDepositMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

//
func (k Keeper) SetLiquidityPoolBatchWithdrawMsg(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, msgIndex uint64, msg types.BatchPoolWithdrawMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchWithdrawMsgIndex(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex, msgIndex), b)
}

// IterateAllLiquidityPoolBatchWithdrawMsgs iterate through all of the LiquidityPoolBatchWithdrawMsgs
func (k Keeper) IterateAllLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, cb func(msg types.BatchPoolWithdrawMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetLiquidityPoolBatchWithdrawMsgsPrefix(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllLiquidityPoolBatchWithdrawMsgs returns all BatchWithdrawMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

//
func (k Keeper) SetLiquidityPoolBatchSwapMsg(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, msgIndex uint64, msg types.BatchPoolSwapMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolSwapMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchSwapMsgIndex(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex, msgIndex), b)
}

// IterateAllLiquidityPoolBatchSwapMsgs iterate through all of the LiquidityPoolBatchSwapMsgs
func (k Keeper) IterateAllLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, cb func(msg types.BatchPoolSwapMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetLiquidityPoolBatchSwapMsgsPrefix(liquidityPoolBatch.PoolID, liquidityPoolBatch.BatchIndex)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllLiquidityPoolBatchSwapMsgs returns all BatchSwapMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolSwapMsg) {
	k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

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

	pool, found := k.GetLiquidityPool(ctx, msg.PoolID)
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
	// TODO: coin sort alphabetical
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

	// TODO: msg.Depositor to escrowModAcc
	batchEscrowAcc := msg.Depositor

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

	// calculate pool token mint amount
	poolCoinAmt := k.GetPoolCoinTotalSupply(ctx, pool).Mul(depositableCoinA).Quo(coinA.Amount)
	poolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
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

	pool, found := k.GetLiquidityPool(ctx, msg.PoolID)
	if !found {
		return types.ErrPoolNotExists
	}

	totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	//params := k.GetParams(ctx)
	for _, reserveCoin := range reserveCoins {
		withdrawAmt := reserveCoin.Amount.Mul(poolCoin.Amount).Quo(totalSupply)
		// TODO: apply fee, (sdk.NewDec(1).Sub(params.LiquidityPoolFeeRate)
		// TODO: to using k.bankKeeper.SendCoinsFromModuleToAccount() with poolKey
		inputs = append(inputs, banktypes.NewInput(pool.ReserveAccount, sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
		outputs = append(outputs, banktypes.NewOutput(msg.Withdrawer, sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
	}

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, msg.Withdrawer, types.ModuleName, poolCoins); err != nil {
		return err
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins); err != nil {
		return err
	}
	// TODO: add events for batch result, each err cases
	return nil
}

func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolID)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before executed on batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolDepositMsg{
		//TxHash: nil,
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}
	// TODO: escrow
	poolBatch.DepositMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchDepositMsg(ctx, poolBatch, poolBatch.DepositMsgIndex, batchPoolMsg)
	return nil
}

func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolID)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before executed on batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolWithdrawMsg{
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}
	// TODO: escrow
	poolBatch.WithdrawMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, poolBatch, poolBatch.WithdrawMsgIndex, batchPoolMsg)
	return nil
}

func (k Keeper) SwapLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgSwap) error {
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolID)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before executed on batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolSwapMsg{
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}
	// TODO: escrow
	poolBatch.WithdrawMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchSwapMsg(ctx, poolBatch, poolBatch.WithdrawMsgIndex, batchPoolMsg)
	return nil
}
