package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"sort"
)

func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	// get All only not processed swap msgs, not executed, not succeed, not toDelete
	swapMsgs := k.GetAllNotProcessedLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	if len(swapMsgs) == 0 {
		return nil
	}
	pool, found := k.GetLiquidityPool(ctx, liquidityPoolBatch.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}
	// set all msgs to executed
	for _, msg := range swapMsgs {
		msg.Executed = true
	}
	k.SetLiquidityPoolBatchSwapMsgs(ctx, pool.PoolId, swapMsgs)

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
	result := types.ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)
	fmt.Println("batch Result before", result)

	// find order match, calculate pool delta with the total x, y amount for the invariant check
	fmt.Println("before XtoY, YtoX", len(XtoY), len(YtoX))
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
	fmt.Println("lastPrice ", lastPrice)

	fmt.Println("result.SwapPrice, X, Y, currentYPriceOverX", result.SwapPrice, X, Y, currentYPriceOverX)
	if beforeXtoYLen-len(matchResultXtoY)+fractionalCntX != (types.MsgList)(XtoY).CountNotMatchedMsgs()+(types.MsgList)(XtoY).CountFractionalMatchedMsgs() {
		panic(beforeXtoYLen)
	}
	if beforeYtoXLen-len(matchResultYtoX)+fractionalCntY != (types.MsgList)(YtoX).CountNotMatchedMsgs()+(types.MsgList)(YtoX).CountFractionalMatchedMsgs() {
		panic(beforeYtoXLen)
	}

	totalAmtX := sdk.ZeroInt()
	totalAmtY := sdk.ZeroInt()

	for _, mr := range matchResultXtoY {
		fmt.Println("matchResultXtoY", mr)
		totalAmtX = totalAmtX.Sub(mr.TransactedCoinAmt)
		totalAmtY = totalAmtY.Add(mr.ExchangedDemandCoinAmt)
	}

	invariantCheckX := totalAmtX
	invariantCheckY := totalAmtY

	totalAmtX = sdk.ZeroInt()
	totalAmtY = sdk.ZeroInt()

	for _, mr := range matchResultYtoX {
		fmt.Println("matchResultYtoX", mr)
		totalAmtY = totalAmtY.Sub(mr.TransactedCoinAmt)
		totalAmtX = totalAmtX.Add(mr.ExchangedDemandCoinAmt)
	}

	invariantCheckX = invariantCheckX.Add(totalAmtX)
	invariantCheckY = invariantCheckY.Add(totalAmtY)

	invariantCheckX = invariantCheckX.Add(poolXdelta)
	invariantCheckY = invariantCheckY.Add(poolYdelta)

	// print the invariant check and validity with swap, match result
	if invariantCheckX.IsZero() && invariantCheckY.IsZero() {
		fmt.Println("swap execution invariant check: True")
	} else {
		fmt.Println("swap execution invariant check: False", invariantCheckX, invariantCheckY)
		panic(invariantCheckX)
	}

	if result.MatchType == 1 {
		fmt.Println("matchType: ", "ExactMatch")
	} else if result.MatchType == 2 {
		fmt.Println("matchType: ", "No Match")
	} else if result.MatchType == 3 {
		fmt.Println("matchType: ", "FractionalMatch")
	}

	fmt.Println("swapPrice: ", result.SwapPrice)
	fmt.Println("matchResultXtoY: ", matchResultXtoY)
	fmt.Println("matchResultYtoX: ", matchResultYtoX)
	fmt.Println("matched totalAmtX, totalAmtY", totalAmtX, totalAmtY)
	fmt.Println("poolXdelta, poolYdelta", poolXdelta, poolYdelta, poolXdelta2, poolYdelta2)

	if !poolXdelta.Add(decimalErrorX).Equal(poolXdelta2) || !poolYdelta.Add(decimalErrorY).Equal(poolYdelta2) {
		panic(poolXdelta)
	}

	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY, true)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	fmt.Println("orderbook after batch")
	orderBookValidity := types.CheckValidityOrderBook(orderBookExecuted, lastPrice)
	fmt.Println("orderBookValidity:", orderBookValidity)
	if !orderBookValidity {
		fmt.Println(orderBookValidity, "ErrOrderBookInvalidity")
		for _, v := range orderBookExecuted {
			fmt.Println(v)
		}
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
	// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapLiquidityPool
	if err := k.TransactAndRefundSwapLiquidityPool(ctx, swapMsgs, matchResultMap, pool); err != nil {
		panic(err)
		return err
	}

	//TODO: emit event per msg
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
	return nil
}

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
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
			decimalErrorX = decimalErrorX.Add(sdk.OneInt())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors to pool
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()
			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			matchedIndexMapXtoY[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeed = true
			match.BatchMsg.ToDelete = false
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
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
			decimalErrorY = decimalErrorY.Add(sdk.OneInt())
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors one to pool
			if match.BatchMsg.RemainingOfferCoin.Amount.Equal(sdk.OneInt()) {
				match.BatchMsg.RemainingOfferCoin.Amount = sdk.ZeroInt()

			}
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
		} else {
			// fractional match
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			matchedIndexMapYtoX[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeed = true
			match.BatchMsg.ToDelete = false
			fractionalCntY += 1
		}
	}

	poolXdelta = poolXdelta.Add(decimalErrorX)
	poolYdelta = poolYdelta.Add(decimalErrorY)

	X = X.Add(poolXdelta.ToDec())
	Y = Y.Add(poolYdelta.ToDec())

	return XtoY, YtoX, X, Y, poolXdelta, poolYdelta, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY
}
