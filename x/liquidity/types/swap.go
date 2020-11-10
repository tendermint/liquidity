package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

const (
	// Price Directions
	Increase = 1
	Decrease = 2
	Stay     = 3

	// Match Types
	ExactMatch      = 1
	NoMatch         = 2
	FractionalMatch = 3

	OrderLifeSpanHeight = 0

	DirectionXtoY = 1
	DirectionYtoX = 2

	DefaultSwapType = 0
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
	//sort.Sort(orderBook)
	sort.Slice(orderBook, func(i, j int) bool {
		return orderBook[i].OrderPrice.LT(orderBook[j].OrderPrice)
	})
}

func (orderBook OrderBook) Reverse() {
	//sort.Reverse(orderBook)
	sort.Slice(orderBook, func(i, j int) bool {
		return orderBook[i].OrderPrice.GT(orderBook[j].OrderPrice)
	})
}

func MinDec(a, b sdk.Dec) sdk.Dec {
	if a.LTE(b) {
		return a
	} else {
		return b
	}
}

func MaxDec(a, b sdk.Dec) sdk.Dec {
	if a.GTE(b) {
		return a
	} else {
		return b
	}
}

type OrderMap map[sdk.Dec]OrderByPrice

func (orderMap OrderMap) SortOrderBook() (orderBook OrderBook) {
	orderPriceList := make([]sdk.Dec, 0, len(orderMap))
	for k := range orderMap {
		orderPriceList = append(orderPriceList, k)
	}

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

type BatchResult struct {
	MatchType   int
	SwapPrice   sdk.Dec
	EX          sdk.Dec
	EY          sdk.Dec
	OriginalEX  sdk.Int
	OriginalEY  sdk.Int
	PoolX       sdk.Dec
	PoolY       sdk.Dec
	TransactAmt sdk.Dec
}

func NewBatchResult() BatchResult {
	return BatchResult{
		SwapPrice:   sdk.ZeroDec(),
		EX:          sdk.ZeroDec(),
		EY:          sdk.ZeroDec(),
		OriginalEX:  sdk.ZeroInt(),
		OriginalEY:  sdk.ZeroInt(),
		PoolX:       sdk.ZeroDec(),
		PoolY:       sdk.ZeroDec(),
		TransactAmt: sdk.ZeroDec(),
	}
}

type MatchResult struct {
	OrderHeight       int64
	OrderCancelHeight int64
	OrderMsgIndex     uint64
	OrderPrice        sdk.Dec
	OrderAmt          sdk.Int
	MatchedAmt        sdk.Int
	RefundAmt         sdk.Int
	ResidualAmt       sdk.Int
	ReceiveAmt        sdk.Int
	FeeAmt            sdk.Int
}

func ComputePriceDirection(X, Y, currentPrice sdk.Dec, orderBook OrderBook) (result BatchResult) {
	result = NewBatchResult()
	orderBook.Sort()
	priceDirection := GetPriceDirection(currentPrice, orderBook)

	if priceDirection == Stay {
		fmt.Println("priceDirection: stay")
		return CalculateMatchStay(currentPrice, orderBook)
	} else { // Increase, Decrease
		if priceDirection == Decrease {
			orderBook.Reverse()
			fmt.Println("priceDirection: decrease")
		} else {
			fmt.Println("priceDirection: increase")
		}
		return CalculateMatch(priceDirection, X, Y, currentPrice, orderBook)
	}
}

func CheckValidityOrderBook(orderBook OrderBook, currentPrice sdk.Dec) bool {
	orderBook.Reverse()
	maxBuyOrderPrice := sdk.ZeroDec()
	minSellOrderPrice := sdk.NewDec(1000000000000) // TODO: fix naive logic
	for _, order := range orderBook {
		if order.BuyOrderAmt.IsPositive() && order.OrderPrice.GT(maxBuyOrderPrice) {
			maxBuyOrderPrice = order.OrderPrice
		}
		if order.SellOrderAmt.IsPositive() && (order.OrderPrice.LT(minSellOrderPrice)) {
			minSellOrderPrice = order.OrderPrice
		}
		fmt.Println(order)
	}
	//if minSellOrderPrice.Equal(sdk.NewDec(1000000000000)){
	//	minSellOrderPrice = sdk.ZeroDec()
	//}
	// TODO: fix naive error rate
	oneOverWithErr, _ := sdk.NewDecFromStr("1.001")
	oneUnderWithErr, _ := sdk.NewDecFromStr("0.999")
	fmt.Println(maxBuyOrderPrice.GT(minSellOrderPrice), maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr), minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr))
	fmt.Println(maxBuyOrderPrice, minSellOrderPrice, currentPrice)
	fmt.Println(maxBuyOrderPrice.Quo(currentPrice), minSellOrderPrice.Quo(currentPrice))
	if maxBuyOrderPrice.GT(minSellOrderPrice) || maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr) || minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr) {
		//fmt.Println(maxBuyOrderPrice.TruncateInt().Quo(currentPrice.TruncateInt()), minSellOrderPrice.TruncateInt().Quo(currentPrice.TruncateInt()))
		return false
	} else {
		return true
	}
}

func ClearOrders(XtoY, YtoX []BatchPoolSwapMsg) ([]BatchPoolSwapMsg, []BatchPoolSwapMsg) {
	// TODO: add clear logic for orderCancelHeight
	newI := 0
	for _, order := range XtoY {
		if !order.Msg.OfferCoin.Amount.IsZero() {
			XtoY[newI] = order
			newI += 1
		}
	}
	XtoY = XtoY[:newI]

	newI = 0
	for _, order := range YtoX {
		if !order.Msg.OfferCoin.Amount.IsZero() {
			YtoX[newI] = order
			newI += 1
		}
	}
	YtoX = YtoX[:newI]
	return XtoY, YtoX
}

func CalculateMatchStay(currentPrice sdk.Dec, orderBook OrderBook) (r BatchResult) {
	r = NewBatchResult()
	r.SwapPrice = currentPrice
	r.OriginalEX, r.OriginalEY = GetExecutableAmt(r.SwapPrice, orderBook)
	r.EX = r.OriginalEX.ToDec()
	r.EY = r.OriginalEY.ToDec()

	if r.EX.Add(r.PoolX).Equal(sdk.ZeroDec()) || r.EY.Add(r.PoolY).Equal(sdk.ZeroDec()) {
		r.MatchType = NoMatch
	} else if r.EX.Equal(r.SwapPrice.Mul(r.EY)) {
		r.MatchType = ExactMatch
	} else {
		r.MatchType = FractionalMatch
		if r.EX.GT(r.SwapPrice.Mul(r.EY)) {
			r.EX = r.SwapPrice.Mul(r.EY)
		} else if r.EX.LT(r.SwapPrice.Mul(r.EY)) {
			r.EY = r.EX.Quo(r.SwapPrice)
		}
	}
	return
}

// TODO: need to debugging
func FindOrderMatch(direction int, swapList []BatchPoolSwapMsg, executableAmt, swapPrice, swapFeeRate sdk.Dec, height int64) (matchResultList []MatchResult, swapListExecuted []BatchPoolSwapMsg, poolXdelta, poolYdelta sdk.Int) {
	swapFeeRate = sdk.ZeroDec() // TODO: temporary zero for simulation
	fmt.Println("FindOrderMatch", direction, executableAmt, swapPrice, swapFeeRate, height, swapList)

	poolXdelta = sdk.ZeroInt()
	poolYdelta = sdk.ZeroInt()

	if direction == DirectionXtoY {
		sort.SliceStable(swapList, func(i, j int) bool {
			return swapList[i].Msg.OrderPrice.GT(swapList[j].Msg.OrderPrice)
		})
	} else if direction == DirectionYtoX {
		sort.SliceStable(swapList, func(i, j int) bool {
			return swapList[i].Msg.OrderPrice.LT(swapList[j].Msg.OrderPrice)
		})
	}

	matchAmt := sdk.ZeroInt()
	accumMatchAmt := sdk.ZeroInt()
	var matchedOrderMsgIndexList []uint64
	var matchOrderList []BatchPoolSwapMsg
	matchedIndexMap := make(map[uint64]sdk.Coin)

	lenSwapList := len(swapList)
	for i, order := range swapList {
		var breakFlag, appendFlag bool

		// include the order in matchAmt, matchOrderList
		if (direction == DirectionXtoY && order.Msg.OrderPrice.GTE(swapPrice)) || // TODO: GTE nil pointer error, swapPrice nil
			(direction == DirectionYtoX && order.Msg.OrderPrice.LTE(swapPrice)) {
			matchAmt = matchAmt.Add(order.Msg.OfferCoin.Amount)
			matchOrderList = append(matchOrderList, order)
		}

		// case check
		if lenSwapList > i+1 { // check next order exist
			if swapList[i+1].Msg.OrderPrice == order.Msg.OrderPrice { // check next orderPrice is same
				breakFlag = false
				appendFlag = false
			} else { // next orderPrice is new
				appendFlag = true
				if (direction == DirectionXtoY && swapList[i+1].Msg.OrderPrice.GTE(swapPrice)) ||
					(direction == DirectionYtoX && swapList[i+1].Msg.OrderPrice.LTE(swapPrice)) { // check next price is matchable
					breakFlag = false
				} else { // next orderPrice is unmatchable
					breakFlag = true
				}
			}
		} else { // next order does not exist
			breakFlag = true
			appendFlag = true
		}

		var fractionalMatchRatio sdk.Dec
		if appendFlag {
			if matchAmt.IsPositive() {
				if accumMatchAmt.Add(matchAmt).ToDec().GTE(executableAmt) {
					fractionalMatchRatio = executableAmt.Sub(accumMatchAmt.ToDec()).Quo(matchAmt.ToDec())
				} else {
					fractionalMatchRatio = sdk.OneDec()
				}
				if fractionalMatchRatio.IsPositive() {
					for _, matchOrder := range matchOrderList {
						orderAmt := matchOrder.Msg.OfferCoin.Amount.ToDec()
						matchResult := MatchResult{
							OrderHeight:       height,
							OrderCancelHeight: height + OrderLifeSpanHeight,
							OrderMsgIndex:     matchOrder.MsgIndex,
							OrderPrice:        matchOrder.Msg.OrderPrice,
							OrderAmt:          matchOrder.Msg.OfferCoin.Amount,
							MatchedAmt:        orderAmt.Mul(fractionalMatchRatio).TruncateInt(),
							RefundAmt:         orderAmt.Mul(sdk.OneDec().Sub(fractionalMatchRatio)).TruncateInt(),
							ReceiveAmt:        orderAmt.Mul(fractionalMatchRatio).Quo(swapPrice).TruncateInt(),
							FeeAmt:            orderAmt.Mul(fractionalMatchRatio).Quo(swapPrice).Mul(swapFeeRate).TruncateInt(),
						}
						matchResult.ResidualAmt = matchResult.OrderAmt.Sub(matchResult.MatchedAmt).Sub(matchResult.RefundAmt)

						if matchOrder.Msg.OfferCoin.Amount.Sub(matchResult.MatchedAmt).LT(sdk.OneInt()) {
							// full match
							matchedOrderMsgIndexList = append(matchedOrderMsgIndexList, matchOrder.MsgIndex)
						} else {
							// fractional match
							matchOrder.Msg.OfferCoin = matchOrder.Msg.OfferCoin.Sub(sdk.NewCoin(matchOrder.Msg.OfferCoin.Denom, matchResult.MatchedAmt))
							matchedIndexMap[matchOrder.MsgIndex] = matchOrder.Msg.OfferCoin
						}

						matchResultList = append(matchResultList, matchResult)
						if direction == DirectionXtoY {
							poolXdelta = poolXdelta.Add(matchResult.MatchedAmt)
							poolYdelta = poolYdelta.Sub(matchResult.ReceiveAmt)
						} else if direction == DirectionYtoX {
							poolXdelta = poolXdelta.Sub(matchResult.ReceiveAmt)
							poolYdelta = poolYdelta.Add(matchResult.MatchedAmt)
						}
					}
				}
			}
			// update accumMatchAmt and initiate matchAmt and matchOrderList
			accumMatchAmt = accumMatchAmt.Add(matchAmt)
			matchAmt = sdk.ZeroInt()
			matchOrderList = matchOrderList[:0]
		}

		if breakFlag {
			break
		}
	}
	if len(matchedOrderMsgIndexList) > 0 {
		newI := 0
		for _, order := range swapList {
			if val, ok := matchedIndexMap[order.MsgIndex]; ok {
				order.Msg.OfferCoin = val
			}
			removeFlag := false
			for _, i := range matchedOrderMsgIndexList {
				if i == order.MsgIndex {
					removeFlag = true
					break
				}
			}
			if !removeFlag {
				swapList[newI] = order
				newI += 1
			}

		}
		swapListExecuted = swapList[:newI]
	}
	return
}

// TODO: find and fix decimal errors
func CalculateSwap(direction int, X, Y, orderPrice, lastOrderPrice sdk.Dec, orderBook OrderBook) (r BatchResult) {
	r = NewBatchResult()
	r.OriginalEX, r.OriginalEY = GetExecutableAmt(lastOrderPrice.Add(orderPrice).Quo(sdk.NewDec(2)), orderBook)
	r.EX = r.OriginalEX.ToDec()
	r.EY = r.OriginalEY.ToDec()

	//r.SwapPrice = X.Add(r.EX).Quo(Y.Add(r.EY)) // legacy constant product model
	r.SwapPrice = X.Add(r.EX.MulInt64(2)).Quo(Y.Add(r.EY.MulInt64(2))) // newSwapPriceModel

	if direction == Increase {
		//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))  // legacy constant product model
		r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)) // newSwapPriceModel
		if lastOrderPrice.LT(r.SwapPrice) && r.SwapPrice.LT(orderPrice) && !r.PoolY.IsNegative() {
			if r.EX.IsZero() && r.EY.IsZero() {
				r.MatchType = NoMatch
			} else {
				r.MatchType = ExactMatch
			}
		}
	} else if direction == Decrease {
		//r.PoolX = X.Sub(Y.Mul(r.SwapPrice))   // legacy constant product model
		r.PoolX = X.Sub(r.SwapPrice.Mul(Y)).QuoInt64(2) // newSwapPriceModel
		if orderPrice.LT(r.SwapPrice) && r.SwapPrice.LT(lastOrderPrice) && !r.PoolX.IsNegative() {
			if r.EX.IsZero() && r.EY.IsZero() {
				r.MatchType = NoMatch
			} else {
				r.MatchType = ExactMatch
			}
		}
	}

	if r.MatchType == 0 {
		r.OriginalEX, r.OriginalEY = GetExecutableAmt(lastOrderPrice.Add(orderPrice).Quo(sdk.NewDec(2)), orderBook)
		r.EX = r.OriginalEX.ToDec()
		r.EY = r.OriginalEY.ToDec()
		r.SwapPrice = orderPrice
		if direction == Increase {
			//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))  // legacy constant product model
			r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)) // newSwapPriceModel
			r.EX = MinDec(r.EX, r.EY.Add(r.PoolY).Mul(r.SwapPrice))
			r.EY = MaxDec(MinDec(r.EY, r.EX.Quo(r.SwapPrice).Sub(r.PoolY)), sdk.ZeroDec())
		} else if direction == Decrease {
			//r.PoolX = X.Sub(Y.Mul(r.SwapPrice)) // legacy constant product model
			r.PoolX = X.Sub(r.SwapPrice.Mul(Y)).QuoInt64(2) // newSwapPriceModel
			r.EX = MinDec(r.EY, r.EX.Add(r.PoolX).Quo(r.SwapPrice))
			r.EY = MaxDec(MinDec(r.EX, r.EY.Mul(r.SwapPrice).Sub(r.PoolX)), sdk.ZeroDec())
		}
		r.MatchType = FractionalMatch
	}

	if direction == Increase {
		//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))
		r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)) // newSwapPriceModel
		if r.SwapPrice.LT(X.Quo(Y)) || r.PoolY.IsNegative() {
			r.TransactAmt = sdk.ZeroDec()
		} else {
			r.TransactAmt = MinDec(r.EX, r.EY.Add(r.PoolY).Mul(r.SwapPrice))
		}
	} else if direction == Decrease {
		//r.PoolX = X.Sub(Y.Mul(r.SwapPrice))
		if r.SwapPrice.GT(X.Quo(Y)) || r.PoolX.IsNegative() {
			r.TransactAmt = sdk.ZeroDec()
		} else {
			r.TransactAmt = MinDec(r.EY, r.EX.Add(r.PoolX).Quo(r.SwapPrice))
		}
	}

	return
}

func CalculateMatch(direction int, X, Y, currentPrice sdk.Dec, orderBook OrderBook) (result BatchResult) {
	result = NewBatchResult()
	lastOrderPrice := currentPrice
	var matchScenarioList []BatchResult
	for _, order := range orderBook {
		if (direction == Increase && order.OrderPrice.LT(currentPrice)) ||
			(direction == Decrease && order.OrderPrice.GT(currentPrice)) {
			continue
		} else {
			orderPrice := order.OrderPrice
			r := CalculateSwap(direction, X, Y, orderPrice, lastOrderPrice, orderBook)
			matchScenarioList = append(matchScenarioList, r)
			lastOrderPrice = orderPrice
		}
	}
	var maxScenario BatchResult
	maxScenario.TransactAmt = sdk.ZeroDec()
	for _, s := range matchScenarioList {
		if s.MatchType == ExactMatch && s.TransactAmt.TruncateInt().IsPositive() {
			maxScenario = s
			break
		} else if s.TransactAmt.TruncateInt().GT(maxScenario.TransactAmt.TruncateInt()) {
			maxScenario = s
		}
	}
	return maxScenario
}

func GetOrderMap(swapMsgs []BatchPoolSwapMsg, denomX, denomY string) (OrderMap, []BatchPoolSwapMsg, []BatchPoolSwapMsg) {
	orderMap := make(OrderMap)
	var XtoY []BatchPoolSwapMsg // buying Y from X
	var YtoX []BatchPoolSwapMsg // selling Y for X
	for _, m := range swapMsgs {
		if m.Msg.OfferCoin.Denom == denomX { // buying Y from X
			XtoY = append(XtoY, m)
			if _, ok := orderMap[m.Msg.OrderPrice]; ok {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice].BuyOrderAmt.Add(m.Msg.OfferCoin.Amount),
					orderMap[m.Msg.OrderPrice].SellOrderAmt}
			} else {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{m.Msg.OrderPrice, m.Msg.OfferCoin.Amount, sdk.ZeroInt()}
			}
		} else if m.Msg.OfferCoin.Denom == denomY { // selling Y for X
			YtoX = append(YtoX, m)
			if _, ok := orderMap[m.Msg.OrderPrice]; ok {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice].BuyOrderAmt,
					orderMap[m.Msg.OrderPrice].SellOrderAmt.Add(m.Msg.OfferCoin.Amount)}
			} else {
				orderMap[m.Msg.OrderPrice] = OrderByPrice{m.Msg.OrderPrice, sdk.ZeroInt(), m.Msg.OfferCoin.Amount}
			}
		} else {
			//return sdk.ErrInvalidDenom
		}
	}
	return orderMap, XtoY, YtoX
}

func GetPriceDirection(currentPrice sdk.Dec, orderBook OrderBook) int {
	buyAmtOverCurrentPrice := sdk.ZeroDec()
	buyAmtAtCurrentPrice := sdk.ZeroDec()
	sellAmtUnderCurrentPrice := sdk.ZeroDec()
	sellAmtAtCurrentPrice := sdk.ZeroDec()

	for _, order := range orderBook {
		if order.OrderPrice.GT(currentPrice) {
			buyAmtOverCurrentPrice = buyAmtOverCurrentPrice.Add(order.BuyOrderAmt.ToDec())
		} else if order.OrderPrice.Equal(currentPrice) {
			buyAmtAtCurrentPrice = buyAmtAtCurrentPrice.Add(order.BuyOrderAmt.ToDec())
			sellAmtAtCurrentPrice = sellAmtAtCurrentPrice.Add(order.SellOrderAmt.ToDec())
		} else if order.OrderPrice.LT(currentPrice) {
			sellAmtUnderCurrentPrice = sellAmtUnderCurrentPrice.Add(order.SellOrderAmt.ToDec())
		}
	}
	// TODO: verify Dec, Int math logic
	if buyAmtOverCurrentPrice.Sub(currentPrice.Mul(sellAmtUnderCurrentPrice.Add(sellAmtAtCurrentPrice))).IsPositive() {
		return Increase
	} else if currentPrice.Mul(sellAmtUnderCurrentPrice).Sub(buyAmtOverCurrentPrice.Add(buyAmtAtCurrentPrice)).IsPositive() {
		return Decrease
	} else {
		return Stay
	}
}

func GetExecutableAmt(swapPrice sdk.Dec, orderBook OrderBook) (executableBuyAmtX, executableSellAmtY sdk.Int) {
	executableBuyAmtX = sdk.ZeroInt()
	executableSellAmtY = sdk.ZeroInt()
	for _, order := range orderBook {
		if order.OrderPrice.GTE(swapPrice) {
			executableBuyAmtX = executableBuyAmtX.Add(order.BuyOrderAmt)
		}
		if order.OrderPrice.LTE(swapPrice) {
			executableSellAmtY = executableSellAmtY.Add(order.SellOrderAmt)
		}
	}
	return
}
