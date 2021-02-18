package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/x/liquidity/simulation"
)

func TestParamChanges(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)

	expected := []struct {
		composedKey string
		key         string
		simValue    string
		subspace    string
	}{
		{"liquidity/KeyMinInitDepositToPool", "KeyMinInitDepositToPool", "\"12498081\"", "liquidity"},
		{"liquidity/KeyInitPoolCoinMintAmount", "KeyInitPoolCoinMintAmount", "\"40727887\"", "liquidity"},
		{"liquidity/KeySwapFeeRate", "KeySwapFeeRate", "\"0.461190000000000000\"", "liquidity"},
		{"liquidity/KeyWithdrawFeeRate", "KeyWithdrawFeeRate", "\"0.934590000000000000\"", "liquidity"},
		{"liquidity/KeyMaxOrderAmountRatio", "KeyMaxOrderAmountRatio", "\"0.112010000000000000\"", "liquidity"},
		{"liquidity/KeyUnitBatchSize", "KeyUnitBatchSize", "\"9\"", "liquidity"},
	}

	paramChanges := simulation.ParamChanges(r)

	require.Len(t, paramChanges, 6)

	for i, p := range paramChanges {
		require.Equal(t, expected[i].composedKey, p.ComposedKey())
		require.Equal(t, expected[i].key, p.Key())
		require.Equal(t, expected[i].simValue, p.SimValue()(r))
		require.Equal(t, expected[i].subspace, p.Subspace())
	}
}
