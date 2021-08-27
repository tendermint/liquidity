package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/tendermint/liquidity/x/liquidity/simulation"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestDecodeLiquidityStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	pool := types.Pool{
		Id:                    uint64(1),
		TypeId:                uint32(1),
		ReserveCoinDenoms:     []string{"dzkiv", "imwo"},
		ReserveAccountAddress: "cosmos1ldxcd2qkjnhhu4avkpt378aupqazex6qh0eg20",
		PoolCoinDenom:         "poolFB4D86A81694EF7E57ACB0571F1FBC083A2C9B403D6127DFA7B2D9212EA85D72",
	}
	batch := types.NewPoolBatch(1, 1)
	// depositMsgState := types.DepositMsgState{
	// 	MsgHeight:   int64(50),
	// 	MsgIndex:    uint64(1),
	// 	Executed:    true,
	// 	Succeeded:   true,
	// 	ToBeDeleted: true,
	// 	Msg:         &types.MsgDepositWithinBatch{PoolId: uint64(1)},
	// }
	// withdrawMsgState := types.WithdrawMsgState{
	// 	MsgHeight:   int64(50),
	// 	MsgIndex:    uint64(1),
	// 	Executed:    true,
	// 	Succeeded:   true,
	// 	ToBeDeleted: true,
	// 	Msg:         &types.MsgWithdrawWithinBatch{PoolId: uint64(1)},
	// }
	// swapMsgState := types.SwapMsgState{
	// 	MsgHeight:   int64(50),
	// 	MsgIndex:    uint64(1),
	// 	Executed:    true,
	// 	Succeeded:   true,
	// 	ToBeDeleted: true,
	// 	Msg:         &types.MsgSwapWithinBatch{PoolId: uint64(1)},
	// }
	depositBatch := types.MsgDepositWithinBatch{}
	withdrawBatch := types.MsgWithdrawWithinBatch{}
	swapBatch := types.MsgSwapWithinBatch{}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.PoolKeyPrefix, Value: cdc.MustMarshal(&pool)},
			{Key: types.PoolByReserveAccIndexKeyPrefix, Value: cdc.MustMarshal(&pool)},
			{Key: types.PoolBatchIndexKeyPrefix, Value: cdc.MustMarshal(&batch)},
			{Key: types.PoolBatchKeyPrefix, Value: cdc.MustMarshal(&batch)},
			{Key: types.PoolBatchDepositMsgStateIndexKeyPrefix, Value: cdc.MustMarshal(&depositBatch)},
			{Key: types.PoolBatchWithdrawMsgStateIndexKeyPrefix, Value: cdc.MustMarshal(&withdrawBatch)},
			{Key: types.PoolBatchSwapMsgStateIndexKeyPrefix, Value: cdc.MustMarshal(&swapBatch)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Pool", fmt.Sprintf("%v\n%v", pool, pool)},
		{"PoolByReserveAccIndex", fmt.Sprintf("%v\n%v", pool, pool)},
		{"PoolBatchIndex", fmt.Sprintf("%v\n%v", batch, batch)},
		{"PoolBatchKey", fmt.Sprintf("%v\n%v", batch, batch)},
		{"PoolBatchDepositMsgStateIndex PoolBatch", fmt.Sprintf("%v\n%v", depositBatch, depositBatch)},
		{"PoolBatchWithdrawMsgStateIndex", fmt.Sprintf("%v\n%v", withdrawBatch, withdrawBatch)},
		{"PoolBatchSwapMsgStateIndex", fmt.Sprintf("%v\n%v", swapBatch, swapBatch)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}

// func TestDecodeMsgStateArray(t *testing.T) {
// 	cdc := simapp.MakeTestEncodingConfig().Marshaler

// 	state := []types.DepositMsgState{
// 		types.DepositMsgState{
// 			MsgHeight:   int64(50),
// 			MsgIndex:    uint64(1),
// 			Executed:    true,
// 			Succeeded:   true,
// 			ToBeDeleted: true,
// 			Msg: &types.MsgDepositWithinBatch{
// 				PoolId:           uint64(25),
// 				DepositorAddress: "cosmos1p63k4hmsw3z5frfpxestkp690wc8whcus70dff",
// 				DepositCoins: sdk.NewCoins(
// 					sdk.NewInt64Coin("mklp", 96071022),
// 					sdk.NewInt64Coin("uzme", 78171341),
// 				),
// 			},
// 		},
// 	}

// 	bz, err := cdc.Marshal(state)

// }
