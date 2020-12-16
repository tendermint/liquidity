package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// read form kvstore and return a specific liquidityPool
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

// store to kvstore a specific liquidityPool
func (k Keeper) SetLiquidityPool(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPool(k.cdc, liquidityPool)
	store.Set(types.GetLiquidityPoolKey(liquidityPool.PoolId), b)
}

// delete from kvstore a specific liquidityPool
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

// return next liquidity pool id for new pool, using index of latest pool id
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

// read form kvstore and return a specific liquidityPool indexed by given reserve account
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

// Set Index by ReserveAcc for liquidity Pool duplication check
func (k Keeper) SetLiquidityPoolByReserveAccIndex(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPool(k.cdc, liquidityPool)
	store.Set(types.GetLiquidityPoolByReserveAccIndexKey(liquidityPool.GetReserveAccount()), b)
}

// Set Liquidity Pool with set global pool id index +1 and index by reserveAcc
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

// set index for liquidity pool batch, it should be increase after batch executed
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

// return next batch index, with set index increased
func (k Keeper) GetNextBatchIndexWithUpdate(ctx sdk.Context, poolId uint64) (batchIndex uint64) {
	batchIndex = k.GetLiquidityPoolBatchIndex(ctx, poolId)
	batchIndex += 1
	k.SetLiquidityPoolBatchIndex(ctx, poolId, batchIndex)
	return
}

// Get All batches of the all existed liquidity pools
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

// Delete batch of the liquidity pool, it used for test case
func (k Keeper) DeleteLiquidityPoolBatch(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetLiquidityPoolBatchKey(liquidityPoolBatch.PoolId)
	store.Delete(batchKey)
}

// set batch of the liquidity pool, with current state
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

// set deposit batch msg of the liquidity pool batch, with current state
func (k Keeper) SetLiquidityPoolBatchDepositMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolDepositMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolDepositMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
}

// set deposit batch msgs of the liquidity pool batch, with current state using pointers
func (k Keeper) SetLiquidityPoolBatchDepositMsgsByPointer(ctx sdk.Context, poolId uint64, msgList []*types.BatchPoolDepositMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolDepositMsg(k.cdc, *msg)
		store.Set(types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

// set deposit batch msgs of the liquidity pool batch, with current state
func (k Keeper) SetLiquidityPoolBatchDepositMsgs(ctx sdk.Context, poolId uint64, msgList []types.BatchPoolDepositMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolDepositMsg(k.cdc, msg)
		store.Set(types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

//func (k Keeper) DeleteLiquidityPoolBatchDepositMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
//	store := ctx.KVStore(k.storeKey)
//	batchKey := types.GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msgIndex)
//	store.Delete(batchKey)
//}

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

// IterateAllBatchDepositMsgs iterate through all of the BatchDepositMsgs of all batches
func (k Keeper) IterateAllBatchDepositMsgs(ctx sdk.Context, cb func(msg types.BatchPoolDepositMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.LiquidityPoolBatchDepositMsgIndexKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllBatchDepositMsgs returns all BatchDepositMsgs for all batches.
func (k Keeper) GetAllBatchDepositMsgs(ctx sdk.Context) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllBatchDepositMsgs(ctx, func(msg types.BatchPoolDepositMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllLiquidityPoolBatchDepositMsgs returns all BatchDepositMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolDepositMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllToDeleteLiquidityPoolBatchDepositMsgs returns all Not toDelete BatchDepositMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllNotToDeleteLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolDepositMsg) bool {
		if !msg.ToBeDeleted {
			msgs = append(msgs, msg)
		}
		return false
	})
	return msgs
}

// GetAllRemainingLiquidityPoolBatchDepositMsgs returns All only remaining BatchDepositMsgs after endblock , executed but not toDelete
func (k Keeper) GetAllRemainingLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []*types.BatchPoolDepositMsg) {
	k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolDepositMsg) bool {
		if msg.Executed && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// delete deposit batch msgs of the liquidity pool batch which has state ToBeDeleted
func (k Keeper) DeleteAllReadyLiquidityPoolBatchDepositMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidityPoolBatchDepositMsgsPrefix(liquidityPoolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if msg.ToBeDeleted {
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

// set withdraw batch msg of the liquidity pool batch, with current state
func (k Keeper) SetLiquidityPoolBatchWithdrawMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolWithdrawMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
}

// set withdraw batch msgs of the liquidity pool batch, with current state using pointers
func (k Keeper) SetLiquidityPoolBatchWithdrawMsgsByPointer(ctx sdk.Context, poolId uint64, msgList []*types.BatchPoolWithdrawMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, *msg)
		store.Set(types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

// set withdraw batch msgs of the liquidity pool batch, with current state
func (k Keeper) SetLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, poolId uint64, msgList []types.BatchPoolWithdrawMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, msg)
		store.Set(types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

//func (k Keeper) DeleteLiquidityPoolBatchWithdrawMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
//	store := ctx.KVStore(k.storeKey)
//	batchKey := types.GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msgIndex)
//	store.Delete(batchKey)
//}

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

// IterateAllBatchWithdrawMsgs iterate through all of the BatchPoolWithdrawMsg of all batches
func (k Keeper) IterateAllBatchWithdrawMsgs(ctx sdk.Context, cb func(msg types.BatchPoolWithdrawMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.LiquidityPoolBatchWithdrawMsgIndexKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllBatchWithdrawMsgs returns all BatchWithdrawMsgs for all batches
func (k Keeper) GetAllBatchWithdrawMsgs(ctx sdk.Context) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllBatchWithdrawMsgs(ctx, func(msg types.BatchPoolWithdrawMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllLiquidityPoolBatchWithdrawMsgs returns all BatchWithdrawMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllToDeleteLiquidityPoolBatchWithdrawMsgs returns all Not to delete BatchWithdrawMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllNotToDeleteLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		if !msg.ToBeDeleted {
			msgs = append(msgs, msg)
		}
		return false
	})
	return msgs
}

// GetAllRemainingLiquidityPoolBatchWithdrawMsgs returns All only remaining BatchWithdrawMsgs after endblock, executed but not toDelete
func (k Keeper) GetAllRemainingLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []*types.BatchPoolWithdrawMsg) {
	k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		if msg.Executed && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// delete withdraw batch msgs of the liquidity pool batch which has state ToBeDeleted
func (k Keeper) DeleteAllReadyLiquidityPoolBatchWithdrawMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidityPoolBatchWithdrawMsgsPrefix(liquidityPoolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if msg.ToBeDeleted {
			store.Delete(iterator.Key())
		}
	}
}

// return a specific GetLiquidityPoolBatchSwapMsg, not used currently
//func (k Keeper) GetLiquidityPoolBatchSwapMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolSwapMsg, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msgIndex)
//
//	value := store.Get(key)
//	if value == nil {
//		return msg, false
//	}
//
//	msg = types.MustUnmarshalBatchPoolSwapMsg(k.cdc, value)
//	return msg, true
//}

// set swap batch msg of the liquidity pool batch, with current state
func (k Keeper) SetLiquidityPoolBatchSwapMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolSwapMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolSwapMsg(k.cdc, msg)
	store.Set(types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
}

// Delete swap batch msg of the liquidity pool batch, it used for test case
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

// IterateAllBatchSwapMsgs iterate through all of the BatchPoolSwapMsg of all batches
func (k Keeper) IterateAllBatchSwapMsgs(ctx sdk.Context, cb func(msg types.BatchPoolSwapMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.LiquidityPoolBatchSwapMsgIndexKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllBatchSwapMsgs returns all BatchSwapMsgs of all batches
func (k Keeper) GetAllBatchSwapMsgs(ctx sdk.Context) (msgs []types.BatchPoolSwapMsg) {
	k.IterateAllBatchSwapMsgs(ctx, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// delete swap batch msgs of the liquidity pool batch which has state ToBeDeleted
func (k Keeper) DeleteAllReadyLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidityPoolBatchSwapMsgsPrefix(liquidityPoolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if msg.ToBeDeleted {
			store.Delete(iterator.Key())
		}
	}
}

// GetAllLiquidityPoolBatchSwapMsgsAsPointer returns all BatchSwapMsgs pointer indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, &msg)
		return false
	})
	return msgs
}

// GetAllLiquidityPoolBatchSwapMsgs returns all BatchSwapMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []types.BatchPoolSwapMsg) {
	k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllNotProcessedLiquidityPoolBatchSwapMsgs returns All only not processed swap msgs, not executed with not succeed and not toDelete BatchSwapMsgs indexed by the liquidityPoolBatch
func (k Keeper) GetAllNotProcessedLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if !msg.Executed && !msg.Succeeded && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// GetAllRemainingLiquidityPoolBatchSwapMsgs returns All only remaining after endblock swap msgs, executed but not toDelete
func (k Keeper) GetAllRemainingLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if msg.Executed && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// GetAllNotToDeleteLiquidityPoolBatchSwapMsgs returns All only not to delete swap msgs
func (k Keeper) GetAllNotToDeleteLiquidityPoolBatchSwapMsgs(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// set swap batch msgs of the liquidity pool batch, with current state using pointers
func (k Keeper) SetLiquidityPoolBatchSwapMsgPointers(ctx sdk.Context, poolId uint64, msgList []*types.BatchPoolSwapMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolSwapMsg(k.cdc, *msg)
		store.Set(types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

// set swap batch msgs of the liquidity pool batch, with current state
func (k Keeper) SetLiquidityPoolBatchSwapMsgs(ctx sdk.Context, poolId uint64, msgList []types.BatchPoolSwapMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolSwapMsg(k.cdc, msg)
		store.Set(types.GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}
