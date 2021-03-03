package types_test

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

func PointerListToValueList(pointerList []*interface{}) (valueList []interface{}) {
	for _, i := range pointerList {
		valueList = append(valueList, *i)
	}
	return valueList
}

func TestSwapScenario(t *testing.T) {
	// init test app and context
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	X := params.MinInitDepositToPool
	Y := params.MinInitDepositToPool

	// init addresses for the test
	addrs := app.AddTestAddrs(simapp, ctx, 20, params.LiquidityPoolCreationFee)

	// Create pool
	// The create pool msg is not run in batch, but is processed immediately.
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// In case of deposit, withdraw, and swap msg, unlike other normal tx msgs,
	// collect them in the batch and perform an execution at once at the endblock.

	// add a deposit to pool and run batch execution on endblock
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:2], poolId, true)

	// next block, reinitialize batch and increase batchIndex at beginBlocker,
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// Create swap msg for test purposes and put it in the batch.
	price, _ := sdk.NewDecFromStr("1.1")
	priceY, _ := sdk.NewDecFromStr("1.2")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000))}
	offerCoinListY := []sdk.Coin{sdk.NewCoin(denomY, sdk.NewInt(5000))}
	orderPriceList := []sdk.Dec{price}
	orderPriceListY := []sdk.Dec{priceY}
	orderAddrList := addrs[1:2]
	orderAddrListY := addrs[2:3]
	_, batch := app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	_, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	_, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	_, batch = app.TestSwapPool(t, simapp, ctx, offerCoinListY, orderPriceListY, orderAddrListY, poolId, false)

	// Set the execution status flag of messages to true.
	msgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, batch)
	for _, msg := range msgs {
		msg.Executed = true
	}
	simapp.LiquidityKeeper.SetLiquidityPoolBatchSwapMsgPointers(ctx, poolId, msgs)

	// Generate an orderbook by arranging swap messages in order price
	orderMap, XtoY, YtoX := types.GetOrderMap(msgs, denomX, denomY, false)
	orderBook := orderMap.SortOrderBook()
	currentPrice := X.Quo(Y).ToDec()
	require.Equal(t, orderMap[orderPriceList[0].String()].BuyOfferAmt, offerCoinList[0].Amount.MulRaw(3))
	require.Equal(t, orderMap[orderPriceList[0].String()].OrderPrice, orderPriceList[0])

	require.Equal(t, 3, len(XtoY))
	require.Equal(t, 1, len(YtoX))
	require.Equal(t, 3, len(orderMap[orderPriceList[0].String()].MsgList))
	require.Equal(t, 1, len(orderMap[orderPriceListY[0].String()].MsgList))
	require.Equal(t, 3, len(orderBook[0].MsgList))
	require.Equal(t, 1, len(orderBook[1].MsgList))

	require.Equal(t, len(orderBook), orderBook.Len())

	fmt.Println(orderBook, currentPrice)
	fmt.Println(XtoY, YtoX)

	types.ValidateStateAndExpireOrders(XtoY, ctx.BlockHeight(), false)
	types.ValidateStateAndExpireOrders(YtoX, ctx.BlockHeight(), false)

	currentYPriceOverX := X.Quo(Y).ToDec()

	// The price and coins of swap messages in orderbook are calculated
	// to derive match result with the price direction.
	result := types.MatchOrderbook(X.ToDec(), Y.ToDec(), currentYPriceOverX, orderBook)
	require.NotEqual(t, types.NoMatch, result.MatchType)

	matchResultXtoY, _, poolXDeltaXtoY, poolYDeltaXtoY := types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX,
		result.SwapPrice, ctx.BlockHeight())
	matchResultYtoX, _, poolXDeltaYtoX, poolYDeltaYtoX := types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY,
		result.SwapPrice, ctx.BlockHeight())

	XtoY, YtoX, XDec, YDec, poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY :=
		simapp.LiquidityKeeper.UpdateState(X.ToDec(), Y.ToDec(), XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	require.Equal(t, 0, (types.MsgList)(XtoY).CountNotMatchedMsgs())
	require.Equal(t, 0, (types.MsgList)(XtoY).CountFractionalMatchedMsgs())
	require.Equal(t, 1, (types.MsgList)(YtoX).CountNotMatchedMsgs())
	require.Equal(t, 0, (types.MsgList)(YtoX).CountFractionalMatchedMsgs())
	require.Equal(t, 3, len(XtoY))
	require.Equal(t, 1, len(YtoX))

	fmt.Println(matchResultXtoY)
	fmt.Println(poolXDeltaXtoY)
	fmt.Println(poolYDeltaXtoY)

	fmt.Println(poolXDeltaYtoX, poolYDeltaYtoX)
	fmt.Println(poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY)
	fmt.Println(decimalErrorX, decimalErrorY)
	fmt.Println(XDec, YDec)

	// Verify swap result by creating an orderbook with remaining messages that have been matched and not transacted.
	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY, true)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	lastPrice := XDec.Quo(YDec)
	fmt.Println("lastPrice", lastPrice)
	fmt.Println("X", XDec)
	fmt.Println("Y", YDec)
	require.True(t, types.CheckValidityOrderBook(orderBookExecuted, lastPrice))

	require.Equal(t, 0, (types.MsgList)(orderMapExecuted[orderPriceList[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 1, (types.MsgList)(orderMapExecuted[orderPriceListY[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 1, (types.MsgList)(orderBookExecuted[0].MsgList).CountNotMatchedMsgs())

	types.ValidateStateAndExpireOrders(XtoY, ctx.BlockHeight(), true)
	types.ValidateStateAndExpireOrders(YtoX, ctx.BlockHeight(), true)

	orderMapCleared, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY, true)
	orderBookCleared := orderMapCleared.SortOrderBook()
	require.True(t, types.CheckValidityOrderBook(orderBookCleared, lastPrice))

	require.Equal(t, 0, (types.MsgList)(orderMapCleared[orderPriceList[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 0, (types.MsgList)(orderMapCleared[orderPriceListY[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 0, len(orderBookCleared))

	// next block
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// test genesisState with export, init
	genesis := simapp.LiquidityKeeper.ExportGenesis(ctx)
	simapp.LiquidityKeeper.InitGenesis(ctx, *genesis)
	err := types.ValidateGenesis(*genesis)
	require.NoError(t, err)
	genesisNew := simapp.LiquidityKeeper.ExportGenesis(ctx)
	err = types.ValidateGenesis(*genesisNew)
	require.NoError(t, err)
	require.Equal(t, genesis, genesisNew)
	for _, record := range genesisNew.LiquidityPoolRecords {
		err = record.Validate()
		require.NoError(t, err)
	}

	// validate genesis fail case
	batch.DepositMsgIndex = 0
	simapp.LiquidityKeeper.SetLiquidityPoolBatch(ctx, batch)
	genesisNew = simapp.LiquidityKeeper.ExportGenesis(ctx)
	err = types.ValidateGenesis(*genesisNew)
	require.Error(t, err, types.ErrBadBatchMsgIndex)
	batch.WithdrawMsgIndex = 0
	simapp.LiquidityKeeper.SetLiquidityPoolBatch(ctx, batch)
	genesisNew = simapp.LiquidityKeeper.ExportGenesis(ctx)
	err = types.ValidateGenesis(*genesisNew)
	require.Error(t, err, types.ErrBadBatchMsgIndex)
	batch.SwapMsgIndex = 20
	simapp.LiquidityKeeper.SetLiquidityPoolBatch(ctx, batch)
	genesisNew = simapp.LiquidityKeeper.ExportGenesis(ctx)
	err = types.ValidateGenesis(*genesisNew)
	require.Error(t, err, types.ErrBadBatchMsgIndex)
}

func TestMaxOrderRatio(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := params.MinInitDepositToPool
	Y := params.MinInitDepositToPool

	addrs := app.AddTestAddrs(simapp, ctx, 20, params.LiquidityPoolCreationFee)
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:2], poolId, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	maxOrderRatio := params.MaxOrderAmountRatio

	// Success case, not exceed GetMaxOrderRatio orders
	priceBuy, _ := sdk.NewDecFromStr("1.1")
	priceSell, _ := sdk.NewDecFromStr("1.2")

	offerCoin := sdk.NewCoin(denomX, sdk.NewInt(1000))
	offerCoinY := sdk.NewCoin(denomY, sdk.NewInt(1000))

	app.SaveAccountWithFee(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin), offerCoin)
	app.SaveAccountWithFee(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY), offerCoinY)

	msgBuy := types.NewMsgSwap(addrs[1], poolId, DefaultSwapType, offerCoin, DenomY, priceBuy, params.SwapFeeRate)
	msgSell := types.NewMsgSwap(addrs[2], poolId, DefaultSwapType, offerCoinY, DenomX, priceSell, params.SwapFeeRate)

	_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy, 0)
	require.NoError(t, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell, 0)
	require.NoError(t, err)

	// Fail case, exceed GetMaxOrderRatio orders
	offerCoin = sdk.NewCoin(denomX, X)
	offerCoinY = sdk.NewCoin(denomY, Y)

	app.SaveAccountWithFee(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin), offerCoin)
	app.SaveAccountWithFee(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY), offerCoinY)

	msgBuy = types.NewMsgSwap(addrs[1], poolId, DefaultSwapType, offerCoin, DenomY, priceBuy, params.SwapFeeRate)
	msgSell = types.NewMsgSwap(addrs[2], poolId, DefaultSwapType, offerCoinY, DenomX, priceSell, params.SwapFeeRate)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy, 0)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell, 0)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

	// Success case, same GetMaxOrderRatio orders
	offerCoin = sdk.NewCoin(denomX, X.ToDec().Mul(maxOrderRatio).TruncateInt())
	offerCoinY = sdk.NewCoin(denomY, Y.ToDec().Mul(maxOrderRatio).TruncateInt())

	app.SaveAccountWithFee(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin), offerCoin)
	app.SaveAccountWithFee(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY), offerCoinY)

	msgBuy = types.NewMsgSwap(addrs[1], poolId, DefaultSwapType, offerCoin, DenomY, priceBuy, params.SwapFeeRate)
	msgSell = types.NewMsgSwap(addrs[2], poolId, DefaultSwapType, offerCoinY, DenomX, priceSell, params.SwapFeeRate)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy, 0)
	require.NoError(t, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell, 0)
	require.NoError(t, err)

	// Success case, same GetMaxOrderRatio orders
	offerCoin = sdk.NewCoin(denomX, X.ToDec().Mul(maxOrderRatio).TruncateInt().AddRaw(1))
	offerCoinY = sdk.NewCoin(denomY, Y.ToDec().Mul(maxOrderRatio).TruncateInt().AddRaw(1))

	offerCoin = sdk.NewCoin(denomX, params.MinInitDepositToPool.Quo(sdk.NewInt(2)))
	offerCoinY = sdk.NewCoin(denomY, params.MinInitDepositToPool.Quo(sdk.NewInt(10)))
	app.SaveAccountWithFee(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin), offerCoin)
	app.SaveAccountWithFee(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY), offerCoinY)

	msgBuy = types.NewMsgSwap(addrs[1], poolId, DefaultSwapType, offerCoin, DenomY, priceBuy, params.SwapFeeRate)
	msgSell = types.NewMsgSwap(addrs[2], poolId, DefaultSwapType, offerCoinY, DenomX, priceSell, params.SwapFeeRate)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy, 0)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell, 0)
	require.NoError(t, err)
}

func TestOrderBookSort(t *testing.T) {
	orderMap := make(types.OrderMap)
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")
	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[b.String()] = types.OrderByPrice{
		OrderPrice:   b,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[c.String()] = types.OrderByPrice{
		OrderPrice:   c,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()
	fmt.Println(orderBook)

	res := orderBook.Less(0, 1)
	require.True(t, res)
	res = orderBook.Less(1, 2)
	require.True(t, res)
	res = orderBook.Less(2, 1)
	require.False(t, res)

	orderBook.Swap(1, 2)
	fmt.Println(orderBook)
	require.Equal(t, c, orderBook[1].OrderPrice)
	require.Equal(t, b, orderBook[2].OrderPrice)

	orderBook.Sort()
	fmt.Println(orderBook)
	require.Equal(t, a, orderBook[0].OrderPrice)
	require.Equal(t, b, orderBook[1].OrderPrice)
	require.Equal(t, c, orderBook[2].OrderPrice)

	orderBook.Reverse()
	fmt.Println(orderBook)
	require.Equal(t, a, orderBook[2].OrderPrice)
	require.Equal(t, b, orderBook[1].OrderPrice)
	require.Equal(t, c, orderBook[0].OrderPrice)
}

func TestMinMaxDec(t *testing.T) {
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")

	require.Equal(t, a, types.MinDec(a, b))
	require.Equal(t, a, types.MinDec(a, c))
	require.Equal(t, b, types.MaxDec(a, b))
	require.Equal(t, c, types.MaxDec(a, c))
	require.Equal(t, a, types.MaxDec(a, a))
	require.Equal(t, a, types.MinDec(a, a))
}

func TestMaxInt(t *testing.T) {
	a := sdk.NewInt(1)
	b := sdk.NewInt(2)
	c := sdk.NewInt(3)

	require.Equal(t, a, types.MinInt(a, b))
	require.Equal(t, a, types.MinInt(a, c))
	require.Equal(t, b, types.MaxInt(a, b))
	require.Equal(t, c, types.MaxInt(a, c))
	require.Equal(t, a, types.MaxInt(a, a))
	require.Equal(t, a, types.MinInt(a, a))
}

func TestGetExecutableAmt(t *testing.T) {
	orderMap := make(types.OrderMap)
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")
	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.NewInt(30000000),
	}
	orderMap[b.String()] = types.OrderByPrice{
		OrderPrice:   b,
		BuyOfferAmt:  sdk.NewInt(90000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[c.String()] = types.OrderByPrice{
		OrderPrice:   c,
		BuyOfferAmt:  sdk.NewInt(50000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	executableBuyAmtX, executableSellAmtY := types.GetExecutableAmt(b, orderBook)
	require.Equal(t, sdk.NewInt(140000000), executableBuyAmtX)
	require.Equal(t, sdk.NewInt(30000000), executableSellAmtY)
}

func TestGetPriceDirection(t *testing.T) {
	// increase case
	orderMap := make(types.OrderMap)
	a, _ := sdk.NewDecFromStr("1")
	b, _ := sdk.NewDecFromStr("1.1")
	c, _ := sdk.NewDecFromStr("1.2")
	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.NewInt(40000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[b.String()] = types.OrderByPrice{
		OrderPrice:   b,
		BuyOfferAmt:  sdk.NewInt(40000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[c.String()] = types.OrderByPrice{
		OrderPrice:   c,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.NewInt(20000000),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()
	currentYPriceOverX, _ := sdk.NewDecFromStr("1.0")
	result := types.GetPriceDirection(currentYPriceOverX, orderBook)
	require.Equal(t, types.Increase, result)

	// decrease case
	orderMap = make(types.OrderMap)
	a, _ = sdk.NewDecFromStr("0.7")
	b, _ = sdk.NewDecFromStr("0.9")
	c, _ = sdk.NewDecFromStr("0.8")
	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.NewInt(20000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[b.String()] = types.OrderByPrice{
		OrderPrice:   b,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.NewInt(40000000),
	}
	orderMap[c.String()] = types.OrderByPrice{
		OrderPrice:   c,
		BuyOfferAmt:  sdk.NewInt(10000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook = orderMap.SortOrderBook()
	currentYPriceOverX, _ = sdk.NewDecFromStr("1.0")
	result = types.GetPriceDirection(currentYPriceOverX, orderBook)
	require.Equal(t, types.Decrease, result)

	// stay case
	orderMap = make(types.OrderMap)
	a, _ = sdk.NewDecFromStr("1.0")

	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.NewInt(50000000),
		SellOfferAmt: sdk.NewInt(50000000),
	}
	orderBook = orderMap.SortOrderBook()
	currentYPriceOverX, _ = sdk.NewDecFromStr("1.0")
	result = types.GetPriceDirection(currentYPriceOverX, orderBook)
	require.Equal(t, types.Stay, result)
}

func TestComputePriceDirection(t *testing.T) {
	// increase case
	orderMap := make(types.OrderMap)
	a, _ := sdk.NewDecFromStr("1")
	b, _ := sdk.NewDecFromStr("1.1")
	c, _ := sdk.NewDecFromStr("1.2")
	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.NewInt(40000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[b.String()] = types.OrderByPrice{
		OrderPrice:   b,
		BuyOfferAmt:  sdk.NewInt(40000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[c.String()] = types.OrderByPrice{
		OrderPrice:   c,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.NewInt(20000000),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	X := orderMap[a.String()].BuyOfferAmt.ToDec().Add(orderMap[b.String()].BuyOfferAmt.ToDec())
	Y := orderMap[c.String()].SellOfferAmt.ToDec()

	currentYPriceOverX := X.Quo(Y)
	direction := types.GetPriceDirection(currentYPriceOverX, orderBook)
	result := types.MatchOrderbook(X, Y, currentYPriceOverX, orderBook)
	require.Equal(t, types.CalculateMatch(direction, X, Y, currentYPriceOverX, orderBook), result)

	// decrease case
	orderMap = make(types.OrderMap)
	a, _ = sdk.NewDecFromStr("0.7")
	b, _ = sdk.NewDecFromStr("0.9")
	c, _ = sdk.NewDecFromStr("0.8")
	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.NewInt(20000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	orderMap[b.String()] = types.OrderByPrice{
		OrderPrice:   b,
		BuyOfferAmt:  sdk.ZeroInt(),
		SellOfferAmt: sdk.NewInt(40000000),
	}
	orderMap[c.String()] = types.OrderByPrice{
		OrderPrice:   c,
		BuyOfferAmt:  sdk.NewInt(10000000),
		SellOfferAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook = orderMap.SortOrderBook()

	X = orderMap[a.String()].BuyOfferAmt.ToDec().Add(orderMap[c.String()].BuyOfferAmt.ToDec())
	Y = orderMap[b.String()].SellOfferAmt.ToDec()

	currentYPriceOverX = X.Quo(Y)
	direction = types.GetPriceDirection(currentYPriceOverX, orderBook)
	result = types.MatchOrderbook(X, Y, currentYPriceOverX, orderBook)
	require.Equal(t, types.CalculateMatch(direction, X, Y, currentYPriceOverX, orderBook), result)

	// stay case
	orderMap = make(types.OrderMap)
	a, _ = sdk.NewDecFromStr("1.0")

	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice:   a,
		BuyOfferAmt:  sdk.NewInt(50000000),
		SellOfferAmt: sdk.NewInt(50000000),
	}
	orderBook = orderMap.SortOrderBook()

	X = orderMap[a.String()].BuyOfferAmt.ToDec()
	Y = orderMap[a.String()].SellOfferAmt.ToDec()
	currentYPriceOverX = X.Quo(Y)

	result = types.MatchOrderbook(X, Y, currentYPriceOverX, orderBook)
	require.Equal(t, types.CalculateMatchStay(currentYPriceOverX, orderBook), result)
}

// Match Stay case with fractional match type
func TestCalculateMatchStayEdgeCase(t *testing.T) {
	currentPrice, err := sdk.NewDecFromStr("1.844380246375231658")
	require.NoError(t, err)
	var orderbook types.OrderBook
	orderbookEdgeCase := `[{"OrderPrice":"1.827780824157854573","BuyOfferAmt":"12587364000","SellOfferAmt":"6200948000","MsgList":[{"msg_index":12,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2097894000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2097894000"},"demand_coin_denom":"denomY","order_price":"1.827780824157854573"}},{"msg_index":16,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"4669506000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg44npvhm","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"4669506000"},"demand_coin_denom":"denomY","order_price":"1.827780824157854573"}},{"msg_index":23,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"609066000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfzwk37gt","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"609066000"},"demand_coin_denom":"denomY","order_price":"1.827780824157854573"}},{"msg_index":39,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5210898000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfckxufsg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5210898000"},"demand_coin_denom":"denomY","order_price":"1.827780824157854573"}},{"msg_index":56,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1284220000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg44npvhm","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1284220000"},"demand_coin_denom":"denomX","order_price":"1.827780824157854573"}},{"msg_index":78,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1981368000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfhft040s","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1981368000"},"demand_coin_denom":"denomX","order_price":"1.827780824157854573"}},{"msg_index":85,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2935360000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c2yrhrufk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2935360000"},"demand_coin_denom":"denomX","order_price":"1.827780824157854573"}}]},{"OrderPrice":"1.829625204404229805","BuyOfferAmt":"9203664000","SellOfferAmt":"6971480000","MsgList":[{"msg_index":18,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5210898000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cghxkq0yk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5210898000"},"demand_coin_denom":"denomY","order_price":"1.829625204404229805"}},{"msg_index":36,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3992766000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf46wwkua","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3992766000"},"demand_coin_denom":"denomY","order_price":"1.829625204404229805"}},{"msg_index":44,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3155512000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3155512000"},"demand_coin_denom":"denomX","order_price":"1.829625204404229805"}},{"msg_index":55,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"513688000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg5g94e2f","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"513688000"},"demand_coin_denom":"denomX","order_price":"1.829625204404229805"}},{"msg_index":61,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3302280000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfqansamx","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3302280000"},"demand_coin_denom":"denomX","order_price":"1.829625204404229805"}}]},{"OrderPrice":"1.831469584650605036","BuyOfferAmt":"18001284000","SellOfferAmt":"2311596000","MsgList":[{"msg_index":21,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3248352000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfqansamx","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3248352000"},"demand_coin_denom":"denomY","order_price":"1.831469584650605036"}},{"msg_index":32,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5007876000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf34yvsn8","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5007876000"},"demand_coin_denom":"denomY","order_price":"1.831469584650605036"}},{"msg_index":33,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5955312000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfjmhexac","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5955312000"},"demand_coin_denom":"denomY","order_price":"1.831469584650605036"}},{"msg_index":34,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3789744000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfnxpdnq2","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3789744000"},"demand_coin_denom":"denomY","order_price":"1.831469584650605036"}},{"msg_index":65,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2311596000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfyjejm5u","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2311596000"},"demand_coin_denom":"denomX","order_price":"1.831469584650605036"}}]},{"OrderPrice":"1.833313964896980268","BuyOfferAmt":"12113646000","SellOfferAmt":"4806652000","MsgList":[{"msg_index":6,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"6632052000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg9qjf5zg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"6632052000"},"demand_coin_denom":"denomY","order_price":"1.833313964896980268"}},{"msg_index":28,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5481594000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf8u28d6r","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5481594000"},"demand_coin_denom":"denomY","order_price":"1.833313964896980268"}},{"msg_index":41,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"660456000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"660456000"},"demand_coin_denom":"denomX","order_price":"1.833313964896980268"}},{"msg_index":64,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2421672000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfrnq9t4e","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2421672000"},"demand_coin_denom":"denomX","order_price":"1.833313964896980268"}},{"msg_index":73,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1724524000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfjmhexac","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1724524000"},"demand_coin_denom":"denomX","order_price":"1.833313964896980268"}}]},{"OrderPrice":"1.835158345143355500","BuyOfferAmt":"0","SellOfferAmt":"6421100000","MsgList":[{"msg_index":47,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2715208000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxwpuzvh","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2715208000"},"demand_coin_denom":"denomX","order_price":"1.835158345143355500"}},{"msg_index":58,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2678516000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cghxkq0yk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2678516000"},"demand_coin_denom":"denomX","order_price":"1.835158345143355500"}},{"msg_index":82,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1027376000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c2p3t40m7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1027376000"},"demand_coin_denom":"denomX","order_price":"1.835158345143355500"}}]},{"OrderPrice":"1.837002725389730731","BuyOfferAmt":"9135990000","SellOfferAmt":"3852660000","MsgList":[{"msg_index":13,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"744414000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgj52kuk7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"744414000"},"demand_coin_denom":"denomY","order_price":"1.837002725389730731"}},{"msg_index":19,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5143224000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgcemnnmw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5143224000"},"demand_coin_denom":"denomY","order_price":"1.837002725389730731"}},{"msg_index":22,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"541392000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfpq9ygx5","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"541392000"},"demand_coin_denom":"denomY","order_price":"1.837002725389730731"}},{"msg_index":35,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2706960000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf58c6rp0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2706960000"},"demand_coin_denom":"denomY","order_price":"1.837002725389730731"}},{"msg_index":48,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2274904000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8nhgh39","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2274904000"},"demand_coin_denom":"denomX","order_price":"1.837002725389730731"}},{"msg_index":51,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1394296000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgs80hl9n","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1394296000"},"demand_coin_denom":"denomX","order_price":"1.837002725389730731"}},{"msg_index":80,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"183460000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfetsgud6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"183460000"},"demand_coin_denom":"denomX","order_price":"1.837002725389730731"}}]},{"OrderPrice":"1.838847105636105963","BuyOfferAmt":"6226008000","SellOfferAmt":"2715208000","MsgList":[{"msg_index":5,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"6226008000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyayapl6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"6226008000"},"demand_coin_denom":"denomY","order_price":"1.838847105636105963"}},{"msg_index":43,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2715208000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2715208000"},"demand_coin_denom":"denomX","order_price":"1.838847105636105963"}}]},{"OrderPrice":"1.840691485882481195","BuyOfferAmt":"6496704000","SellOfferAmt":"3155512000","MsgList":[{"msg_index":8,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"6496704000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8nhgh39","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"6496704000"},"demand_coin_denom":"denomY","order_price":"1.840691485882481195"}},{"msg_index":81,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3155512000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c2qvap6xv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3155512000"},"demand_coin_denom":"denomX","order_price":"1.840691485882481195"}}]},{"OrderPrice":"1.842535866128856426","BuyOfferAmt":"0","SellOfferAmt":"1137452000","MsgList":[{"msg_index":45,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1137452000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyayapl6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1137452000"},"demand_coin_denom":"denomX","order_price":"1.842535866128856426"}}]},{"OrderPrice":"1.844380246375231658","BuyOfferAmt":"15700368000","SellOfferAmt":"2274904000","MsgList":[{"msg_index":14,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1759524000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgnfuzftv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1759524000"},"demand_coin_denom":"denomY","order_price":"1.844380246375231658"}},{"msg_index":24,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1624176000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfrnq9t4e","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1624176000"},"demand_coin_denom":"denomY","order_price":"1.844380246375231658"}},{"msg_index":25,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3248352000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfyjejm5u","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3248352000"},"demand_coin_denom":"denomY","order_price":"1.844380246375231658"}},{"msg_index":29,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"4263462000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfgr8539m","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"4263462000"},"demand_coin_denom":"denomY","order_price":"1.844380246375231658"}},{"msg_index":31,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"4804854000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfsgjc9w4","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"4804854000"},"demand_coin_denom":"denomY","order_price":"1.844380246375231658"}},{"msg_index":59,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1651140000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgcemnnmw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1651140000"},"demand_coin_denom":"denomX","order_price":"1.844380246375231658"}},{"msg_index":62,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"623764000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfpq9ygx5","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"623764000"},"demand_coin_denom":"denomX","order_price":"1.844380246375231658"}}]},{"OrderPrice":"1.846224626621606890","BuyOfferAmt":"19963830000","SellOfferAmt":"3338972000","MsgList":[{"msg_index":11,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"6429030000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgs80hl9n","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"6429030000"},"demand_coin_denom":"denomY","order_price":"1.846224626621606890"}},{"msg_index":20,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5143224000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgeyd8xxu","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5143224000"},"demand_coin_denom":"denomY","order_price":"1.846224626621606890"}},{"msg_index":27,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2300916000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfxpunc83","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2300916000"},"demand_coin_denom":"denomY","order_price":"1.846224626621606890"}},{"msg_index":38,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"6090660000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfhft040s","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"6090660000"},"demand_coin_denom":"denomY","order_price":"1.846224626621606890"}},{"msg_index":42,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"660456000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"660456000"},"demand_coin_denom":"denomX","order_price":"1.846224626621606890"}},{"msg_index":68,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2678516000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf8u28d6r","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2678516000"},"demand_coin_denom":"denomX","order_price":"1.846224626621606890"}}]},{"OrderPrice":"1.848069006867982121","BuyOfferAmt":"0","SellOfferAmt":"3302280000","MsgList":[{"msg_index":46,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2201520000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg9qjf5zg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2201520000"},"demand_coin_denom":"denomX","order_price":"1.848069006867982121"}},{"msg_index":70,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1100760000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cff73qycf","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1100760000"},"demand_coin_denom":"denomX","order_price":"1.848069006867982121"}}]},{"OrderPrice":"1.849913387114357353","BuyOfferAmt":"2233242000","SellOfferAmt":"10420528000","MsgList":[{"msg_index":4,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2233242000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2233242000"},"demand_coin_denom":"denomY","order_price":"1.849913387114357353"}},{"msg_index":54,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"917300000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgnfuzftv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"917300000"},"demand_coin_denom":"denomX","order_price":"1.849913387114357353"}},{"msg_index":60,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3485740000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgeyd8xxu","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3485740000"},"demand_coin_denom":"denomX","order_price":"1.849913387114357353"}},{"msg_index":63,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"697148000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfzwk37gt","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"697148000"},"demand_coin_denom":"denomX","order_price":"1.849913387114357353"}},{"msg_index":66,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2421672000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf900xwfw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2421672000"},"demand_coin_denom":"denomX","order_price":"1.849913387114357353"}},{"msg_index":84,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1357604000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c2rzw5vgn","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1357604000"},"demand_coin_denom":"denomX","order_price":"1.849913387114357353"}},{"msg_index":87,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1541064000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c2xsjzl6m","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1541064000"},"demand_coin_denom":"denomX","order_price":"1.849913387114357353"}}]},{"OrderPrice":"1.851757767360732585","BuyOfferAmt":"23550552000","SellOfferAmt":"1577756000","MsgList":[{"msg_index":1,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5075550000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5075550000"},"demand_coin_denom":"denomY","order_price":"1.851757767360732585"}},{"msg_index":7,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"4128114000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxwpuzvh","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"4128114000"},"demand_coin_denom":"denomY","order_price":"1.851757767360732585"}},{"msg_index":9,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"4940202000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cggv6mtwa","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"4940202000"},"demand_coin_denom":"denomY","order_price":"1.851757767360732585"}},{"msg_index":15,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3113004000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg5g94e2f","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3113004000"},"demand_coin_denom":"denomY","order_price":"1.851757767360732585"}},{"msg_index":26,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"6293682000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf900xwfw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"6293682000"},"demand_coin_denom":"denomY","order_price":"1.851757767360732585"}},{"msg_index":67,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"146768000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfxpunc83","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"146768000"},"demand_coin_denom":"denomX","order_price":"1.851757767360732585"}},{"msg_index":71,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1430988000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfsgjc9w4","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1430988000"},"demand_coin_denom":"denomX","order_price":"1.851757767360732585"}}]},{"OrderPrice":"1.853602147607107816","BuyOfferAmt":"3519048000","SellOfferAmt":"5577184000","MsgList":[{"msg_index":10,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3519048000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgf3v07n0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3519048000"},"demand_coin_denom":"denomY","order_price":"1.853602147607107816"}},{"msg_index":52,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"403612000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"403612000"},"demand_coin_denom":"denomX","order_price":"1.853602147607107816"}},{"msg_index":53,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"770532000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgj52kuk7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"770532000"},"demand_coin_denom":"denomX","order_price":"1.853602147607107816"}},{"msg_index":72,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"146768000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf34yvsn8","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"146768000"},"demand_coin_denom":"denomX","order_price":"1.853602147607107816"}},{"msg_index":74,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3155512000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfnxpdnq2","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3155512000"},"demand_coin_denom":"denomX","order_price":"1.853602147607107816"}},{"msg_index":75,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"183460000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf58c6rp0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"183460000"},"demand_coin_denom":"denomX","order_price":"1.853602147607107816"}},{"msg_index":76,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"917300000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf46wwkua","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"917300000"},"demand_coin_denom":"denomX","order_price":"1.853602147607107816"}}]},{"OrderPrice":"1.855446527853483048","BuyOfferAmt":"5752290000","SellOfferAmt":"1357604000","MsgList":[{"msg_index":3,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"3654396000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"3654396000"},"demand_coin_denom":"denomY","order_price":"1.855446527853483048"}},{"msg_index":17,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2097894000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgkmq56ey","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2097894000"},"demand_coin_denom":"denomY","order_price":"1.855446527853483048"}},{"msg_index":49,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1357604000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cggv6mtwa","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1357604000"},"demand_coin_denom":"denomX","order_price":"1.855446527853483048"}}]},{"OrderPrice":"1.857290908099858280","BuyOfferAmt":"2774634000","SellOfferAmt":"4256272000","MsgList":[{"msg_index":37,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2774634000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfk5amqjz","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2774634000"},"demand_coin_denom":"denomY","order_price":"1.857290908099858280"}},{"msg_index":50,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2128136000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgf3v07n0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2128136000"},"demand_coin_denom":"denomX","order_price":"1.857290908099858280"}},{"msg_index":77,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"256844000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfk5amqjz","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"256844000"},"demand_coin_denom":"denomX","order_price":"1.857290908099858280"}},{"msg_index":83,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1871292000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c2zlcqe4p","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1871292000"},"demand_coin_denom":"denomX","order_price":"1.857290908099858280"}}]},{"OrderPrice":"1.859135288346233511","BuyOfferAmt":"10760166000","SellOfferAmt":"5283648000","MsgList":[{"msg_index":2,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1421154000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1421154000"},"demand_coin_denom":"denomY","order_price":"1.859135288346233511"}},{"msg_index":30,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"4331136000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cff73qycf","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"4331136000"},"demand_coin_denom":"denomY","order_price":"1.859135288346233511"}},{"msg_index":40,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"5007876000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfetsgud6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"5007876000"},"demand_coin_denom":"denomY","order_price":"1.859135288346233511"}},{"msg_index":57,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1137452000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgkmq56ey","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1137452000"},"demand_coin_denom":"denomX","order_price":"1.859135288346233511"}},{"msg_index":69,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"293536000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfgr8539m","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"293536000"},"demand_coin_denom":"denomX","order_price":"1.859135288346233511"}},{"msg_index":79,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3302280000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfckxufsg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3302280000"},"demand_coin_denom":"denomX","order_price":"1.859135288346233511"}},{"msg_index":86,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"550380000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5c297phf5y","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"550380000"},"demand_coin_denom":"denomX","order_price":"1.859135288346233511"}}]}]`
	json.Unmarshal([]byte(orderbookEdgeCase), &orderbook)
	r := types.CalculateMatchStay(currentPrice, orderbook)
	require.Equal(t, types.FractionalMatch, r.MatchType)
	// stay case with fractional
}

// Match Stay case with no match type
func TestCalculateNoMatchEdgeCase(t *testing.T) {
	currentPrice, err := sdk.NewDecFromStr("1.007768598527187219")
	require.NoError(t, err)
	var orderbook types.OrderBook
	orderbookEdgeCase := `[{"OrderPrice":"1.007768598527187219","BuyOfferAmt":"0","SellOfferAmt":"417269600","MsgList":[{"msg_index":1,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"417269600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"417269600"},"demand_coin_denom":"denomX","order_price":"1.007768598527187219"}}]},{"OrderPrice":"1.011799672921295968","BuyOfferAmt":"0","SellOfferAmt":"2190665400","MsgList":[{"msg_index":2,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2190665400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2190665400"},"demand_coin_denom":"denomX","order_price":"1.011799672921295968"}}]}]`
	json.Unmarshal([]byte(orderbookEdgeCase), &orderbook)
	r := types.CalculateMatchStay(currentPrice, orderbook)
	require.Equal(t, types.NoMatch, r.MatchType)
	// stay case with fractional
}

// Reproduce GetOrderMapEdgeCase, selling Y for X case, ErrInvalidDenom case
func TestGetOrderMapEdgeCase(t *testing.T) {
	onlyNotMatched := false
	var swapMsgs []*types.BatchPoolSwapMsg
	swapMsgsJson := `[{"msg_index":1,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"19228500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"19228500"},"demand_coin_denom":"denomY","order_price":"0.027506527499265415"}},{"msg_index":2,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"141009000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"141009000"},"demand_coin_denom":"denomY","order_price":"0.027341323129900457"}},{"msg_index":3,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"23501500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"23501500"},"demand_coin_denom":"denomY","order_price":"0.027616663745508720"}},{"msg_index":4,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"200831000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"200831000"},"demand_coin_denom":"denomY","order_price":"0.027589129683947893"}},{"msg_index":5,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"160237500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyayapl6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"160237500"},"demand_coin_denom":"denomY","order_price":"0.027313789068339631"}},{"msg_index":6,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"175193000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg9qjf5zg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"175193000"},"demand_coin_denom":"denomY","order_price":"0.027478993437704589"}},{"msg_index":7,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"183739000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxwpuzvh","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"183739000"},"demand_coin_denom":"denomY","order_price":"0.027699265930191198"}},{"msg_index":8,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"32047500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8nhgh39","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"32047500"},"demand_coin_denom":"denomY","order_price":"0.027451459376143762"}},{"msg_index":9,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"111098000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cggv6mtwa","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"111098000"},"demand_coin_denom":"denomY","order_price":"0.027286255006778805"}},{"msg_index":10,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"166647000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgf3v07n0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"166647000"},"demand_coin_denom":"denomY","order_price":"0.027341323129900457"}},{"msg_index":11,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"98279000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgs80hl9n","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"98279000"},"demand_coin_denom":"denomY","order_price":"0.027368857191461284"}},{"msg_index":12,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"8546000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"8546000"},"demand_coin_denom":"denomY","order_price":"0.027396391253022110"}},{"msg_index":13,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"87596500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgj52kuk7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"87596500"},"demand_coin_denom":"denomY","order_price":"0.027451459376143762"}},{"msg_index":14,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"111098000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgnfuzftv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"111098000"},"demand_coin_denom":"denomY","order_price":"0.027478993437704589"}},{"msg_index":15,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"38457000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg5g94e2f","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"38457000"},"demand_coin_denom":"denomY","order_price":"0.027451459376143762"}},{"msg_index":16,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"153828000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg44npvhm","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"153828000"},"demand_coin_denom":"denomY","order_price":"0.027616663745508720"}},{"msg_index":17,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"70504500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgkmq56ey","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"70504500"},"demand_coin_denom":"denomY","order_price":"0.027451459376143762"}},{"msg_index":18,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"47003000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cghxkq0yk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"47003000"},"demand_coin_denom":"denomY","order_price":"0.027396391253022110"}},{"msg_index":19,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"132463000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgcemnnmw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"132463000"},"demand_coin_denom":"denomY","order_price":"0.027726799991752025"}},{"msg_index":20,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"66231500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgeyd8xxu","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"66231500"},"demand_coin_denom":"denomY","order_price":"0.027561595622387067"}},{"msg_index":21,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"119644000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfqansamx","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"119644000"},"demand_coin_denom":"denomY","order_price":"0.027506527499265415"}},{"msg_index":22,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"17092000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfpq9ygx5","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"17092000"},"demand_coin_denom":"denomY","order_price":"0.027341323129900457"}},{"msg_index":23,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"209377000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfzwk37gt","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"209377000"},"demand_coin_denom":"denomY","order_price":"0.027478993437704589"}},{"msg_index":24,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"207240500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfrnq9t4e","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"207240500"},"demand_coin_denom":"denomY","order_price":"0.027396391253022110"}},{"msg_index":25,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"155964500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfyjejm5u","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"155964500"},"demand_coin_denom":"denomY","order_price":"0.027423925314582936"}},{"msg_index":26,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"194421500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf900xwfw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"194421500"},"demand_coin_denom":"denomY","order_price":"0.027286255006778805"}},{"msg_index":27,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"102552000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfxpunc83","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"102552000"},"demand_coin_denom":"denomY","order_price":"0.027368857191461284"}},{"msg_index":28,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"151691500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf8u28d6r","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"151691500"},"demand_coin_denom":"denomY","order_price":"0.027478993437704589"}},{"msg_index":29,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"113234500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfgr8539m","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"113234500"},"demand_coin_denom":"denomY","order_price":"0.027368857191461284"}},{"msg_index":30,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"117507500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cff73qycf","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"117507500"},"demand_coin_denom":"denomY","order_price":"0.027423925314582936"}},{"msg_index":31,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"141009000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfsgjc9w4","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"141009000"},"demand_coin_denom":"denomY","order_price":"0.027423925314582936"}},{"msg_index":32,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"200831000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf34yvsn8","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"200831000"},"demand_coin_denom":"denomY","order_price":"0.027534061560826241"}},{"msg_index":33,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"141009000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfjmhexac","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"141009000"},"demand_coin_denom":"denomY","order_price":"0.027726799991752025"}},{"msg_index":34,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"98279000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfnxpdnq2","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"98279000"},"demand_coin_denom":"denomY","order_price":"0.027478993437704589"}},{"msg_index":35,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"76914000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf58c6rp0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"76914000"},"demand_coin_denom":"denomY","order_price":"0.027423925314582936"}},{"msg_index":36,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"23501500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf46wwkua","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"23501500"},"demand_coin_denom":"denomY","order_price":"0.027754334053312851"}},{"msg_index":37,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4733282800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4733282800"},"demand_coin_denom":"denomX","order_price":"0.027699265930191198"}},{"msg_index":38,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3957334800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3957334800"},"demand_coin_denom":"denomX","order_price":"0.027478993437704589"}},{"msg_index":39,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2483033600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2483033600"},"demand_coin_denom":"denomX","order_price":"0.027589129683947893"}},{"msg_index":40,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"5509230800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"5509230800"},"demand_coin_denom":"denomX","order_price":"0.027561595622387067"}},{"msg_index":41,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2327844000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyayapl6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2327844000"},"demand_coin_denom":"denomX","order_price":"0.027423925314582936"}},{"msg_index":42,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4733282800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg9qjf5zg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4733282800"},"demand_coin_denom":"denomX","order_price":"0.027451459376143762"}},{"msg_index":43,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"7061126800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxwpuzvh","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"7061126800"},"demand_coin_denom":"denomX","order_price":"0.027726799991752025"}},{"msg_index":44,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4655688000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8nhgh39","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4655688000"},"demand_coin_denom":"denomX","order_price":"0.027589129683947893"}},{"msg_index":45,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3026197200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cggv6mtwa","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3026197200"},"demand_coin_denom":"denomX","order_price":"0.027589129683947893"}},{"msg_index":46,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"7293911200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgf3v07n0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"7293911200"},"demand_coin_denom":"denomX","order_price":"0.027616663745508720"}},{"msg_index":47,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4810877600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgs80hl9n","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4810877600"},"demand_coin_denom":"denomX","order_price":"0.027534061560826241"}},{"msg_index":48,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4345308800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4345308800"},"demand_coin_denom":"denomX","order_price":"0.027451459376143762"}},{"msg_index":49,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"5509230800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgj52kuk7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"5509230800"},"demand_coin_denom":"denomX","order_price":"0.027368857191461284"}},{"msg_index":50,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4190119200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgnfuzftv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4190119200"},"demand_coin_denom":"denomX","order_price":"0.027451459376143762"}},{"msg_index":51,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"543163600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg5g94e2f","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"543163600"},"demand_coin_denom":"denomX","order_price":"0.027286255006778805"}},{"msg_index":52,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4578093200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg44npvhm","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4578093200"},"demand_coin_denom":"denomX","order_price":"0.027506527499265415"}},{"msg_index":53,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"6517963200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgkmq56ey","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"6517963200"},"demand_coin_denom":"denomX","order_price":"0.027368857191461284"}},{"msg_index":54,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"4190119200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cghxkq0yk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"4190119200"},"demand_coin_denom":"denomX","order_price":"0.027368857191461284"}},{"msg_index":55,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1939870000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgcemnnmw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1939870000"},"demand_coin_denom":"denomX","order_price":"0.027754334053312851"}},{"msg_index":56,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1163922000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgeyd8xxu","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1163922000"},"demand_coin_denom":"denomX","order_price":"0.027478993437704589"}},{"msg_index":57,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"5897204800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfqansamx","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"5897204800"},"demand_coin_denom":"denomX","order_price":"0.027644197807069546"}},{"msg_index":58,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"155189600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfpq9ygx5","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"155189600"},"demand_coin_denom":"denomX","order_price":"0.027671731868630372"}},{"msg_index":59,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2250249200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfzwk37gt","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2250249200"},"demand_coin_denom":"denomX","order_price":"0.027286255006778805"}},{"msg_index":60,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2948602400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfrnq9t4e","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2948602400"},"demand_coin_denom":"denomX","order_price":"0.027286255006778805"}},{"msg_index":61,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"7449100800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfyjejm5u","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"7449100800"},"demand_coin_denom":"denomX","order_price":"0.027313789068339631"}},{"msg_index":62,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"6129989200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf900xwfw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"6129989200"},"demand_coin_denom":"denomX","order_price":"0.027341323129900457"}},{"msg_index":63,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3491766000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfxpunc83","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3491766000"},"demand_coin_denom":"denomX","order_price":"0.027534061560826241"}},{"msg_index":64,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"6362773600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf8u28d6r","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"6362773600"},"demand_coin_denom":"denomX","order_price":"0.027726799991752025"}},{"msg_index":65,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"7138721600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfgr8539m","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"7138721600"},"demand_coin_denom":"denomX","order_price":"0.027534061560826241"}},{"msg_index":66,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3724550400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cff73qycf","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3724550400"},"demand_coin_denom":"denomX","order_price":"0.027616663745508720"}},{"msg_index":67,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3103792000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfsgjc9w4","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3103792000"},"demand_coin_denom":"denomX","order_price":"0.027589129683947893"}},{"msg_index":68,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"232784400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf34yvsn8","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"232784400"},"demand_coin_denom":"denomX","order_price":"0.027478993437704589"}},{"msg_index":69,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"6052394400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfjmhexac","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"6052394400"},"demand_coin_denom":"denomX","order_price":"0.027478993437704589"}},{"msg_index":70,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"5121256800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfnxpdnq2","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"5121256800"},"demand_coin_denom":"denomX","order_price":"0.027644197807069546"}}]`
	json.Unmarshal([]byte(swapMsgsJson), &swapMsgs)
	orderMap, XtoY, YtoX := types.GetOrderMap(swapMsgs, DenomX, DenomY, onlyNotMatched)
	require.NotZero(t, len(orderMap))
	require.NotNil(t, XtoY)
	require.NotNil(t, YtoX)

	// ErrInvalidDenom case
	require.Panics(t, func() {
		types.GetOrderMap(swapMsgs, "12421miklfdjnfiasdjfidosa8381813818---", DenomY, onlyNotMatched)
	})
}

// TODO: update as half-half fee version
// Reproduce next orderPrice is new on FindOrderMatch
//func TestFindOrderMatchEdgeCaseX(t *testing.T) {
//
//	direction := types.DirectionXtoY
//	executableAmt := sdk.NewInt(17329883339)
//	swapPrice, _ := sdk.NewDecFromStr("1.123759863025281136")
//	swapFeeRate, _ := sdk.NewDecFromStr("0.003000000000000000")
//	height := int64(0)
//
//	var swapMsgs []*types.BatchPoolSwapMsg
//	swapMsgsJson := `[{"msg_index":23,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1599684000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfzwk37gt","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1599684000"},"demand_coin_denom":"denomY","order_price":"1.136158417181026465"}},{"msg_index":1,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"26661400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"26661400"},"demand_coin_denom":"denomY","order_price":"1.135031275894140526"}},{"msg_index":14,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"773180600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgnfuzftv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"773180600"},"demand_coin_denom":"denomY","order_price":"1.135031275894140526"}},{"msg_index":4,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1519699800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1519699800"},"demand_coin_denom":"denomY","order_price":"1.133904134607254587"}},{"msg_index":18,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1653006800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cghxkq0yk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1653006800"},"demand_coin_denom":"denomY","order_price":"1.132776993320368648"}},{"msg_index":15,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"399921000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg5g94e2f","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"399921000"},"demand_coin_denom":"denomY","order_price":"1.131649852033482709"}},{"msg_index":26,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2372864600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf900xwfw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2372864600"},"demand_coin_denom":"denomY","order_price":"1.131649852033482709"}},{"msg_index":11,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1039794600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgs80hl9n","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1039794600"},"demand_coin_denom":"denomY","order_price":"1.130522710746596770"}},{"msg_index":12,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"613212200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"613212200"},"demand_coin_denom":"denomY","order_price":"1.130522710746596770"}},{"msg_index":17,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1573022600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgkmq56ey","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1573022600"},"demand_coin_denom":"denomY","order_price":"1.130522710746596770"}},{"msg_index":21,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"159968400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfqansamx","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"159968400"},"demand_coin_denom":"denomY","order_price":"1.130522710746596770"}},{"msg_index":8,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2132912000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8nhgh39","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2132912000"},"demand_coin_denom":"denomY","order_price":"1.127141286885938953"}},{"msg_index":13,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1386392800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgj52kuk7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1386392800"},"demand_coin_denom":"denomY","order_price":"1.127141286885938953"}},{"msg_index":28,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"426582400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf8u28d6r","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"426582400"},"demand_coin_denom":"denomY","order_price":"1.127141286885938953"}},{"msg_index":22,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"559889400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfpq9ygx5","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"559889400"},"demand_coin_denom":"denomY","order_price":"1.124887004312167075"}},{"msg_index":9,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1199763000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cggv6mtwa","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1199763000"},"demand_coin_denom":"denomY","order_price":"1.123759863025281136"}},{"msg_index":2,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"826503400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"826503400"},"demand_coin_denom":"denomY","order_price":"1.122632721738395197"}},{"msg_index":10,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"799842000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgf3v07n0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"799842000"},"demand_coin_denom":"denomY","order_price":"1.122632721738395197"}},{"msg_index":16,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1279747200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg44npvhm","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1279747200"},"demand_coin_denom":"denomY","order_price":"1.122632721738395197"}},{"msg_index":27,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2612817200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfxpunc83","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2612817200"},"demand_coin_denom":"denomY","order_price":"1.121505580451509258"}},{"msg_index":5,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1839636600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyayapl6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1839636600"},"demand_coin_denom":"denomY","order_price":"1.120378439164623319"}},{"msg_index":19,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"719857800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgcemnnmw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"719857800"},"demand_coin_denom":"denomY","order_price":"1.120378439164623319"}},{"msg_index":20,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"773180600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgeyd8xxu","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"773180600"},"demand_coin_denom":"denomY","order_price":"1.120378439164623319"}},{"msg_index":24,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"533228000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfrnq9t4e","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"533228000"},"demand_coin_denom":"denomY","order_price":"1.120378439164623319"}},{"msg_index":6,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1599684000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg9qjf5zg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1599684000"},"demand_coin_denom":"denomY","order_price":"1.119251297877737380"}},{"msg_index":7,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"2639478600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxwpuzvh","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"2639478600"},"demand_coin_denom":"denomY","order_price":"1.119251297877737380"}},{"msg_index":25,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1972943600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfyjejm5u","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1972943600"},"demand_coin_denom":"denomY","order_price":"1.119251297877737380"}},{"msg_index":3,"executed":true,"exchanged_offer_coin":{"denom":"denomX","amount":"0"},"remaining_offer_coin":{"denom":"denomX","amount":"1519699800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomX","amount":"1519699800"},"demand_coin_denom":"denomY","order_price":"1.118124156590851441"}}]`
//	json.Unmarshal([]byte(swapMsgsJson), &swapMsgs)
//	matchResultList, swapListExecuted, _, _ := types.FindOrderMatch(direction, swapMsgs, executableAmt, swapPrice, swapFeeRate, height)
//	require.Equal(t, 16, len(matchResultList))
//	require.Equal(t, 16, len(swapListExecuted))
//}

// TODO: update as half-half fee version
// Reproduce DirectionYtoX case on FindOrderMatch
//func TestFindOrderMatchEdgeCaseY(t *testing.T) {
//	direction := types.DirectionYtoX
//	executableAmt := sdk.NewInt(11376981000)
//	swapPrice, _ := sdk.NewDecFromStr("1.508781582653605557")
//	swapFeeRate, _ := sdk.NewDecFromStr("0.003000000000000000")
//	height := int64(0)
//
//	var swapMsgs []*types.BatchPoolSwapMsg
//	swapMsgsJson := `[{"msg_index":40,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"3413094300"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"3413094300"},"demand_coin_denom":"denomX","order_price":"1.487805473262194953"}},{"msg_index":42,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1462754700"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1462754700"},"demand_coin_denom":"denomX","order_price":"1.502788408541773956"}},{"msg_index":43,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"5092553400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"5092553400"},"demand_coin_denom":"denomX","order_price":"1.505784995597689756"}},{"msg_index":41,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1408578600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1408578600"},"demand_coin_denom":"denomX","order_price":"1.507283289125647657"}}]`
//	json.Unmarshal([]byte(swapMsgsJson), &swapMsgs)
//	matchResultList, swapListExecuted, _, _ := types.FindOrderMatch(direction, swapMsgs, executableAmt, swapPrice, swapFeeRate, height)
//	require.Equal(t, 4, len(matchResultList))
//	require.Equal(t, 4, len(swapListExecuted))
//}

// // TODO: update as half-half fee version
// Reproduce negative value of fractional rate case on FindOrderMatch
//func TestFindOrderMatchEdgeCaseNegativeFractionalRate(t *testing.T) {
//	direction := types.DirectionYtoX
//	executableAmt := sdk.NewInt(25835513200)
//	swapPrice, _ := sdk.NewDecFromStr("2.256160794978905238")
//	swapFeeRate, _ := sdk.NewDecFromStr("0.003000000000000000")
//	height := int64(0)
//
//	var swapMsgs []*types.BatchPoolSwapMsg
//	swapMsgsJson := `[{"msg_index":65,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1608716900"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgkmq56ey","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1608716900"},"demand_coin_denom":"denomX","order_price":"2.235855347824095091"}},{"msg_index":71,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2329037900"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfzwk37gt","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2329037900"},"demand_coin_denom":"denomX","order_price":"2.235855347824095091"}},{"msg_index":81,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"696310300"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfjmhexac","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"696310300"},"demand_coin_denom":"denomX","order_price":"2.235855347824095091"}},{"msg_index":49,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1296577800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1296577800"},"demand_coin_denom":"denomX","order_price":"2.240367669414052901"}},{"msg_index":59,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1392620600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgs80hl9n","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1392620600"},"demand_coin_denom":"denomX","order_price":"2.240367669414052901"}},{"msg_index":67,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2040909500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgcemnnmw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2040909500"},"demand_coin_denom":"denomX","order_price":"2.240367669414052901"}},{"msg_index":77,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"960428000"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfgr8539m","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"960428000"},"demand_coin_denom":"denomX","order_price":"2.240367669414052901"}},{"msg_index":83,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1632727600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf58c6rp0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1632727600"},"demand_coin_denom":"denomX","order_price":"2.240367669414052901"}},{"msg_index":55,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1584706200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgxwpuzvh","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1584706200"},"demand_coin_denom":"denomX","order_price":"2.242623830209031807"}},{"msg_index":60,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1368609900"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg36er2cp","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1368609900"},"demand_coin_denom":"denomX","order_price":"2.242623830209031807"}},{"msg_index":72,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1872834600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfrnq9t4e","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1872834600"},"demand_coin_denom":"denomX","order_price":"2.242623830209031807"}},{"msg_index":70,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"552246100"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfpq9ygx5","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"552246100"},"demand_coin_denom":"denomX","order_price":"2.244879991004010712"}},{"msg_index":52,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"696310300"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgrua237l","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"696310300"},"demand_coin_denom":"denomX","order_price":"2.247136151798989617"}},{"msg_index":63,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1224545700"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg5g94e2f","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1224545700"},"demand_coin_denom":"denomX","order_price":"2.247136151798989617"}},{"msg_index":76,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2016898800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf8u28d6r","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2016898800"},"demand_coin_denom":"denomX","order_price":"2.247136151798989617"}},{"msg_index":79,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"168074900"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfsgjc9w4","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"168074900"},"demand_coin_denom":"denomX","order_price":"2.247136151798989617"}},{"msg_index":61,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2208984400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgj52kuk7","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2208984400"},"demand_coin_denom":"denomX","order_price":"2.251648473388947428"}},{"msg_index":66,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"264117700"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cghxkq0yk","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"264117700"},"demand_coin_denom":"denomX","order_price":"2.251648473388947428"}},{"msg_index":75,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"96042800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfxpunc83","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"96042800"},"demand_coin_denom":"denomX","order_price":"2.251648473388947428"}},{"msg_index":53,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1824813200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgyayapl6","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1824813200"},"demand_coin_denom":"denomX","order_price":"2.253904634183926333"}},{"msg_index":64,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"192085600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg44npvhm","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"192085600"},"demand_coin_denom":"denomX","order_price":"2.256160794978905238"}},{"msg_index":69,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1152513600"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfqansamx","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1152513600"},"demand_coin_denom":"denomX","order_price":"2.256160794978905238"}},{"msg_index":78,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2016898800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cff73qycf","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2016898800"},"demand_coin_denom":"denomX","order_price":"2.256160794978905238"}},{"msg_index":68,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"552246100"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgeyd8xxu","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"552246100"},"demand_coin_denom":"denomX","order_price":"2.258416955773884143"}},{"msg_index":58,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"96042800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgf3v07n0","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"96042800"},"demand_coin_denom":"denomX","order_price":"2.260673116568863048"}},{"msg_index":73,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2232995100"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfyjejm5u","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2232995100"},"demand_coin_denom":"denomX","order_price":"2.260673116568863048"}},{"msg_index":51,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"600267500"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgzpt7yrd","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"600267500"},"demand_coin_denom":"denomX","order_price":"2.262929277363841954"}},{"msg_index":82,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2208984400"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cfnxpdnq2","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2208984400"},"demand_coin_denom":"denomX","order_price":"2.262929277363841954"}},{"msg_index":62,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"816363800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgnfuzftv","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"816363800"},"demand_coin_denom":"denomX","order_price":"2.265185438158820859"}},{"msg_index":54,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2329037900"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg9qjf5zg","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2329037900"},"demand_coin_denom":"denomX","order_price":"2.267441598953799764"}},{"msg_index":74,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"96042800"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf900xwfw","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"96042800"},"demand_coin_denom":"denomX","order_price":"2.267441598953799764"}},{"msg_index":80,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"216096300"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cf34yvsn8","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"216096300"},"demand_coin_denom":"denomX","order_price":"2.267441598953799764"}},{"msg_index":56,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"624278200"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cg8nhgh39","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"624278200"},"demand_coin_denom":"denomX","order_price":"2.269697759748778669"}},{"msg_index":50,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"1848823900"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgp0ctjdj","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"1848823900"},"demand_coin_denom":"denomX","order_price":"2.274210081338736480"}},{"msg_index":57,"executed":true,"exchanged_offer_coin":{"denom":"denomY","amount":"0"},"remaining_offer_coin":{"denom":"denomY","amount":"2232995100"},"msg":{"swap_requester_address":"cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cggv6mtwa","pool_id":1,"pool_type_index":1,"offer_coin":{"denom":"denomY","amount":"2232995100"},"demand_coin_denom":"denomX","order_price":"2.274210081338736480"}}]`
//	json.Unmarshal([]byte(swapMsgsJson), &swapMsgs)
//	_, _, poolXdelta, poolYdelta := types.FindOrderMatch(direction, swapMsgs, executableAmt, swapPrice, swapFeeRate, height)
//	require.Equal(t, poolXdelta, sdk.NewInt(-65243456686))
//	require.Equal(t, poolYdelta, sdk.NewInt(29004925600))
//	//require.Panics(t, func() {
//	//	types.FindOrderMatch(direction, swapMsgs, executableAmt, swapPrice, swapFeeRate, height)
//	//} )
//}

func TestCheckValidityOrderBook(t *testing.T) {
	currentPrice := sdk.MustNewDecFromStr("1.0")
	for _, testCase := range []struct {
		buyPrice  string
		sellPrice string
		valid     bool
	}{
		{
			buyPrice:  "0.99",
			sellPrice: "1.01",
			valid:     true,
		},
		{
			// maxBuyOrderPrice > minSellOrderPrice
			buyPrice:  "1.01",
			sellPrice: "0.99",
			valid:     false,
		},
		{
			buyPrice:  "1.1",
			sellPrice: "1.2",
			valid:     true,
		},
		{
			// maxBuyOrderPrice/currentPrice > 1.10
			buyPrice:  "1.11",
			sellPrice: "1.2",
			valid:     false,
		},
		{
			buyPrice:  "0.8",
			sellPrice: "0.9",
			valid:     true,
		},
		{
			// minSellOrderPrice/currentPrice < 0.90
			buyPrice:  "0.8",
			sellPrice: "0.89",
			valid:     false,
		},
	} {
		buyPrice := sdk.MustNewDecFromStr(testCase.buyPrice)
		sellPrice := sdk.MustNewDecFromStr(testCase.sellPrice)
		orderMap := types.OrderMap{
			buyPrice.String(): {
				OrderPrice:   buyPrice,
				BuyOfferAmt:  sdk.OneInt(),
				SellOfferAmt: sdk.ZeroInt(),
			},
			sellPrice.String(): {
				OrderPrice:   sellPrice,
				BuyOfferAmt:  sdk.ZeroInt(),
				SellOfferAmt: sdk.OneInt(),
			},
		}
		orderBook := orderMap.SortOrderBook()
		require.Equal(t, testCase.valid, types.CheckValidityOrderBook(orderBook, currentPrice))
	}
}
