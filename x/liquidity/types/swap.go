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
	MsgList      []*BatchPoolSwapMsg
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

func (msgList MsgList) CountNotMatchedMsgs() int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && !m.Succeed {
			cnt++
		}
	}
	return cnt
}

func (msgList MsgList) CountFractionalMatchedMsgs() int {
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

// make orderbook by sort orderMap
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
			MsgList:      orderMap[k.String()].MsgList,
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
	// TODO: add swapPrice
}

// The price and coins of swap messages in orderbook are calculated
// to derive match result with the price direction.
func MatchOrderbook(X, Y, currentPrice sdk.Dec, orderBook OrderBook) (result BatchResult) {
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
	}

	// TODO: fix naive error rate
	oneOverWithErr, _ := sdk.NewDecFromStr("1.002")
	oneUnderWithErr, _ := sdk.NewDecFromStr("0.998")
	if maxBuyOrderPrice.GT(minSellOrderPrice) ||
		maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr) ||
		minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr) {

		fmt.Println(maxBuyOrderPrice.GT(minSellOrderPrice),
			maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr),
			minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr), maxBuyOrderPrice, minSellOrderPrice, currentPrice)
		return false

	} else {
		return true
	}
}

func ValidateStateAndExpireOrders(msgList []*BatchPoolSwapMsg, currentHeight int64, expireThisHeight bool) []*BatchPoolSwapMsg {
	for _, order := range msgList {
		if !order.Executed {
			panic("not executed")
			continue
		}
		if order.RemainingOfferCoin.IsZero() {
			if !order.Succeed || !order.ToDelete {
				panic("broken state consistency for not matched order")
			}
			order.Succeed = true
			order.ToDelete = true
			continue
		}
		// set toDelete, expired msgs
		if currentHeight > order.OrderExpiryHeight {
			if order.Succeed || !order.ToDelete {
				panic("broken state consistency for fractional matched order")
			}
			order.Succeed = false
			order.ToDelete = true
			continue
		}
		if expireThisHeight && currentHeight == order.OrderExpiryHeight {
			order.ToDelete = true
		}
	}
	return msgList
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
							BatchMsg:          matchOrder,
						}
						if matchOrder != matchResult.BatchMsg {
							panic("not matched msg pointer ")
						}
						if direction == DirectionXtoY {
							// TODO: offer-FeeAmt for exchanged
							matchResult.FeeAmt = matchResult.TransactedCoinAmt.ToDec().Mul(swapFeeRate).TruncateInt()
							matchResult.ExchangedDemandCoinAmt = matchResult.TransactedCoinAmt.Sub(matchResult.FeeAmt).ToDec().Quo(swapPrice).TruncateInt()
						} else if direction == DirectionYtoX {
							// TODO: offer-FeeAmt for exchanged
							matchResult.FeeAmt = matchResult.TransactedCoinAmt.ToDec().Mul(swapFeeRate).TruncateInt()
							matchResult.ExchangedDemandCoinAmt = matchResult.TransactedCoinAmt.Sub(matchResult.FeeAmt).ToDec().Mul(swapPrice).TruncateInt()
						}
						if matchResult.TransactedCoinAmt.GT(matchResult.OfferCoinAmt) ||
							(matchResult.FeeAmt.GT(matchResult.OfferCoinAmt) && matchResult.FeeAmt.GT(sdk.OneInt())) {
							panic(matchResult.TransactedCoinAmt)
						}
						matchResultList = append(matchResultList, matchResult)
						swapListExecuted = append(swapListExecuted, matchOrder)
						if direction == DirectionXtoY {
							poolXdelta = poolXdelta.Add(matchResult.TransactedCoinAmt)
							poolYdelta = poolYdelta.Sub(matchResult.ExchangedDemandCoinAmt)
						} else if direction == DirectionYtoX {
							poolXdelta = poolXdelta.Sub(matchResult.ExchangedDemandCoinAmt)
							poolYdelta = poolYdelta.Add(matchResult.TransactedCoinAmt)
						}
					} else {
						fmt.Println("fractional ratio is negative", fractionalMatchRatio)
						fmt.Println(accumMatchAmt, matchAmt, executableAmt)
						// TODO: check stop or pass
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
			// TODO: need to re-check on milestone 2
			if (direction == Increase && r.PoolY.ToDec().Sub(r.EX.ToDec().Quo(r.SwapPrice)).GTE(sdk.OneDec())) ||
				(direction == Decrease && r.PoolX.ToDec().Sub(r.EY.ToDec().Mul(r.SwapPrice)).GTE(sdk.OneDec())) {
				continue
			}
			matchScenarioList = append(matchScenarioList, r)
			lastOrderPrice = orderPrice
		}
	}
	maxScenario = NewBatchResult()
	maxScenario.TransactAmt = sdk.ZeroInt()
	for _, s := range matchScenarioList {
		//if s.MatchType == ExactMatch && s.TransactAmt.IsPositive() {
		//	maxScenario = s
		//	break
		//} else if s.TransactAmt.GT(maxScenario.TransactAmt) {
		//	maxScenario = s
		//}
		MEX, MEY := GetMustExecutableAmt(s.SwapPrice, orderBook)
		fmt.Println("Scenario, MEX, MEY", s, MEX, MEY)
		if s.EX.GTE(MEX) && s.EY.GTE(MEY) {
			if s.MatchType == ExactMatch && s.TransactAmt.IsPositive() {
				maxScenario = s
				break
			} else if s.TransactAmt.GT(maxScenario.TransactAmt) {
				maxScenario = s
			}
		}
	}
	r := maxScenario

	tmpInvariant := r.EX.Add(r.PoolX).ToDec().Sub(r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice))
	if tmpInvariant.GT(r.SwapPrice) && tmpInvariant.GT(sdk.OneDec()) {
		fmt.Println(tmpInvariant.GT(r.SwapPrice), tmpInvariant.GT(sdk.OneDec()))
		fmt.Println(tmpInvariant, r.SwapPrice)
		panic("maxScenario CalculateSwap")
	}
	return maxScenario
}

// TODO: WIP new validity, Fee
func CheckSwapPrice(matchResultXtoY, matchResultYtoX []MatchResult, swapPrice sdk.Dec) bool {
	for _, m := range matchResultXtoY {
		if m.TransactedCoinAmt.Sub(m.FeeAmt).ToDec().Quo(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()).Abs().GT(sdk.OneDec()) {
			fmt.Println(swapPrice, m)
			fmt.Println(m.TransactedCoinAmt.ToDec().Quo(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()))
			fmt.Println(m.TransactedCoinAmt.Sub(m.FeeAmt).ToDec().Quo(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()))
			return false
		}
	}
	for _, m := range matchResultYtoX {
		if m.TransactedCoinAmt.Sub(m.FeeAmt).ToDec().Mul(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()).Abs().GT(sdk.OneDec()) {
			fmt.Println(swapPrice, m)
			fmt.Println(m.TransactedCoinAmt.ToDec().Mul(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()))
			fmt.Println(m.TransactedCoinAmt.Sub(m.FeeAmt).ToDec().Mul(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()))
			return false
		}
	}
	return true
}

// TODO: WIP new validity
func GetMustExecutableAmt(swapPrice sdk.Dec, orderBook OrderBook) (mustExecutableBuyAmtX, mustExecutableSellAmtY sdk.Int) {
	mustExecutableBuyAmtX = sdk.ZeroInt()
	mustExecutableSellAmtY = sdk.ZeroInt()
	for _, order := range orderBook {
		if order.OrderPrice.GT(swapPrice) {
			mustExecutableBuyAmtX = mustExecutableBuyAmtX.Add(order.BuyOfferAmt)
		}
		if order.OrderPrice.LT(swapPrice) {
			mustExecutableSellAmtY = mustExecutableSellAmtY.Add(order.SellOfferAmt)
		}
	}
	return
}

// TODO: WIP new validity
func CheckValidityMustExecutable(orderBook OrderBook, swapPrice sdk.Dec) bool {
	MEX, MEY := GetMustExecutableAmt(swapPrice, orderBook)
	if MEX.GT(sdk.NewInt(1000)) || MEY.GT(sdk.NewInt(1000)) {
		fmt.Println("CheckValidityMustExecutable False", MEX, MEY, swapPrice)
		return false
	} else {
		return true
	}
}


// make orderMap key as swap price, value as Buy, Sell Amount from swap msgs,  with split as Buy XtoY, Sell YtoX msg list.
func GetOrderMap(swapMsgs []*BatchPoolSwapMsg, denomX, denomY string, onlyNotMatched bool) (OrderMap, []*BatchPoolSwapMsg, []*BatchPoolSwapMsg) {
	orderMap := make(OrderMap)
	var XtoY []*BatchPoolSwapMsg // buying Y from X
	var YtoX []*BatchPoolSwapMsg // selling Y for X
	for _, m := range swapMsgs {
		if onlyNotMatched && (m.ToDelete || m.RemainingOfferCoin.IsZero()) {
			continue
		}
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
