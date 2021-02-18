package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"sort"
)

// Execute Swap of the pool batch, Collect swap messages in batch for transact the same price for each batch and run them on endblock.
func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) (uint64, error) {
	// get All only not processed swap msgs, not executed, not succeed, not toDelete
	swapMsgs := k.GetAllNotProcessedLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	if len(swapMsgs) == 0 {
		return 0, nil
	}
	pool, found := k.GetLiquidityPool(ctx, liquidityPoolBatch.PoolId)
	if !found {
		return 0, types.ErrPoolNotExists
	}
	// set all msgs to executed
	for _, msg := range swapMsgs {
		msg.Executed = true
	}
	k.SetLiquidityPoolBatchSwapMsgPointers(ctx, pool.PoolId, swapMsgs)

	params := k.GetParams(ctx)
	currentHeight := ctx.BlockHeight()
	invariantCheckFlag := true // temporary flag for test

	swapMsgs = types.ValidateStateAndExpireOrders(swapMsgs, currentHeight, false)

	// get reserve Coin from the liquidity pool
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	// get current pool pair and price
	X := reserveCoins[0].Amount.ToDec()
	Y := reserveCoins[1].Amount.ToDec()
	currentYPriceOverX := X.QuoTruncate(Y)

	denomX := reserveCoins[0].Denom
	denomY := reserveCoins[1].Denom

	// make orderMap, orderbook by sort orderMap
	orderMap, XtoY, YtoX := types.GetOrderMap(swapMsgs, denomX, denomY, false)
	orderBook := orderMap.SortOrderBook()

	// check orderbook validity and compute batchResult(direction, swapPrice, ..)
	result := types.MatchOrderbook(X, Y, currentYPriceOverX, orderBook)
	resultDec := types.MatchOrderbookDec(X, Y, currentYPriceOverX, orderBook)
	types.BatchResultDecimalDelta(result, resultDec)

	// find order match, calculate pool delta with the total x, y amount for the invariant check
	var matchResultXtoY, matchResultYtoX []types.MatchResult
	var matchResultXtoYDec, matchResultYtoXDec []types.MatchResultDec
	poolXdelta := sdk.ZeroInt()
	poolYdelta := sdk.ZeroInt()
	poolXdeltaDec := sdk.ZeroDec()
	poolYdeltaDec := sdk.ZeroDec()
	if result.MatchType != types.NoMatch {
		var poolXDeltaXtoY, poolXDeltaYtoX, poolYDeltaYtoX, poolYDeltaXtoY sdk.Int
		matchResultXtoY, _, poolXDeltaXtoY, poolYDeltaXtoY = types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX,
			result.SwapPrice, params.SwapFeeRate, currentHeight)
		matchResultYtoX, _, poolXDeltaYtoX, poolYDeltaYtoX = types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY,
			result.SwapPrice, params.SwapFeeRate, currentHeight)
		poolXdelta = poolXDeltaXtoY.Add(poolXDeltaYtoX)
		poolYdelta = poolYDeltaXtoY.Add(poolYDeltaYtoX)

	}
	if resultDec.MatchType != types.NoMatch {
		var poolXDeltaXtoYDec, poolXDeltaYtoXDec, poolYDeltaYtoXDec, poolYDeltaXtoYDec sdk.Dec

		matchResultXtoYDec, _, poolXDeltaXtoYDec, poolYDeltaXtoYDec = types.FindOrderMatchDec(types.DirectionXtoY, XtoY, resultDec.EX,
			resultDec.SwapPrice, params.SwapFeeRate, currentHeight)
		matchResultYtoXDec, _, poolXDeltaYtoXDec, poolYDeltaYtoXDec = types.FindOrderMatchDec(types.DirectionYtoX, YtoX, resultDec.EY,
			resultDec.SwapPrice, params.SwapFeeRate, currentHeight)
		poolXdeltaDec = poolXDeltaXtoYDec.Add(poolXDeltaYtoXDec)
		poolYdeltaDec = poolYDeltaXtoYDec.Add(poolYDeltaYtoXDec)

	}
	types.MatchResultDecimalDelta(matchResultXtoY, matchResultXtoYDec)
	types.MatchResultDecimalDelta(matchResultYtoX, matchResultYtoXDec)
	//fmt.Println("pool_delta", "poolXdelta", "poolYdelta")
	//fmt.Println("pool_delta", poolXdeltaDec.Sub(poolXdelta.ToDec()), poolYdeltaDec.Sub(poolYdelta.ToDec()))
	//if poolXdeltaDec.Sub(poolXdelta.ToDec()).Abs().GT(types.DecimalErrThresholdAmount) ||
	//	poolYdeltaDec.Sub(poolYdelta.ToDec()).Abs().GT(types.DecimalErrThresholdAmount) {
	//	fmt.Println("pool_delta Threshold", poolXdelta, poolXdeltaDec, poolYdelta, poolYdeltaDec)
	//}

	//XtoY, YtoX, X, Y, poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY :=
	//	k.UpdateState(X, Y, XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	XtoY, YtoX, X, Y, _, _, _, _, decimalErrorX, decimalErrorY :=
		k.UpdateStateDec(X, Y, XtoY, YtoX, matchResultXtoYDec, matchResultYtoXDec)

	fmt.Println("UpdateState", "poolXdelta", "poolYdelta", "decimalErrorX", "decimalErrorY")
	fmt.Println("UpdateState", poolXdeltaDec.Sub(poolXdelta.ToDec()), poolYdeltaDec.Sub(poolYdelta.ToDec()), decimalErrorX, decimalErrorY)
	//fmt.Println("UpdateState", poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY)
	lastPrice := X.QuoTruncate(Y)

	//if invariantCheckFlag {
	//	beforeXtoYLen := len(XtoY)
	//	beforeYtoXLen := len(YtoX)
	//	//if beforeXtoYLen-len(matchResultXtoY)+fractionalCntX != (types.MsgList)(XtoY).CountNotMatchedMsgs()+(types.MsgList)(XtoY).CountFractionalMatchedMsgs() {
	//	//	panic(beforeXtoYLen)
	//	//}
	//	//if beforeYtoXLen-len(matchResultYtoX)+fractionalCntY != (types.MsgList)(YtoX).CountNotMatchedMsgs()+(types.MsgList)(YtoX).CountFractionalMatchedMsgs() {
	//	//	panic(beforeYtoXLen)
	//	//}
	//	if beforeXtoYLen-len(matchResultXtoYDec)+fractionalCntX != (types.MsgList)(XtoY).CountNotMatchedMsgs()+(types.MsgList)(XtoY).CountFractionalMatchedMsgs() {
	//		panic(beforeXtoYLen)
	//	}
	//	if beforeYtoXLen-len(matchResultYtoXDec)+fractionalCntY != (types.MsgList)(YtoX).CountNotMatchedMsgs()+(types.MsgList)(YtoX).CountFractionalMatchedMsgs() {
	//		panic(beforeYtoXLen)
	//	}
	//
	//	totalAmtX := sdk.ZeroDec()
	//	totalAmtY := sdk.ZeroDec()
	//
	//	for _, mr := range matchResultXtoYDec {
	//		totalAmtX = totalAmtX.Sub(mr.TransactedCoinAmt)
	//		totalAmtY = totalAmtY.Add(mr.ExchangedDemandCoinAmt)
	//	}
	//
	//	invariantCheckX := totalAmtX
	//	invariantCheckY := totalAmtY
	//
	//	totalAmtX = sdk.ZeroDec()
	//	totalAmtY = sdk.ZeroDec()
	//
	//	for _, mr := range matchResultYtoX {
	//		totalAmtY = totalAmtY.Sub(mr.TransactedCoinAmt.ToDec())
	//		totalAmtX = totalAmtX.Add(mr.ExchangedDemandCoinAmt.ToDec())
	//	}
	//
	//	invariantCheckX = invariantCheckX.Add(totalAmtX)
	//	invariantCheckY = invariantCheckY.Add(totalAmtY)
	//
	//	invariantCheckX = invariantCheckX.Add(poolXdelta.ToDec())
	//	invariantCheckY = invariantCheckY.Add(poolYdelta.ToDec())
	//
	//	// print the invariant check and validity with swap, match result
	//	if invariantCheckX.IsZero() && invariantCheckY.IsZero() {
	//	} else {
	//		panic(invariantCheckX)
	//	}
	//
	//	//if !poolXdelta.Add(decimalErrorX).Equal(poolXdelta2) || !poolYdelta.Add(decimalErrorY).Equal(poolYdelta2) {
	//	//	panic(poolXdelta)
	//	//}
	//	if !poolXdeltaDec.Add(decimalErrorX).Equal(poolXdelta2) || !poolYdeltaDec.Add(decimalErrorY).Equal(poolYdelta2) {
	//		panic(poolXdelta)
	//	}
	//
	//	validitySwapPrice := types.CheckSwapPrice(matchResultXtoY, matchResultYtoX, result.SwapPrice)
	//	if !validitySwapPrice {
	//		panic("validitySwapPrice")
	//	}
	//}

	XtoY = types.ValidateStateAndExpireOrders(XtoY, currentHeight, false)
	YtoX = types.ValidateStateAndExpireOrders(YtoX, currentHeight, false)

	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY, true)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	orderBookValidity := types.CheckValidityOrderBook(orderBookExecuted, lastPrice)
	if !orderBookValidity {
		fmt.Println(orderBookValidity, "ErrOrderBookInvalidity")
		panic(types.ErrOrderBookInvalidity)
	}

	XtoY = types.ValidateStateAndExpireOrders(XtoY, currentHeight, true)
	YtoX = types.ValidateStateAndExpireOrders(YtoX, currentHeight, true)

	// Make index map for match result
	matchResultMap := make(map[uint64]types.MatchResult)
	for _, msg := range matchResultXtoY {
		if _, ok := matchResultMap[msg.OrderMsgIndex]; ok {
			panic("duplicatedMatchOrder")
		}
		matchResultMap[msg.OrderMsgIndex] = msg
		if msg.OrderMsgIndex != matchResultMap[msg.OrderMsgIndex].OrderMsgIndex {
			panic("map broken1")
		}
	}
	for _, msg := range matchResultYtoX {
		if _, ok := matchResultMap[msg.OrderMsgIndex]; ok {
			panic("duplicatedMatchOrder")
		}
		matchResultMap[msg.OrderMsgIndex] = msg
		if msg.OrderMsgIndex != matchResultMap[msg.OrderMsgIndex].OrderMsgIndex {
			panic("map broken1")
		}
	}
	// Make index map for match result
	matchResultMapDec := make(map[uint64]types.MatchResultDec)
	for _, msg := range matchResultXtoYDec {
		if _, ok := matchResultMapDec[msg.OrderMsgIndex]; ok {
			panic("duplicatedMatchOrder")
		}
		matchResultMapDec[msg.OrderMsgIndex] = msg
		if msg.OrderMsgIndex != matchResultMapDec[msg.OrderMsgIndex].OrderMsgIndex {
			panic("map broken1")
		}
	}
	for _, msg := range matchResultYtoXDec {
		if _, ok := matchResultMapDec[msg.OrderMsgIndex]; ok {
			panic("duplicatedMatchOrder")
		}
		matchResultMapDec[msg.OrderMsgIndex] = msg
		if msg.OrderMsgIndex != matchResultMapDec[msg.OrderMsgIndex].OrderMsgIndex {
			panic("map broken1")
		}
	}

	executedMsgCount := uint64(len(swapMsgs))

	//if invariantCheckFlag {
	//	if len(matchResultXtoY)+len(matchResultYtoX) != len(matchResultMap) {
	//		panic("match result map err")
	//	}
	//
	//	for k, v := range matchResultMap {
	//		if k != v.OrderMsgIndex {
	//			panic("broken map consistency")
	//		}
	//	}
	//
	//	// compare swapMsgs state with XtoY, YtoX
	//	notMatchedCount := 0
	//	for k, v := range matchResultMap {
	//		if k != v.OrderMsgIndex {
	//			panic("broken map consistency2")
	//		}
	//	}
	//	for _, msg := range swapMsgs {
	//		for _, msgAfter := range XtoY {
	//			if msg.MsgIndex == msgAfter.MsgIndex {
	//				if *(msg) != *(msgAfter) || msg != msgAfter {
	//					panic("msg not matched")
	//				} else {
	//					break
	//				}
	//			}
	//		}
	//		for _, msgAfter := range YtoX {
	//			if msg.MsgIndex == msgAfter.MsgIndex {
	//				if *(msg) != *(msgAfter) || msg != msgAfter {
	//					panic("msg not matched")
	//				} else {
	//					break
	//				}
	//			}
	//		}
	//		if msgAfter, ok := matchResultMap[msg.MsgIndex]; ok {
	//			if msg.MsgIndex == msgAfter.BatchMsg.MsgIndex {
	//				if *(msg) != *(msgAfter.BatchMsg) || msg != msgAfter.BatchMsg {
	//					panic("msg not matched")
	//				} else {
	//					break
	//				}
	//				// TODO: check for half-half-fee
	//				if !msgAfter.OfferCoinFeeAmt.IsPositive() {
	//					panic(msgAfter.OfferCoinFeeAmt)
	//				}
	//			} else {
	//				panic("fail msg pointer consistency")
	//			}
	//		} else {
	//			// not matched
	//			notMatchedCount++
	//		}
	//	}
	//
	//	// invariant check, swapPrice check
	//	switch result.PriceDirection {
	//	// check whether the calculated swapPrice is actually increased from last pool price
	//	case types.Increase:
	//		if !result.SwapPrice.GT(currentYPriceOverX) {
	//			panic("invariant check fail swapPrice Increase")
	//		}
	//	// check whether the calculated swapPrice is actually decreased from last pool price
	//	case types.Decrease:
	//		if !result.SwapPrice.LT(currentYPriceOverX) {
	//			panic("invariant check fail swapPrice Decrease")
	//		}
	//	// check whether the calculated swapPrice is actually equal to last pool price
	//	case types.Stay:
	//		if !result.SwapPrice.Equal(currentYPriceOverX) {
	//			panic("invariant check fail swapPrice Stay")
	//		}
	//	}
	//
	//	// invariant check, execution validity check
	//	for _, batchMsg := range swapMsgs {
	//		// check whether every executed orders have order price which is not "unexecutable"
	//		if _, ok := matchResultMap[batchMsg.MsgIndex]; ok {
	//			if !batchMsg.Executed || !batchMsg.Succeeded {
	//				panic("batchMsg consistency error, matched but not succeeded")
	//			}
	//
	//			if batchMsg.Msg.OfferCoin.Denom == denomX {
	//				// buy orders having equal or higher order price than found swapPrice
	//				if !batchMsg.Msg.OrderPrice.GTE(result.SwapPrice) {
	//					panic("execution validity failed, executed but unexecutable")
	//				}
	//			} else {
	//				// sell orders having equal or lower order price than found swapPrice
	//				if !batchMsg.Msg.OrderPrice.LTE(result.SwapPrice) {
	//					panic("execution validity failed, executed but unexecutable")
	//				}
	//			}
	//
	//		} else {
	//			// check whether every unexecuted orders have order price which is not "executable"
	//			if batchMsg.Executed && batchMsg.Succeeded {
	//				panic("batchMsg consistency error, not matched but succeeded")
	//			}
	//
	//			if batchMsg.Msg.OfferCoin.Denom == denomX {
	//				// buy orders having equal or lower order price than found swapPrice
	//				if !batchMsg.Msg.OrderPrice.LTE(result.SwapPrice) {
	//					panic("execution validity failed, unexecuted but executable")
	//				}
	//			} else {
	//				// sell orders having equal or higher order price than found swapPrice
	//				if !batchMsg.Msg.OrderPrice.GTE(result.SwapPrice) {
	//					panic("execution validity failed, unexecuted but executable")
	//				}
	//			}
	//		}
	//	}
	//}
	// invariantCheck for matchResultMapDec
	if invariantCheckFlag {
		if len(matchResultXtoYDec)+len(matchResultYtoXDec) != len(matchResultMapDec) {
			panic("match result map err")
		}

		for k, v := range matchResultMapDec {
			if k != v.OrderMsgIndex {
				panic("broken map consistency")
			}
		}

		// compare swapMsgs state with XtoY, YtoX
		notMatchedCount := 0
		for k, v := range matchResultMapDec {
			if k != v.OrderMsgIndex {
				panic("broken map consistency2")
			}
		}
		for _, msg := range swapMsgs {
			for _, msgAfter := range XtoY {
				if msg.MsgIndex == msgAfter.MsgIndex {
					if *(msg) != *(msgAfter) || msg != msgAfter {
						panic("msg not matched")
					} else {
						break
					}
				}
			}
			for _, msgAfter := range YtoX {
				if msg.MsgIndex == msgAfter.MsgIndex {
					if *(msg) != *(msgAfter) || msg != msgAfter {
						panic("msg not matched")
					} else {
						break
					}
				}
			}
			if msgAfter, ok := matchResultMapDec[msg.MsgIndex]; ok {
				if msg.MsgIndex == msgAfter.BatchMsg.MsgIndex {
					if *(msg) != *(msgAfter.BatchMsg) || msg != msgAfter.BatchMsg {
						panic("msg not matched")
					} else {
						break
					}
					// TODO: check for half-half-fee
					if !msgAfter.OfferCoinFeeAmt.IsPositive() {
						panic(msgAfter.OfferCoinFeeAmt)
					}
				} else {
					panic("fail msg pointer consistency")
				}
			} else {
				// not matched
				notMatchedCount++
			}
		}

		// invariant check, swapPrice check
		switch resultDec.PriceDirection {
		// check whether the calculated swapPrice is actually increased from last pool price
		case types.Increase:
			if !resultDec.SwapPrice.GT(currentYPriceOverX) {
				panic("invariant check fail swapPrice Increase")
			}
		// check whether the calculated swapPrice is actually decreased from last pool price
		case types.Decrease:
			if !resultDec.SwapPrice.LT(currentYPriceOverX) {
				panic("invariant check fail swapPrice Decrease")
			}
		// check whether the calculated swapPrice is actually equal to last pool price
		case types.Stay:
			if !resultDec.SwapPrice.Equal(currentYPriceOverX) {
				panic("invariant check fail swapPrice Stay")
			}
		}

		// invariant check, execution validity check
		for _, batchMsg := range swapMsgs {
			// check whether every executed orders have order price which is not "unexecutable"
			if _, ok := matchResultMapDec[batchMsg.MsgIndex]; ok {
				if !batchMsg.Executed || !batchMsg.Succeeded {
					panic("batchMsg consistency error, matched but not succeeded")
				}

				if batchMsg.Msg.OfferCoin.Denom == denomX {
					// buy orders having equal or higher order price than found swapPrice
					if !batchMsg.Msg.OrderPrice.GTE(resultDec.SwapPrice) {
						fmt.Println(batchMsg.Msg.OrderPrice, resultDec.SwapPrice)
						panic("execution validity failed, executed but unexecutable")
					}
				} else {
					// sell orders having equal or lower order price than found swapPrice
					if !batchMsg.Msg.OrderPrice.LTE(resultDec.SwapPrice) {
						fmt.Println(batchMsg.Msg.OrderPrice, resultDec.SwapPrice)
						panic("execution validity failed, executed but unexecutable")
					}
				}

			} else {
				// check whether every unexecuted orders have order price which is not "executable"
				if batchMsg.Executed && batchMsg.Succeeded {
					panic("batchMsg consistency error, not matched but succeeded")
				}

				if batchMsg.Msg.OfferCoin.Denom == denomX {
					// buy orders having equal or lower order price than found swapPrice
					if !batchMsg.Msg.OrderPrice.LTE(resultDec.SwapPrice) {
						fmt.Println(batchMsg.Msg.OrderPrice, resultDec.SwapPrice)
						fmt.Println("execution validity failed, unexecuted but executable")
						//panic("execution validity failed, unexecuted but executable")
					}
				} else {
					// sell orders having equal or higher order price than found swapPrice
					if !batchMsg.Msg.OrderPrice.GTE(resultDec.SwapPrice) {
						fmt.Println(batchMsg.Msg.OrderPrice, resultDec.SwapPrice)
						fmt.Println("execution validity failed, unexecuted but executable")
						//panic("execution validity failed, unexecuted but executable")
					}
				}
			}
		}
	}

	//fmt.Println("########", result.SwapPrice, result.PriceDirection, result.TransactAmt, result.MatchType, result.EX, result.EY)
	//fmt.Println("########", resultDec.SwapPrice, resultDec.PriceDirection, resultDec.TransactAmt, resultDec.MatchType, resultDec.EX, resultDec.EY)
	// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapLiquidityPool
	if err := k.TransactAndRefundSwapLiquidityPoolDec(ctx, swapMsgs, matchResultMapDec, pool, resultDec); err != nil {
		panic(err)
		return executedMsgCount, err
	}
	//if err := k.TransactAndRefundSwapLiquidityPool(ctx, swapMsgs, matchResultMap, pool, result); err != nil {
	//	panic(err)
	//	return executedMsgCount, err
	//}

	return executedMsgCount, nil
}

// Update Buy, Sell swap batch messages using the result of match.
func (k Keeper) UpdateState(X, Y sdk.Dec, XtoY, YtoX []*types.BatchPoolSwapMsg, matchResultXtoY, matchResultYtoX []types.MatchResult) (
	[]*types.BatchPoolSwapMsg, []*types.BatchPoolSwapMsg, sdk.Dec, sdk.Dec, sdk.Int, sdk.Int, int, int, sdk.Int, sdk.Int) {
	sort.SliceStable(XtoY, func(i, j int) bool {
		return XtoY[i].Msg.OrderPrice.GT(XtoY[j].Msg.OrderPrice)
	})
	sort.SliceStable(YtoX, func(i, j int) bool {
		return YtoX[i].Msg.OrderPrice.LT(YtoX[j].Msg.OrderPrice)
	})

	poolXdelta := sdk.ZeroInt()
	poolYdelta := sdk.ZeroInt()
	matchedIndexMapXtoY := make(map[uint64]sdk.Coin)
	matchedIndexMapYtoX := make(map[uint64]sdk.Coin)
	fractionalCntX := 0
	fractionalCntY := 0

	// Variables to accumulate and offset the values of int 1 caused by decimal error
	decimalErrorX := sdk.ZeroInt()
	decimalErrorY := sdk.ZeroInt()

	for _, match := range matchResultXtoY {
		poolXdelta = poolXdelta.Add(match.TransactedCoinAmt)
		poolYdelta = poolYdelta.Sub(match.ExchangedDemandCoinAmt)
		if match.BatchMsg.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Equal(match.TransactedCoinAmt) {
			// full match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt))

			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt)
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 1")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) {
			// TODO: add testcase for coverage
			decimalErrorX = decimalErrorX.Add(sdk.OneInt())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt)
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()
			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 2")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt)
			matchedIndexMapXtoY[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeeded = true
			match.BatchMsg.ToBeDeleted = false
			fractionalCntX += 1
		}
	}
	for _, match := range matchResultYtoX {
		poolXdelta = poolXdelta.Sub(match.ExchangedDemandCoinAmt)
		poolYdelta = poolYdelta.Add(match.TransactedCoinAmt)
		if match.BatchMsg.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Equal(match.TransactedCoinAmt) {
			// full match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt)
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 3")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) {
			// TODO: add testcase for coverage
			decimalErrorY = decimalErrorY.Add(sdk.OneInt())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt)
			// TODO: verify RemainingOfferCoin about decimal errors one to pool
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()

			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 4")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt)
			matchedIndexMapYtoX[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeeded = true
			match.BatchMsg.ToBeDeleted = false
			fractionalCntY += 1
		}
	}

	// Offset accumulated decimal error values
	poolXdelta = poolXdelta.Add(decimalErrorX)
	poolYdelta = poolYdelta.Add(decimalErrorY)

	X = X.Add(poolXdelta.ToDec())
	Y = Y.Add(poolYdelta.ToDec())

	return XtoY, YtoX, X, Y, poolXdelta, poolYdelta, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY
}

func (k Keeper) UpdateStateDec(X, Y sdk.Dec, XtoY, YtoX []*types.BatchPoolSwapMsg, matchResultXtoY, matchResultYtoX []types.MatchResultDec) (
	[]*types.BatchPoolSwapMsg, []*types.BatchPoolSwapMsg, sdk.Dec, sdk.Dec, sdk.Dec, sdk.Dec, int, int, sdk.Dec, sdk.Dec) {
	sort.SliceStable(XtoY, func(i, j int) bool {
		return XtoY[i].Msg.OrderPrice.GT(XtoY[j].Msg.OrderPrice)
	})
	sort.SliceStable(YtoX, func(i, j int) bool {
		return YtoX[i].Msg.OrderPrice.LT(YtoX[j].Msg.OrderPrice)
	})

	poolXdelta := sdk.ZeroDec()
	poolYdelta := sdk.ZeroDec()
	matchedIndexMapXtoY := make(map[uint64]sdk.Coin)
	matchedIndexMapYtoX := make(map[uint64]sdk.Coin)
	fractionalCntX := 0
	fractionalCntY := 0

	// Variables to accumulate and offset the values of int 1 caused by decimal error
	decimalErrorX := sdk.ZeroDec()
	decimalErrorY := sdk.ZeroDec()

	for _, match := range matchResultXtoY {
		poolXdelta = poolXdelta.Add(match.TransactedCoinAmt)
		poolYdelta = poolYdelta.Sub(match.ExchangedDemandCoinAmt)
		if match.BatchMsg.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt.TruncateInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Equal(match.TransactedCoinAmt.TruncateInt()) {
			// full match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt.TruncateInt()))

			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt.TruncateInt())
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt.TruncateInt())
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 1")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt.TruncateInt()).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt.TruncateInt()).Equal(sdk.OneInt()) {
			// TODO: add testcase for coverage
			decimalErrorX = decimalErrorX.Add(sdk.OneDec())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt.TruncateInt()))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt.TruncateInt())
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt.TruncateInt())
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()
			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 2")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt.TruncateInt()))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt.TruncateInt())
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt.TruncateInt())
			matchedIndexMapXtoY[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeeded = true
			match.BatchMsg.ToBeDeleted = false
			fractionalCntX += 1
		}
	}
	for _, match := range matchResultYtoX {
		poolXdelta = poolXdelta.Sub(match.ExchangedDemandCoinAmt)
		poolYdelta = poolYdelta.Add(match.TransactedCoinAmt)
		if match.BatchMsg.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt.TruncateInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Equal(match.TransactedCoinAmt.TruncateInt()) {
			// full match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt.TruncateInt()))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt.TruncateInt())
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt.TruncateInt())
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 3")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt.TruncateInt()).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt.TruncateInt()).Equal(sdk.OneInt()) {
			// TODO: add testcase for coverage
			decimalErrorY = decimalErrorY.Add(sdk.OneDec())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt.TruncateInt()))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt.TruncateInt())
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt.TruncateInt())
			// TODO: verify RemainingOfferCoin about decimal errors one to pool
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()

			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) ||
				match.BatchMsg.OfferCoinFeeReserve.IsGTE(sdk.NewCoin(match.BatchMsg.OfferCoinFeeReserve.Denom, sdk.NewInt(2))) {
				panic("remaining not matched 4")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt.TruncateInt()))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt.TruncateInt())
			match.BatchMsg.OfferCoinFeeReserve = types.CoinSafeSubAmount(match.BatchMsg.OfferCoinFeeReserve, match.OfferCoinFeeAmt.TruncateInt())
			matchedIndexMapYtoX[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeeded = true
			match.BatchMsg.ToBeDeleted = false
			fractionalCntY += 1
		}
	}

	// Offset accumulated decimal error values
	poolXdelta = poolXdelta.Add(decimalErrorX)
	poolYdelta = poolYdelta.Add(decimalErrorY)

	X = X.Add(poolXdelta)
	Y = Y.Add(poolYdelta)

	return XtoY, YtoX, X, Y, poolXdelta, poolYdelta, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY
}
