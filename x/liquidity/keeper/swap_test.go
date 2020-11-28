package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"math/rand"
	"testing"
	"time"
)

func TestSimulationSwapExecution(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println("test count", i+1)
		TestSwapExecution(t)
	}
}

func TestSwapExecution(t *testing.T) {
	// TODO: to with simulation, invariants, ransim
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	// define test denom X, Y for Liquidity Pool
	denomX := "denomX"
	denomY := "denomY"
	denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)
	denoms := []string{denomX, denomY}

	// get random X, Y amount for create pool
	param := simapp.LiquidityKeeper.GetParams(ctx)
	X, Y := app.GetRandPoolAmt(r, param.MinInitDepositToPool)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	fmt.Println("-------------------------------------------------------")
	fmt.Println("X/Y", X.ToDec().Quo(Y.ToDec()), "X", X, "Y", Y)

	// set pool creator account, balance for deposit
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)
	app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeIndex := types.DefaultPoolTypeIndex
	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)
	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	poolId := lpList[0].PoolId
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	var XtoY []*types.MsgSwap // buying Y from X
	var YtoX []*types.MsgSwap // selling Y for X

	// make random orders, set buyer, seller accounts for the orders
	XtoY, YtoX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 50, 50)
	buyerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
	sellerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

	for i, msg := range XtoY {
		app.SaveAccount(simapp, ctx, buyerAccs[i], sdk.NewCoins(msg.OfferCoin))
		msg.SwapRequesterAddress = buyerAccs[i].String()
		msg.PoolId = poolId
		msg.PoolTypeIndex = poolTypeIndex
	}
	for i, msg := range YtoX {
		app.SaveAccount(simapp, ctx, sellerAccs[i], sdk.NewCoins(msg.OfferCoin))
		msg.SwapRequesterAddress = sellerAccs[i].String()
		msg.PoolId = poolId
		msg.PoolTypeIndex = poolTypeIndex
	}

	// begin block, delete and init pool batch
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	//simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}
	for _, msg := range YtoX {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}

	// verify pool batch
	//batchIndex := simapp.LiquidityKeeper.GetLiquidityPoolBatchIndex(ctx, poolId)
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolId)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// end block, swap execution
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	//err = simapp.LiquidityKeeper.SwapExecution(ctx, liquidityPoolBatch)
	require.NoError(t, err)
}

func TestGetRandomOrders(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// get random X, Y amount for create pool
	X, Y := app.GetRandPoolAmt(r, types.DefaultMinInitDepositToPool)
	XtoY, YtoX := app.GetRandomSizeOrders("denomX", "denomY", X, Y, r, 50, 50)
	fmt.Println(XtoY, YtoX)
	require.Equal(t, X.ToDec().MulInt64(2).TruncateInt(), X.MulRaw(2))
	require.Equal(t, X.ToDec().MulInt64(2), X.MulRaw(2).ToDec())
	require.NotEqual(t, X.ToDec().MulInt64(2).TruncateInt(), X.MulRaw(2).ToDec())
}
