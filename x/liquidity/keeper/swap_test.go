package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"math/rand"
	"testing"
	"time"
)

func getRandPoolAmt(r *rand.Rand) (X, Y sdk.Int){
	X = sdk.NewInt(int64(r.Float32()*100000000))
	Y = sdk.NewInt(int64(r.Float32()*100000000))
	return
}


func TestSwapExecution(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	simapp, ctx := createTestInput()

	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	poolTypeIndex := uint32(0)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))

	denomX := "uXXX"
	denomY := "uYYY"
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
	poolID := lpList[0].PoolID
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(0), poolID)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	var XtoY []*types.MsgSwap // buying Y from X
	var YtoX []*types.MsgSwap // selling Y for X
	// TODO: make random orders, set to batch
	// GetRandomOrders()

	// TODO: init SwapRequesterAcc list, set balance
	for _, msg := range XtoY {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}
	for _, msg := range YtoX {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}

	simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)
	batchIndex := simapp.LiquidityKeeper.GetLiquidityPoolBatchIndex(ctx, poolID)
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolID, batchIndex)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// TODO: SwapExecution()
	//err = simapp.LiquidityKeeper.SwapExecution(ctx, liquidityPoolBatch)
	//require.NoError(t, err)

	// TODO: invariant check
}

func randFloats(min, max float64) float64 {
	return min + rand.Float64() * (max - min)
}

func randRange(r *rand.Rand, min, max int) sdk.Int {
	return sdk.NewInt(int64(r.Intn(max-min) + min))
}

func GetRandomOrders(X, Y sdk.Int, r *rand.Rand) (XtoY, YtoX []*types.MsgSwap){
	currentPrice := X.ToDec().Quo(Y.ToDec())
	XtoYnewSize := int(r.Float32()*200)
	YtoXnewSize := int(r.Float32()*200)

	for i := 0; i < XtoYnewSize; i++ {
		fmt.Println(sdk.NewDecFromIntWithPrec(randRange(r, 991, 1009),3), sdk.NewDecFromIntWithPrec(randRange(r, 1, 100),4))
	}

	for i := 0; i < XtoYnewSize; i++ {
		randFloats(0.1, 0.9, )
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(randRange(r, 991, 1009),3))
		orderAmt := X.ToDec().Mul(sdk.NewDecFromIntWithPrec(randRange(r, 1, 100),4))
		fmt.Println(orderPrice, orderAmt)
	}

	for i := 0; i < YtoXnewSize; i++ {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(randRange(r, 991, 1009),3))
		orderAmt := Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(randRange(r, 1, 100),4))
		fmt.Println(orderPrice, orderAmt)
	}
	return
}

func TestGetRandomOrders(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	X, Y := getRandPoolAmt(r)
	GetRandomOrders(X, Y, r)
}