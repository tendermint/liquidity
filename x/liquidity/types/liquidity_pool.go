package types

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LiquidityPoolLegacy struct {
	PoolID            uint64         // index of this liquidity pool
	PoolTypeIndex     uint32         // pool type of this liquidity pool
	ReserveCoinDenoms []string       // list of reserve coin denoms for this liquidity pool
	ReserveAccount    sdk.AccAddress // module account address for this liquidity pool to store reserve coins
	PoolCoinDenom     string         // denom of pool coin for this liquidity pool
	SwapFeeRate       sdk.Dec        // swap fee rate for every executed swap on this liquidity pool
	PoolFeeRate       sdk.Dec        // liquidity pool fee rate for swaps consumed liquidity from this liquidity pool
	BatchSize         uint32         // size of each batch as a number of block heights  // TODO: set default Param
	//LastBatchIndex     uint64         // index of the last batch of this liquidity pool  // TODO: separate
}

// need to validate alphabetical ordering of ReserveCoinDenoms when New() and Store
// Denominations can be 3 ~ 128 characters long and support letters, followed by either
// a letter, a number or a separator ('/').
// reDnmString = `[a-zA-Z][a-zA-Z0-9/]{2,127}`
func (lp LiquidityPool) GetPoolKey() string {
	return GetPoolKey(lp.ReserveCoinDenoms, lp.PoolTypeIndex)
}

func GetPoolKey(reserveCoinDenoms []string, poolTypeIndex uint32) string {
	return strings.Join(append(reserveCoinDenoms, strconv.FormatUint(uint64(poolTypeIndex), 10)), "-")
}

// NewLiquidityPool creates a new liquidityPool object
func NewLiquidityPool() LiquidityPool {
	return LiquidityPool{}
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

func (lp LiquidityPool) GetReserveAccount() sdk.AccAddress { return lp.ReserveAccount }
func (lp LiquidityPool) GetPoolCoinDenom() string          { return lp.PoolCoinDenom }
func (lp LiquidityPool) GetPoolID() uint64                 { return lp.PoolID }

// String returns a human readable string representation of a LiquidityPool.
//func (lp LiquidityPool) String() string {
//	out, _ := yaml.Marshal(lp)
//	return string(out)
//}

// LiquidityPools is a collection of liquidityPools
type LiquidityPools []LiquidityPool

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
