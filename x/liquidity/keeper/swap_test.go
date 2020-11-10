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
	// TODO: need to set range for avoid min deposit amt errors
	X = sdk.NewInt(int64(r.Float32()*100000000000))
	Y = sdk.NewInt(int64(r.Float32()*100000000000))
	return
}

func TestSimulationSwapExecution(t *testing.T){
	for i:=0; i<100; i++ {
		TestSwapExecution(t)
	}
}

func TestSwapExecution(t *testing.T) {
	// set simapp, random
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX := "denomX"
	denomY := "denomY"
	denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)
	denoms := []string{denomX, denomY}

	// get random X, Y amount for create pool
	X, Y := getRandPoolAmt(r)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	fmt.Println("-------------------------------------------------------")
	fmt.Println("X/Y", X.ToDec().Quo(Y.ToDec()), "X", X, "Y", Y)

	// set pool creator account, balance for deposit
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)  // pool creator
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeIndex := uint32(0)
	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)
	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	poolID := lpList[0].PoolID
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(0), poolID)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	var XtoY []*types.MsgSwap // buying Y from X
	var YtoX []*types.MsgSwap // selling Y for X

	// make random orders, set buyer, seller accounts for the orders
	XtoY, YtoX = GetRandomOrders(denomX, denomY, X, Y, r)
	buyerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
	sellerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

	for i, msg := range XtoY {
		app.SaveAccount(simapp, ctx, buyerAccs[i], sdk.NewCoins(msg.OfferCoin))
		msg.SwapRequester = buyerAccs[i]
		msg.PoolID = poolID
		msg.PoolTypeIndex = poolTypeIndex
	}
	for i, msg := range YtoX {
		app.SaveAccount(simapp, ctx, sellerAccs[i], sdk.NewCoins(msg.OfferCoin))
		msg.SwapRequester = sellerAccs[i]
		msg.PoolID = poolID
		msg.PoolTypeIndex = poolTypeIndex
	}

	// begin block, delete and init pool batch
	simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}
	for _, msg := range YtoX {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}

	// verify pool batch
	batchIndex := simapp.LiquidityKeeper.GetLiquidityPoolBatchIndex(ctx, poolID)
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolID, batchIndex)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// end block, swap execution
	err = simapp.LiquidityKeeper.SwapExecution(ctx, liquidityPoolBatch)
	require.NoError(t, err)
}

func randFloats(min, max float64) float64 {
	return min + rand.Float64() * (max - min)
}

func randRange(r *rand.Rand, min, max int) sdk.Int {
	return sdk.NewInt(int64(r.Intn(max-min) + min))
}

func GetRandomOrders(denomX, denomY string, X, Y sdk.Int, r *rand.Rand) (XtoY, YtoX []*types.MsgSwap){
	currentPrice := X.ToDec().Quo(Y.ToDec())

	XtoYnewSize := int(r.Int31n(20))  // 0~19
	YtoXnewSize := int(r.Int31n(20))  // 0~19

	for i := 0; i < XtoYnewSize; i++ {
		randFloats(0.1, 0.9, )
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(randRange(r, 991, 1009),3))
		orderAmt := X.ToDec().Mul(sdk.NewDecFromIntWithPrec(randRange(r, 1, 100),4))
		orderCoin := sdk.NewCoin(denomX, orderAmt.RoundInt())

		XtoY = append(XtoY, &types.MsgSwap{
			OfferCoin:     orderCoin,
			OrderPrice:    orderPrice,
		})
	}

	for i := 0; i < YtoXnewSize; i++ {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(randRange(r, 991, 1009),3))
		orderAmt := Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(randRange(r, 1, 100),4))
		orderCoin := sdk.NewCoin(denomY, orderAmt.RoundInt())

		YtoX = append(YtoX, &types.MsgSwap{
			OfferCoin:     orderCoin,
			OrderPrice:    orderPrice,
		})
	}
	return
}

func TestGetRandomOrders(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	X, Y := getRandPoolAmt(r)
	XtoY, YtoX := GetRandomOrders("denomX", "denomY", X, Y, r)
	fmt.Println(XtoY, YtoX)
	require.Equal(t, X.ToDec().MulInt64(2).TruncateInt(), X.MulRaw(2))
	require.Equal(t, X.ToDec().MulInt64(2), X.MulRaw(2).ToDec())
	require.NotEqual(t, X.ToDec().MulInt64(2).TruncateInt(), X.MulRaw(2).ToDec())
}