package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"sort"
)

func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	params := k.GetParams(ctx)
	pool, found := k.GetLiquidityPool(ctx, liquidityPoolBatch.PoolId)
	if !found {
		return types.ErrPoolNotExists
	}
	currentHeight := ctx.BlockHeight()

	// get All only not processed swap msgs, not executed, not succeed, not toDelete
	swapMsgs := k.GetAllNotProcessedPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	if len(swapMsgs) == 0 {
		return nil
	}

	// TODO: add validate MsgSwap
	// set all msgs to executed
	for _, msg := range swapMsgs {
		msg.Executed = true
	}
	k.SetLiquidityPoolBatchSwapMsgs(ctx, pool.PoolId, swapMsgs)

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
	orderMap, XtoY, YtoX := types.GetOrderMap(swapMsgs, denomX, denomY)
	orderBook := orderMap.SortOrderBook()

	// check orderbook validity and compute batchResult(direction, swapPrice, ..)
	fmt.Println("orderbook before batch")
	//orderBookValidity := types.CheckValidityOrderBook(orderBook, currentYPriceOverX)
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
	//XtoY, YtoX = types.ClearOrders(XtoY, YtoX, currentHeight, false)

	fmt.Println("result.SwapPrice, X, Y, currentYPriceOverX", result.SwapPrice, X, Y, currentYPriceOverX)
	if beforeXtoYLen-len(matchResultXtoY)+fractionalCntX != (types.MsgList)(XtoY).LenRemainingMsgs()+(types.MsgList)(XtoY).LenFractionalMsgs() {
		fmt.Println("!! match invariant Fail X")
		fmt.Println(beforeXtoYLen-len(matchResultXtoY)+fractionalCntX, (types.MsgList)(XtoY).LenRemainingMsgs(), (types.MsgList)(XtoY).LenFractionalMsgs())
		fmt.Println(beforeXtoYLen, len(matchResultXtoY), fractionalCntX, (types.MsgList)(XtoY).LenRemainingMsgs(), len(XtoY))
		fmt.Println(XtoY)
		panic(beforeXtoYLen)
	}
	if beforeYtoXLen-len(matchResultYtoX)+fractionalCntY != (types.MsgList)(YtoX).LenRemainingMsgs()+(types.MsgList)(YtoX).LenFractionalMsgs() {
		fmt.Println("!! match invariant Fail Y")
		fmt.Println(beforeYtoXLen-len(matchResultYtoX)+fractionalCntY, (types.MsgList)(YtoX).LenRemainingMsgs())
		fmt.Println(beforeYtoXLen, len(matchResultYtoX), fractionalCntY, (types.MsgList)(YtoX).LenRemainingMsgs(), len(YtoX))
		fmt.Println(YtoX)
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
		// TODO: verify after batch
		panic(poolXdelta)
	}

	XtoY, YtoX = types.ClearOrders(XtoY, YtoX, currentHeight, true)
	if (types.MsgList)(XtoY).LenRemainingMsgs() != len(XtoY) {
		fmt.Println((types.MsgList)(XtoY).LenRemainingMsgs(), len(XtoY), (types.MsgList)(XtoY).LenFractionalMsgs())
		panic("not clear orders invariants")
	}
	if (types.MsgList)(YtoX).LenRemainingMsgs() != len(YtoX) {
		fmt.Println((types.MsgList)(YtoX).LenRemainingMsgs(), len(YtoX), (types.MsgList)(YtoX).LenFractionalMsgs())
		panic("not clear orders invariants")
	}

	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	fmt.Println("orderbook after batch")
	orderBookValidity := types.CheckValidityOrderBook(orderBookExecuted, lastPrice)
	fmt.Println("after orderBookValidity", orderBookValidity)
	if !orderBookValidity {
		fmt.Println(orderBookValidity, "ErrOrderBookInvalidity", orderBookExecuted)
		panic(orderBookValidity)
	}

	// TODO: updateState, KV Set, with escrow, emit event
	// TODO: check order expiry height, set toDelete flag
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
	return nil
}

// TODO: keeper, err, set kv, test code
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
	var matchedOrderMsgIndexListXtoY []uint64
	var matchedOrderMsgIndexListYtoX []uint64
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
			// TODO: verify set batch msg
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.Msg.OfferCoin.Sub(match.BatchMsg.ExchangedOfferCoin)
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				// TODO: add verify batchSwapMsg
				fmt.Println(match)
				fmt.Println(match.BatchMsg.RemainingOfferCoin, match.BatchMsg.ExchangedOfferCoin.Amount, match.BatchMsg.Msg.OfferCoin)
				fmt.Println(match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
					GT(match.BatchMsg.Msg.OfferCoin.Amount),
					match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())))
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
			matchedOrderMsgIndexListXtoY = append(matchedOrderMsgIndexListXtoY, match.BatchMsg.MsgIndex)
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
			decimalErrorX = decimalErrorX.Add(sdk.OneInt())
			//poolXdelta = poolXdelta.Add(sdk.OneInt())
			// TODO: verify set batch msg
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				// TODO: add verify batchSwapMsg
				fmt.Println(match)
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
			matchedOrderMsgIndexListXtoY = append(matchedOrderMsgIndexListXtoY, match.BatchMsg.MsgIndex)
		} else {
			// fractional match
			// TODO: verify msg edit
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.Msg.OfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			matchedIndexMapXtoY[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeed = true
			match.BatchMsg.ToDelete = false
			fractionalCntX += 1
		}
		// TODO: check break logic
		//break
	}
	if len(matchedOrderMsgIndexListXtoY) > 0 {
		newI := 0
		for _, order := range XtoY {
			if _, ok := matchedIndexMapXtoY[order.MsgIndex]; ok {
				//order.Succeed = true
				// already updated RemainingOfferCoin
				//order.Msg.OfferCoin = val
				// TODO: verify batch
			}
			removeFlag := false
			for _, i := range matchedOrderMsgIndexListXtoY {
				if i == order.MsgIndex {
					//order.Succeed = true
					//order.ToDelete = true
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
		poolXdelta = poolXdelta.Sub(match.ExchangedDemandCoinAmt)
		poolYdelta = poolYdelta.Add(match.TransactedCoinAmt)
		if match.BatchMsg.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Equal(match.TransactedCoinAmt) {
			// full match
			// TODO: verify set batch msg
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.Msg.OfferCoin.Sub(match.BatchMsg.ExchangedOfferCoin)
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				// TODO: add verify batchSwapMsg
				fmt.Println(match)
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
			matchedOrderMsgIndexListYtoX = append(matchedOrderMsgIndexListYtoX, match.BatchMsg.MsgIndex)
		} else if match.BatchMsg.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) ||
			match.BatchMsg.RemainingOfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
			decimalErrorY = decimalErrorY.Add(sdk.OneInt())
			//poolXdelta = poolXdelta.Add(sdk.OneInt())
			// TODO: verify set batch msg
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(
				sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			// TODO: verify RemainingOfferCoin about deciaml errors
			if match.BatchMsg.RemainingOfferCoin.Amount.Add(match.BatchMsg.ExchangedOfferCoin.Amount).
				GT(match.BatchMsg.Msg.OfferCoin.Amount) ||
				!match.BatchMsg.RemainingOfferCoin.Equal(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, sdk.ZeroInt())) {
				// TODO: add verify batchSwapMsg
				fmt.Println(match)
				panic("remaining not matched")
			} else {
				match.BatchMsg.Succeed = true
				match.BatchMsg.ToDelete = true
			}
			matchedOrderMsgIndexListYtoX = append(matchedOrderMsgIndexListYtoX, match.BatchMsg.MsgIndex)
		} else {
			// fractional match
			// TODO: verify msg edit
			match.BatchMsg.ExchangedOfferCoin = match.BatchMsg.ExchangedOfferCoin.Add(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			match.BatchMsg.RemainingOfferCoin = match.BatchMsg.RemainingOfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			//match.BatchMsg.RemainingOfferCoin = match.BatchMsg.Msg.OfferCoin.Sub(sdk.NewCoin(match.BatchMsg.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
			matchedIndexMapYtoX[match.BatchMsg.MsgIndex] = match.BatchMsg.RemainingOfferCoin
			match.BatchMsg.Succeed = true
			match.BatchMsg.ToDelete = false
			fractionalCntY += 1
		}
		// TODO: check break logic
		//break
	}
	if len(matchedOrderMsgIndexListYtoX) > 0 {
		// TODO: set toDelete without delete
		newI := 0
		for _, order := range YtoX {
			if _, ok := matchedIndexMapYtoX[order.MsgIndex]; ok {
				//order.Succeed = true
				// already updated RemainingOfferCoin
				//order.Msg.OfferCoin = val
				// TODO: verify batch
			}
			removeFlag := false
			for _, i := range matchedOrderMsgIndexListYtoX {
				if i == order.MsgIndex {
					//order.Succeed = true
					//order.ToDelete = true
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

