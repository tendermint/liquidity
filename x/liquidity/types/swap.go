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

	//OrderLifeSpanHeight = 0

	DirectionXtoY = 1
	DirectionYtoX = 2
)

type OrderByPrice struct {
	OrderPrice   sdk.Dec
	BuyOfferAmt  sdk.Int
	SellOfferAmt sdk.Int
	MsgList []*BatchPoolSwapMsg
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

type MsgList []*BatchPoolSwapMsg

func (msgList MsgList) LenRemainingMsgs() int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && !m.Succeed {
			cnt++
		}
	}
	return cnt
}

func (msgList MsgList) LenFractionalMsgs() int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && m.Succeed && !m.ToDelete {
			cnt++
		}
	}
	return cnt
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
			BuyOfferAmt:  orderMap[k.String()].BuyOfferAmt,
			SellOfferAmt: orderMap[k.String()].SellOfferAmt,
			MsgList: orderMap[k.String()].MsgList,
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
	OrderHeight            int64
	OrderExpiryHeight      int64
	OrderMsgIndex          uint64
	OrderPrice             sdk.Dec
	OfferCoinAmt           sdk.Int
	TransactedCoinAmt      sdk.Int
	ExchangedDemandCoinAmt sdk.Int
	FeeAmt                 sdk.Int
	BatchMsg               *BatchPoolSwapMsg
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
		if order.BuyOfferAmt.IsPositive() && order.OrderPrice.GT(maxBuyOrderPrice) {
			maxBuyOrderPrice = order.OrderPrice
		}
		if order.SellOfferAmt.IsPositive() && (order.OrderPrice.LT(minSellOrderPrice)) {
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

func ClearOrders(XtoY, YtoX []*BatchPoolSwapMsg, currentHeight int64, clearThisHeight bool) ([]*BatchPoolSwapMsg, []*BatchPoolSwapMsg) {
	for _, order := range XtoY {
		if order.RemainingOfferCoin.IsZero() {  // TODO: verify
			order.Succeed = true
			order.ToDelete = true
			// TODO: set exchangedAmt, remainingAmt, without fix msg
		}
		// set toDelete, expired msgs
		if order.OrderExpiryHeight > currentHeight ||
			(clearThisHeight && order.OrderExpiryHeight >= currentHeight ){
			order.Succeed = false
			order.ToDelete = true
		}
	}
	for _, order := range YtoX {
		if order.RemainingOfferCoin.IsZero() {  // TODO: verify
			order.Succeed = true
			order.ToDelete = true
		}
		// set toDelete, expired msgs
		if order.OrderExpiryHeight > currentHeight ||
			(clearThisHeight && order.OrderExpiryHeight >= currentHeight ){
			order.Succeed = false
			order.ToDelete = true
		}
	}
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

// Find matched orders and set status for msgs
func FindOrderMatch(direction int, swapList []*BatchPoolSwapMsg, executableAmt sdk.Int,
	swapPrice, swapFeeRate sdk.Dec, height int64) (
	matchResultList []MatchResult, swapListExecuted []*BatchPoolSwapMsg, poolXdelta, poolYdelta sdk.Int) {

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
	var matchOrderList []*BatchPoolSwapMsg
	//matchedIndexMap := make(map[uint64]sdk.Coin)

	fmt.Println("executableAmt", executableAmt)
	if executableAmt.IsZero() {
		return
	}

	lenSwapList := len(swapList)
	for i, order := range swapList {
		var breakFlag, appendFlag bool

		// include the matched order in matchAmt, matchOrderList
		if (direction == DirectionXtoY && order.Msg.OrderPrice.GTE(swapPrice)) ||
			(direction == DirectionYtoX && order.Msg.OrderPrice.LTE(swapPrice)) {
			matchAmt = matchAmt.Add(order.Msg.OfferCoin.Amount)
			// TODO: set on state update?
			//order.Succeed = true
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
						offerAmt := matchOrder.Msg.OfferCoin.Amount.ToDec()
						matchResult := MatchResult{
							OrderHeight:       height,
							OrderExpiryHeight: height + CancelOrderLifeSpan,
							OrderMsgIndex:     matchOrder.MsgIndex,
							OrderPrice:        matchOrder.Msg.OrderPrice,
							OfferCoinAmt:      matchOrder.Msg.OfferCoin.Amount,
							TransactedCoinAmt: offerAmt.Mul(fractionalMatchRatio).Ceil().TruncateInt(),
							BatchMsg: matchOrder,
						}
						if direction == DirectionXtoY {
							// TODO: verify exchanged
							matchResult.ExchangedDemandCoinAmt = offerAmt.Mul(fractionalMatchRatio).Quo(swapPrice).TruncateInt()
							matchResult.FeeAmt = offerAmt.Mul(fractionalMatchRatio).Quo(swapPrice).Mul(swapFeeRate).TruncateInt()
						} else if direction == DirectionYtoX {
							// TODO: verify exchanged
							matchResult.ExchangedDemandCoinAmt = offerAmt.Mul(fractionalMatchRatio).Mul(swapPrice).TruncateInt()
							matchResult.FeeAmt = offerAmt.Mul(fractionalMatchRatio).Mul(swapPrice).Mul(swapFeeRate).TruncateInt()
						}
						// TODO: need to verify logic
						if matchResult.TransactedCoinAmt.GT(matchResult.OfferCoinAmt) ||
							(matchResult.FeeAmt.GT(matchResult.ExchangedDemandCoinAmt) && matchResult.FeeAmt.GT(sdk.OneInt())) {
							fmt.Println("panic(matchResult.TransactedCoinAmt)", matchResult,
								offerAmt.Mul(fractionalMatchRatio).Mul(swapPrice).Mul(swapFeeRate))
							panic(matchResult.TransactedCoinAmt)
						}
						// set exchangedAmt and remainingAmt on batch pointer
						// TODO: set on state update?
						//matchOrder.ExchangedOfferCoin = matchOrder.ExchangedOfferCoin.Add(
						//	sdk.NewCoin(matchOrder.Msg.OfferCoin.Denom, matchResult.TransactedCoinAmt))
						//// TODO: verify RemainingOfferCoin about deciaml errors
						//matchOrder.RemainingOfferCoin = matchOrder.Msg.OfferCoin.Sub(matchOrder.ExchangedOfferCoin)
						matchResultList = append(matchResultList, matchResult)
						if direction == DirectionXtoY {
							poolXdelta = poolXdelta.Add(matchResult.TransactedCoinAmt)
							poolYdelta = poolYdelta.Sub(matchResult.ExchangedDemandCoinAmt)
						} else if direction == DirectionYtoX {
							poolXdelta = poolXdelta.Sub(matchResult.ExchangedDemandCoinAmt)
							poolYdelta = poolYdelta.Add(matchResult.TransactedCoinAmt)
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
func GetOrderMap(swapMsgs []*BatchPoolSwapMsg, denomX, denomY string) (OrderMap, []*BatchPoolSwapMsg, []*BatchPoolSwapMsg) {
	orderMap := make(OrderMap)
	var XtoY []*BatchPoolSwapMsg // buying Y from X
	var YtoX []*BatchPoolSwapMsg // selling Y for X
	for _, m := range swapMsgs {
		if m.Msg.OfferCoin.Denom == denomX { // buying Y from X
			XtoY = append(XtoY, m)
			if _, ok := orderMap[m.Msg.OrderPrice.String()]; ok {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice.String()].BuyOfferAmt.Add(m.Msg.OfferCoin.Amount),
					orderMap[m.Msg.OrderPrice.String()].SellOfferAmt,
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			} else {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{m.Msg.OrderPrice,
					m.Msg.OfferCoin.Amount, sdk.ZeroInt(),
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			}
		} else if m.Msg.OfferCoin.Denom == denomY { // selling Y for X
			YtoX = append(YtoX, m)
			if _, ok := orderMap[m.Msg.OrderPrice.String()]; ok {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice.String()].BuyOfferAmt,
					orderMap[m.Msg.OrderPrice.String()].SellOfferAmt.Add(m.Msg.OfferCoin.Amount),
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			} else {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{m.Msg.OrderPrice,
					sdk.ZeroInt(), m.Msg.OfferCoin.Amount,
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			}
		} else {
			panic("ErrInvalidDenom")
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
			buyAmtOverCurrentPrice = buyAmtOverCurrentPrice.Add(order.BuyOfferAmt.ToDec())
		} else if order.OrderPrice.Equal(currentPrice) {
			buyAmtAtCurrentPrice = buyAmtAtCurrentPrice.Add(order.BuyOfferAmt.ToDec())
			sellAmtAtCurrentPrice = sellAmtAtCurrentPrice.Add(order.SellOfferAmt.ToDec())
		} else if order.OrderPrice.LT(currentPrice) {
			sellAmtUnderCurrentPrice = sellAmtUnderCurrentPrice.Add(order.SellOfferAmt.ToDec())
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
			executableBuyAmtX = executableBuyAmtX.Add(order.BuyOfferAmt)
		}
		if order.OrderPrice.LTE(swapPrice) {
			executableSellAmtY = executableSellAmtY.Add(order.SellOfferAmt)
		}
	}
	return
}
