package keeper

import "C"
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"sort"
)

// TODO: move to types/swap.go
const (
	// Price Directions
	Increase = 1
	Decrease = -1
	Stay     = 0

	// Match Types
	ExactMatch      = 1
	FractionalMatch = 0
)

type OrderByPrice struct {
	OrderPrice   sdk.Dec
	BuyOrderAmt  sdk.Int
	SellOrderAmt sdk.Int
}
type OrderBook []OrderByPrice

// Len implements sort.Interface for OrderBook
func (orderBook OrderBook) Len() int { return len(orderBook) }

// Less implements sort.Interface for OrderBook
func (orderBook OrderBook) Less(i, j int) bool {
	return orderBook[i].OrderPrice.LT(orderBook[j].OrderPrice)
}

// Swap implements sort.Interface for OrderBook
func (orderBook OrderBook) Swap(i, j int) { orderBook[i], orderBook[j] = orderBook[j], orderBook[i] }

func (orderBook OrderBook) Sort() {
	sort.Sort(orderBook)
}

func (orderBook OrderBook) Reverse() {
	sort.Reverse(orderBook)
}

type OrderMap map[sdk.Dec]OrderByPrice

// TODO: testcode
func (orderMap OrderMap) SortOrderBook() (orderBook OrderBook) {
	orderPriceList := make([]sdk.Dec, 0, len(orderMap))
	for k := range orderMap {
		orderPriceList = append(orderPriceList, k)
	}

	// TODO: verify sort
	sort.Slice(orderPriceList, func(i, j int) bool {
		return orderPriceList[i].LT(orderPriceList[j])
	})

	for _, k := range orderPriceList {
		orderBook = append(orderBook, OrderByPrice{
			OrderPrice:   k,
			BuyOrderAmt:  orderMap[k].BuyOrderAmt,
			SellOrderAmt: orderMap[k].SellOrderAmt,
		})
	}
	return orderBook
}

// TODO: WIP https://repl.it/@HyungyeonLee/batchExecution#main.py
// TODO: testcode
func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	pool, found := k.GetLiquidityPool(ctx, liquidityPoolBatch.PoolID)
	if !found {
		return types.ErrPoolNotExists
	}
	//totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort() // TODO: validate alphabetical

	X := reserveCoins[0].Amount.ToDec()
	Y := reserveCoins[1].Amount.ToDec()
	currentYPriceOverX := X.Quo(Y)
	var XtoY []*types.MsgSwap // buying Y from X
	var YtoX []*types.MsgSwap // selling Y for X
	//var orderBook OrderBook
	//orderMap := make(map[sdk.Dec]OrderByPrice)
	orderMap := make(OrderMap)

	var sumOfBuy sdk.Coin
	var sumOfSell sdk.Coin
	swapMsgs := k.GetAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	for _, m := range swapMsgs {
		if m.Msg.OfferCoin.Denom == sumOfSell.Denom {
			sumOfSell = sumOfSell.Add(m.Msg.OfferCoin)
			YtoX = append(YtoX, m.Msg)
			if _, ok := orderMap[m.Msg.OrderPrice]; ok {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice].BuyOrderAmt,
					orderMap[m.Msg.OrderPrice].SellOrderAmt.Add(m.Msg.OfferCoin.Amount)}
			} else {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{m.Msg.OrderPrice, sdk.ZeroInt(), m.Msg.OfferCoin.Amount}
			}
		} else if m.Msg.DemandCoin.Denom == sumOfBuy.Denom {
			sumOfBuy = sumOfBuy.Add(m.Msg.DemandCoin)
			XtoY = append(XtoY, m.Msg)
			if _, ok := orderMap[m.Msg.OrderPrice]; ok {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice].BuyOrderAmt.Add(m.Msg.OfferCoin.Amount),
					orderMap[m.Msg.OrderPrice].SellOrderAmt}
			} else {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{m.Msg.OrderPrice, m.Msg.OfferCoin.Amount, sdk.ZeroInt()}
			}
		} else {
			return types.ErrInvalidDenom
		}
	}

	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	//orderPriceList := make([]sdk.Dec, 0, len(orderMap))
	//for k := range orderMap {
	//	orderPriceList = append(orderPriceList, k)
	//}
	//
	//// TODO: verify sort
	//sort.Slice(orderPriceList, func(i, j int) bool {
	//	return orderPriceList[i].LT(orderPriceList[j])
	//})
	//
	//for _, k := range orderPriceList {
	//	orderBook = append(orderBook, OrderByPrice{
	//		OrderPrice:   k,
	//		BuyOrderAmt:  orderMap[k].BuyOrderAmt,
	//		SellOrderAmt: orderMap[k].SellOrderAmt,
	//	})
	//}
	// TODO: verify sorted by orderPrice
	sort.SliceStable(XtoY, func(i, j int) bool {
		return XtoY[i].OrderPrice.GT(XtoY[j].OrderPrice)
	})
	sort.SliceStable(YtoX, func(i, j int) bool {
		return YtoX[i].OrderPrice.LT(YtoX[j].OrderPrice)
	})

	priceDirection := GetPriceDirection(currentYPriceOverX, orderBook)
	if priceDirection == Decrease {
		orderBook.Reverse()
	}

	var matchType int
	var poolX, poolY, swapPrice, lastOrderPrice, decEX, decEY sdk.Dec
	var EX, EY, originalEX, originalEY sdk.Int

	// price does not change
	if priceDirection == Stay {
		swapPrice := currentYPriceOverX
		originalEX, originalEY = GetExecutableAmt(swapPrice, orderBook)
		poolY = sdk.ZeroDec()
		EX := originalEX.ToDec()
		EY := originalEY.ToDec()

		// fractionalMatch
		if EX.Equal(swapPrice.Mul(EY)) {
			matchType = ExactMatch
		} else {
			matchType = FractionalMatch
			if EX.GT(swapPrice.Mul(EY)) {
				EX = swapPrice.Mul(EY)
			} else if EX.GT(swapPrice.Mul(EY)) {
				EY = EX.Quo(swapPrice)
			}
		}

	// price increases or decrease
	} else {
		lastOrderPrice = currentYPriceOverX
		var lastExecutableBuyAmtX, lastExecutableSellAmtY, executableBuyAmtX, executableSellAmtY sdk.Int
		var swapPrice sdk.Dec

		// iterate orderbook from current price to upwards(increase)/downwards(decrease)
		for _, order := range orderBook {
			orderPrice := order.OrderPrice

			// calculate executable amounts in X coins
			executableBuyAmtX, executableSellAmtY = GetExecutableAmt(orderPrice, orderBook)

			// #### variables #####
			// swapPrice : swap price for this batch
			// PoolY : Y coins provided by the liquidity pool
			// EX : X coins provided by users
			// EY : Y coins provided by users

			// #### equations for price increase #####
			// 1) swap equation : EX = EY*SwapPrice + PoolY*SwapPrice --> PoolY = EX/SwapPrice - EY
			// 2) constant product equation : X*Y = (X+SwapPrice*PoolY)*(Y-PoolY) --> PoolY = Y - X/SwapPrice
			// 3) 1) & 2) : EX/SwapPrice - EY = Y - X/SwapPrice --> SwapPrice = (X + EX)/(Y + EY)
			// 4) 1) & 3) : PoolY = (Y*EX - X*EY)/(X + EX)

			// #### equations for price decrease #####
			// 1) swap equation : EY = EX/SwapPrice + PoolX/SwapPrice --> PoolX = EY*SwapPrice - EX
			// 2) constant product equation : X*Y = (X-PoolX)*(Y-PoolX/SwapPrice) --> PoolX = X - Y*SwapPrice
			// 3) 1) & 2) : EY*SwapPrice - EX = X - Y*SwapPrice --> SwapPrice = (X + EX)/(Y + EY)
			// 4) 1) & 3) : PoolX = (X*EY - Y*EX)/(Y + EY)

			// simulation) check whether all executable EX/EY are matched

			if priceDirection == Increase {
				EX = executableBuyAmtX
				EY = executableBuyAmtX
				swapPrice = orderPrice
				poolY = Y.Sub(X.Quo(swapPrice))

				//check all EX are matched
				if EX.ToDec().Sub(EY.ToDec().Mul(swapPrice)).Sub(poolY.Mul(swapPrice)).IsNegative() {
					// check whether exactMatch is possible in last price range (lastOrderPrice ~ orderPrice)
					EY = lastExecutableSellAmtY
					swapPrice = X.Add(EX.ToDec()).Quo(Y.Add(EY.ToDec()))
					poolY = Y.Mul(EX.ToDec()).Sub(X.Mul(EY.ToDec())).Quo(X.Add(EX.ToDec()))
					if lastOrderPrice.LT(swapPrice) && swapPrice.LT(orderPrice) && poolY.GTE(sdk.ZeroDec()) {
						matchType = ExactMatch
						break
					}
				} else { // exactMatch is not found --> fractionalMatch
					matchType = FractionalMatch
					EX = lastExecutableBuyAmtX
					EY = lastExecutableSellAmtY
					swapPrice = lastOrderPrice
					poolY = Y.Sub(X.Quo(swapPrice))
					break
				}
			} else if priceDirection == Decrease {
				EX = executableBuyAmtX
				EY = executableBuyAmtX
				swapPrice := orderPrice
				poolX = X.Sub(Y.Quo(swapPrice))

				// check all EY are matched
				if EY.ToDec().Sub(EX.ToDec().Quo(swapPrice)).Sub(poolX.Quo(swapPrice)).IsNegative() {
					// check whether exactMatch is possible in last price range (orderPrice ~ lastOrderPrice)
					EX = lastExecutableBuyAmtX
					swapPrice = X.Add(EX.ToDec()).Quo(Y.Add(EY.ToDec()))
					poolX = X.Mul(EY.ToDec()).Sub(Y.Mul(EY.ToDec())).Quo(Y.Add(EY.ToDec()))
					// check swapPrice within given price range
					if orderPrice.LT(swapPrice) && swapPrice.LT(lastOrderPrice) && poolX.GTE(sdk.ZeroDec()) {
						matchType = ExactMatch  // all orders are exactly matched
						break
					}
				} else { // exactMatch is not found --> fractionalMatch
					matchType = FractionalMatch
					EX = lastExecutableBuyAmtX
					EY = lastExecutableSellAmtY
					swapPrice = lastOrderPrice
					poolX = X.Sub(Y.Quo(swapPrice))
					break
				}
			}

			// update last variables
			lastOrderPrice = orderPrice
			lastExecutableBuyAmtX = executableBuyAmtX
			lastExecutableSellAmtY = executableSellAmtY

		}

		originalEX = EX
		originalEY = EY
		// fractional match for EX
		if matchType == FractionalMatch {
			if priceDirection == Increase {
				frac := EY.ToDec().Mul(swapPrice).Add(poolY.Mul(swapPrice))
				if EX.ToDec().LT(frac) {
					decEX = EX.ToDec()
				} else {
					decEX = frac
				}
			} else if priceDirection == Decrease {
				frac := EX.ToDec().Mul(swapPrice).Add(poolX.Quo(swapPrice))
				if EY.ToDec().LT(frac) {
					decEY = EY.ToDec()
				} else {
					decEY = frac
				}
			}
		}
	}
	var EXFractionalRatio, EYFractionalRatio sdk.Dec
	if originalEX.IsZero() {
		EXFractionalRatio = sdk.NewDec(1)
	} else {
		EXFractionalRatio = EX.ToDec().Quo(originalEX.ToDec())
	}

	if originalEY.IsZero() {
		EYFractionalRatio = sdk.NewDec(1)
	} else {
		EYFractionalRatio = EY.ToDec().Quo(originalEY.ToDec())
	}

	ctx.Logger().Info("Swap Execution Result",
		"priceDirection", priceDirection,
		"currentYPriceOverX", currentYPriceOverX,
		"swapPrice", swapPrice,
		"priceChangeRatio", swapPrice.Quo(currentYPriceOverX.Sub(sdk.NewDec(1))),
		"matchType", matchType,
		"originEx", originalEX,
		"EX", EX,
		"decEX", decEX,
		"EXFractionalRatio", EXFractionalRatio,
		"originalEY", originalEY,
		"EY", EY,
		"decEY", decEY,
		"EYFractionalRatio", EYFractionalRatio,
		)

	buySellAmtCheck := EX.ToDec().Sub(EY.ToDec().Mul(swapPrice)).Sub(poolY.Mul(swapPrice))
	if priceDirection == Increase {
		ctx.Logger().Info("Swap Direction Increase",
			"PoolY", poolY,
			"buy-sell amount check(should be zero): ", buySellAmtCheck,
			"constant product check(should be zero): ", X.Mul(Y).Sub(X.Add(swapPrice.Mul(poolY)).Mul(Y.Sub(poolY))),
			)
	} else if priceDirection == Decrease {
		ctx.Logger().Info("Swap Direction Decrease",
			"PoolX", poolX,
			"buy-sell amount check(should be zero): ", EY.ToDec().Sub(EX.ToDec().Quo(swapPrice)).Sub(poolX.Quo(swapPrice)),
			"constant product check(should be zero): ", X.Mul(Y).Sub(X.Sub(poolX).Mul(Y.Add(poolX.Quo(swapPrice)))),
			)
	} else if priceDirection == Stay {
		ctx.Logger().Info("Swap Direction Stay",
			"Pool", nil,
			"buy-sell amount check(should be zero): ", buySellAmtCheck,
			"constant product check(should be zero): ", nil,
		)
	}

	// finding orders to be matched
	// TODO: WIP
	//for _, order := range XtoY {
	//	if order.OrderPrice == lastOrderPrice
	//}




	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
	return nil
}

// TODO: testcode
func GetPriceDirection(currentYPriceOverX sdk.Dec, orderBook OrderBook) int {
	buyAmtOverCurrentPrice := sdk.ZeroInt()
	buyAmtAtCurrentPrice := sdk.ZeroInt()
	sellAmtUnderCurrentPrice := sdk.ZeroInt()
	sellAmtAtCurrentPrice := sdk.ZeroInt()

	for _, order := range orderBook {
		if order.OrderPrice.GT(currentYPriceOverX) {
			buyAmtOverCurrentPrice = buyAmtOverCurrentPrice.Add(order.BuyOrderAmt)
		} else if order.OrderPrice.Equal(currentYPriceOverX) {
			buyAmtAtCurrentPrice = buyAmtAtCurrentPrice.Add(order.BuyOrderAmt)
			sellAmtAtCurrentPrice = sellAmtAtCurrentPrice.Add(order.SellOrderAmt)
		} else if order.OrderPrice.LT(currentYPriceOverX) {
			sellAmtUnderCurrentPrice = sellAmtUnderCurrentPrice.Add(order.SellOrderAmt)
		} else {
			// TODO: err
		}
	}
	// TODO: verify Dec, Int math logic
	if buyAmtOverCurrentPrice.ToDec().GT(currentYPriceOverX.MulInt(sellAmtUnderCurrentPrice.Add(sellAmtAtCurrentPrice))) {
		return Increase
	} else if currentYPriceOverX.MulInt(sellAmtAtCurrentPrice).GT(buyAmtOverCurrentPrice.Add(buyAmtAtCurrentPrice).ToDec()) {
		return Decrease
	} else {
		return Stay
	}
}

// TODO: testcode
func GetExecutableAmt(swapPrice sdk.Dec, orderBook OrderBook) (executableBuyAmtX, executableSellAmtY sdk.Int) {
	for _, order := range orderBook {
		if order.OrderPrice.GTE(swapPrice) {
			executableBuyAmtX = executableBuyAmtX.Add(order.BuyOrderAmt)
		}
		if order.OrderPrice.GTE(swapPrice) {
			executableSellAmtY = executableSellAmtY.Add(order.SellOrderAmt)
		}
	}
	return
}
