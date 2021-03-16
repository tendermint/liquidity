package keeper_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestSimulationSwapExecution(t *testing.T) {
	for i := 0; i < 50; i++ {
		if i%10 == 0 {
			fmt.Println("TestSimulationSwapExecution count", i)
		}
		TestSwapExecution(t)
	}
	for i := 0; i < 10; i++ {
		if i%10 == 0 {
			fmt.Println("TestSimulationSwapExecutionFindEdgeCase count", i)
		}
		TestSimulationSwapExecutionFindEdgeCase(t)
	}
}

func TestSimulationSwapExecutionFindEdgeCase(t *testing.T) {
	simapp, ctx := createTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)

	// define test denom X, Y for Liquidity Pool
	denomX := "denomX"
	denomY := "denomY"
	denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// get random X, Y amount for create pool
	param := simapp.LiquidityKeeper.GetParams(ctx)
	X, Y := app.GetRandPoolAmt(r, param.MinInitDepositToPool)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

	// set pool creator account, balance for deposit
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)
	app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeId := types.DefaultPoolTypeId
	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	for i := 0; i < 20; i++ {
		ctx = ctx.WithBlockHeight(int64(i))
		testSwapEdgeCases(t, simapp, ctx, X, Y, depositBalance, addrs)
	}
}

func TestSwapExecution(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	// define test denom X, Y for Liquidity Pool
	denomX := "denomX"
	denomY := "denomY"
	denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

	// get random X, Y amount for create pool
	X, Y := app.GetRandPoolAmt(r, params.MinInitDepositToPool)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

	// set pool creator account, balance for deposit
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)
	app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeId := types.DefaultPoolTypeId
	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolId := lpList[0].Id
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	var XtoY []*types.MsgSwapWithinBatch // buying Y from X
	var YtoX []*types.MsgSwapWithinBatch // selling Y for X

	// make random orders, set buyer, seller accounts for the orders
	XtoY, YtoX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 250, 250)
	buyerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
	sellerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

	for i, msg := range XtoY {
		app.SaveAccountWithFee(simapp, ctx, buyerAccs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
		msg.SwapRequesterAddress = buyerAccs[i].String()
		msg.PoolId = poolId
		msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
	}
	for i, msg := range YtoX {
		app.SaveAccountWithFee(simapp, ctx, sellerAccs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
		msg.SwapRequesterAddress = sellerAccs[i].String()
		msg.PoolId = poolId
		msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
	}

	// begin block, delete and init pool batch
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
		require.NoError(t, err)
	}
	for _, msg := range YtoX {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
		require.NoError(t, err)
	}

	// verify pool batch
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// end block, swap execution
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}

func testSwapEdgeCases(t *testing.T, simapp *app.LiquidityApp, ctx sdk.Context, X, Y sdk.Int, depositBalance sdk.Coins, addrs []sdk.AccAddress) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	//simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	denomX := depositBalance[0].Denom
	denomY := depositBalance[1].Denom

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolId := lpList[0].Id
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	var XtoY []*types.MsgSwapWithinBatch // buying Y from X
	var YtoX []*types.MsgSwapWithinBatch // selling Y for X

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)

	remainingSwapMsgs := simapp.LiquidityKeeper.GetAllNotProcessedPoolBatchSwapMsgStates(ctx, batch)
	if ctx.BlockHeight() == 0 || len(remainingSwapMsgs) == 0 {
		// make random orders, set buyer, seller accounts for the orders
		XtoY, YtoX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 100, 100)
		buyerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
		sellerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

		for i, msg := range XtoY {
			app.SaveAccountWithFee(simapp, ctx, buyerAccs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = buyerAccs[i].String()
			msg.PoolId = poolId
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
		for i, msg := range YtoX {
			app.SaveAccountWithFee(simapp, ctx, sellerAccs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = sellerAccs[i].String()
			msg.PoolId = poolId
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
	}

	// begin block, delete and init pool batch
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, int64(r.Intn(4)))
		require.NoError(t, err)
	}
	for _, msg := range YtoX {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, int64(r.Intn(4)))
		require.NoError(t, err)
	}

	// verify pool batch
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// end block, swap execution
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
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

func TestBadSwapExecution(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	simapp, ctx := app.CreateTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)
	denomX, denomY := types.AlphabeticalDenomPair("denomX", "denomY")

	// add pool creator account
	X, Y := app.GetRandPoolAmt(r, params.MinInitDepositToPool)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	creatorAddr := app.AddRandomTestAddr(simapp, ctx, deposit.Add(params.LiquidityPoolCreationFee...))
	balanceX := simapp.BankKeeper.GetBalance(ctx, creatorAddr, denomX)
	balanceY := simapp.BankKeeper.GetBalance(ctx, creatorAddr, denomY)
	creatorBalance := sdk.NewCoins(balanceX, balanceY)
	require.Equal(t, deposit, creatorBalance)

	// create pool
	createPoolMsg := types.NewMsgCreatePool(creatorAddr, types.DefaultPoolTypeId, creatorBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createPoolMsg)
	require.NoError(t, err)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	offerCoin := sdk.NewCoin(denomX, sdk.NewInt(10000))
	offerCoinFee := types.GetOfferCoinFee(offerCoin, params.SwapFeeRate)
	testAddr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins(offerCoin.Add(offerCoinFee)))

	currentPrice := X.ToDec().Quo(Y.ToDec())
	swapMsg := types.NewMsgSwapWithinBatch(testAddr, 0, types.DefaultSwapTypeId, offerCoin, denomY, currentPrice, params.SwapFeeRate)
	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, swapMsg, 0)
	require.ErrorIs(t, err, types.ErrPoolNotExists)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}
