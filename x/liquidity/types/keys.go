package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Routes, Keys for liquidity module
const (
	// ModuleName is the name of the module.
	ModuleName = "liquidity"

	// RouterKey is the message route for the liquidity module.
	RouterKey = ModuleName

	// StoreKey is the default store key for the liquidity module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the liquidity module.
	QuerierRoute = ModuleName

	// PoolCoinDenomPrefix is the prefix used for liquidity pool coin representation.
	PoolCoinDenomPrefix = "pool"
)

// prefix key of liquidity states for indexing when kvstore
var (
	// param key for global Liquidity Pool IDs
	GlobalLiquidityPoolIdKey = []byte("globalLiquidityPoolId")

	LiquidityPoolKeyPrefix               = []byte{0x11}
	LiquidityPoolByReserveIndexKeyPrefix = []byte{0x12}

	LiquidityPoolBatchIndexKeyPrefix = []byte{0x21} // LastLiquidityPoolBatchIndex
	LiquidityPoolBatchKeyPrefix      = []byte{0x22}

	LiquidityPoolBatchDepositMsgIndexKeyPrefix  = []byte{0x31}
	LiquidityPoolBatchWithdrawMsgIndexKeyPrefix = []byte{0x32}
	LiquidityPoolBatchSwapMsgIndexKeyPrefix     = []byte{0x33}
)

// return kv indexing key of the pool
func GetLiquidityPoolKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

// return kv indexing key of the pool indexed by reserve account
func GetLiquidityPoolByReserveAccIndexKey(reserveAcc sdk.AccAddress) []byte {
	return append(LiquidityPoolByReserveIndexKeyPrefix, reserveAcc.Bytes()...)
}

// return kv indexing key of the latest index value of the pool batch
func GetLiquidityPoolBatchIndexKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchIndexKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

// return kv indexing key of the pool batch indexed by pool id
func GetLiquidityPoolBatchKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// Get prefix of the deposit batch messages that given pool for iteration
func GetLiquidityPoolBatchDepositMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// Get prefix of the withdraw batch messages that given pool for iteration
func GetLiquidityPoolBatchWithdrawMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// Get prefix of the swap batch messages that given pool for iteration
func GetLiquidityPoolBatchSwapMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// return kv indexing key of the latest index value of the msg index
func GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

// return kv indexing key of the latest index value of the msg index
func GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

// return kv indexing key of the latest index value of the msg index
func GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}
