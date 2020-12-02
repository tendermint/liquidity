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
	QuerierRoute = ModuleName
)

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

func GetLiquidityPoolKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

func GetLiquidityPoolByReserveAccIndexKey(reserveAcc sdk.AccAddress) []byte {
	return append(LiquidityPoolByReserveIndexKeyPrefix, reserveAcc.Bytes()...)
}

func GetLiquidityPoolBatchIndexKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchIndexKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

func GetLiquidityPoolBatchKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

func GetLiquidityPoolBatchDepositMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

func GetLiquidityPoolBatchWithdrawMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

func GetLiquidityPoolBatchSwapMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

func GetLiquidityPoolBatchDepositMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

func GetLiquidityPoolBatchWithdrawMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

func GetLiquidityPoolBatchSwapMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}
