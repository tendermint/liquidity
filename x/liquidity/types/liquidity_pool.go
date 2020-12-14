package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"
)

// Calculate unique Pool key of the liquidity pool
// need to validate alphabetical ordering of ReserveCoinDenoms when New() and Store
// Denominations can be 3 ~ 128 characters long and support letters, followed by either
// a letter, a number or a separator ('/').
// reDnmString = `[a-zA-Z][a-zA-Z0-9/]{2,127}`.
func (lp LiquidityPool) GetPoolKey() string {
	return GetPoolKey(lp.ReserveCoinDenoms, lp.PoolTypeIndex)
}

// Validate each constraint of the liquidity pool
func (lp LiquidityPool) Validate() error {
	if lp.PoolId == 0 {
		return ErrPoolNotExists
	}
	if lp.PoolTypeIndex == 0 {
		return ErrPoolTypeNotExists
	}
	if lp.ReserveCoinDenoms == nil || len(lp.ReserveCoinDenoms) == 0 {
		return ErrNumOfReserveCoinDenoms
	}
	if uint32(len(lp.ReserveCoinDenoms)) > MaxReserveCoinNum || uint32(len(lp.ReserveCoinDenoms)) < MinReserveCoinNum {
		return ErrNumOfReserveCoinDenoms
	}
	sortedDenomA, sortedDenomB := AlphabeticalDenomPair(lp.ReserveCoinDenoms[0], lp.ReserveCoinDenoms[1])
	if sortedDenomA != lp.ReserveCoinDenoms[0] || sortedDenomB != lp.ReserveCoinDenoms[1] {
		return ErrBadOrderingReserveCoinDenoms
	}
	if lp.ReserveAccountAddress == "" {
		return ErrEmptyReserveAccountAddress
	}
	//addr, err := sdk.AccAddressFromBech32(lp.ReserveAccountAddress)
	//if err != nil || lp.GetReserveAccount().Equals(addr) {
	//	return ErrBadReserveAccountAddress
	//}
	if lp.ReserveAccountAddress != GetPoolReserveAcc(lp.GetPoolKey()).String() {
		return ErrBadReserveAccountAddress
	}
	if lp.PoolCoinDenom == "" {
		return ErrEmptyPoolCoinDenom
	}
	if lp.PoolCoinDenom != lp.GetPoolKey() {
		return ErrBadPoolCoinDenom
	}
	return nil
}

// Calculate unique Pool key of the liquidity pool
func GetPoolKey(reserveCoinDenoms []string, poolTypeIndex uint32) string {
	return strings.Join(append(reserveCoinDenoms, strconv.FormatUint(uint64(poolTypeIndex), 10)), "-")
}

func NewLiquidityPoolBatch(poolId, batchIndex uint64) LiquidityPoolBatch {
	return LiquidityPoolBatch{
		PoolId:           poolId,
		BatchIndex:       batchIndex,
		BeginHeight:      0,
		DepositMsgIndex:  1,
		WithdrawMsgIndex: 1,
		SwapMsgIndex:     1,
		Executed:         false,
	}
}

// MustMarshalLiquidityPool returns the liquidityPool bytes. Panics if fails
func MustMarshalLiquidityPool(cdc codec.BinaryMarshaler, liquidityPool LiquidityPool) []byte {
	return cdc.MustMarshalBinaryBare(&liquidityPool)
}

// MustUnmarshalLiquidityPool return the unmarshaled liquidityPool from bytes.
// Panics if fails.
func MustUnmarshalLiquidityPool(cdc codec.BinaryMarshaler, value []byte) LiquidityPool {
	liquidityPool, err := UnmarshalLiquidityPool(cdc, value)
	if err != nil {
		panic(err)
	}

	return liquidityPool
}

// return the liquidityPool
func UnmarshalLiquidityPool(cdc codec.BinaryMarshaler, value []byte) (liquidityPool LiquidityPool, err error) {
	err = cdc.UnmarshalBinaryBare(value, &liquidityPool)
	return liquidityPool, err
}

// return sdk.AccAddress object of he address saved as string because of protobuf
func (lp LiquidityPool) GetReserveAccount() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(lp.ReserveAccountAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// return pool coin denom of the liquidity poool
func (lp LiquidityPool) GetPoolCoinDenom() string { return lp.PoolCoinDenom }

// return pool id of the liquidity poool
func (lp LiquidityPool) GetPoolId() uint64 { return lp.PoolId }

// LiquidityPools is a collection of liquidityPools
type LiquidityPools []LiquidityPool

// LiquidityPoolsBatch is a collection of liquidityPoolBatch
type LiquidityPoolsBatch []LiquidityPoolBatch

// get string of list of liquidity pool
func (lps LiquidityPools) String() (out string) {
	for _, del := range lps {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// MustMarshalLiquidityPoolBatch returns the LiquidityPoolBatch bytes. Panics if fails
func MustMarshalLiquidityPoolBatch(cdc codec.BinaryMarshaler, liquidityPoolBatch LiquidityPoolBatch) []byte {
	return cdc.MustMarshalBinaryBare(&liquidityPoolBatch)
}

// return the liquidityPoolBatch
func UnmarshalLiquidityPoolBatch(cdc codec.BinaryMarshaler, value []byte) (liquidityPoolBatch LiquidityPoolBatch, err error) {
	err = cdc.UnmarshalBinaryBare(value, &liquidityPoolBatch)
	return liquidityPoolBatch, err
}

// MustUnmarshalLiquidityPool return the unmarshaled LiquidityPoolBatch from bytes.
// Panics if fails.
func MustUnmarshalLiquidityPoolBatch(cdc codec.BinaryMarshaler, value []byte) LiquidityPoolBatch {
	liquidityPoolBatch, err := UnmarshalLiquidityPoolBatch(cdc, value)
	if err != nil {
		panic(err)
	}

	return liquidityPoolBatch
}

// MustMarshalBatchPoolDepositMsg returns the BatchPoolDepositMsg bytes. Panics if fails
func MustMarshalBatchPoolDepositMsg(cdc codec.BinaryMarshaler, msg BatchPoolDepositMsg) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the BatchPoolDepositMsg
func UnmarshalBatchPoolDepositMsg(cdc codec.BinaryMarshaler, value []byte) (msg BatchPoolDepositMsg, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalBatchPoolDepositMsg return the unmarshaled BatchPoolDepositMsg from bytes.
// Panics if fails.
func MustUnmarshalBatchPoolDepositMsg(cdc codec.BinaryMarshaler, value []byte) BatchPoolDepositMsg {
	msg, err := UnmarshalBatchPoolDepositMsg(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalBatchPoolWithdrawMsg returns the BatchPoolWithdrawMsg bytes. Panics if fails
func MustMarshalBatchPoolWithdrawMsg(cdc codec.BinaryMarshaler, msg BatchPoolWithdrawMsg) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the BatchPoolWithdrawMsg
func UnmarshalBatchPoolWithdrawMsg(cdc codec.BinaryMarshaler, value []byte) (msg BatchPoolWithdrawMsg, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalBatchPoolWithdrawMsg return the unmarshaled BatchPoolWithdrawMsg from bytes.
// Panics if fails.
func MustUnmarshalBatchPoolWithdrawMsg(cdc codec.BinaryMarshaler, value []byte) BatchPoolWithdrawMsg {
	msg, err := UnmarshalBatchPoolWithdrawMsg(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalBatchPoolSwapMsg returns the BatchPoolSwapMsg bytes. Panics if fails
func MustMarshalBatchPoolSwapMsg(cdc codec.BinaryMarshaler, msg BatchPoolSwapMsg) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the UnmarshalBatchPoolSwapMsg
func UnmarshalBatchPoolSwapMsg(cdc codec.BinaryMarshaler, value []byte) (msg BatchPoolSwapMsg, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalBatchPoolSwapMsg return the unmarshaled BatchPoolSwapMsg from bytes.
// Panics if fails.
func MustUnmarshalBatchPoolSwapMsg(cdc codec.BinaryMarshaler, value []byte) BatchPoolSwapMsg {
	msg, err := UnmarshalBatchPoolSwapMsg(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}
