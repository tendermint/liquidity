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
		case bytes.Equal(kvA.Key[:1], types.PoolKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.PoolByReserveAccIndexKeyPrefix):
			var lpA, lpB types.Pool
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpB)
			return fmt.Sprintf("%v\n%v", lpA, lpB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchIndexKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.PoolBatchKeyPrefix):
			var lpbA, lpbB types.PoolBatch
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpbA)
			cdc.MustUnmarshalBinaryBare(kvA.Value, &lpbB)
			return fmt.Sprintf("%v\n%v", lpbA, lpbB)

		default:
			panic(fmt.Sprintf("invalid liquidity key prefix %X", kvA.Key[:1]))
		}
	}
}
