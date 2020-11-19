package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO: template, delete
func TestGetLiquidityPoolBatchIndexKey(t *testing.T) {
	poolId := uint64(1)
	res := GetLiquidityPoolBatchIndexKey(poolId)
	require.NotNil(t, res)
}
