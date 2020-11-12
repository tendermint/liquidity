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
	GlobalLiquidityPoolIdKey = []byte("globalLiquidityPoolId")

	LiquidityPoolKeyPrefix               = []byte{0x11}
	LiquidityPoolByReserveIndexKeyPrefix = []byte{0x12}

	LiquidityPoolBatchIndexKeyPrefix = []byte{0x21} // LastLiquidityPoolBatchIndex
	LiquidityPoolBatchKeyPrefix      = []byte{0x22}

	LiquidityPoolBatchDepositMsgIndexKeyPrefix  = []byte{0x31}
	LiquidityPoolBatchWithdrawMsgIndexKeyPrefix = []byte{0x32}
	LiquidityPoolBatchSwapMsgIndexKeyPrefix     = []byte{0x33}
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

// TODO: check msgList struck is right? need to like UnbondingDelegations, need to sorted MessageList?
//store.Delete(unbondingTimesliceIterator.Key())

func GetLiquidityPoolBatchIndexKey(poolID uint64) []byte {
	key := make([]byte, 9)
	key[0] = LiquidityPoolBatchIndexKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolID))
	return key
}

func GetLiquidityPoolBatchKey(poolID uint64, batchIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:], sdk.Uint64ToBigEndian(batchIndex))
	return key
}

func GetLiquidityPoolBatchDepositMsgsPrefix(poolID, batchIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:], sdk.Uint64ToBigEndian(batchIndex))
	return key
}

func GetLiquidityPoolBatchWithdrawMsgsPrefix(poolID, batchIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:], sdk.Uint64ToBigEndian(batchIndex))
	return key
}

func GetLiquidityPoolBatchSwapMsgsPrefix(poolID, batchIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = LiquidityPoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:], sdk.Uint64ToBigEndian(batchIndex))
	return key
}

func GetLiquidityPoolBatchDepositMsgIndex(poolID, batchIndex, msgIndex uint64) []byte {
	key := make([]byte, 25)
	key[0] = LiquidityPoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:17], sdk.Uint64ToBigEndian(batchIndex))
	copy(key[17:], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

func GetLiquidityPoolBatchWithdrawMsgIndex(poolID, batchIndex, msgIndex uint64) []byte {
	key := make([]byte, 25)
	key[0] = LiquidityPoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:17], sdk.Uint64ToBigEndian(batchIndex))
	copy(key[17:], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

func GetLiquidityPoolBatchSwapMsgIndex(poolID, batchIndex, msgIndex uint64) []byte {
	key := make([]byte, 25)
	key[0] = LiquidityPoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolID))
	copy(key[9:17], sdk.Uint64ToBigEndian(batchIndex))
	copy(key[17:], sdk.Uint64ToBigEndian(msgIndex))
	return key
}
