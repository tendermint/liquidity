package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

// liquidity module const types for swap
const (
	// Price Directions
	Increase = 1
	Decrease = 2
	Stay     = 3

	// Match Types
	ExactMatch      = 1
	NoMatch         = 2
	FractionalMatch = 3

	// Order Directions
	DirectionXtoY = 1
	DirectionYtoX = 2
)

// Type of order map to index at price, having the pointer list of the swap batch message.
type OrderByPrice struct {
	OrderPrice   sdk.Dec
	BuyOfferAmt  sdk.Int
	SellOfferAmt sdk.Int
	MsgList      []*BatchPoolSwapMsg
}

// list of orderByPrice
type OrderBook []OrderByPrice

// Len implements sort.Interface for OrderBook
func (orderBook OrderBook) Len() int { return len(orderBook) }

// Less implements sort.Interface for OrderBook
func (orderBook OrderBook) Less(i, j int) bool {
	return orderBook[i].OrderPrice.LT(orderBook[j].OrderPrice)
}

// Swap implements sort.Interface for OrderBook
func (orderBook OrderBook) Swap(i, j int) { orderBook[i], orderBook[j] = orderBook[j], orderBook[i] }

// increasing sort orderbook by order price
func (orderBook OrderBook) Sort() {
	//sort.Sort(orderBook)
	sort.Slice(orderBook, func(i, j int) bool {
		return orderBook[i].OrderPrice.LT(orderBook[j].OrderPrice)
	})
}

// decreasing sort orderbook by order price
func (orderBook OrderBook) Reverse() {
	//sort.Reverse(orderBook)
	sort.Slice(orderBook, func(i, j int) bool {
		return orderBook[i].OrderPrice.GT(orderBook[j].OrderPrice)
	})
}

// The pointer list of the swap batch message.
type MsgList []*BatchPoolSwapMsg

// Get number of not matched messages on the list.
func (msgList MsgList) CountNotMatchedMsgs() int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && !m.Succeeded {
			cnt++
		}
	}
	return cnt
}

// Get number of fractional matched messages on the list.
func (msgList MsgList) CountFractionalMatchedMsgs() int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && m.Succeeded && !m.ToBeDeleted {
			cnt++
		}
	}
	return cnt
}

// Return minimum Decimal
func MinDec(a, b sdk.Dec) sdk.Dec {
	if a.LTE(b) {
		return a
	} else {
		return b
	}
}

// Return maximum Decimal
func MaxDec(a, b sdk.Dec) sdk.Dec {
	if a.GTE(b) {
		return a
	} else {
		return b
	}
}

// Return minimum Int
func MinInt(a, b sdk.Int) sdk.Int {
	if a.LTE(b) {
		return a
	} else {
		return b
	}
}

// Return Maximum Int
func MaxInt(a, b sdk.Int) sdk.Int {
	if a.GTE(b) {
		return a
	} else {
		return b
	}
}

// Order map type indexed by order price at price
type OrderMap map[string]OrderByPrice

// Make orderbook by sort orderMap.
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

// struct of swap matching result of the batch
type BatchResult struct {
	MatchType      int
	PriceDirection int
	SwapPrice      sdk.Dec
	EX             sdk.Int
	EY             sdk.Int
	OriginalEX     sdk.Int
	OriginalEY     sdk.Int
	PoolX          sdk.Int
	PoolY          sdk.Int
	TransactAmt    sdk.Int
}

// return of zero object, to avoid nil
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

// struct of swap matching result of each Batch swap message
type MatchResult struct {
	OrderHeight            int64
	OrderExpiryHeight      int64
	OrderMsgIndex          uint64
	OrderPrice             sdk.Dec
	OfferCoinAmt           sdk.Int
	TransactedCoinAmt      sdk.Int
	ExchangedDemandCoinAmt sdk.Int
	OfferCoinFeeAmt        sdk.Int
	ExchangedCoinFeeAmt    sdk.Int
	BatchMsg               *BatchPoolSwapMsg
}

// The price and coins of swap messages in orderbook are calculated
// to derive match result with the price direction.
func MatchOrderbook(X, Y, currentPrice sdk.Dec, orderBook OrderBook) (result BatchResult) {
	result = NewBatchResult()
	orderBook.Sort()
	priceDirection := GetPriceDirection(currentPrice, orderBook)

	if priceDirection == Stay {
		return CalculateMatchStay(currentPrice, orderBook)
	} else { // Increase, Decrease
		if priceDirection == Decrease {
			orderBook.Reverse()
		}
		return CalculateMatch(priceDirection, X, Y, currentPrice, orderBook)
	}
}

// Check orderbook validity
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
	oneOverWithErr, _ := sdk.NewDecFromStr("1.10")
	oneUnderWithErr, _ := sdk.NewDecFromStr("0.90")

	if maxBuyOrderPrice.GT(minSellOrderPrice) ||
		maxBuyOrderPrice.Quo(currentPrice).GT(oneOverWithErr) ||
		minSellOrderPrice.Quo(currentPrice).LT(oneUnderWithErr) {
		return false
	} else {
		return true
	}
}

//check validity state of the batch swap messages, and set to delete state to height timeout expired order
func ValidateStateAndExpireOrders(msgList []*BatchPoolSwapMsg, currentHeight int64, expireThisHeight bool) []*BatchPoolSwapMsg {
	for _, order := range msgList {
		if !order.Executed {
			panic("not executed")
		}
		if order.RemainingOfferCoin.IsZero() {
			if !order.Succeeded || !order.ToBeDeleted {
				panic("broken state consistency for not matched order")
			}
			order.Succeeded = true
			order.ToBeDeleted = true
			continue
		}
		// set toDelete, expired msgs
		if currentHeight > order.OrderExpiryHeight {
			if order.Succeeded || !order.ToBeDeleted {
				panic("broken state consistency for fractional matched order")
			}
			order.Succeeded = false
			order.ToBeDeleted = true
			continue
		}
		if expireThisHeight && currentHeight == order.OrderExpiryHeight {
			order.ToBeDeleted = true
		}
	}
	return msgList
}

// Calculate results for orderbook matching with unchanged price case
func CalculateMatchStay(currentPrice sdk.Dec, orderBook OrderBook) (r BatchResult) {
	r = NewBatchResult()
	r.SwapPrice = currentPrice
	r.OriginalEX, r.OriginalEY = GetExecutableAmt(r.SwapPrice, orderBook)
	r.EX = r.OriginalEX
	r.EY = r.OriginalEY
	r.PriceDirection = Stay

	if r.EX.Add(r.PoolX).Equal(sdk.ZeroInt()) || r.EY.Add(r.PoolY).Equal(sdk.ZeroInt()) {
		r.MatchType = NoMatch
		// Normalization to an integrator for easy determination of exactMatch
	} else if r.EX.Equal(r.SwapPrice.MulInt(r.EY).TruncateInt()) {
		r.MatchType = ExactMatch
	} else {
		// Decimal Error, When calculating the Executable value, conservatively Truncated decimal
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
	var matchOrderList []*BatchPoolSwapMsg

	if executableAmt.IsZero() {
		return
	}

	lenSwapList := len(swapList)
	for i, order := range swapList {
		var breakFlag, appendFlag bool

		// include the matched order in matchAmt, matchOrderList
		if (direction == DirectionXtoY && order.Msg.OrderPrice.GTE(swapPrice)) ||
			(direction == DirectionYtoX && order.Msg.OrderPrice.LTE(swapPrice)) {
			matchAmt = matchAmt.Add(order.RemainingOfferCoin.Amount)
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
					if fractionalMatchRatio.GT(sdk.NewDec(1)) {
						panic("Invariant Check: fractionalMatchRatio between 0 and 1")
					}
				} else {
					fractionalMatchRatio = sdk.OneDec()
				}
				for _, matchOrder := range matchOrderList {
					if !fractionalMatchRatio.IsPositive() {
						fractionalMatchRatio = sdk.OneDec()
						continue
					}
					offerAmt := matchOrder.RemainingOfferCoin.Amount.ToDec()
					matchResult := MatchResult{
						OrderHeight:       height,
						OrderExpiryHeight: height + CancelOrderLifeSpan,
						OrderMsgIndex:     matchOrder.MsgIndex,
						OrderPrice:        matchOrder.Msg.OrderPrice,
						OfferCoinAmt:      matchOrder.RemainingOfferCoin.Amount,
						// TransactedCoinAmt is a value that should not be lost, so Ceil it conservatively considering the decimal error.
						TransactedCoinAmt: offerAmt.Mul(fractionalMatchRatio).Ceil().TruncateInt(),
						BatchMsg:          matchOrder,
					}
					if matchOrder != matchResult.BatchMsg {
						panic("not matched msg pointer ")
					}
					// Fee, Exchanged amount are values that should not be overmeasured, so it is lowered conservatively considering the decimal error.
					if direction == DirectionXtoY {
						matchResult.OfferCoinFeeAmt = matchResult.BatchMsg.OfferCoinFeeReserve.Amount.ToDec().Mul(fractionalMatchRatio).TruncateInt()
						matchResult.ExchangedDemandCoinAmt = matchResult.TransactedCoinAmt.ToDec().Quo(swapPrice).TruncateInt()
						matchResult.ExchangedCoinFeeAmt = matchResult.OfferCoinFeeAmt.ToDec().Quo(swapPrice).TruncateInt()
					} else if direction == DirectionYtoX {
						matchResult.OfferCoinFeeAmt = matchResult.BatchMsg.OfferCoinFeeReserve.Amount.ToDec().Mul(fractionalMatchRatio).TruncateInt()
						matchResult.ExchangedDemandCoinAmt = matchResult.TransactedCoinAmt.ToDec().Mul(swapPrice).TruncateInt()
						matchResult.ExchangedCoinFeeAmt = matchResult.OfferCoinFeeAmt.ToDec().Mul(swapPrice).TruncateInt()
					}
					// Check for differences above maximum decimal error
					if matchResult.TransactedCoinAmt.GT(matchResult.OfferCoinAmt) ||
						(matchResult.OfferCoinFeeAmt.GT(matchResult.OfferCoinAmt) && matchResult.OfferCoinFeeAmt.GT(sdk.OneInt())) {
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

	// Normalization to an integrator for easy determination of exactMatch. this decimal error will be minimize
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
		// When calculating the Pool value, conservatively Truncated decimal, so Ceil it to reduce the decimal error
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

	// Round to an integer to minimize decimal errors.
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
			// Check to see if it exceeds a value that can be a decimal error
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
		MEX, MEY := GetMustExecutableAmt(s.SwapPrice, orderBook)
		if s.EX.GTE(MEX) && s.EY.GTE(MEY) {
			if s.MatchType == ExactMatch && s.TransactAmt.IsPositive() {
				maxScenario = s
				break
			} else if s.TransactAmt.GT(maxScenario.TransactAmt) {
				maxScenario = s
			}
		}
	}

	// Invariant Check
	r := maxScenario
	tmpInvariant := r.EX.Add(r.PoolX).ToDec().Sub(r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice))
	if tmpInvariant.GT(r.SwapPrice) && tmpInvariant.GT(sdk.OneDec()) {
		panic("maxScenario CalculateSwap")
	}

	maxScenario.PriceDirection = direction
	return maxScenario
}

// Check swap price validity using list of match result.
func CheckSwapPrice(matchResultXtoY, matchResultYtoX []MatchResult, swapPrice sdk.Dec) bool {
	if len(matchResultXtoY) == 0 && len(matchResultYtoX) == 0 {
		return true
	}
	// Check if it is greater than a value that can be a decimal error
	for _, m := range matchResultXtoY {
		if m.TransactedCoinAmt.ToDec().Quo(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()).Abs().GT(sdk.OneDec()) {
			return false
		}
	}
	for _, m := range matchResultYtoX {
		if m.TransactedCoinAmt.ToDec().Mul(swapPrice).Sub(m.ExchangedDemandCoinAmt.ToDec()).Abs().GT(sdk.OneDec()) {
			return false
		}
	}
	if swapPrice.IsZero() {
		return false
	}
	return true
}

// Check swap executable amount validity of the orderbook
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

// make orderMap key as swap price, value as Buy, Sell Amount from swap msgs,  with split as Buy XtoY, Sell YtoX msg list.
func GetOrderMap(swapMsgs []*BatchPoolSwapMsg, denomX, denomY string, onlyNotMatched bool) (OrderMap, []*BatchPoolSwapMsg, []*BatchPoolSwapMsg) {
	orderMap := make(OrderMap)
	var XtoY []*BatchPoolSwapMsg // buying Y from X
	var YtoX []*BatchPoolSwapMsg // selling Y for X
	for _, m := range swapMsgs {
		if onlyNotMatched && (m.ToBeDeleted || m.RemainingOfferCoin.IsZero()) {
			continue
		}
		if m.Msg.OfferCoin.Denom == denomX { // buying Y from X
			XtoY = append(XtoY, m)
			if _, ok := orderMap[m.Msg.OrderPrice.String()]; ok {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice.String()].BuyOfferAmt.Add(m.RemainingOfferCoin.Amount), // TODO: feeX half
					orderMap[m.Msg.OrderPrice.String()].SellOfferAmt,
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			} else {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{m.Msg.OrderPrice,
					m.RemainingOfferCoin.Amount, sdk.ZeroInt(),
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			}
		} else if m.Msg.OfferCoin.Denom == denomY { // selling Y for X
			YtoX = append(YtoX, m)
			if _, ok := orderMap[m.Msg.OrderPrice.String()]; ok {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{
					m.Msg.OrderPrice,
					orderMap[m.Msg.OrderPrice.String()].BuyOfferAmt,
					orderMap[m.Msg.OrderPrice.String()].SellOfferAmt.Add(m.RemainingOfferCoin.Amount),
					append(orderMap[m.Msg.OrderPrice.String()].MsgList, m),
				}
			} else {
				orderMap[m.Msg.OrderPrice.String()] = OrderByPrice{m.Msg.OrderPrice,
					sdk.ZeroInt(), m.RemainingOfferCoin.Amount,
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
