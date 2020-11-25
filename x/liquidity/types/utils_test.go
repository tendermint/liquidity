package types_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

func TestAlphabeticalDenomPair(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"
	afterDenomA, afterDenomB := types.AlphabeticalDenomPair(denomA, denomB)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)

	afterDenomA, afterDenomB = types.AlphabeticalDenomPair(denomB, denomA)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)
}

func TestStringInSlice(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"
	denomC := "uCoinC"
	denoms := []string{denomA, denomB}
	require.True(t, types.StringInSlice(denomA, denoms))
	require.True(t, types.StringInSlice(denomB, denoms))
	require.False(t, types.StringInSlice(denomC, denoms))
}

func TestGetPoolReserveAcc(t *testing.T) {
	reserveAcc := types.GetPoolReserveAcc("denomX-denomY-1")
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos18gvpvm3lrzx6rs6yq5c6klnye2t5qumm3v3re8", reserveAcc.String())
	require.Equal(t, "cosmos18gvpvm3lrzx6rs6yq5c6klnye2t5qumm3v3re8", types.GetPoolCoinDenom(reserveAcc))
}
