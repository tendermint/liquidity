package types

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Calculate unique Pool key of the liquidity pool
// need to validate alphabetical ordering of ReserveCoinDenoms when New() and Store
// Denominations can be 3 ~ 128 characters long and support letters, followed by either
// a letter, a number or a separator ('/').
// reDnmString = `[a-zA-Z][a-zA-Z0-9/]{2,127}`.
func (pool Pool) Name() string {
	return PoolName(pool.ReserveCoinDenoms, pool.TypeId)
}

// Validate each constraint of the liquidity pool
func (pool Pool) Validate() error {
	if pool.Id == 0 {
		return ErrPoolNotExists
	}
	if pool.TypeId == 0 {
		return ErrPoolTypeNotExists
	}
	if pool.ReserveCoinDenoms == nil || len(pool.ReserveCoinDenoms) == 0 {
		return ErrNumOfReserveCoinDenoms
	}
	if uint32(len(pool.ReserveCoinDenoms)) > MaxReserveCoinNum || uint32(len(pool.ReserveCoinDenoms)) < MinReserveCoinNum {
		return ErrNumOfReserveCoinDenoms
	}
	sortedDenomA, sortedDenomB := AlphabeticalDenomPair(pool.ReserveCoinDenoms[0], pool.ReserveCoinDenoms[1])
	if sortedDenomA != pool.ReserveCoinDenoms[0] || sortedDenomB != pool.ReserveCoinDenoms[1] {
		return ErrBadOrderingReserveCoinDenoms
	}
	if pool.ReserveAccountAddress == "" {
		return ErrEmptyReserveAccountAddress
	}
	//addr, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	//if err != nil || pool.GetReserveAccount().Equals(addr) {
	//	return ErrBadReserveAccountAddress
	//}
	if pool.ReserveAccountAddress != GetPoolReserveAcc(pool.Name()).String() {
		return ErrBadReserveAccountAddress
	}
	if pool.PoolCoinDenom == "" {
		return ErrEmptyPoolCoinDenom
	}
	if pool.PoolCoinDenom != pool.Name() {
		return ErrBadPoolCoinDenom
	}
	return nil
}

// Calculate unique Pool key of the liquidity pool
func PoolName(reserveCoinDenoms []string, poolTypeId uint32) string {
	return strings.Join(append(SortDenoms(reserveCoinDenoms), strconv.FormatUint(uint64(poolTypeId), 10)), "/")
}

// NewPoolBatch creates a new PoolBatch object.
func NewPoolBatch(poolId, batchIndex uint64) PoolBatch {
	return PoolBatch{
		PoolId:           poolId,
		Index:            batchIndex,
		BeginHeight:      0,
		DepositMsgIndex:  1,
		WithdrawMsgIndex: 1,
		SwapMsgIndex:     1,
		Executed:         false,
	}
}

// GetPoolBatchResponse returns a PoolBatchResponse object skipped pool_id, It used for result of queries
func GetPoolBatchResponse(poolBatch PoolBatch) PoolBatchResponse {
	return PoolBatchResponse{
		Index:            poolBatch.Index,
		BeginHeight:      poolBatch.BeginHeight,
		DepositMsgIndex:  poolBatch.DepositMsgIndex,
		WithdrawMsgIndex: poolBatch.WithdrawMsgIndex,
		SwapMsgIndex:     poolBatch.SwapMsgIndex,
		Executed:         poolBatch.Executed,
	}
}

// MustMarshalPool returns the liquidityPool bytes. Panics if fails
func MustMarshalPool(cdc codec.BinaryMarshaler, liquidityPool Pool) []byte {
	return cdc.MustMarshalBinaryBare(&liquidityPool)
}

// MustUnmarshalPool return the unmarshalled liquidityPool from bytes.
// Panics if fails.
func MustUnmarshalPool(cdc codec.BinaryMarshaler, value []byte) Pool {
	liquidityPool, err := UnmarshalPool(cdc, value)
	if err != nil {
		panic(err)
	}

	return liquidityPool
}

// return the liquidityPool
func UnmarshalPool(cdc codec.BinaryMarshaler, value []byte) (liquidityPool Pool, err error) {
	err = cdc.UnmarshalBinaryBare(value, &liquidityPool)
	return liquidityPool, err
}

// return sdk.AccAddress object of the address saved as string because of protobuf
func (pool Pool) GetReserveAccount() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// return pool coin denom of the liquidity pool
func (pool Pool) GetPoolCoinDenom() string { return pool.PoolCoinDenom }

// return pool id of the liquidity pool
func (pool Pool) GetPoolId() uint64 { return pool.Id }

// Pools is a collection of liquidityPools
type Pools []Pool

// get string of list of liquidity pool
func (pools Pools) String() (out string) {
	for _, del := range pools {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// MustMarshalPoolBatch returns the PoolBatch bytes. Panics if fails
func MustMarshalPoolBatch(cdc codec.BinaryMarshaler, poolBatch PoolBatch) []byte {
	return cdc.MustMarshalBinaryBare(&poolBatch)
}

// return the poolBatch
func UnmarshalPoolBatch(cdc codec.BinaryMarshaler, value []byte) (poolBatch PoolBatch, err error) {
	err = cdc.UnmarshalBinaryBare(value, &poolBatch)
	return poolBatch, err
}

// MustUnmarshalPool return the unmarshalled PoolBatch from bytes.
// Panics if fails.
func MustUnmarshalPoolBatch(cdc codec.BinaryMarshaler, value []byte) PoolBatch {
	poolBatch, err := UnmarshalPoolBatch(cdc, value)
	if err != nil {
		panic(err)
	}

	return poolBatch
}

// MustMarshalDepositMsgState returns the DepositMsgState bytes. Panics if fails
func MustMarshalDepositMsgState(cdc codec.BinaryMarshaler, msg DepositMsgState) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the DepositMsgState
func UnmarshalDepositMsgState(cdc codec.BinaryMarshaler, value []byte) (msg DepositMsgState, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalDepositMsgState return the unmarshalled DepositMsgState from bytes.
// Panics if fails.
func MustUnmarshalDepositMsgState(cdc codec.BinaryMarshaler, value []byte) DepositMsgState {
	msg, err := UnmarshalDepositMsgState(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalWithdrawMsgState returns the WithdrawMsgState bytes. Panics if fails
func MustMarshalWithdrawMsgState(cdc codec.BinaryMarshaler, msg WithdrawMsgState) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the WithdrawMsgState
func UnmarshalWithdrawMsgState(cdc codec.BinaryMarshaler, value []byte) (msg WithdrawMsgState, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalWithdrawMsgState return the unmarshalled WithdrawMsgState from bytes.
// Panics if fails.
func MustUnmarshalWithdrawMsgState(cdc codec.BinaryMarshaler, value []byte) WithdrawMsgState {
	msg, err := UnmarshalWithdrawMsgState(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalSwapMsgState returns the SwapMsgState bytes. Panics if fails
func MustMarshalSwapMsgState(cdc codec.BinaryMarshaler, msg SwapMsgState) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the UnmarshalSwapMsgState
func UnmarshalSwapMsgState(cdc codec.BinaryMarshaler, value []byte) (msg SwapMsgState, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalSwapMsgState return the unmarshalled SwapMsgState from bytes.
// Panics if fails.
func MustUnmarshalSwapMsgState(cdc codec.BinaryMarshaler, value []byte) SwapMsgState {
	msg, err := UnmarshalSwapMsgState(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}
