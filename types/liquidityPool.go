package types

import (
	yaml "gopkg.in/yaml.v2"
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
func (lp LiquidityPool) getPoolKey() string {
	return GetPoolKey(lp.ReserveCoinDenoms, lp.PoolTypeIndex)
}

func GetPoolKey(reserveCoinDenoms []string, poolTypeIndex uint32) string {
	return strings.Join(append(reserveCoinDenoms, string(poolTypeIndex)), "-")
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
func (lp LiquidityPool) String() string {
	out, _ := yaml.Marshal(lp)
	return string(out)
}

// LiquidityPools is a collection of liquidityPools
type LiquidityPools []LiquidityPool

func (lps LiquidityPools) String() (out string) {
	for _, del := range lps {
		out += del.String() + "\n"
	}

	return strings.TrimSpace(out)
}
