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
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PoolKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.PoolByReserveAccIndexKeyPrefix):
			var lpA, lpB types.Pool
			cdc.MustUnmarshal(kvA.Value, &lpA)
			cdc.MustUnmarshal(kvA.Value, &lpB)
			return fmt.Sprintf("%v\n%v", lpA, lpB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchIndexKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.PoolBatchKeyPrefix):
			var lpbA, lpbB types.PoolBatch
			cdc.MustUnmarshal(kvA.Value, &lpbA)
			cdc.MustUnmarshal(kvA.Value, &lpbB)
			return fmt.Sprintf("%v\n%v", lpbA, lpbB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchDepositMsgStateIndexKeyPrefix):
			var lpbA, lpbB types.DepositMsgState
			lpbA = types.MustUnmarshalDepositMsgState(cdc, kvA.Value)
			lpbB = types.MustUnmarshalDepositMsgState(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", lpbA, lpbB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchWithdrawMsgStateIndexKeyPrefix):
			var lpbA, lpbB types.WithdrawMsgState
			lpbA = types.MustUnmarshalWithdrawMsgState(cdc, kvA.Value)
			lpbB = types.MustUnmarshalWithdrawMsgState(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", lpbA, lpbB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchSwapMsgStateIndexKeyPrefix):
			var lpbA, lpbB types.SwapMsgState
			lpbA = types.MustUnmarshalSwapMsgState(cdc, kvA.Value)
			lpbB = types.MustUnmarshalSwapMsgState(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", lpbA, lpbB)

		default:
			panic(fmt.Sprintf("invalid liquidity key prefix %X", kvA.Key[:1]))
		}
	}
}
