package keeper_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

func TestLiquidityPool(t *testing.T) {
	app, ctx := createTestInput()
	lp := types.LiquidityPool{
		PoolID:         1,
		PoolTypeIndex:  1,
		ReserveAccount: nil,
		PoolCoinDenom:  "poolCoin",
	}
	lp.ReserveCoinDenoms = append(lp.ReserveCoinDenoms, "a")
	lp.ReserveCoinDenoms = append(lp.ReserveCoinDenoms, "b")
	app.LiquidityKeeper.SetLiquidityPool(ctx, lp)

	lpGet, found := app.LiquidityKeeper.GetLiquidityPool(ctx, 1)
	require.True(t, found)
	require.Equal(t, lp, lpGet)
}

// TODO: WIP
func TestCreateLiquidityPool(t *testing.T) {
	app, ctx := createTestInput()


}
