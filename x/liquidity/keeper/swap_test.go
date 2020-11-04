package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"math/rand"
	"testing"
)

func getRandPoolAmt(r *rand.Rand) (X, Y sdk.Int){
	X = sdk.NewInt(r.Int63()).Mul(sdk.NewInt(1000000))
	Y = sdk.NewInt(r.Int63()).Mul(sdk.NewInt(1000000))
	return
}


func TestSwapExecution(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)

	simapp, ctx := createTestInput()

	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	poolTypeIndex := uint32(0)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))

	denomX := "X"
	denomY := "Y"
	denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

	denoms := []string{denomX, denomY}

	X, Y := getRandPoolAmt(r)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)

	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	poolID := simapp.LiquidityKeeper.GetNextLiquidityPoolID(ctx)-1
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(0), lpList[0].PoolID)
	require.Equal(t, poolID, lpList[0].PoolID)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextLiquidityPoolID(ctx))
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)



	// TODO: set order msgs to batch
	//var XtoY []*types.MsgSwap // buying Y from X
	//var YtoX []*types.MsgSwap // selling Y for X

	// TODO: SwapExecution()

	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolID)
	require.False(t, found)

	err = simapp.LiquidityKeeper.SwapExecution(ctx, liquidityPoolBatch)
	require.NoError(t, err)

	// TODO: invariant check


}