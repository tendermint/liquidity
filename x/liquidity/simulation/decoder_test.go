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
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	liquidityPool := types.LiquidityPool{}
	liquidityPool.PoolId = 1
	liquidityPoolBatch := types.NewLiquidityPoolBatch(1, 1)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.LiquidityPoolKeyPrefix, Value: cdc.MustMarshalBinaryBare(&liquidityPool)},
			{Key: types.LiquidityPoolBatchKeyPrefix, Value: cdc.MustMarshalBinaryBare(&liquidityPoolBatch)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"LiquidityPool", fmt.Sprintf("%v\n%v", liquidityPool, liquidityPool)},
		{"LiquidityPoolBatch", fmt.Sprintf("%v\n%v", liquidityPoolBatch, liquidityPoolBatch)},
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
