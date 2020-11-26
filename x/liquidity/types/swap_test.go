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

func TestOrderMap(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
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
	msgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgs(ctx, batch)
	orderMap, XtoY, YtoX := types.GetOrderMap(msgs, denomX, denomY)
	orderBook := orderMap.SortOrderBook()
	currentPrice := X.Quo(Y).ToDec()
	require.Equal(t, orderMap[orderPriceList[0].String()].BuyOfferAmt, offerCoinList[0].Amount.MulRaw(3))
	require.Equal(t, orderMap[orderPriceList[0].String()].OrderPrice, orderPriceList[0])

	require.Equal(t,3, len(XtoY))
	require.Equal(t, 1, len(YtoX))
	require.Equal(t,3, len(orderMap[orderPriceList[0].String()].MsgList))
	require.Equal(t,1, len(orderMap[orderPriceListY[0].String()].MsgList))
	require.Equal(t,3, len(orderBook[0].MsgList))
	require.Equal(t,1, len(orderBook[1].MsgList))

	fmt.Println(orderBook, currentPrice)
	fmt.Println(XtoY, YtoX)

	clearedXtoY, clearedYtoX := types.ClearOrders(XtoY, YtoX, ctx.BlockHeight(), false)
	require.Equal(t, XtoY, clearedXtoY)
	require.Equal(t, YtoX, clearedYtoX)

	require.False(t, types.CheckValidityOrderBook(orderBook, currentPrice))

	currentYPriceOverX := X.Quo(Y).ToDec()
	//direction := types.GetPriceDirection(currentYPriceOverX, orderBook)
	result := types.ComputePriceDirection(X.ToDec(), Y.ToDec(), currentYPriceOverX, orderBook)

	require.NotEqual(t, types.NoMatch, result.MatchType)

	matchResultXtoY, _, poolXDeltaXtoY, poolYDeltaXtoY := types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX,
		result.SwapPrice, sdk.ZeroDec(), ctx.BlockHeight())
	matchResultYtoX, _, poolXDeltaYtoX, poolYDeltaYtoX := types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY,
		result.SwapPrice, sdk.ZeroDec(), ctx.BlockHeight())

	XtoY, YtoX, XDec, YDec, poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY :=
		simapp.LiquidityKeeper.UpdateState(X.ToDec(), Y.ToDec(), XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	clearedXtoY, clearedYtoX = types.ClearOrders(XtoY, YtoX, ctx.BlockHeight(), true)
	require.Equal(t, 0, (types.MsgList)(clearedXtoY).LenRemainingMsgs())
	require.Equal(t, 0, (types.MsgList)(clearedXtoY).LenFractionalMsgs())
	require.Equal(t, 0, (types.MsgList)(clearedYtoX).LenRemainingMsgs())
	require.Equal(t, 0, (types.MsgList)(clearedYtoX).LenFractionalMsgs())
	require.Equal(t,1, len(clearedYtoX))
	require.Equal(t,0, len(clearedXtoY))

	fmt.Println(matchResultXtoY)
	fmt.Println(poolXDeltaXtoY)
	fmt.Println(poolYDeltaXtoY)

	fmt.Println(poolXDeltaYtoX, poolYDeltaYtoX)
	fmt.Println(poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY)
	fmt.Println(decimalErrorX, decimalErrorY)
	fmt.Println(XDec, YDec)
	// TODO: detailed assertion
	// TODO: debug Ydec 999970003, poolYdelta2, poolYDeltaXtoY -29997


	orderMapExecuted, _, _ := types.GetOrderMap(append(clearedXtoY, clearedYtoX...), denomX, denomY)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	lastPrice := XDec.Quo(YDec)
	require.True(t, types.CheckValidityOrderBook(orderBookExecuted, lastPrice))

	require.Equal(t,0, (types.MsgList)(orderMapExecuted[orderPriceList[0].String()].MsgList).LenRemainingMsgs())
	require.Equal(t,0, (types.MsgList)(orderMapExecuted[orderPriceListY[0].String()].MsgList).LenRemainingMsgs())
	require.Equal(t,0, (types.MsgList)(orderBookExecuted[0].MsgList).LenRemainingMsgs())

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
		OrderPrice: a,
		BuyOfferAmt: sdk.NewInt(50000000),
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
	result := types.ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)
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
	result = types.ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)
	require.Equal(t, types.CalculateMatch(direction, X, Y, currentYPriceOverX, orderBook), result)

	// stay case
	orderMap = make(types.OrderMap)
	a, _ = sdk.NewDecFromStr("1.0")

	orderMap[a.String()] = types.OrderByPrice{
		OrderPrice: a,
		BuyOfferAmt: sdk.NewInt(50000000),
		SellOfferAmt: sdk.NewInt(50000000),
	}
	orderBook = orderMap.SortOrderBook()

	X = orderMap[a.String()].BuyOfferAmt.ToDec()
	Y = orderMap[a.String()].SellOfferAmt.ToDec()
	currentYPriceOverX = X.Quo(Y)

	result = types.ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)
	require.Equal(t, types.CalculateMatchStay(currentYPriceOverX, orderBook), result)
}
