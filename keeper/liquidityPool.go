package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/liquidity/types"
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
