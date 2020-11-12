package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderBookSort(t *testing.T) {
	orderMap := make(OrderMap)
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")
	orderMap[a] = OrderByPrice{
		OrderPrice: a,
		BuyOrderAmt: sdk.ZeroInt(),
		SellOrderAmt: sdk.ZeroInt(),
	}
	orderMap[b] = OrderByPrice{
		OrderPrice: b,
		BuyOrderAmt: sdk.ZeroInt(),
		SellOrderAmt: sdk.ZeroInt(),
	}
	orderMap[c] = OrderByPrice{
		OrderPrice: c,
		BuyOrderAmt: sdk.ZeroInt(),
		SellOrderAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()
	fmt.Println(orderBook)

	res := orderBook.Less(0,1)
	require.True(t, res)
	res = orderBook.Less(1,2)
	require.True(t, res)
	res = orderBook.Less(2,1)
	require.False(t, res)

	orderBook.Swap(1,2)
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

	require.Equal(t, a, MinDec(a, b))
	require.Equal(t, a, MinDec(a, c))
	require.Equal(t, b, MaxDec(a, b))
	require.Equal(t, c, MaxDec(a, c))
	require.Equal(t, a, MaxDec(a, a))
	require.Equal(t, a, MinDec(a, a))
}

func TestGetExecutableAmt(t *testing.T) {
	orderMap := make(OrderMap)
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")
	orderMap[a] = OrderByPrice{
		OrderPrice:   a,
		BuyOrderAmt:  sdk.ZeroInt(),
		SellOrderAmt: sdk.NewInt(30000000),
	}
	orderMap[b] = OrderByPrice{
		OrderPrice:   b,
		BuyOrderAmt:  sdk.NewInt(90000000),
		SellOrderAmt: sdk.ZeroInt(),
	}
	orderMap[c] = OrderByPrice{
		OrderPrice:   c,
		BuyOrderAmt:  sdk.NewInt(50000000),
		SellOrderAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	executableBuyAmtX, executableSellAmtY := GetExecutableAmt(b, orderBook)
	require.Equal(t, sdk.NewInt(140000000), executableBuyAmtX)
	require.Equal(t, sdk.NewInt(30000000), executableSellAmtY)
}


// TODO: WIP
func TestGetPriceDirection(t *testing.T) {

	// decrease case
	orderMap := make(OrderMap)
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")
	orderMap[a] = OrderByPrice{
		OrderPrice:   a,
		BuyOrderAmt:  sdk.ZeroInt(),
		SellOrderAmt: sdk.NewInt(30000000),
	}
	orderMap[b] = OrderByPrice{
		OrderPrice:   b,
		BuyOrderAmt:  sdk.NewInt(90000000),
		SellOrderAmt: sdk.ZeroInt(),
	}
	orderMap[c] = OrderByPrice{
		OrderPrice:   c,
		BuyOrderAmt:  sdk.NewInt(50000000),
		SellOrderAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	// increase
	X := sdk.NewInt(10000000).ToDec()
	Y := sdk.NewInt(50000000).ToDec()
	currentYPriceOverX := X.Quo(Y)
	require.Equal(t, currentYPriceOverX, b)
	result := GetPriceDirection(currentYPriceOverX, orderBook)
	require.Equal(t, Increase, result)

	// decrease
	X = sdk.NewInt(100000000).ToDec()
	Y = sdk.NewInt(50000000).ToDec()
	currentYPriceOverX = X.Quo(Y)
	result = GetPriceDirection(currentYPriceOverX, orderBook)
	require.Equal(t, Decrease, result)

	// TODO: stay case
}


// TODO: WIP
func TestComputePriceDirection(t *testing.T) {

	// decrease case
	orderMap := make(OrderMap)
	a, _ := sdk.NewDecFromStr("0.1")
	b, _ := sdk.NewDecFromStr("0.2")
	c, _ := sdk.NewDecFromStr("0.3")
	orderMap[a] = OrderByPrice{
		OrderPrice:   a,
		BuyOrderAmt:  sdk.ZeroInt(),
		SellOrderAmt: sdk.NewInt(3),
	}
	orderMap[b] = OrderByPrice{
		OrderPrice:   b,
		BuyOrderAmt:  sdk.NewInt(9),
		SellOrderAmt: sdk.ZeroInt(),
	}
	orderMap[c] = OrderByPrice{
		OrderPrice:   c,
		BuyOrderAmt:  sdk.NewInt(5),
		SellOrderAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	X := sdk.NewInt(100000000).ToDec()
	Y := sdk.NewInt(50000000).ToDec()
	currentYPriceOverX := X.Quo(Y)
	result := ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)

	fmt.Println(X, Y, currentYPriceOverX)
	fmt.Println(result)


	// increase case
	orderMap[c] = OrderByPrice{
		OrderPrice:   c,
		BuyOrderAmt:  sdk.ZeroInt(),
		SellOrderAmt: sdk.NewInt(3),
	}
	orderMap[b] = OrderByPrice{
		OrderPrice:   b,
		BuyOrderAmt:  sdk.NewInt(9),
		SellOrderAmt: sdk.ZeroInt(),
	}
	orderMap[a] = OrderByPrice{
		OrderPrice:   a,
		BuyOrderAmt:  sdk.NewInt(5),
		SellOrderAmt: sdk.ZeroInt(),
	}
	// make orderbook to sort orderMap
	orderBook = orderMap.SortOrderBook()

	X = sdk.NewInt(100000000).ToDec()
	Y = sdk.NewInt(50000000).ToDec()
	currentYPriceOverX = X.Quo(Y)
	result = ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)

	fmt.Println(X, Y, currentYPriceOverX)
	fmt.Println(result)
}