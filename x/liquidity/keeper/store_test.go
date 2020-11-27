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

func TestGetAllLiquidityPoolBatchSwapMsgs(t *testing.T) {
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
	X, Y := app.GetRandPoolAmt(r)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

	// set pool creator account, balance for deposit
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))
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

	var XtoY []*types.MsgSwap // buying Y from X
	var YtoX []*types.MsgSwap // selling Y for X

	// make random orders, set buyer, seller accounts for the orders
	XtoY, YtoX = app.GetRandomOrders(denomX, denomY, X, Y, r, 11, 11)
	buyerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
	sellerAccs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

	poolId := uint64(1)
	pool, found := simapp.LiquidityKeeper.GetLiquidityPool(ctx, poolId)
	require.True(t, found)

	poolBatch, found := simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolId)
	require.Equal(t, uint64(1), poolBatch.SwapMsgIndex)

	for i, msg := range XtoY {
		app.SaveAccount(simapp, ctx, buyerAccs[i], sdk.NewCoins(msg.OfferCoin))
		msg.SwapRequesterAddress = buyerAccs[i].String()
		msg.PoolId = pool.PoolId
		msg.PoolTypeIndex = poolTypeIndex
	}
	for i, msg := range YtoX {
		app.SaveAccount(simapp, ctx, sellerAccs[i], sdk.NewCoins(msg.OfferCoin))
		msg.SwapRequesterAddress = sellerAccs[i].String()
		msg.PoolId = pool.PoolId
		msg.PoolTypeIndex = poolTypeIndex
	}

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY[:10] {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}
	for _, msg := range YtoX[:10] {
		simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg)
	}

	msgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgs(ctx, poolBatch)
	require.Equal(t, 20, len(msgs))

	simapp.LiquidityKeeper.IterateAllLiquidityPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if msg.MsgIndex%2 == 1 {
			simapp.LiquidityKeeper.DeleteLiquidityPoolBatchSwapMsg(ctx, msg.Msg.PoolId, msg.MsgIndex)
		}
		return false
	})

	msgs = simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgs(ctx, poolBatch)
	require.Equal(t, 10, len(msgs))

	poolBatch, found = simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolId)
	require.Equal(t, uint64(21), poolBatch.SwapMsgIndex)

	poolBatch.SwapMsgIndex = uint64(18446744073709551610)
	simapp.LiquidityKeeper.SetLiquidityPoolBatch(ctx, poolBatch)

	simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, XtoY[10])
	simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, YtoX[10])

	msgs = simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgs(ctx, poolBatch)
	require.Equal(t, 12, len(msgs))
	require.Equal(t, XtoY[10], msgs[10].Msg)
	require.Equal(t, YtoX[10], msgs[11].Msg)
	fmt.Println(msgs)
}

func TestGetAllNotProcessedPoolBatchSwapMsgs(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, true)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	price, _ := sdk.NewDecFromStr("1.1")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000)), sdk.NewCoin(denomX, sdk.NewInt(10000)), sdk.NewCoin(denomX, sdk.NewInt(10000))}
	orderPriceList := []sdk.Dec{price, price, price}
	orderAddrList := addrs[1:4]
	batchMsgs, _ := app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	batchMsgs2, batch := app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	require.Equal(t, 3, len(batchMsgs))
	for _, msg := range batchMsgs2 {
		msg.Executed = true
		msg.Succeed = true
		msg.ToDelete = true
	}
	require.Equal(t, 3, len(batchMsgs2))
	simapp.LiquidityKeeper.SetLiquidityPoolBatchSwapMsgs(ctx, poolId, batchMsgs2)

	resultMsgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgs(ctx, batch)
	resultProccessedMsgs := simapp.LiquidityKeeper.GetAllNotProcessedLiquidityPoolBatchSwapMsgs(ctx, batch)
	require.Equal(t, 6, len(resultMsgs))
	require.Equal(t, 3, len(resultProccessedMsgs))

}
