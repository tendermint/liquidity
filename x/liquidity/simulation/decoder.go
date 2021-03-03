package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding liquidity type.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.LiquidityPoolKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.LiquidityPoolByReserveIndexKeyPrefix):
			var lpA, lpB types.LiquidityPool
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpB)
			return fmt.Sprintf("%v\n%v", lpA, lpB)

		case bytes.Equal(kvA.Key[:1], types.LiquidityPoolBatchIndexKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.LiquidityPoolBatchKeyPrefix):
			var lpbA, lpbB types.LiquidityPoolBatch
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpbA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpbB)
			return fmt.Sprintf("%v\n%v", lpbA, lpbB)

		default:
			panic(fmt.Sprintf("invalid liquidity key prefix %X", kvA.Key[:1]))
		}
	}
}
