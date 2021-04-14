package types

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PoolName returns unique name of the pool consists of given reserve coin denoms and type id.
func PoolName(reserveCoinDenoms []string, poolTypeId uint32) string {
	return strings.Join(append(SortDenoms(reserveCoinDenoms), strconv.FormatUint(uint64(poolTypeId), 10)), "/")
}

// Name returns the pool's name.
func (pool Pool) Name() string {
	return PoolName(pool.ReserveCoinDenoms, pool.TypeId)
}

// Validate validates Pool.
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

// MustMarshalPool returns the Pool bytes. Panics if fails.
func MustMarshalPool(cdc codec.BinaryMarshaler, liquidityPool Pool) []byte {
	return cdc.MustMarshalBinaryBare(&liquidityPool)
}

// MustUnmarshalPool returns the Pool from bytes. Panics if fails.
func MustUnmarshalPool(cdc codec.BinaryMarshaler, value []byte) Pool {
	liquidityPool, err := UnmarshalPool(cdc, value)
	if err != nil {
		panic(err)
	}
	return liquidityPool
}

// UnmarshalPool returns the Pool from bytes.
func UnmarshalPool(cdc codec.BinaryMarshaler, value []byte) (liquidityPool Pool, err error) {
	err = cdc.UnmarshalBinaryBare(value, &liquidityPool)
	return liquidityPool, err
}

// GetReserveAccount returns sdk.AccAddress of the pool's reserve account.
func (pool Pool) GetReserveAccount() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetPoolCoinDenom returns the pool coin's denom.
func (pool Pool) GetPoolCoinDenom() string { return pool.PoolCoinDenom }

// GetPoolId returns id of the pool.
func (pool Pool) GetPoolId() uint64 { return pool.Id }

// Pools is a collection of pools.
type Pools []Pool

func (pools Pools) String() (out string) {
	for _, del := range pools {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// MustMarshalPoolBatch returns the PoolBatch bytes. Panics if fails.
func MustMarshalPoolBatch(cdc codec.BinaryMarshaler, poolBatch PoolBatch) []byte {
	return cdc.MustMarshalBinaryBare(&poolBatch)
}

// UnmarshalPoolBatch returns the PoolBatch from bytes.
func UnmarshalPoolBatch(cdc codec.BinaryMarshaler, value []byte) (poolBatch PoolBatch, err error) {
	err = cdc.UnmarshalBinaryBare(value, &poolBatch)
	return poolBatch, err
}

// MustUnmarshalPoolBatch returns the PoolBatch from bytes. Panics if fails.
func MustUnmarshalPoolBatch(cdc codec.BinaryMarshaler, value []byte) PoolBatch {
	poolBatch, err := UnmarshalPoolBatch(cdc, value)
	if err != nil {
		panic(err)
	}
	return poolBatch
}

// MustMarshalDepositMsgState returns the DepositMsgState bytes. Panics if fails.
func MustMarshalDepositMsgState(cdc codec.BinaryMarshaler, msg DepositMsgState) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// UnmarshalDepositMsgState returns the DepositMsgState from bytes.
func UnmarshalDepositMsgState(cdc codec.BinaryMarshaler, value []byte) (msg DepositMsgState, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalDepositMsgState returns the DepositMsgState from bytes. Panics if fails.
func MustUnmarshalDepositMsgState(cdc codec.BinaryMarshaler, value []byte) DepositMsgState {
	msg, err := UnmarshalDepositMsgState(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalWithdrawMsgState returns the WithdrawMsgState bytes. Panics if fails.
func MustMarshalWithdrawMsgState(cdc codec.BinaryMarshaler, msg WithdrawMsgState) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// UnmarshalWithdrawMsgState returns the WithdrawMsgState from bytes.
func UnmarshalWithdrawMsgState(cdc codec.BinaryMarshaler, value []byte) (msg WithdrawMsgState, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalWithdrawMsgState returns the WithdrawMsgState from bytes. Panics if fails.
func MustUnmarshalWithdrawMsgState(cdc codec.BinaryMarshaler, value []byte) WithdrawMsgState {
	msg, err := UnmarshalWithdrawMsgState(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalSwapMsgState returns the SwapMsgState bytes. Panics if fails.
func MustMarshalSwapMsgState(cdc codec.BinaryMarshaler, msg SwapMsgState) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// UnmarshalSwapMsgState returns the UnmarshalSwapMsgState from bytes.
func UnmarshalSwapMsgState(cdc codec.BinaryMarshaler, value []byte) (msg SwapMsgState, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalSwapMsgState returns the SwapMsgState from bytes. Panics if fails.
func MustUnmarshalSwapMsgState(cdc codec.BinaryMarshaler, value []byte) SwapMsgState {
	msg, err := UnmarshalSwapMsgState(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}
