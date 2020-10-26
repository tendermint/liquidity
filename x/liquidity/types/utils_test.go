package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlphabeticalDenomPair(t *testing.T) {
	denomA := "A"
	denomB := "B"
	afterDenomA, afterDenomB := AlphabeticalDenomPair(denomA, denomB)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)

	afterDenomA, afterDenomB = AlphabeticalDenomPair(denomB, denomA)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)
}
