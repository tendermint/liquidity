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

func MinInt(a, b sdk.Int) sdk.Int {
	if a.LTE(b) {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b sdk.Int) sdk.Int {
	if a.GTE(b) {
		return a
	} else {
		return b
	}
}

type OrderMap map[string]OrderByPrice

// make orderbook by sort orderMap, increased
func (orderMap OrderMap) SortOrderBook() (orderBook OrderBook) {
	orderPriceList := make([]sdk.Dec, 0, len(orderMap))
	for _, v := range orderMap {
		orderPriceList = append(orderPriceList, v.OrderPrice)
	}

	sort.Slice(orderPriceList, func(i, j int) bool {
		return orderPriceList[i].LT(orderPriceList[j])
	})

	for _, k := range orderPriceList {
		orderBook = append(orderBook, OrderByPrice{
			OrderPrice:   k,
			BuyOrderAmt:  orderMap[k.String()].BuyOrderAmt,
			SellOrderAmt: orderMap[k.String()].SellOrderAmt,
		})
	}
	return orderBook
}

type BatchResult struct {
	MatchType   int
	SwapPrice   sdk.Dec
	EX          sdk.Int
	EY          sdk.Int
	OriginalEX  sdk.Int
	OriginalEY  sdk.Int
	PoolX       sdk.Int
	PoolY       sdk.Int
	TransactAmt sdk.Int
}

func NewBatchResult() BatchResult {
	return BatchResult{
		SwapPrice:   sdk.ZeroDec(),
		EX:          sdk.ZeroInt(),
		EY:          sdk.ZeroInt(),
		OriginalEX:  sdk.ZeroInt(),
		OriginalEY:  sdk.ZeroInt(),
		PoolX:       sdk.ZeroInt(),
		PoolY:       sdk.ZeroInt(),
		TransactAmt: sdk.ZeroInt(),
	}
}

type MatchResult struct {
	OrderHeight       int64
	OrderCancelHeight int64
	OrderMsgIndex     uint64
	OrderPrice        sdk.Dec
	OrderAmt          sdk.Int
	MatchedAmt        sdk.Int
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

// check orderbook validity
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

	// TODO: fix naive error rate
	oneOverWithErr, _ := sdk.NewDecFromStr("1.001")
	oneUnderWithErr, _ := sdk.NewDecFromStr("0.999")
	if maxBuyOrderPrice.GT(minSellOrderPrice) ||
		maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr) ||
		minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr) {
		fmt.Println(maxBuyOrderPrice.GT(minSellOrderPrice), maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr),
			minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr))
		fmt.Println(maxBuyOrderPrice, minSellOrderPrice, currentPrice)
		fmt.Println(maxBuyOrderPrice.Quo(currentPrice), minSellOrderPrice.Quo(currentPrice))
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
	r.EX = r.OriginalEX
	r.EY = r.OriginalEY

	if r.EX.Add(r.PoolX).Equal(sdk.ZeroInt()) || r.EY.Add(r.PoolY).Equal(sdk.ZeroInt()) {
		r.MatchType = NoMatch
	} else if r.EX.Equal(r.SwapPrice.MulInt(r.EY).TruncateInt()) {
		r.MatchType = ExactMatch
	} else {
		r.MatchType = FractionalMatch
		if r.EX.GT(r.SwapPrice.MulInt(r.EY).TruncateInt()) {
			r.EX = r.SwapPrice.MulInt(r.EY).TruncateInt()
		} else if r.EX.LT(r.SwapPrice.MulInt(r.EY).TruncateInt()) {
			r.EY = r.EX.ToDec().Quo(r.SwapPrice).TruncateInt()
		}
	}
	if !r.EX.Add(r.PoolX).Equal(r.EY.Add(r.PoolY)) {
		fmt.Println("!! CalculateMatchStay", r, r.EX.Add(r.PoolX), r.EY.Add(r.PoolY))
	}
	return
}

func UpdateState(X, Y sdk.Dec, XtoY, YtoX []BatchPoolSwapMsg, matchResultXtoY, matchResultYtoX []MatchResult) (
	[]BatchPoolSwapMsg, []BatchPoolSwapMsg, sdk.Dec, sdk.Dec, sdk.Int, sdk.Int, int, int, sdk.Int, sdk.Int) {
	sort.SliceStable(XtoY, func(i, j int) bool {
		return XtoY[i].Msg.OrderPrice.GT(XtoY[j].Msg.OrderPrice)
	})
	sort.SliceStable(YtoX, func(i, j int) bool {
		return YtoX[i].Msg.OrderPrice.LT(YtoX[j].Msg.OrderPrice)
	})

	poolXdelta := sdk.ZeroInt()
	poolYdelta := sdk.ZeroInt()
	var matchedOrderMsgIndexListXtoY []uint64
	var matchedOrderMsgIndexListYtoX []uint64
	matchedIndexMapXtoY := make(map[uint64]sdk.Coin)
	matchedIndexMapYtoX := make(map[uint64]sdk.Coin)
	fractionalCntX := 0
	fractionalCntY := 0
	decimalErrorX := sdk.ZeroInt()
	decimalErrorY := sdk.ZeroInt()

	for _, match := range matchResultXtoY {
		for _, order := range XtoY {
			if match.OrderMsgIndex == order.MsgIndex {
				poolXdelta = poolXdelta.Add(match.MatchedAmt)
				poolYdelta = poolYdelta.Sub(match.ReceiveAmt)
				if order.Msg.OfferCoin.Amount.Equal(match.MatchedAmt) {
					// full match
					matchedOrderMsgIndexListXtoY = append(matchedOrderMsgIndexListXtoY, order.MsgIndex)
				} else if order.Msg.OfferCoin.Amount.Sub(match.MatchedAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
					decimalErrorX = decimalErrorX.Add(sdk.OneInt())
					//poolXdelta = poolXdelta.Add(sdk.OneInt())
					matchedOrderMsgIndexListXtoY = append(matchedOrderMsgIndexListXtoY, order.MsgIndex)
				} else {
					// fractional match
					order.Msg.OfferCoin = order.Msg.OfferCoin.Sub(sdk.NewCoin(order.Msg.OfferCoin.Denom, match.MatchedAmt))
					matchedIndexMapXtoY[order.MsgIndex] = order.Msg.OfferCoin
					fractionalCntX += 1
				}
				break
			}
		}
	}
	if len(matchedOrderMsgIndexListXtoY) > 0 {
		newI := 0
		for _, order := range XtoY {
			if val, ok := matchedIndexMapXtoY[order.MsgIndex]; ok {
				order.Msg.OfferCoin = val
			}
			removeFlag := false
			for _, i := range matchedOrderMsgIndexListXtoY {
				if i == order.MsgIndex {
					removeFlag = true
					break
				}
			}
			if !removeFlag {
				XtoY[newI] = order
				newI += 1
			}
			removeFlag = false

		}
		XtoY = XtoY[:newI]
	}
	for _, match := range matchResultYtoX {
		for _, order := range YtoX {
			if match.OrderMsgIndex == order.MsgIndex {
				poolXdelta = poolXdelta.Sub(match.ReceiveAmt)
				poolYdelta = poolYdelta.Add(match.MatchedAmt)
				if order.Msg.OfferCoin.Amount.Equal(match.MatchedAmt) {
					// full match
					matchedOrderMsgIndexListYtoX = append(matchedOrderMsgIndexListYtoX, order.MsgIndex)
				} else if order.Msg.OfferCoin.Amount.Sub(match.MatchedAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
					decimalErrorY = decimalErrorY.Add(sdk.OneInt())
					//poolYdelta = poolYdelta.Add(sdk.OneInt())
					matchedOrderMsgIndexListYtoX = append(matchedOrderMsgIndexListYtoX, order.MsgIndex)
				} else {
					// fractional match
					order.Msg.OfferCoin = order.Msg.OfferCoin.Sub(sdk.NewCoin(order.Msg.OfferCoin.Denom, match.MatchedAmt))
					matchedIndexMapYtoX[order.MsgIndex] = order.Msg.OfferCoin
					fractionalCntY += 1
				}
				break
			}
		}
	}
	if len(matchedOrderMsgIndexListYtoX) > 0 {
		newI := 0
		for _, order := range YtoX {
			if val, ok := matchedIndexMapYtoX[order.MsgIndex]; ok {
				order.Msg.OfferCoin = val
			}
			removeFlag := false
			for _, i := range matchedOrderMsgIndexListYtoX {
				if i == order.MsgIndex {
					removeFlag = true
					break
				}
			}
			if !removeFlag {
				YtoX[newI] = order
				newI += 1
			}
			removeFlag = false

		}
		YtoX = YtoX[:newI]
	}

	poolXdelta = poolXdelta.Add(decimalErrorX)
	poolYdelta = poolYdelta.Add(decimalErrorY)

	X = X.Add(poolXdelta.ToDec())
	Y = Y.Add(poolYdelta.ToDec())

	return XtoY, YtoX, X, Y, poolXdelta, poolYdelta, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY
}

func FindOrderMatch(direction int, swapList []BatchPoolSwapMsg, executableAmt sdk.Int,
	swapPrice, swapFeeRate sdk.Dec, height int64) (
	matchResultList []MatchResult, swapListExecuted []BatchPoolSwapMsg, poolXdelta, poolYdelta sdk.Int) {

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
	//var matchedOrderMsgIndexList []uint64
	var matchOrderList []BatchPoolSwapMsg
	//matchedIndexMap := make(map[uint64]sdk.Coin)

	fmt.Println("executableAmt", executableAmt)
	if executableAmt.IsZero() {
		return
	}

	lenSwapList := len(swapList)
	for i, order := range swapList {
		var breakFlag, appendFlag bool

		// include the order in matchAmt, matchOrderList
		if (direction == DirectionXtoY && order.Msg.OrderPrice.GTE(swapPrice)) ||
			(direction == DirectionYtoX && order.Msg.OrderPrice.LTE(swapPrice)) {
			matchAmt = matchAmt.Add(order.Msg.OfferCoin.Amount)
			matchOrderList = append(matchOrderList, order)
		}

		// case check
		if lenSwapList > i+1 { // check next order exist
			if swapList[i+1].Msg.OrderPrice.Equal(order.Msg.OrderPrice) { // check next orderPrice is same
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
				if accumMatchAmt.Add(matchAmt).GTE(executableAmt) {
					fractionalMatchRatio = executableAmt.Sub(accumMatchAmt).ToDec().Quo(matchAmt.ToDec())
				} else {
					fractionalMatchRatio = sdk.OneDec()
				}
				if fractionalMatchRatio.GT(sdk.OneDec()) {
					fmt.Println("!!! fractionalMatchRatio.GT(sdk.OneDec())", fractionalMatchRatio,
						executableAmt.Sub(accumMatchAmt), matchAmt)
				}
				for _, matchOrder := range matchOrderList {
					if fractionalMatchRatio.IsPositive() {
						orderAmt := matchOrder.Msg.OfferCoin.Amount.ToDec()
						matchResult := MatchResult{
							OrderHeight:       height,
							OrderCancelHeight: height + OrderLifeSpanHeight,
							OrderMsgIndex:     matchOrder.MsgIndex,
							OrderPrice:        matchOrder.Msg.OrderPrice,
							OrderAmt:          matchOrder.Msg.OfferCoin.Amount,
							MatchedAmt:        orderAmt.Mul(fractionalMatchRatio).Ceil().TruncateInt(),
						}
						if direction == DirectionXtoY {
							matchResult.ReceiveAmt = orderAmt.Mul(fractionalMatchRatio).Quo(swapPrice).TruncateInt()
							matchResult.FeeAmt = orderAmt.Mul(fractionalMatchRatio).Quo(swapPrice).Mul(swapFeeRate).TruncateInt()
						} else if direction == DirectionYtoX {
							matchResult.ReceiveAmt = orderAmt.Mul(fractionalMatchRatio).Mul(swapPrice).TruncateInt()
							matchResult.FeeAmt = orderAmt.Mul(fractionalMatchRatio).Mul(swapPrice).Mul(swapFeeRate).TruncateInt()
						}
						// TODO: need to verify logic
						if matchResult.MatchedAmt.GT(matchResult.OrderAmt) ||
							(matchResult.FeeAmt.GT(matchResult.ReceiveAmt) && matchResult.FeeAmt.GT(sdk.OneInt())) {
							fmt.Println("panic(matchResult.MatchedAmt)", matchResult,
								orderAmt.Mul(fractionalMatchRatio).Mul(swapPrice).Mul(swapFeeRate))
							panic(matchResult.MatchedAmt)
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
	return
}

// Calculates the batch results with the processing logic for each direction
func CalculateSwap(direction int, X, Y, orderPrice, lastOrderPrice sdk.Dec, orderBook OrderBook) (r BatchResult) {
	r = NewBatchResult()
	r.OriginalEX, r.OriginalEY = GetExecutableAmt(lastOrderPrice.Add(orderPrice).Quo(sdk.NewDec(2)), orderBook)
	r.EX = r.OriginalEX
	r.EY = r.OriginalEY

	//r.SwapPrice = X.Add(r.EX).Quo(Y.Add(r.EY)) // legacy constant product model
	r.SwapPrice = X.Add(r.EX.MulRaw(2).ToDec()).Quo(Y.Add(r.EY.MulRaw(2).ToDec())) // newSwapPriceModel

	if direction == Increase {
		//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))  // legacy constant product model
		r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)).TruncateInt() // newSwapPriceModel
		if lastOrderPrice.LT(r.SwapPrice) && r.SwapPrice.LT(orderPrice) && !r.PoolY.IsNegative() {
			if r.EX.IsZero() && r.EY.IsZero() {
				r.MatchType = NoMatch
			} else {
				r.MatchType = ExactMatch
			}
		}
	} else if direction == Decrease {
		//r.PoolX = X.Sub(Y.Mul(r.SwapPrice))   // legacy constant product model
		r.PoolX = X.Sub(r.SwapPrice.Mul(Y)).QuoInt64(2).TruncateInt() // newSwapPriceModel
		if orderPrice.LT(r.SwapPrice) && r.SwapPrice.LT(lastOrderPrice) && r.PoolX.GTE(sdk.ZeroInt()) {
			if r.EX.IsZero() && r.EY.IsZero() {
				r.MatchType = NoMatch
			} else {
				r.MatchType = ExactMatch
			}
		}
	}

	if r.MatchType == 0 {
		r.OriginalEX, r.OriginalEY = GetExecutableAmt(orderPrice, orderBook)
		r.EX = r.OriginalEX
		r.EY = r.OriginalEY
		r.SwapPrice = orderPrice
		if direction == Increase {
			//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))  // legacy constant product model
			r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)).TruncateInt() // newSwapPriceModel
			r.EX = MinDec(r.EX.ToDec(), r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice)).Ceil().TruncateInt()
			r.EY = MaxDec(MinDec(r.EY.ToDec(), r.EX.ToDec().Quo(r.SwapPrice).Sub(r.PoolY.ToDec())), sdk.ZeroDec()).Ceil().TruncateInt()
		} else if direction == Decrease {
			//r.PoolX = X.Sub(Y.Mul(r.SwapPrice)) // legacy constant product model
			r.PoolX = X.Sub(r.SwapPrice.Mul(Y)).QuoInt64(2).TruncateInt() // newSwapPriceModel
			r.EY = MinDec(r.EY.ToDec(), r.EX.Add(r.PoolX).ToDec().Quo(r.SwapPrice)).Ceil().TruncateInt()
			r.EX = MaxDec(MinDec(r.EX.ToDec(), r.EY.ToDec().Mul(r.SwapPrice).Sub(r.PoolX.ToDec())), sdk.ZeroDec()).Ceil().TruncateInt()
		}
		r.MatchType = FractionalMatch
	}

	if direction == Increase {
		if r.SwapPrice.LT(X.Quo(Y)) || r.PoolY.IsNegative() {
			r.TransactAmt = sdk.ZeroInt()
		} else {
			r.TransactAmt = MinInt(r.EX, r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice).RoundInt())
		}
	} else if direction == Decrease {
		if r.SwapPrice.GT(X.Quo(Y)) || r.PoolX.LT(sdk.ZeroInt()) {
			r.TransactAmt = sdk.ZeroInt()
		} else {
			r.TransactAmt = MinInt(r.EY, r.EX.Add(r.PoolX).ToDec().Quo(r.SwapPrice).RoundInt())
		}
	}

	//if r.EX.Add(r.PoolX).ToDec().Sub(r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice)).GT(sdk.OneDec()) {
	//	fmt.Println("!! CalculateSwap invariant check fail", r, r.EX.Add(r.PoolX).ToDec(), r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice))
	//}
	return
}

// Calculates the batch results with the logic for each direction
func CalculateMatch(direction int, X, Y, currentPrice sdk.Dec, orderBook OrderBook) (maxScenario BatchResult) {
	lastOrderPrice := currentPrice
	var matchScenarioList []BatchResult
	for _, order := range orderBook {
		if (direction == Increase && order.OrderPrice.LT(currentPrice)) ||
			(direction == Decrease && order.OrderPrice.GT(currentPrice)) {
			continue
		} else {
			orderPrice := order.OrderPrice
			r := CalculateSwap(direction, X, Y, orderPrice, lastOrderPrice, orderBook)
			// TODO: need to re-check on v2
			if (direction == Increase && r.PoolY.ToDec().Sub(r.EX.ToDec().Quo(r.SwapPrice)).GTE(sdk.OneDec())) ||
				(direction == Decrease && r.PoolX.ToDec().Sub(r.EY.ToDec().Mul(r.SwapPrice)).GTE(sdk.OneDec())) {
				//fmt.Println("!! CalculateMatch, cant cover case", r)
				continue
			}
			matchScenarioList = append(matchScenarioList, r)
			lastOrderPrice = orderPrice
		}
	}
	maxScenario = NewBatchResult()
	maxScenario.TransactAmt = sdk.ZeroInt()
	for _, s := range matchScenarioList {
		if s.MatchType == ExactMatch && s.TransactAmt.IsPositive() {
			maxScenario = s
			break
		} else if s.TransactAmt.GT(maxScenario.TransactAmt) {
			maxScenario = s
		}
	}
	r := maxScenario

	// TODO: verify logic
	tmpInvariant := r.EX.Add(r.PoolX).ToDec().Sub(r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice))
	if tmpInvariant.GT(r.SwapPrice) && tmpInvariant.GT(sdk.OneDec()) {
		fmt.Println("!! maxScenario CalculateSwap ", r, r.EX.Add(r.PoolX).ToDec(),
			r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice),
			r.EX.Add(r.PoolX).ToDec().Sub(r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice)).Quo(r.SwapPrice), tmpInvariant)
		panic("maxScenario CalculateSwap")
	}
	return maxScenario
}

// make orderMap key as swap price, value as Buy, Sell Amount from swap msgs,  with split as Buy XtoY, Sell YtoX msg list
func GetOrderMap(swapMsgs []BatchPoolSwapMsg, denomX, denomY string) (OrderMap, []BatchPoolSwapMsg, []BatchPoolSwapMsg) {
	orderMap := make(OrderMap)
	var XtoY []BatchPoolSwapMsg // buying Y from X
	var YtoX []BatchPoolSwapMsg // selling Y for X
	for _, m := range swapMsgs {
		if m.Msg.OfferCoin.Denom == denomX { // buying Y from X
			XtoY = append(XtoY, m)
			if _, ok := orderMap[m.Msg.OrderPrice.String()]; ok {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice.String()].BuyOrderAmt.Add(m.Msg.OfferCoin.Amount),
					orderMap[m.Msg.OrderPrice.String()].SellOrderAmt}
			} else {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{m.Msg.OrderPrice,
					m.Msg.OfferCoin.Amount, sdk.ZeroInt()}
			}
		} else if m.Msg.OfferCoin.Denom == denomY { // selling Y for X
			YtoX = append(YtoX, m)
			if _, ok := orderMap[m.Msg.OrderPrice.String()]; ok {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice.String()].BuyOrderAmt,
					orderMap[m.Msg.OrderPrice.String()].SellOrderAmt.Add(m.Msg.OfferCoin.Amount)}
			} else {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{m.Msg.OrderPrice,
					sdk.ZeroInt(), m.Msg.OfferCoin.Amount}
			}
		} else {
			//return sdk.ErrInvalidDenom
		}
	}
	return orderMap, XtoY, YtoX
}

// Get Price direction of the orderbook with current Price
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

	if buyAmtOverCurrentPrice.Sub(currentPrice.Mul(sellAmtUnderCurrentPrice.Add(sellAmtAtCurrentPrice))).IsPositive() {
		return Increase
	} else if currentPrice.Mul(sellAmtUnderCurrentPrice).Sub(buyAmtOverCurrentPrice.Add(buyAmtAtCurrentPrice)).IsPositive() {
		return Decrease
	} else {
		return Stay
	}
}

// calculate the executable amount of the orderbook for each X, Y
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
