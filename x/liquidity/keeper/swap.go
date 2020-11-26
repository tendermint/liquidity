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

	// TODO: get all past queued orders
	// get All swap msgs
	swapMsgs := k.GetAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	if len(swapMsgs) == 0 {
		return nil
	}

	// TODO: add validate MsgSwap

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
	orderBookValidity := types.CheckValidityOrderBook(orderBook, currentYPriceOverX)
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
			result.SwapPrice, params.SwapFeeRate, ctx.BlockHeight())
		matchResultYtoX, _, poolXDeltaYtoX, poolYDeltaYtoX = types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY,
			result.SwapPrice, params.SwapFeeRate, ctx.BlockHeight())
		poolXdelta = poolXDeltaXtoY.Add(poolXDeltaYtoX)
		poolYdelta = poolYDeltaXtoY.Add(poolYDeltaYtoX)
	}

	XtoY, YtoX, X, Y, poolXdelta2, poolYdelta2, fractionalCntX, fractionalCntY, decimalErrorX, decimalErrorY :=
		k.UpdateState(X, Y, XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	lastPrice := X.Quo(Y)
	fmt.Println("lastPrice ", lastPrice)

	fmt.Println("result.SwapPrice, X, Y, currentYPriceOverX", result.SwapPrice, X, Y, currentYPriceOverX)
	if beforeXtoYLen-len(matchResultXtoY)+fractionalCntX != len(XtoY) {
		fmt.Println("!! match invariant Fail X")
		panic(beforeXtoYLen)
	}
	if beforeYtoXLen-len(matchResultYtoX)+fractionalCntY != len(YtoX) {
		fmt.Println("!! match invariant Fail Y")
		panic(beforeYtoXLen)
	}

	totalAmtX := sdk.ZeroInt()
	totalAmtY := sdk.ZeroInt()

	for _, mr := range matchResultXtoY {
		fmt.Println("matchResultXtoY", mr)
		totalAmtX = totalAmtX.Sub(mr.TransactedCoinAmt)
		totalAmtY = totalAmtY.Add(mr.ExchangedCoinAmt)
	}

	invariantCheckX := totalAmtX
	invariantCheckY := totalAmtY

	totalAmtX = sdk.ZeroInt()
	totalAmtY = sdk.ZeroInt()

	for _, mr := range matchResultYtoX {
		fmt.Println("matchResultYtoX", mr)
		totalAmtY = totalAmtY.Sub(mr.TransactedCoinAmt)
		totalAmtX = totalAmtX.Add(mr.ExchangedCoinAmt)
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

	XtoY, YtoX = types.ClearOrders(XtoY, YtoX)

	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	fmt.Println("orderbook after batch")
	orderBookValidity = types.CheckValidityOrderBook(orderBookExecuted, lastPrice)
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
		for _, order := range XtoY {
			if match.OrderMsgIndex == order.MsgIndex {
				poolXdelta = poolXdelta.Add(match.TransactedCoinAmt)
				poolYdelta = poolYdelta.Sub(match.ExchangedCoinAmt)
				if order.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt) {
					// full match
					matchedOrderMsgIndexListXtoY = append(matchedOrderMsgIndexListXtoY, order.MsgIndex)
				} else if order.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
					decimalErrorX = decimalErrorX.Add(sdk.OneInt())
					//poolXdelta = poolXdelta.Add(sdk.OneInt())
					matchedOrderMsgIndexListXtoY = append(matchedOrderMsgIndexListXtoY, order.MsgIndex)
				} else {
					// fractional match
					order.Msg.OfferCoin = order.Msg.OfferCoin.Sub(sdk.NewCoin(order.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
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
				poolXdelta = poolXdelta.Sub(match.ExchangedCoinAmt)
				poolYdelta = poolYdelta.Add(match.TransactedCoinAmt)
				if order.Msg.OfferCoin.Amount.Equal(match.TransactedCoinAmt) {
					// full match
					matchedOrderMsgIndexListYtoX = append(matchedOrderMsgIndexListYtoX, order.MsgIndex)
				} else if order.Msg.OfferCoin.Amount.Sub(match.TransactedCoinAmt).Equal(sdk.OneInt()) { // TODO: need to verify logic
					decimalErrorY = decimalErrorY.Add(sdk.OneInt())
					//poolYdelta = poolYdelta.Add(sdk.OneInt())
					matchedOrderMsgIndexListYtoX = append(matchedOrderMsgIndexListYtoX, order.MsgIndex)
				} else {
					// fractional match
					order.Msg.OfferCoin = order.Msg.OfferCoin.Sub(sdk.NewCoin(order.Msg.OfferCoin.Denom, match.TransactedCoinAmt))
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

