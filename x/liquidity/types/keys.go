package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module.
	ModuleName = "liquidity"

	// RouterKey is the message route for the liquidity module.
	RouterKey = ModuleName

	// StoreKey is the default store key for the liquidity module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the liquidity module.
	QuerierRoute = StoreKey
)

var (

	// param key for global Liquidity Pool IDs
	GlobalLiquidityPoolIDKey = []byte("globalLiquidityPoolID")

	LiquidityPoolKeyPrefix               = []byte{0x11}
	LiquidityPoolByReserveIndexKeyPrefix = []byte{0x12}

	LiquidityPoolBatchIndexKeyPrefix = []byte{0x21} // LastLiquidityPoolBatchIndex
	LiquidityPoolBatchKeyPrefix      = []byte{0x21}
)

func GetLiquidityPoolKey(poolID uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolID))
	return key
}

func GetLiquidityPoolByReserveAccIndexKey(reserveAcc sdk.AccAddress) []byte {
	return append(LiquidityPoolByReserveIndexKeyPrefix, reserveAcc.Bytes()...)
}

func GetLiquidityPoolBatchIndex(poolID uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchIndexKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolID))
	return key
}

func GetLiquidityPoolBatch(poolID uint64, batchIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:], sdk.Uint64ToBigEndian(batchIndex))
	return key
}
