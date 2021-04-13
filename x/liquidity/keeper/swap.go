package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Execute Swap of the pool batch, Collect swap messages in batch for transact the same price for each batch and run them on endblock.
func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.PoolBatch) (uint64, error) {
	// get all swap message batch states that are not executed, not succeeded, and not to be deleted.
	swapMsgStates := k.GetAllNotProcessedPoolBatchSwapMsgStates(ctx, liquidityPoolBatch)
	if len(swapMsgStates) == 0 {
		return 0, nil
	}

	pool, found := k.GetPool(ctx, liquidityPoolBatch.PoolId)
	if !found {
		return 0, types.ErrPoolNotExists
	}

	// set executed states of all messages to true
	for _, sms := range swapMsgStates {
		sms.Executed = true
	}
	k.SetPoolBatchSwapMsgStatesByPointer(ctx, pool.Id, swapMsgStates)

	currentHeight := ctx.BlockHeight()
	types.ValidateStateAndExpireOrders(swapMsgStates, currentHeight, false)

	// get reserve coins from the liquidity pool and calculate the current pool price (p = x / y)
	reserveCoins := k.GetReserveCoins(ctx, pool)

	X := reserveCoins[0].Amount.ToDec()
	Y := reserveCoins[1].Amount.ToDec()
	currentPoolPrice := X.Quo(Y)
	denomX := reserveCoins[0].Denom
	denomY := reserveCoins[1].Denom

	// make orderMap, orderbook by sort orderMap
	orderMap, XtoY, YtoX := types.MakeOrderMap(swapMsgStates, denomX, denomY, false)
	orderBook := orderMap.SortOrderBook()

	// check orderbook validity and compute batchResult(direction, swapPrice, ..)
	result, found := orderBook.Match(X, Y)

	executedMsgCount := uint64(len(swapMsgStates))

	if !found {
		err := k.RefundSwaps(ctx, pool, swapMsgStates)
		return executedMsgCount, err
	}

	// find order match, calculate pool delta with the total x, y amounts for the invariant check
	var matchResultXtoY, matchResultYtoX []types.MatchResult

	poolXDelta := sdk.ZeroDec()
	poolYDelta := sdk.ZeroDec()

	if result.MatchType != types.NoMatch {
		var poolXDeltaXtoY, poolXDeltaYtoX, poolYDeltaYtoX, poolYDeltaXtoY sdk.Dec
		matchResultXtoY, poolXDeltaXtoY, poolYDeltaXtoY = types.FindOrderMatch(types.DirectionXtoY, XtoY, result.EX, result.SwapPrice, currentHeight)
		matchResultYtoX, poolXDeltaYtoX, poolYDeltaYtoX = types.FindOrderMatch(types.DirectionYtoX, YtoX, result.EY, result.SwapPrice, currentHeight)
		poolXDelta = poolXDeltaXtoY.Add(poolXDeltaYtoX)
		poolYDelta = poolYDeltaXtoY.Add(poolYDeltaYtoX)
	}

	XtoY, YtoX, X, Y, poolXDelta2, poolYDelta2, decimalErrorX, decimalErrorY := types.UpdateSwapMsgStates(X, Y, XtoY, YtoX, matchResultXtoY, matchResultYtoX)

	lastPrice := X.Quo(Y)

	if invariantCheckFlag {
		SwapMatchingInvariants(XtoY, YtoX, matchResultXtoY, matchResultYtoX)
		SwapPriceInvariants(matchResultXtoY, matchResultYtoX, poolXDelta, poolYDelta, poolXDelta2, poolYDelta2, decimalErrorX, decimalErrorY, result)
	}

	types.ValidateStateAndExpireOrders(XtoY, currentHeight, false)
	types.ValidateStateAndExpireOrders(YtoX, currentHeight, false)

	orderMapExecuted, _, _ := types.MakeOrderMap(append(XtoY, YtoX...), denomX, denomY, true)
	orderBookExecuted := orderMapExecuted.SortOrderBook()
	if !orderBookExecuted.Validate(lastPrice) {
		return executedMsgCount, types.ErrOrderBookInvalidity
	}

	types.ValidateStateAndExpireOrders(XtoY, currentHeight, true)
	types.ValidateStateAndExpireOrders(YtoX, currentHeight, true)

	// make index map for match result
	matchResultMap := make(map[uint64]types.MatchResult)
	for _, match := range append(matchResultXtoY, matchResultYtoX...) {
		if _, ok := matchResultMap[match.SwapMsgState.MsgIndex]; ok {
			return executedMsgCount, fmt.Errorf("duplicate match order")
		}
		matchResultMap[match.SwapMsgState.MsgIndex] = match
	}

	if invariantCheckFlag {
		SwapPriceDirection(currentPoolPrice, result)
		SwapMsgStatesInvariants(matchResultXtoY, matchResultYtoX, matchResultMap, swapMsgStates, XtoY, YtoX)
		SwapOrdersExecutionStateInvariants(matchResultMap, swapMsgStates, result, denomX)
	}

	// execute transact, refund, expire, send coins with escrow, update state by TransactAndRefundSwapLiquidityPool
	if err := k.TransactAndRefundSwapLiquidityPool(ctx, swapMsgStates, matchResultMap, pool, result); err != nil {
		return executedMsgCount, err
	}

	return executedMsgCount, nil
}
