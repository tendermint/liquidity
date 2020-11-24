package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// return a specific liquidityPool
func (k Keeper) GetLiquidityPool(ctx sdk.Context, poolId uint64) (liquidityPool types.LiquidityPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolKey(poolId)

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
	store.Set(types.GetLiquidityPoolKey(liquidityPool.PoolId), b)
}

func (k Keeper) DeleteLiquidityPool(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	Key := types.GetLiquidityPoolKey(liquidityPool.PoolId)
	store.Delete(Key)
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
func (k Keeper) GetNextLiquidityPoolIdWithUpdate(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	poolId := k.GetNextLiquidityPoolId(ctx)
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: poolId + 1})
	store.Set(types.GlobalLiquidityPoolIdKey, bz)
	return poolId
}

func (k Keeper) GetNextLiquidityPoolId(ctx sdk.Context) uint64 {
	var poolId uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalLiquidityPoolIdKey)
	if bz == nil {
		// initialize the LiquidityPoolId
		poolId = 1
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.UnmarshalBinaryBare(bz, &val)
		if err != nil {
			panic(err)
		}

		poolId = val.GetValue()
	}
	return poolId
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
	store.Set(types.GetLiquidityPoolByReserveAccIndexKey(liquidityPool.GetReserveAccount()), b)
}

func (k Keeper) SetLiquidityPoolAtomic(ctx sdk.Context, liquidityPool types.LiquidityPool) types.LiquidityPool {
	liquidityPool.PoolId = k.GetNextLiquidityPoolIdWithUpdate(ctx)
	k.SetLiquidityPool(ctx, liquidityPool)
	k.SetLiquidityPoolByReserveAccIndex(ctx, liquidityPool)
	return liquidityPool
}

// return a specific GetLiquidityPoolBatchIndexKey
func (k Keeper) GetLiquidityPoolBatchIndex(ctx sdk.Context, poolId uint64) (liquidityPoolBatchIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolBatchIndexKey(poolId)

	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	liquidityPoolBatchIndex = sdk.BigEndianToUint64(bz)
	return liquidityPoolBatchIndex
}

func (k Keeper) SetLiquidityPoolBatchIndex(ctx sdk.Context, poolId, batchIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	b := sdk.Uint64ToBigEndian(batchIndex)
	store.Set(types.GetLiquidityPoolBatchIndexKey(poolId), b)
}

// return a specific liquidityPoolBatch
func (k Keeper) GetLiquidityPoolBatch(ctx sdk.Context, poolId uint64) (liquidityPoolBatch types.LiquidityPoolBatch, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolBatchKey(poolId)

	value := store.Get(key)
	if value == nil {
		return liquidityPoolBatch, false
	}

	liquidityPoolBatch = types.MustUnmarshalLiquidityPoolBatch(k.cdc, value)

	return liquidityPoolBatch, true
}

func (k Keeper) GetNextBatchIndexWithUpdate(ctx sdk.Context, poolId uint64) (batchIndex uint64) {
	batchIndex = k.GetLiquidityPoolBatchIndex(ctx, poolId)
	batchIndex += 1
	k.SetLiquidityPoolBatchIndex(ctx, poolId, batchIndex)
	return
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
	batchKey := types.GetLiquidityPoolBatchKey(liquidityPoolBatch.PoolId)
	store.Delete(batchKey)
}

func (k Keeper) SetLiquidityPoolBatch(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPoolBatch(k.cdc, liquidityPoolBatch)
	store.Set(types.GetLiquidityPoolBatchKey(liquidityPoolBatch.PoolId), b)
}

// return a specific liquidityPoolBatchDepositMsg
func (k Keeper) GetLiquidityPoolBatchDepositMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolDepositMsg, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msgIndex)

	value := store.Get(key)
	if value == nil {
		return msg, false
	}

	msg = types.MustUnmarshalBatchPoolDepositMsg(k.cdc, value)
	return msg, true
}

func (k Keeper) SetLiquidityPoolBatchDepositMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolDepositMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolDepositMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
}

func (k Keeper) DeleteLiquidityPoolBatchDepositMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msgIndex)
	store.Delete(batchKey)
}

// IterateAllLiquidityPoolBatchDepositMsgs iterate through all of the LiquidityPoolBatchDepositMsgs
func (k Keeper) IterateAllLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, cb func(msg types.BatchPoolDepositMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetLiquidityPoolBatchDepositMsgsPrefix(liquidityPoolBatch.PoolId)
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

func (k Keeper) DeleteAllReadyLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidityPoolBatchDepositMsgsPrefix(liquidityPoolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if msg.ToDelete {
			store.Delete(iterator.Key())
		}
	}
}

// return a specific liquidityPoolBatchWithdrawMsg
func (k Keeper) GetLiquidityPoolBatchWithdrawMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolWithdrawMsg, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msgIndex)

	value := store.Get(key)
	if value == nil {
		return msg, false
	}

	msg = types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, value)
	return msg, true
}

func (k Keeper) SetLiquidityPoolBatchWithdrawMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolWithdrawMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
}

func (k Keeper) DeleteLiquidityPoolBatchWithdrawMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msgIndex)
	store.Delete(batchKey)
}

// IterateAllLiquidityPoolBatchWithdrawMsgs iterate through all of the LiquidityPoolBatchWithdrawMsgs
func (k Keeper) IterateAllLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, cb func(msg types.BatchPoolWithdrawMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetLiquidityPoolBatchWithdrawMsgsPrefix(liquidityPoolBatch.PoolId)
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

func (k Keeper) DeleteAllReadyLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidityPoolBatchWithdrawMsgsPrefix(liquidityPoolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if msg.ToDelete {
			store.Delete(iterator.Key())
		}
	}
}

// return a specific liquidityPoolBatchDepositMsg
func (k Keeper) GetLiquidityPoolBatchSwapMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolSwapMsg, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msgIndex)

	value := store.Get(key)
	if value == nil {
		return msg, false
	}

	msg = types.MustUnmarshalBatchPoolSwapMsg(k.cdc, value)
	return msg, true
}

func (k Keeper) SetLiquidityPoolBatchSwapMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolSwapMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolSwapMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
}

func (k Keeper) DeleteLiquidityPoolBatchSwapMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msgIndex)
	store.Delete(batchKey)
}

// IterateAllLiquidityPoolBatchSwapMsgs iterate through all of the LiquidityPoolBatchSwapMsgs
func (k Keeper) IterateAllLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch, cb func(msg types.BatchPoolSwapMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetLiquidityPoolBatchSwapMsgsPrefix(liquidityPoolBatch.PoolId)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

func (k Keeper) DeleteAllReadyLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidityPoolBatchSwapMsgsPrefix(liquidityPoolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if msg.ToDelete {
			store.Delete(iterator.Key())
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
