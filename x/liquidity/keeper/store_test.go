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
	poolTypeIndex := types.DefaultPoolTypeIndex
	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, depositBalance)
	err, _ := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
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
		app.SaveAccountWithFee(simapp, ctx, buyerAccs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
		msg.SwapRequesterAddress = buyerAccs[i].String()
		msg.PoolId = pool.PoolId
		msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
	}
	for i, msg := range YtoX {
		app.SaveAccountWithFee(simapp, ctx, sellerAccs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
		msg.SwapRequesterAddress = sellerAccs[i].String()
		msg.PoolId = pool.PoolId
		msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
	}

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY[:10] {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
		require.NoError(t, err)
	}
	for _, msg := range YtoX[:10] {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
		require.NoError(t, err)
	}

	msgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, poolBatch)
	require.Equal(t, 20, len(msgs))

	simapp.LiquidityKeeper.IterateAllLiquidityPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if msg.MsgIndex%2 == 1 {
			simapp.LiquidityKeeper.DeleteLiquidityPoolBatchSwapMsg(ctx, msg.Msg.PoolId, msg.MsgIndex)
		}
		return false
	})

	msgs = simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, poolBatch)
	require.Equal(t, 10, len(msgs))

	poolBatch, found = simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolId)
	require.Equal(t, uint64(21), poolBatch.SwapMsgIndex)

	poolBatch.SwapMsgIndex = uint64(18446744073709551610)
	simapp.LiquidityKeeper.SetLiquidityPoolBatch(ctx, poolBatch)

	simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, XtoY[10], 0)
	simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, YtoX[10], 0)

	msgs = simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, poolBatch)
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
		msg.Succeeded = true
		msg.ToBeDeleted = true
	}
	require.Equal(t, 3, len(batchMsgs2))
	simapp.LiquidityKeeper.SetLiquidityPoolBatchSwapMsgPointers(ctx, poolId, batchMsgs2)

	resultMsgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, batch)
	resultProccessedMsgs := simapp.LiquidityKeeper.GetAllNotProcessedLiquidityPoolBatchSwapMsgs(ctx, batch)
	require.Equal(t, 6, len(resultMsgs))
	require.Equal(t, 3, len(resultProccessedMsgs))

}

func TestIterateAllBatchMsgs(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	denomA, denomB := types.AlphabeticalDenomPair("denomA", "denomB")

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)
	A := sdk.NewInt(500000000)
	B := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	poolId2 := app.TestCreatePool(t, simapp, ctx, A, B, denomA, denomB, addrs[4])
	batch, found := simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, poolId)
	require.True(t, found)

	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)

	price, _ := sdk.NewDecFromStr("1.1")
	priceY, _ := sdk.NewDecFromStr("1.2")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000))}
	offerCoinListY := []sdk.Coin{sdk.NewCoin(denomY, sdk.NewInt(5000))}

	orderPriceList := []sdk.Dec{price}
	orderPriceListY := []sdk.Dec{priceY}
	orderAddrList := addrs[1:2]
	orderAddrListY := addrs[2:3]

	offerCoinList2 := []sdk.Coin{sdk.NewCoin(denomA, sdk.NewInt(5000))}

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	app.TestDepositPool(t, simapp, ctx, A, B.QuoRaw(10), addrs[4:5], poolId2, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(1000), addrs[4:5], poolId2, false)
	app.TestSwapPool(t, simapp, ctx, offerCoinList2, orderPriceList, addrs[4:5], poolId2, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	app.TestDepositPool(t, simapp, ctx, A, B.QuoRaw(10), addrs[4:5], poolId2, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(1000), addrs[4:5], poolId2, false)
	app.TestSwapPool(t, simapp, ctx, offerCoinList2, orderPriceList, addrs[4:5], poolId2, true)

	// next block,
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[2:3], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[2:3], poolId, false)

	depositMsgsRemaining := simapp.LiquidityKeeper.GetAllRemainingLiquidityPoolBatchDepositMsgs(ctx, batch)
	require.Equal(t, 0, len(depositMsgsRemaining))

	var depositMsgs []types.BatchPoolDepositMsg
	simapp.LiquidityKeeper.IterateAllBatchDepositMsgs(ctx, func(msg types.BatchPoolDepositMsg) bool {
		depositMsgs = append(depositMsgs, msg)
		return false
	})
	require.Equal(t, 4, len(depositMsgs))

	depositMsgs[0].ToBeDeleted = true
	simapp.LiquidityKeeper.SetLiquidityPoolBatchDepositMsgs(ctx, poolId, []types.BatchPoolDepositMsg{depositMsgs[0]})
	depositMsgsNotToDelete := simapp.LiquidityKeeper.GetAllNotToDeleteLiquidityPoolBatchDepositMsgs(ctx, batch)
	require.Equal(t, 3, len(depositMsgsNotToDelete))

	var withdrawMsgs []types.BatchPoolWithdrawMsg
	simapp.LiquidityKeeper.IterateAllBatchWithdrawMsgs(ctx, func(msg types.BatchPoolWithdrawMsg) bool {
		withdrawMsgs = append(withdrawMsgs, msg)
		return false
	})
	withdrawMsgs[0].ToBeDeleted = true
	simapp.LiquidityKeeper.SetLiquidityPoolBatchWithdrawMsgs(ctx, poolId, withdrawMsgs[0:1])

	withdrawMsgsNotToDelete := simapp.LiquidityKeeper.GetAllNotToDeleteLiquidityPoolBatchWithdrawMsgs(ctx, batch)
	require.Equal(t, 4, len(withdrawMsgs))
	require.Equal(t, 3, len(withdrawMsgsNotToDelete))
	require.NotEqual(t, withdrawMsgsNotToDelete, withdrawMsgs)

	app.TestDepositPool(t, simapp, ctx, A, B.QuoRaw(10), addrs[4:5], poolId2, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(1000), addrs[4:5], poolId2, false)

	depositMsgs = simapp.LiquidityKeeper.GetAllBatchDepositMsgs(ctx)
	require.Equal(t, 5, len(depositMsgs))
	withdrawMsgs = simapp.LiquidityKeeper.GetAllBatchWithdrawMsgs(ctx)
	require.Equal(t, 5, len(depositMsgs))

	var depositMsgs2 []types.BatchPoolDepositMsg
	simapp.LiquidityKeeper.IterateAllBatchDepositMsgs(ctx, func(msg types.BatchPoolDepositMsg) bool {
		depositMsgs2 = append(depositMsgs2, msg)
		return false
	})

	var withdrawMsgs2 []types.BatchPoolWithdrawMsg
	simapp.LiquidityKeeper.IterateAllBatchWithdrawMsgs(ctx, func(msg types.BatchPoolWithdrawMsg) bool {
		withdrawMsgs2 = append(withdrawMsgs2, msg)
		return false
	})

	require.Equal(t, 5, len(depositMsgs2))

	require.Equal(t, 5, len(withdrawMsgs2))

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	depositMsgsRemaining = simapp.LiquidityKeeper.GetAllRemainingLiquidityPoolBatchDepositMsgs(ctx, batch)
	require.Equal(t, 0, len(depositMsgsRemaining))

	// next block,
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	var depositMsgs3 []types.BatchPoolDepositMsg
	simapp.LiquidityKeeper.IterateAllBatchDepositMsgs(ctx, func(msg types.BatchPoolDepositMsg) bool {
		depositMsgs3 = append(depositMsgs3, msg)
		return false
	})
	require.Equal(t, 0, len(depositMsgs3))

	var withdrawMsgs3 []types.BatchPoolWithdrawMsg
	simapp.LiquidityKeeper.IterateAllBatchWithdrawMsgs(ctx, func(msg types.BatchPoolWithdrawMsg) bool {
		withdrawMsgs3 = append(withdrawMsgs3, msg)
		return false
	})
	require.Equal(t, 0, len(withdrawMsgs3))

	app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	app.TestSwapPool(t, simapp, ctx, offerCoinListY, orderPriceListY, orderAddrListY, poolId, false)
	app.TestSwapPool(t, simapp, ctx, offerCoinList2, orderPriceList, addrs[4:5], poolId2, false)

	swapMsgsPool1 := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgs(ctx, batch)
	require.Equal(t, 4, len(swapMsgsPool1))

	swapMsg, found := simapp.LiquidityKeeper.GetLiquidityPoolBatchSwapMsg(ctx, batch.PoolId, 1)
	require.True(t, found)
	require.Equal(t, swapMsg, swapMsgsPool1[0])

	var swapMsgsAllPool []types.BatchPoolSwapMsg
	simapp.LiquidityKeeper.IterateAllBatchSwapMsgs(ctx, func(msg types.BatchPoolSwapMsg) bool {
		swapMsgsAllPool = append(swapMsgsAllPool, msg)
		return false
	})
	require.Equal(t, 5, len(swapMsgsAllPool))

	swapMsgsAllPool = simapp.LiquidityKeeper.GetAllBatchSwapMsgs(ctx)
	require.Equal(t, 5, len(swapMsgsAllPool))
	require.Equal(t, swapMsgsPool1, swapMsgsAllPool[:len(swapMsgsPool1)])

	swapMsgsAllPool[1].Executed = true
	simapp.LiquidityKeeper.SetLiquidityPoolBatchSwapMsgs(ctx, poolId, swapMsgsAllPool[1:2])

	reminingSwapMsgs := simapp.LiquidityKeeper.GetAllRemainingLiquidityPoolBatchSwapMsgs(ctx, batch)
	require.Equal(t, 1, len(reminingSwapMsgs))

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	// next block,
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	var swapMsg2 []types.BatchPoolSwapMsg
	simapp.LiquidityKeeper.IterateAllBatchSwapMsgs(ctx, func(msg types.BatchPoolSwapMsg) bool {
		swapMsg2 = append(swapMsg2, msg)
		return false
	})
	require.Equal(t, 0, len(swapMsg2))

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	genesis := simapp.LiquidityKeeper.ExportGenesis(ctx)
	simapp.LiquidityKeeper.InitGenesis(ctx, *genesis)
	genesisNew := simapp.LiquidityKeeper.ExportGenesis(ctx)
	require.Equal(t, genesis, genesisNew)

	simapp.LiquidityKeeper.DeleteLiquidityPoolBatch(ctx, batch)
	batch, found = simapp.LiquidityKeeper.GetLiquidityPoolBatch(ctx, batch.PoolId)
	require.Equal(t, types.LiquidityPoolBatch{}, batch)
	require.False(t, found)
}
