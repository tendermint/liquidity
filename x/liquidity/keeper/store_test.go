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
	require.Equal(t, uint64(0), poolBatch.SwapMsgIndex)

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
	require.Equal(t, uint64(20), poolBatch.SwapMsgIndex)

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
