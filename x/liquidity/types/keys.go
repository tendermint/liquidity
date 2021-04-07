package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the liquidity module
	ModuleName = "liquidity"

	// RouterKey is the message router key for the liquidity module
	RouterKey = ModuleName

	// StoreKey is the default store key for the liquidity module
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the liquidity module
	QuerierRoute = ModuleName

	// PoolCoinDenomPrefix is the prefix used for liquidity pool coin representation
	PoolCoinDenomPrefix = "pool"
)

// prefix key of liquidity states for indexing when kvstore
var (
	// param key for global Liquidity Pool IDs
	GlobalLiquidityPoolIdKey = []byte("globalLiquidityPoolId")

	PoolKeyPrefix                  = []byte{0x11}
	PoolByReserveAccIndexKeyPrefix = []byte{0x12}

	PoolBatchIndexKeyPrefix = []byte{0x21} // Last PoolBatchIndex
	PoolBatchKeyPrefix      = []byte{0x22}

	PoolBatchDepositMsgStateIndexKeyPrefix  = []byte{0x31}
	PoolBatchWithdrawMsgStateIndexKeyPrefix = []byte{0x32}
	PoolBatchSwapMsgStateIndexKeyPrefix     = []byte{0x33}
)

// GetPoolKey returns kv indexing key of the pool
func GetPoolKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

// GetPoolByReserveAccIndexKey returns kv indexing key of the pool indexed by reserve account
func GetPoolByReserveAccIndexKey(reserveAcc sdk.AccAddress) []byte {
	return append(PoolByReserveAccIndexKeyPrefix, reserveAcc.Bytes()...)
}

// GetPoolBatchIndexKey returns kv indexing key of the latest index value of the pool batch
func GetPoolBatchIndexKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// GetPoolBatchKey returns kv indexing key of the pool batch indexed by pool id
func GetPoolBatchKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// GetPoolBatchDepositMsgStatesPrefix returns prefix of deposit message states in the pool's latest batch for iteration
func GetPoolBatchDepositMsgStatesPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchDepositMsgStateIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// GetPoolBatchWithdrawMsgsPrefix returns prefix of withdraw message states in the pool's latest batch for iteration
func GetPoolBatchWithdrawMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchWithdrawMsgStateIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// GetPoolBatchSwapMsgStatesPrefix returns prefix of swap message states in the pool's latest batch for iteration
func GetPoolBatchSwapMsgStatesPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchSwapMsgStateIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// GetPoolBatchDepositMsgStateIndexKey returns kv indexing key of the latest index value of the msg index
func GetPoolBatchDepositMsgStateIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = PoolBatchDepositMsgStateIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

// GetPoolBatchWithdrawMsgStateIndexKey returns kv indexing key of the latest index value of the msg index
func GetPoolBatchWithdrawMsgStateIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = PoolBatchWithdrawMsgStateIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

// GetPoolBatchSwapMsgStateIndexKey returns kv indexing key of the latest index value of the msg index
func GetPoolBatchSwapMsgStateIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = PoolBatchSwapMsgStateIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}
