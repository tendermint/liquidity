package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	pool, found := k.GetLiquidityPool(ctx, liquidityPoolBatch.PoolID)
	if !found {
		return types.ErrPoolNotExists
	}
	//totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	X := reserveCoins[0].Amount.ToDec()
	Y := reserveCoins[1].Amount.ToDec()
	currentYPriceOverX := X.Quo(Y)

	denomX := reserveCoins[0].Denom
	denomY := reserveCoins[1].Denom

	swapMsgs := k.GetAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	orderMap, XtoY, YtoX := types.GetOrderMap(swapMsgs, denomX, denomY)

	// make orderbook to sort orderMap
	orderBook := orderMap.SortOrderBook()

	fmt.Println("orderbook before batch")
	orderBookValidity := types.CheckValidityOrderBook(orderBook, currentYPriceOverX)

	result := types.ComputePriceDirection(X, Y, currentYPriceOverX, orderBook)

	fmt.Println("priceDirection: ", result)

	params := k.GetParams(ctx)
	fmt.Println("before XtoY, YtoX", len(XtoY), len(YtoX))
	matchResultXtoY, XtoY, poolXDeltaXtoY, poolYDeltaXtoY := types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX, result.SwapPrice, params.SwapFeeRate, ctx.BlockHeight())
	matchResultYtoX, YtoX, poolXDeltaYtoX, poolYDeltaYtoX := types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY, result.SwapPrice, params.SwapFeeRate, ctx.BlockHeight())
	poolXdelta := poolXDeltaXtoY.Add(poolXDeltaYtoX)
	poolYdelta := poolYDeltaXtoY.Add(poolYDeltaYtoX)
	fmt.Println(result, matchResultXtoY, matchResultYtoX, poolXdelta, poolYdelta)
	fmt.Println(result.SwapPrice, X, Y, currentYPriceOverX)
	fmt.Println("after XtoY, YtoX", len(XtoY), len(YtoX), len(matchResultXtoY), len(matchResultYtoX))

	totalAmtX := sdk.ZeroInt()
	totalAmtY := sdk.ZeroInt()

	for _, mr := range matchResultXtoY {
		totalAmtX = totalAmtX.Sub(mr.MatchedAmt)
		totalAmtY = totalAmtY.Add(mr.ReceiveAmt)
	}
	fmt.Println("totalAmtX, totalAmtY", totalAmtX, totalAmtY)

	invariantCheckX := totalAmtX
	invariantCheckY := totalAmtY

	totalAmtX = sdk.ZeroInt()
	totalAmtY = sdk.ZeroInt()

	for _, mr := range matchResultYtoX {
		totalAmtY = totalAmtY.Sub(mr.MatchedAmt)
		totalAmtX = totalAmtX.Add(mr.ReceiveAmt)
	}
	fmt.Println("X, Y, poolXdelta, poolYdelta", totalAmtX, totalAmtY, poolXdelta, poolYdelta, invariantCheckX, invariantCheckY)

	invariantCheckX = invariantCheckX.Add(totalAmtX)
	invariantCheckY = invariantCheckY.Add(totalAmtY)

	invariantCheckX = invariantCheckX.Add(poolXdelta)
	invariantCheckY = invariantCheckY.Add(poolYdelta)

	if invariantCheckX.IsZero() && invariantCheckY.IsZero() {
		fmt.Println("swap execution invariant check: True")
	} else {
		fmt.Println("swap execution invariant check: False", invariantCheckX, invariantCheckY)
	}

	// TODO: updateState, cancelEndOfLifeSpanOrders
	XtoY, YtoX = types.ClearOrders(XtoY, YtoX)

	orderMapExecuted, _, _ := types.GetOrderMap(append(XtoY, YtoX...), denomX, denomY)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	fmt.Println("orderbook after batch")
	orderBookValidity = types.CheckValidityOrderBook(orderBookExecuted, currentYPriceOverX)
	fmt.Println("after orderBookValidity", orderBookValidity)
	if !orderBookValidity {
		fmt.Println(orderBookValidity, "ErrOrderBookInvalidity", orderBookExecuted)
		//return types.ErrOrderBookInvalidity
	}

	// TODO: updateState with escrow, emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
	return nil
}

