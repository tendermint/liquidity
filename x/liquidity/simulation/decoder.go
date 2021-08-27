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
			var poolA, poolB types.Pool
			cdc.MustUnmarshal(kvA.Value, &poolA)
			cdc.MustUnmarshal(kvA.Value, &poolB)
			return fmt.Sprintf("%v\n%v", poolA, poolB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchIndexKeyPrefix),
			bytes.Equal(kvA.Key[:1], types.PoolBatchKeyPrefix):
			var batchA, batchB types.PoolBatch
			cdc.MustUnmarshal(kvA.Value, &batchA)
			cdc.MustUnmarshal(kvA.Value, &batchB)
			return fmt.Sprintf("%v\n%v", batchA, batchB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchDepositMsgStateIndexKeyPrefix):
			var msgA, msgB types.MsgDepositWithinBatch
			cdc.MustUnmarshal(kvA.Value, &msgA)
			cdc.MustUnmarshal(kvA.Value, &msgB)
			return fmt.Sprintf("%v\n%v", msgA, msgB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchWithdrawMsgStateIndexKeyPrefix):
			var msgA, msgB types.MsgWithdrawWithinBatch
			cdc.MustUnmarshal(kvA.Value, &msgA)
			cdc.MustUnmarshal(kvA.Value, &msgB)
			return fmt.Sprintf("%v\n%v", msgA, msgB)

		case bytes.Equal(kvA.Key[:1], types.PoolBatchSwapMsgStateIndexKeyPrefix):
			var msgA, msgB types.MsgSwapWithinBatch
			cdc.MustUnmarshal(kvA.Value, &msgA)
			cdc.MustUnmarshal(kvA.Value, &msgB)
			return fmt.Sprintf("%v\n%v", msgA, msgB)

		//
		// panic: proto: wrong wireType = 2 for field "" [recovered]
		// panic: proto: wrong wireType = 2 for field ""
		//
		// case bytes.Equal(kvA.Key[:1], types.PoolBatchDepositMsgStateIndexKeyPrefix):
		// 	var msgStateA, msgStateB types.DepositMsgState
		// 	cdc.MustUnmarshal(kvA.Value, &msgStateA)
		// 	cdc.MustUnmarshal(kvA.Value, &msgStateB)
		// 	return fmt.Sprintf("%v\n%v", msgStateA, msgStateB)

		// case bytes.Equal(kvA.Key[:1], types.PoolBatchWithdrawMsgStateIndexKeyPrefix):
		// 	var msgStateA, msgStateB types.WithdrawMsgState
		// 	cdc.MustUnmarshal(kvA.Value, &msgStateA)
		// 	cdc.MustUnmarshal(kvA.Value, &msgStateB)
		// 	return fmt.Sprintf("%v\n%v", msgStateA, msgStateB)

		// case bytes.Equal(kvA.Key[:1], types.PoolBatchSwapMsgStateIndexKeyPrefix):
		// 	var msgStateA, msgStateB types.SwapMsgState
		// 	cdc.MustUnmarshal(kvA.Value, &msgStateA)
		// 	cdc.MustUnmarshal(kvA.Value, &msgStateB)
		// 	return fmt.Sprintf("%v\n%v", msgStateA, msgStateB)

		default:
			panic(fmt.Sprintf("invalid liquidity key prefix %X", kvA.Key[:1]))
		}
	}
}
