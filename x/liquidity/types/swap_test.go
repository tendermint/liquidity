package types_test

import (
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
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:1], poolId, true)

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
	simapp.LiquidityKeeper.SetLiquidityPoolBatchSwapMsgs(ctx, poolId, *(msgs))

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

	fmt.Println(orderBook, currentPrice)
	fmt.Println(XtoY, YtoX)

	clearedXtoY := types.ValidateStateAndExpireOrders(XtoY, ctx.BlockHeight(), false)
	clearedYtoX := types.ValidateStateAndExpireOrders(YtoX, ctx.BlockHeight(), false)
	require.Equal(t, XtoY, clearedXtoY)
	require.Equal(t, YtoX, clearedYtoX)

	require.False(t, types.CheckValidityOrderBook(orderBook, currentPrice))

	currentYPriceOverX := X.Quo(Y).ToDec()

	// The price and coins of swap messages in orderbook are calculated
	// to derive match result with the price direction.
	result := types.MatchOrderbook(X.ToDec(), Y.ToDec(), currentYPriceOverX, orderBook)
	require.NotEqual(t, types.NoMatch, result.MatchType)

	matchResultXtoY, _, poolXDeltaXtoY, poolYDeltaXtoY := types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX,
		result.SwapPrice, sdk.ZeroDec(), ctx.BlockHeight())
	matchResultYtoX, _, poolXDeltaYtoX, poolYDeltaYtoX := types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY,
		result.SwapPrice, sdk.ZeroDec(), ctx.BlockHeight())

	XtoY, YtoX, XDec, YDec, poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY :=
		simapp.LiquidityKeeper.UpdateState(X.ToDec(), Y.ToDec(), XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	require.Equal(t, 0, (types.MsgList)(clearedXtoY).CountNotMatchedMsgs())
	require.Equal(t, 0, (types.MsgList)(clearedXtoY).CountFractionalMatchedMsgs())
	require.Equal(t, 1, (types.MsgList)(clearedYtoX).CountNotMatchedMsgs())
	require.Equal(t, 0, (types.MsgList)(clearedYtoX).CountFractionalMatchedMsgs())
	require.Equal(t, 3, len(clearedXtoY))
	require.Equal(t, 1, len(clearedYtoX))

	fmt.Println(matchResultXtoY)
	fmt.Println(poolXDeltaXtoY)
	fmt.Println(poolYDeltaXtoY)

	fmt.Println(poolXDeltaYtoX, poolYDeltaYtoX)
	fmt.Println(poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY)
	fmt.Println(decimalErrorX, decimalErrorY)
	fmt.Println(XDec, YDec)

	// Verify swap result by creating an orderbook with remaining messages that have been matched and not transacted.
	orderMapExecuted, _, _ := types.GetOrderMap(append(clearedXtoY, clearedYtoX...), denomX, denomY, true)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	lastPrice := XDec.Quo(YDec)
	fmt.Println("lastPrice", lastPrice)
	fmt.Println("X", XDec)
	fmt.Println("Y", YDec)
	require.True(t, types.CheckValidityOrderBook(orderBookExecuted, lastPrice))

	require.Equal(t, 0, (types.MsgList)(orderMapExecuted[orderPriceList[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 1, (types.MsgList)(orderMapExecuted[orderPriceListY[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 1, (types.MsgList)(orderBookExecuted[0].MsgList).CountNotMatchedMsgs())

	clearedXtoY = types.ValidateStateAndExpireOrders(XtoY, ctx.BlockHeight(), true)
	clearedYtoX = types.ValidateStateAndExpireOrders(YtoX, ctx.BlockHeight(), true)

	orderMapCleared, _, _ := types.GetOrderMap(append(clearedXtoY, clearedYtoX...), denomX, denomY, true)
	orderBookCleared := orderMapCleared.SortOrderBook()
	require.True(t, types.CheckValidityOrderBook(orderBookCleared, lastPrice))

	require.Equal(t, 0, (types.MsgList)(orderMapCleared[orderPriceList[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 0, (types.MsgList)(orderMapCleared[orderPriceListY[0].String()].MsgList).CountNotMatchedMsgs())
	require.Equal(t, 0, len(orderBookCleared))
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
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:1], poolId, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	maxOrderRatio := types.GetMaxOrderRatio()

	// Success case, not exceed GetMaxOrderRatio orders
	priceBuy, _ := sdk.NewDecFromStr("1.1")
	priceSell, _ := sdk.NewDecFromStr("1.2")

	offerCoin := sdk.NewCoin(denomX, sdk.NewInt(100))
	offerCoinY := sdk.NewCoin(denomY, sdk.NewInt(100))

	app.SaveAccount(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin))
	app.SaveAccount(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY))

	msgBuy := types.NewMsgSwap(addrs[1], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoin, DenomY, priceBuy)
	msgSell := types.NewMsgSwap(addrs[2], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoinY, DenomY, priceSell)

	_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy)
	require.NoError(t, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell)
	require.NoError(t, err)

	// Fail case, exceed GetMaxOrderRatio orders
	offerCoin = sdk.NewCoin(denomX, X)
	offerCoinY = sdk.NewCoin(denomY, Y)

	app.SaveAccount(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin))
	app.SaveAccount(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY))

	msgBuy = types.NewMsgSwap(addrs[1], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoin, DenomY, priceBuy)
	msgSell = types.NewMsgSwap(addrs[2], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoinY, DenomY, priceSell)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

	// Success case, same GetMaxOrderRatio orders
	offerCoin = sdk.NewCoin(denomX, X.ToDec().Mul(maxOrderRatio).TruncateInt())
	offerCoinY = sdk.NewCoin(denomY, Y.ToDec().Mul(maxOrderRatio).TruncateInt())

	app.SaveAccount(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin))
	app.SaveAccount(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY))

	msgBuy = types.NewMsgSwap(addrs[1], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoin, DenomY, priceBuy)
	msgSell = types.NewMsgSwap(addrs[2], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoinY, DenomY, priceSell)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy)
	require.NoError(t, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell)
	require.NoError(t, err)

	// Success case, same GetMaxOrderRatio orders
	offerCoin = sdk.NewCoin(denomX, X.ToDec().Mul(maxOrderRatio).TruncateInt().AddRaw(1))
	offerCoinY = sdk.NewCoin(denomY, Y.ToDec().Mul(maxOrderRatio).TruncateInt().AddRaw(1))

	app.SaveAccount(simapp, ctx, addrs[1], sdk.NewCoins(offerCoin))
	app.SaveAccount(simapp, ctx, addrs[2], sdk.NewCoins(offerCoinY))

	msgBuy = types.NewMsgSwap(addrs[1], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoin, DenomY, priceBuy)
	msgSell = types.NewMsgSwap(addrs[2], poolId, DefaultPoolTypeIndex, DefaultSwapType, offerCoinY, DenomY, priceSell)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgBuy)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msgSell)
	require.Equal(t, types.ErrExceededMaxOrderable, err)

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
