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

	swapMsgs = types.ValidateStateAndExpireOrders(swapMsgs, currentHeight, false)

	// get reserve Coin from the liquidity pool
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	// get current pool pair and price
	X := reserveCoins[0].Amount.ToDec()
	Y := reserveCoins[1].Amount.ToDec()
	currentYPriceOverX := X.Quo(Y)

	denomX := reserveCoins[0].Denom
	denomY := reserveCoins[1].Denom

	// make orderMap, orderbook by sort orderMap
	orderMap, XtoY, YtoX := types.GetOrderMap(swapMsgs, denomX, denomY, false)
	orderBook := orderMap.SortOrderBook()

	// check orderbook validity and compute batchResult(direction, swapPrice, ..)
	result := types.MatchOrderbook(X, Y, currentYPriceOverX, orderBook)

	// find order match, calculate pool delta with the total x, y amount for the invariant check
	beforeXtoYLen := len(XtoY)
	beforeYtoXLen := len(YtoX)
	var matchResultXtoY, matchResultYtoX []types.MatchResult
	poolXdelta := sdk.ZeroInt()
	poolYdelta := sdk.ZeroInt()
	if result.MatchType != types.NoMatch {
		var poolXDeltaXtoY, poolXDeltaYtoX, poolYDeltaYtoX, poolYDeltaXtoY sdk.Int
		matchResultXtoY, _, poolXDeltaXtoY, poolYDeltaXtoY = types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX,
			result.SwapPrice, params.SwapFeeRate, currentHeight)
		matchResultYtoX, _, poolXDeltaYtoX, poolYDeltaYtoX = types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY,
			result.SwapPrice, params.SwapFeeRate, currentHeight)
		poolXdelta = poolXDeltaXtoY.Add(poolXDeltaYtoX)
		poolYdelta = poolYDeltaXtoY.Add(poolYDeltaYtoX)
	}

	XtoY, YtoX, X, Y, poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY :=
		k.UpdateState(X, Y, XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	lastPrice := X.Quo(Y)

	if beforeXtoYLen-len(matchResultXtoY)+fractionalCntX != (types.MsgList)(XtoY).CountNotMatchedMsgs()+(types.MsgList)(XtoY).CountFractionalMatchedMsgs() {
		panic(beforeXtoYLen)
	}
	if beforeYtoXLen-len(matchResultYtoX)+fractionalCntY != (types.MsgList)(YtoX).CountNotMatchedMsgs()+(types.MsgList)(YtoX).CountFractionalMatchedMsgs() {
		panic(beforeYtoXLen)
	}

	totalAmtX := sdk.ZeroInt()
	totalAmtY := sdk.ZeroInt()

	for _, mr := range matchResultXtoY {
		totalAmtX = totalAmtX.Sub(mr.TransactedCoinAmt)
		totalAmtY = totalAmtY.Add(mr.ExchangedDemandCoinAmt)
	}

	invariantCheckX := totalAmtX
	invariantCheckY := totalAmtY

	totalAmtX = sdk.ZeroInt()
	totalAmtY = sdk.ZeroInt()

	for _, mr := range matchResultYtoX {
		totalAmtY = totalAmtY.Sub(mr.TransactedCoinAmt)
		totalAmtX = totalAmtX.Add(mr.ExchangedDemandCoinAmt)
	}

	invariantCheckX = invariantCheckX.Add(totalAmtX)
	invariantCheckY = invariantCheckY.Add(totalAmtY)

	invariantCheckX = invariantCheckX.Add(poolXdelta)
	invariantCheckY = invariantCheckY.Add(poolYdelta)

	// print the invariant check and validity with swap, match result
	if invariantCheckX.IsZero() && invariantCheckY.IsZero() {
	} else {
		panic(invariantCheckX)
	}

	if !poolXdelta.Add(decimalErrorX).Equal(poolXdelta2) || !poolYdelta.Add(decimalErrorY).Equal(poolYdelta2) {
		panic(poolXdelta)
	}

	XtoY = types.ValidateStateAndExpireOrders(XtoY, currentHeight, false)
	YtoX = types.ValidateStateAndExpireOrders(YtoX, currentHeight, false)

	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY, true)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	orderBookValidity := types.CheckValidityOrderBook(orderBookExecuted, lastPrice)
	if !orderBookValidity {
		fmt.Println(orderBookValidity, "ErrOrderBookInvalidity")
		panic(types.ErrOrderBookInvalidity)
	}

	// TODO: WIP new validity
	validitySwapPrice := types.CheckSwapPrice(matchResultXtoY, matchResultYtoX, result.SwapPrice)
	if !validitySwapPrice {
		panic("validitySwapPrice")
	}

	// TODO: WIP new validity
	//var validityMustExecutable bool
	//if result.SwapPrice.IsZero() {
	//	validityMustExecutable = types.CheckValidityMustExecutable(orderBookExecuted, lastPrice)
	//} else {
	//	validityMustExecutable = types.CheckValidityMustExecutable(orderBookExecuted, result.SwapPrice)
	//}
	//if !validityMustExecutable {
	//	panic("CheckValidityMustExecutable")
	//}

	// TODO: WIP new validity
	//errThreshold, _ := sdk.NewDecFromStr("0.001")
	//if lastPrice.Quo(result.SwapPrice).Sub(sdk.OneDec()).Abs().GT(errThreshold) {
	//	fmt.Println(lastPrice, lastPrice.Quo(result.SwapPrice.Sub(sdk.OneDec())), lastPrice.Quo(result.SwapPrice.Sub(sdk.OneDec())).Abs())
	//	panic("abs(lastPrice/swapPrice-1) > 0.001")
	//}

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
	if len(matchResultXtoY)+len(matchResultYtoX) != len(matchResultMap) {
		panic("match result map err")
	}

	for k, v := range matchResultMap {
		if k != v.OrderMsgIndex {
			panic("broken map consistency")
		}
	}

	// TODO: separate verify logic, only for simulation
	// compare swapMsgs state with XtoY, YtoX
	notMatchedCount := 0
	for k, v := range matchResultMap {
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
		if msgAfter, ok := matchResultMap[msg.MsgIndex]; ok {
			if msg.MsgIndex == msgAfter.BatchMsg.MsgIndex {
				if *(msg) != *(msgAfter.BatchMsg) || msg != msgAfter.BatchMsg {
					panic("msg not matched")
				} else {
					break
				}
				if !msgAfter.FeeAmt.IsPositive() {
					panic(msgAfter.FeeAmt)
				}
			} else {
				panic("fail msg pointer consistency")
			}
		} else {
			// not matched
			notMatchedCount++
		}
	}
	executedMsgCount := uint64(len(swapMsgs))

	// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapLiquidityPool
	if err := k.TransactAndRefundSwapLiquidityPool(ctx, swapMsgs, matchResultMap, pool); err != nil {
		panic(err)
		return executedMsgCount, err
	}

	//TODO: emit event per msg
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
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
			//match.BatchMsg.RemainingOfferCoin.Sub(
			//		sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))

			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) {
			decimalErrorX = decimalErrorX.Add(sdk.OneInt())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
			//	sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()
			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
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
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
			//	sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) {
			decimalErrorY = decimalErrorY.Add(sdk.OneInt())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.RemainingOfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
			//	sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors one to pool
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()

			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeeded = true
				match.BatchMsg.ToBeDeleted = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = types.CoinSafeSubAmount(match.BatchMsg.RemainingOfferCoin, match.TransactedCoinAmt)
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			matchedIndexMapYtoX[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeeded = true
			match.BatchMsg.ToBeDeleted = false
			fractionalCntY += 1
		}
	}

	poolXdelta = poolXdelta.Add(decimalErrorX)
	poolYdelta = poolYdelta.Add(decimalErrorY)

	X = X.Add(poolXdelta.ToDec())
	Y = Y.Add(poolYdelta.ToDec())

	return XtoY, YtoX, X, Y, poolXdelta, poolYdelta, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY
}
