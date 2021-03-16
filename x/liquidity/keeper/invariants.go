package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// TODO: write invariants, migrate to invariant logics from swap logic
// TODO: reserve coin, batch total result, set last reserve coin and escrow balance, and assert equal with add this batch result
// TODO: remaining orderbook validity check

// register all liquidity invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "escrow-amount",
		LiquidityPoolsEscrowAmountInvariant(k))
}

// AllInvariants runs all invariants of the liquidity module
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := LiquidityPoolsEscrowAmountInvariant(k)(ctx)
		return res, stop
	}
}

// LiquidityPoolsEscrowAmountInvariant checks that outstanding unwithdrawn fees are never negative
func LiquidityPoolsEscrowAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		remainingCoins := sdk.NewCoins()
		batches := k.GetAllPoolBatches(ctx)
		for _, batch := range batches {
			swapMsgs := k.GetAllPoolBatchSwapMsgStatesNotToBeDeleted(ctx, batch)
			for _, msg := range swapMsgs {
				remainingCoins = remainingCoins.Add(msg.RemainingOfferCoin)
			}
			depositMsgs := k.GetAllPoolBatchDepositMsgStatesNotToBeDeleted(ctx, batch)
			for _, msg := range depositMsgs {
				remainingCoins = remainingCoins.Add(msg.Msg.DepositCoins...)
			}
			withdrawMsgs := k.GetAllPoolBatchWithdrawMsgStatesNotToBeDeleted(ctx, batch)
			for _, msg := range withdrawMsgs {
				remainingCoins = remainingCoins.Add(msg.Msg.PoolCoin)
			}
		}

		batchEscrowAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
		escrowAmt := k.bankKeeper.GetAllBalances(ctx, batchEscrowAcc)

		broken := !escrowAmt.IsAllGTE(remainingCoins)

		return sdk.FormatInvariant(types.ModuleName, "batch escrow amount invariant broken",
			"batch escrow amount LT batch remaining amount"), broken
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// These invariants cannot be registered via RegisterInvariants since the module is implemented in batch style.
// We should approach adding these invariant checks in deposit / withdraw / swap batch execution.

var (
	invariantCheckFlag = true // temporary flag for test
)

// MintingPoolCoinsInvariant checks the correct minting amount of pool coins. The difference can be smaller than 1.
func MintingPoolCoinsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// NewPoolTokenAmount / LastPoolTokenSupply = DepositTokenA / LastReserveTokenA
		// NewPoolTokenAmount / LastPoolTokenSupply = DepositTokenB / LastReserveTokenB
		return sdk.FormatInvariant(types.ModuleName, "", ""), false
	}
}

// DepositReserveCoinsInvariant checks the after deposit amounts.
func DepositReserveCoinsInvariant() {
}

// DepositRatioInvariant checks the correct ratio of deposit coin amounts.
func DepositRatioInvariant() {
}

// ImmutablePoolPriceAfterDepositInvariant checks immutable pool price after depositing coins
func ImmutablePoolPriceAfterDepositInvariant() {
}

// BurningPoolCoinsInvariant checks the correct burning amount of pool coins
func BurningPoolCoinsInvariant() {
}

// WithdrawReserveCoinsInvariant checks the after withdraw amounts
func WithdrawReserveCoinsInvariant() {
}

// WithdrawRatioInvariant checks the correct ratio of withdraw coin amounts
func WithdrawRatioInvariant() {
}

// ImmutablePoolPriceAfterWithdrawInvariant checks the immutable pool price after withdrawing coins
func ImmutablePoolPriceAfterWithdrawInvariant() {
}

// SwapPriceInvariants checks the calculated swap price is increased, decreased, or equal from the last pool price
func SwapPriceInvariants(XtoY, YtoX []*types.SwapMsgState, matchResultXtoY, matchResultYtoX []types.MatchResult,
	fractionalCntX, fractionalCntY int, poolXdelta, poolYdelta, poolXdelta2, poolYdelta2, decimalErrorX, decimalErrorY sdk.Dec, result types.BatchResult) {
	beforeXtoYLen := len(XtoY)
	beforeYtoXLen := len(YtoX)
	if beforeXtoYLen-len(matchResultXtoY)+fractionalCntX != types.CountNotMatchedMsgs(XtoY)+types.CountFractionalMatchedMsgs(XtoY) {
		panic(beforeXtoYLen)
	}
	if beforeYtoXLen-len(matchResultYtoX)+fractionalCntY != types.CountNotMatchedMsgs(YtoX)+types.CountFractionalMatchedMsgs(YtoX) {
		panic(beforeYtoXLen)
	}

	totalAmtX := sdk.ZeroDec()
	totalAmtY := sdk.ZeroDec()

	for _, mr := range matchResultXtoY {
		totalAmtX = totalAmtX.Sub(mr.TransactedCoinAmt)
		totalAmtY = totalAmtY.Add(mr.ExchangedDemandCoinAmt)
	}

	invariantCheckX := totalAmtX
	invariantCheckY := totalAmtY

	totalAmtX = sdk.ZeroDec()
	totalAmtY = sdk.ZeroDec()

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

	validitySwapPrice := types.CheckSwapPrice(matchResultXtoY, matchResultYtoX, result.SwapPrice)
	if !validitySwapPrice {
		panic("validitySwapPrice")
	}
}

// OrdersWithNotExecutedStateInvariants checks all executed orders have order price which is not "executable" or not "unexecutable"
func OrdersWithExecutedAndNotExecutedStateInvariants(matchResultXtoY, matchResultYtoX []types.MatchResult, matchResultMap map[uint64]types.MatchResult,
	swapMsgStates []*types.SwapMsgState, XtoY, YtoX []*types.SwapMsgState, result types.BatchResult, currentPoolPrice sdk.Dec, denomX string) {
	if len(matchResultXtoY)+len(matchResultYtoX) != len(matchResultMap) {
		panic("invalid length of match result")
	}

	// compare swapMsgs state with XtoY, YtoX
	for k, v := range matchResultMap {
		if k != v.OrderMsgIndex {
			panic("broken map consistency")
		}
	}

	for _, sms := range swapMsgStates {
		for _, smsXtoY := range XtoY {
			if sms.MsgIndex == smsXtoY.MsgIndex {
				if *(sms) != *(smsXtoY) || sms != smsXtoY {
					panic("swap message state not matched")
				} else {
					break
				}
			}
		}

		for _, smsYtoX := range YtoX {
			if sms.MsgIndex == smsYtoX.MsgIndex {
				if *(sms) != *(smsYtoX) || sms != smsYtoX {
					panic("swap message state not matched")
				} else {
					break
				}
			}
		}

		if msgAfter, ok := matchResultMap[sms.MsgIndex]; ok {
			if sms.MsgIndex == msgAfter.BatchMsg.MsgIndex {
				if *(sms) != *(msgAfter.BatchMsg) || sms != msgAfter.BatchMsg {
					panic("msg not matched")
				} else {
					break
				}
				if !msgAfter.OfferCoinFeeAmt.IsPositive() {
					panic(msgAfter.OfferCoinFeeAmt)
				}
			} else {
				panic("fail msg pointer consistency")
			}
		}
	}

	// checks whether the calculated swapPrice is increased / decreased/ or stayed from the last pool price
	switch result.PriceDirection {
	case types.Increasing:
		if !result.SwapPrice.GTE(currentPoolPrice) {
			panic("invariant check fails due to increase of swap price")
		}
	case types.Decreasing:
		if !result.SwapPrice.LTE(currentPoolPrice) {
			panic("invariant check fails due to decrease of swap price")
		}
	case types.Staying:
		if !result.SwapPrice.Equal(currentPoolPrice) {
			panic("invariant check fails due to stay of swap price")
		}
	}

	// invariant check, execution validity check
	for _, sms := range swapMsgStates {
		if _, ok := matchResultMap[sms.MsgIndex]; ok {
			// checks whether all executed orders have order price which is not "unexecutable"
			if !sms.Executed || !sms.Succeeded {
				panic("swap msg state consistency error, matched but not succeeded")
			}

			if sms.Msg.OfferCoin.Denom == denomX {
				// buy orders having equal or higher order price than found swapPrice
				if !sms.Msg.OrderPrice.GTE(result.SwapPrice) {
					panic("execution validity failed, executed but unexecutable")
				}
			} else {
				// sell orders having equal or lower order price than found swapPrice
				if !sms.Msg.OrderPrice.LTE(result.SwapPrice) {
					panic("execution validity failed, executed but unexecutable")
				}
			}

		} else {
			// check whether every unexecuted orders have order price which is not "executable"
			if sms.Executed && sms.Succeeded {
				panic("sms consistency error, not matched but succeeded")
			}

			if sms.Msg.OfferCoin.Denom == denomX {
				// buy orders having equal or lower order price than found swapPrice
				if !sms.Msg.OrderPrice.LTE(result.SwapPrice) {
					panic("execution validity failed, unexecuted but executable")
				}
			} else {
				// sell orders having equal or higher order price than found swapPrice
				if !sms.Msg.OrderPrice.GTE(result.SwapPrice) {
					panic("execution validity failed, unexecuted but executable")
				}
			}
		}
	}
}
