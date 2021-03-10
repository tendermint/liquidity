package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Type of match
type MatchType int

const (
	ExactMatch MatchType = iota + 1
	NoMatch
	FractionalMatch
)

// Direction of price
type PriceDirection int

const (
	Increasing PriceDirection = iota + 1
	Decreasing
	Staying
)

// Direction of order
type OrderDirection int

const (
	DirectionXtoY OrderDirection = iota + 1
	DirectionYtoX
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

// Get number of not matched messages on the list.
func CountNotMatchedMsgs(msgList []*BatchPoolSwapMsg) int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && !m.Succeeded {
			cnt++
		}
	}
	return cnt
}

// Get number of fractional matched messages on the list.
func CountFractionalMatchedMsgs(msgList []*BatchPoolSwapMsg) int {
	cnt := 0
	for _, m := range msgList {
		if m.Executed && m.Succeeded && !m.ToBeDeleted {
			cnt++
		}
	}
	return cnt
}

// Order map type indexed by order price at price
type OrderMap map[string]OrderByPrice

// Make orderbook by sort orderMap.
func (orderMap OrderMap) SortOrderBook() (orderBook OrderBook) {
	for _, o := range orderMap {
		orderBook = append(orderBook, o)
	}
	orderBook.Sort()
	return orderBook
}

// struct of swap matching result of the batch
type BatchResult struct {
	MatchType      MatchType
	PriceDirection PriceDirection
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
func (orderBook OrderBook) Match(X, Y sdk.Dec) BatchResult {
	currentPrice := X.Quo(Y)
	priceDirection := orderBook.PriceDirection(currentPrice)
	if priceDirection == Staying {
		return orderBook.CalculateMatchStay(currentPrice)
	}
	if priceDirection == Decreasing {
		orderBook.Reverse()
	}
	return orderBook.CalculateMatch(priceDirection, X, Y)
}

// Check orderbook validity
func (orderBook OrderBook) Validate(currentPrice sdk.Dec) bool {
	//orderBook.Reverse() // TODO: is it necessary?
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
	if maxBuyOrderPrice.GT(minSellOrderPrice) ||
		maxBuyOrderPrice.Quo(currentPrice).GT(sdk.MustNewDecFromStr("1.10")) ||
		minSellOrderPrice.Quo(currentPrice).LT(sdk.MustNewDecFromStr("0.90")) {
		return false
	}
	return true
}

// Calculate results for orderbook matching with unchanged price case
func (orderBook OrderBook) CalculateMatchStay(currentPrice sdk.Dec) (r BatchResult) {
	r = NewBatchResult()
	r.SwapPrice = currentPrice
	r.OriginalEX, r.OriginalEY = orderBook.ExecutableAmt(r.SwapPrice)
	r.EX = r.OriginalEX
	r.EY = r.OriginalEY
	r.PriceDirection = Staying

	s := r.SwapPrice.MulInt(r.EY).TruncateInt()
	if r.EX.IsZero() || r.EY.IsZero() {
		r.MatchType = NoMatch
	} else if r.EX.Equal(s) { // Normalization to an integrator for easy determination of exactMatch
		r.MatchType = ExactMatch
	} else {
		// Decimal Error, When calculating the Executable value, conservatively Truncated decimal
		r.MatchType = FractionalMatch
		if r.EX.GT(s) {
			r.EX = s
		} else if r.EX.LT(s) {
			r.EY = r.EX.ToDec().Quo(r.SwapPrice).TruncateInt()
		}
	}
	return
}

// Calculates the batch results with the logic for each direction
func (orderBook OrderBook) CalculateMatch(direction PriceDirection, X, Y sdk.Dec) (maxScenario BatchResult) {
	currentPrice := X.Quo(Y)
	lastOrderPrice := currentPrice
	var matchScenarioList []BatchResult
	for _, order := range orderBook {
		if (direction == Increasing && order.OrderPrice.LT(currentPrice)) ||
			(direction == Decreasing && order.OrderPrice.GT(currentPrice)) {
			continue
		} else {
			orderPrice := order.OrderPrice
			r := orderBook.CalculateSwap(direction, X, Y, orderPrice, lastOrderPrice)
			// Check to see if it exceeds a value that can be a decimal error
			if (direction == Increasing && r.PoolY.ToDec().Sub(r.EX.ToDec().Quo(r.SwapPrice)).GTE(sdk.OneDec())) ||
				(direction == Decreasing && r.PoolX.ToDec().Sub(r.EY.ToDec().Mul(r.SwapPrice)).GTE(sdk.OneDec())) {
				continue
			}
			matchScenarioList = append(matchScenarioList, r)
			lastOrderPrice = orderPrice
		}
	}
	maxScenario = NewBatchResult()
	maxScenario.TransactAmt = sdk.ZeroInt()
	for _, s := range matchScenarioList {
		MEX, MEY := orderBook.MustExecutableAmt(s.SwapPrice)
		if s.EX.GTE(MEX) && s.EY.GTE(MEY) {
			if s.MatchType == ExactMatch && s.TransactAmt.IsPositive() {
				maxScenario = s
				break
			} else if s.TransactAmt.GT(maxScenario.TransactAmt) {
				maxScenario = s
			}
		}
	}
	maxScenario.PriceDirection = direction
	return maxScenario
}

// Calculates the batch results with the processing logic for each direction
func (orderBook OrderBook) CalculateSwap(direction PriceDirection, X, Y, orderPrice, lastOrderPrice sdk.Dec) (r BatchResult) {
	r = NewBatchResult()
	r.OriginalEX, r.OriginalEY = orderBook.ExecutableAmt(lastOrderPrice.Add(orderPrice).Quo(sdk.NewDec(2)))
	r.EX = r.OriginalEX
	r.EY = r.OriginalEY

	//r.SwapPrice = X.Add(r.EX).Quo(Y.Add(r.EY)) // legacy constant product model
	r.SwapPrice = X.Add(r.EX.MulRaw(2).ToDec()).Quo(Y.Add(r.EY.MulRaw(2).ToDec())) // newSwapPriceModel

	// Normalization to an integrator for easy determination of exactMatch. this decimal error will be minimize
	if direction == Increasing {
		//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))  // legacy constant product model
		r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)).TruncateInt() // newSwapPriceModel
		if lastOrderPrice.LT(r.SwapPrice) && r.SwapPrice.LT(orderPrice) && !r.PoolY.IsNegative() {
			if r.EX.IsZero() && r.EY.IsZero() {
				r.MatchType = NoMatch
			} else {
				r.MatchType = ExactMatch
			}
		}
	} else if direction == Decreasing {
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
		r.OriginalEX, r.OriginalEY = orderBook.ExecutableAmt(orderPrice)
		r.EX = r.OriginalEX
		r.EY = r.OriginalEY
		r.SwapPrice = orderPrice
		// When calculating the Pool value, conservatively Truncated decimal, so Ceil it to reduce the decimal error
		if direction == Increasing {
			//r.PoolY = Y.Sub(X.Quo(r.SwapPrice))  // legacy constant product model
			r.PoolY = r.SwapPrice.Mul(Y).Sub(X).Quo(r.SwapPrice.MulInt64(2)).TruncateInt() // newSwapPriceModel
			r.EX = sdk.MinDec(r.EX.ToDec(), r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice)).Ceil().TruncateInt()
			r.EY = sdk.MaxDec(sdk.MinDec(r.EY.ToDec(), r.EX.ToDec().Quo(r.SwapPrice).Sub(r.PoolY.ToDec())), sdk.ZeroDec()).Ceil().TruncateInt()
		} else if direction == Decreasing {
			//r.PoolX = X.Sub(Y.Mul(r.SwapPrice)) // legacy constant product model
			r.PoolX = X.Sub(r.SwapPrice.Mul(Y)).QuoInt64(2).TruncateInt() // newSwapPriceModel
			r.EY = sdk.MinDec(r.EY.ToDec(), r.EX.Add(r.PoolX).ToDec().Quo(r.SwapPrice)).Ceil().TruncateInt()
			r.EX = sdk.MaxDec(sdk.MinDec(r.EX.ToDec(), r.EY.ToDec().Mul(r.SwapPrice).Sub(r.PoolX.ToDec())), sdk.ZeroDec()).Ceil().TruncateInt()
		}
		r.MatchType = FractionalMatch
	}

	// Round to an integer to minimize decimal errors.
	if direction == Increasing {
		if r.SwapPrice.LT(X.Quo(Y)) || r.PoolY.IsNegative() {
			r.TransactAmt = sdk.ZeroInt()
		} else {
			r.TransactAmt = sdk.MinInt(r.EX, r.EY.Add(r.PoolY).ToDec().Mul(r.SwapPrice).RoundInt())
		}
	} else if direction == Decreasing {
		if r.SwapPrice.GT(X.Quo(Y)) || r.PoolX.LT(sdk.ZeroInt()) {
			r.TransactAmt = sdk.ZeroInt()
		} else {
			r.TransactAmt = sdk.MinInt(r.EY, r.EX.Add(r.PoolX).ToDec().Quo(r.SwapPrice).RoundInt())
		}
	}
	return
}

// Get Price direction of the orderbook with current Price
func (orderBook OrderBook) PriceDirection(currentPrice sdk.Dec) PriceDirection {
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

	if buyAmtOverCurrentPrice.GT(currentPrice.Mul(sellAmtUnderCurrentPrice.Add(sellAmtAtCurrentPrice))) {
		return Increasing
	} else if currentPrice.Mul(sellAmtUnderCurrentPrice).GT(buyAmtOverCurrentPrice.Add(buyAmtAtCurrentPrice)) {
		return Decreasing
	}
	return Staying
}

// calculate the executable amount of the orderbook for each X, Y
func (orderBook OrderBook) ExecutableAmt(swapPrice sdk.Dec) (executableBuyAmtX, executableSellAmtY sdk.Int) {
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

// Check swap executable amount validity of the orderbook
func (orderBook OrderBook) MustExecutableAmt(swapPrice sdk.Dec) (mustExecutableBuyAmtX, mustExecutableSellAmtY sdk.Int) {
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

// make orderMap key as swap price, value as Buy, Sell Amount from swap msgs, with split as Buy XtoY, Sell YtoX msg list.
func MakeOrderMap(swapMsgs []*BatchPoolSwapMsg, denomX, denomY string, onlyNotMatched bool) (OrderMap, []*BatchPoolSwapMsg, []*BatchPoolSwapMsg) {
	orderMap := make(OrderMap)
	var XtoY []*BatchPoolSwapMsg // buying Y from X
	var YtoX []*BatchPoolSwapMsg // selling Y for X
	for _, m := range swapMsgs {
		if onlyNotMatched && (m.ToBeDeleted || m.RemainingOfferCoin.IsZero()) {
			continue
		}
		order := OrderByPrice{
			OrderPrice:   m.Msg.OrderPrice,
			BuyOfferAmt:  sdk.ZeroInt(),
			SellOfferAmt: sdk.ZeroInt(),
		}
		orderPriceString := m.Msg.OrderPrice.String()
		switch {
		// buying Y from X
		case m.Msg.OfferCoin.Denom == denomX:
			XtoY = append(XtoY, m)
			if o, ok := orderMap[orderPriceString]; ok {
				order = o
				order.BuyOfferAmt = o.BuyOfferAmt.Add(m.RemainingOfferCoin.Amount) // TODO: feeX half
			} else {
				order.BuyOfferAmt = m.RemainingOfferCoin.Amount
			}
		// selling Y for X
		case m.Msg.OfferCoin.Denom == denomY:
			YtoX = append(YtoX, m)
			if o, ok := orderMap[orderPriceString]; ok {
				order = o
				order.SellOfferAmt = o.SellOfferAmt.Add(m.RemainingOfferCoin.Amount)
			} else {
				order.SellOfferAmt = m.RemainingOfferCoin.Amount
			}
		default:
			panic(ErrInvalidDenom)
		}
		order.MsgList = append(order.MsgList, m)
		orderMap[orderPriceString] = order
	}
	return orderMap, XtoY, YtoX
}

//check validity state of the batch swap messages, and set to delete state to height timeout expired order
func ValidateStateAndExpireOrders(msgList []*BatchPoolSwapMsg, currentHeight int64, expireThisHeight bool) {
	for _, order := range msgList {
		if !order.Executed {
			panic("not executed")
		}
		if order.RemainingOfferCoin.IsZero() {
			if !order.Succeeded || !order.ToBeDeleted {
				panic("broken state consistency for not matched order")
			}
			continue
		}
		// set toDelete, expired msgs
		if currentHeight > order.OrderExpiryHeight {
			if order.Succeeded || !order.ToBeDeleted {
				panic("broken state consistency for fractional matched order")
			}
			continue
		}
		if expireThisHeight && currentHeight == order.OrderExpiryHeight {
			order.ToBeDeleted = true
		}
	}
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

// Find matched orders and set status for msgs
func FindOrderMatch(direction OrderDirection, swapList []*BatchPoolSwapMsg, executableAmt sdk.Int, swapPrice sdk.Dec, height int64) (
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
				if !fractionalMatchRatio.IsPositive() {
					fractionalMatchRatio = sdk.OneDec()
				}
				for _, matchOrder := range matchOrderList {
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
						panic("not matched msg pointer")
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
