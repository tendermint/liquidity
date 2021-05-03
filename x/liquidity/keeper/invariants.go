package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// RegisterInvariants registers all liquidity invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "escrow-amount",
		LiquidityPoolsEscrowAmountInvariant(k))
}

// AllInvariants runs all invariants of the liquidity module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := LiquidityPoolsEscrowAmountInvariant(k)(ctx)
		return res, stop
	}
}

// LiquidityPoolsEscrowAmountInvariant checks that outstanding unwithdrawn fees are never negative.
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
// These invariants cannot be registered via RegisterInvariants since the module uses per-block batch execution.
// We should approach adding these invariant checks inside actual logics of deposit / withdraw / swap.

var (
	invariantCheckFlag = true                     // TODO: better way to handle below invariant checks?
	diffThreshold      = sdk.NewDecWithPrec(2, 1) // 20%
)

func diff(a, b sdk.Dec) sdk.Dec {
	return a.Sub(b).Abs().Quo(b)
}

// MintingPoolCoinsInvariant checks the correct ratio of minting amount of pool coins.
func MintingPoolCoinsInvariant(poolCoinTotalSupply, mintPoolCoin, depositCoinA, depositCoinB, lastReserveCoinA, lastReserveCoinB, refundedCoinA, refundedCoinB sdk.Dec) {
	if !refundedCoinA.IsZero() {
		depositCoinA = depositCoinA.Sub(refundedCoinA)
	}

	if !refundedCoinB.IsZero() {
		depositCoinB = depositCoinB.Sub(refundedCoinB)
	}

	poolCoinRatio := mintPoolCoin.Quo(poolCoinTotalSupply)
	depositCoinARatio := depositCoinA.Quo(lastReserveCoinA)
	depositCoinBRatio := depositCoinB.Quo(lastReserveCoinB)

	// NewPoolCoinAmount / LastPoolCoinSupply <= AfterRefundedDepositCoinA / LastReserveCoinA
	// NewPoolCoinAmount / LastPoolCoinSupply <= AfterRefundedDepositCoinA / LastReserveCoinB
	if depositCoinARatio.LT(poolCoinRatio) || depositCoinBRatio.LT(poolCoinRatio) {
		panic("invariant check fails due to incorrect ratio of pool coins")
	}
}

// DepositReserveCoinsInvariant checks after deposit amounts.
func DepositReserveCoinsInvariant(lastReserveCoinA, lastReserveCoinB, depositCoinA, depositCoinB, afterReserveCoinA, afterReserveCoinB, refundedCoinA, refundedCoinB sdk.Dec) {
	if !refundedCoinA.IsZero() {
		depositCoinA = depositCoinA.Sub(refundedCoinA)
	}

	if !refundedCoinB.IsZero() {
		depositCoinB = depositCoinB.Sub(refundedCoinB)
	}

	// AfterDepositReserveCoinA = LastReserveCoinA + AfterRefundedDepositCoinA
	// AfterDepositReserveCoinB = LastReserveCoinB + AfterRefundedDepositCoinA
	if !afterReserveCoinA.Equal(lastReserveCoinA.Add(depositCoinA)) ||
		!afterReserveCoinB.Equal(lastReserveCoinB.Add(depositCoinB)) {
		panic("invariant check fails due to incorrect deposit amounts")
	}
}

// DepositRatioInvariant checks the correct ratio of deposit coin amounts.
func DepositRatioInvariant(depositCoinA, depositCoinB, refundedCoinA, refundedCoinB, lastReserveCoinA, lastReserveCoinB sdk.Dec) {
	if !refundedCoinA.IsZero() {
		depositCoinA = depositCoinA.Sub(refundedCoinA)
	}

	if !refundedCoinB.IsZero() {
		depositCoinB = depositCoinB.Sub(refundedCoinB)
	}

	depositCoinRatio := depositCoinA.Quo(depositCoinB)
	lastReserveCoinRatio := lastReserveCoinA.Quo(lastReserveCoinB)

	// AfterRefundedDepositCoinA / AfterRefundedDepositCoinA = LastReserveCoinA / LastReserveCoinB
	if depositCoinA.GTE(sdk.NewDec(1000)) && depositCoinB.GTE(sdk.NewDec(1000)) &&
		lastReserveCoinA.GTE(sdk.NewDec(1000)) && lastReserveCoinB.GTE(sdk.NewDec(1000)) &&
		diff(depositCoinRatio, lastReserveCoinRatio).GT(diffThreshold) {
		panic("invariant check fails due to incorrect deposit ratio")
	}
}

// ImmutablePoolPriceAfterDepositInvariant checks immutable pool price after depositing coins.
func ImmutablePoolPriceAfterDepositInvariant(lastReserveCoinRatio, afterReserveCoinRatio sdk.Dec) {
	// LastReserveCoinA / LastReserveCoinB = AfterDepositReserveCoinA / AfterDepositReserveCoinB
	if diff(lastReserveCoinRatio, afterReserveCoinRatio).GT(diffThreshold) {
		panic("invariant check fails due to incorrect pool price ratio")
	}
}

// BurningPoolCoinsInvariant checks the correct burning amount of pool coins.
func BurningPoolCoinsInvariant(burnedPoolCoin, withdrawCoinA, withdrawCoinB, reserveCoinA, reserveCoinB, lastPoolCoinSupply, withdrawProportion sdk.Dec) {
	burningPoolCoinRatio := burnedPoolCoin.Quo(lastPoolCoinSupply)
	if burningPoolCoinRatio.Equal(sdk.OneDec()) {
		return
	}

	withdrawCoinARatio := withdrawCoinA.Quo(withdrawProportion).Quo(reserveCoinA)
	withdrawCoinBRatio := withdrawCoinB.Quo(withdrawProportion).Quo(reserveCoinB)

	// BurnedPoolCoinAmount / LastPoolCoinSupply >= (WithdrawCoinA+WithdrawFeeCoinA) / LastReserveCoinA
	// BurnedPoolCoinAmount / LastPoolCoinSupply >= (WithdrawCoinB+WithdrawFeeCoinB) / LastReserveCoinB
	if withdrawCoinARatio.GT(burningPoolCoinRatio) || withdrawCoinBRatio.GT(burningPoolCoinRatio) {
		panic("invariant check fails due to incorrect ratio of burning pool coins")
	}
}

// WithdrawReserveCoinsInvariant checks the after withdraw amounts.
func WithdrawReserveCoinsInvariant(withdrawCoinA, withdrawCoinB, reserveCoinA, reserveCoinB,
	afterReserveCoinA, afterReserveCoinB, afterPoolCoinTotalSupply, lastPoolCoinSupply, burnedPoolCoin sdk.Dec) {
	// AfterWithdrawReserveCoinA = LastReserveCoinA - WithdrawCoinA
	if !afterReserveCoinA.Equal(reserveCoinA.Sub(withdrawCoinA)) {
		panic("invariant check fails due to incorrect withdraw coin A amount")
	}

	// AfterWithdrawReserveCoinB = LastReserveCoinB - WithdrawCoinB
	if !afterReserveCoinB.Equal(reserveCoinB.Sub(withdrawCoinB)) {
		panic("invariant check fails due to incorrect withdraw coin B amount")
	}

	// AfterWithdrawPoolCoinSupply = LastPoolCoinSupply - BurnedPoolCoinAmount
	if !afterPoolCoinTotalSupply.Equal(lastPoolCoinSupply.Sub(burnedPoolCoin)) {
		panic("invariant check fails due to incorrect total supply")
	}
}

// WithdrawAmountInvariant checks the correct ratio of withdraw coin amounts.
func WithdrawAmountInvariant(withdrawCoinA, withdrawCoinB, reserveCoinA, reserveCoinB, burnedPoolCoin, poolCoinSupply, withdrawFeeRate sdk.Dec) {
	ratio := burnedPoolCoin.Quo(poolCoinSupply).Mul(sdk.OneDec().Sub(withdrawFeeRate))
	idealWithdrawCoinA := reserveCoinA.Mul(ratio)
	idealWithdrawCoinB := reserveCoinB.Mul(ratio)
	diffA := idealWithdrawCoinA.Sub(withdrawCoinA).Abs()
	diffB := idealWithdrawCoinB.Sub(withdrawCoinB).Abs()
	if !burnedPoolCoin.Equal(poolCoinSupply) {
		if diffA.GTE(sdk.OneDec()) {
			panic(fmt.Sprintf("withdraw coin amount %v differs too much from %v", withdrawCoinA, idealWithdrawCoinA))
		}
		if diffB.GTE(sdk.OneDec()) {
			panic(fmt.Sprintf("withdraw coin amount %v differs too much from %v", withdrawCoinB, idealWithdrawCoinB))
		}
	}
}

// TODO: add invariant check for withdrawed coin ratio

// ImmutablePoolPriceAfterWithdrawInvariant checks the immutable pool price after withdrawing coins.
func ImmutablePoolPriceAfterWithdrawInvariant(reserveCoinA, reserveCoinB, withdrawCoinA, withdrawCoinB, afterReserveCoinA, afterReserveCoinB sdk.Dec) {
	// TestReinitializePool tests a scenario where after reserve coins are zero
	if !afterReserveCoinA.IsZero() && !afterReserveCoinB.IsZero() {
		reserveCoinA = reserveCoinA.Sub(withdrawCoinA)
		reserveCoinB = reserveCoinB.Sub(withdrawCoinB)

		reserveCoinRatio := reserveCoinA.Quo(reserveCoinB)
		afterReserveCoinRatio := afterReserveCoinA.Quo(afterReserveCoinB)

		// LastReserveCoinA / LastReserveCoinB = AfterWithdrawReserveCoinA / AfterWithdrawReserveCoinB
		if reserveCoinA.GTE(sdk.NewDec(1000)) && reserveCoinB.GTE(sdk.NewDec(1000)) &&
			withdrawCoinA.GTE(sdk.NewDec(1000)) && withdrawCoinB.GTE(sdk.NewDec(1000)) &&
			diff(reserveCoinRatio, afterReserveCoinRatio).GT(diffThreshold) {
			panic("invariant check fails due to incorrect pool price ratio")
		}
	}
}

// SwapMatchingInvariants checks swap matching results of both X to Y and Y to X cases.
func SwapMatchingInvariants(XtoY, YtoX []*types.SwapMsgState, matchResultXtoY, matchResultYtoX []types.MatchResult) {
	beforeMatchingXtoYLen := len(XtoY)
	beforeMatchingYtoXLen := len(YtoX)
	afterMatchingXtoYLen := len(matchResultXtoY)
	afterMatchingYtoXLen := len(matchResultYtoX)

	notMatchedXtoYLen := beforeMatchingXtoYLen - afterMatchingXtoYLen
	notMatchedYtoXLen := beforeMatchingYtoXLen - afterMatchingYtoXLen

	if notMatchedXtoYLen != types.CountNotMatchedMsgs(XtoY) {
		panic("invariant check fails due to invalid XtoY match length")
	}

	if notMatchedYtoXLen != types.CountNotMatchedMsgs(YtoX) {
		panic("invariant check fails due to invalid YtoX match length")
	}
}

// SwapPriceInvariants checks swap price invariants.
func SwapPriceInvariants(matchResultXtoY, matchResultYtoX []types.MatchResult, poolXDelta, poolYDelta, poolXDelta2, poolYDelta2,
	decimalErrorX, decimalErrorY sdk.Dec, result types.BatchResult) {
	invariantCheckX := sdk.ZeroDec()
	invariantCheckY := sdk.ZeroDec()

	for _, m := range matchResultXtoY {
		invariantCheckX = invariantCheckX.Sub(m.TransactedCoinAmt)
		invariantCheckY = invariantCheckY.Add(m.ExchangedDemandCoinAmt)
	}

	for _, m := range matchResultYtoX {
		invariantCheckY = invariantCheckY.Sub(m.TransactedCoinAmt)
		invariantCheckX = invariantCheckX.Add(m.ExchangedDemandCoinAmt)
	}

	invariantCheckX = invariantCheckX.Add(poolXDelta2)
	invariantCheckY = invariantCheckY.Add(poolYDelta2)

	if !invariantCheckX.IsZero() && !invariantCheckY.IsZero() {
		panic(fmt.Errorf("invariant check fails due to invalid swap price: %s", invariantCheckX.String()))
	}

	if !poolXDelta.Add(decimalErrorX).Equal(poolXDelta2) || !poolYDelta.Add(decimalErrorY).Equal(poolYDelta2) {
		panic(fmt.Errorf("invariant check fails due to invalid swap price: %s", poolXDelta.String()))
	}

	validitySwapPrice := types.CheckSwapPrice(matchResultXtoY, matchResultYtoX, result.SwapPrice)
	if !validitySwapPrice {
		panic("invariant check fails due to invalid swap price")
	}
}

// SwapPriceDirection checks whether the calculated swap price is increased, decreased, or stayed from the last pool price.
func SwapPriceDirection(currentPoolPrice sdk.Dec, batchResult types.BatchResult) {
	switch batchResult.PriceDirection {
	case types.Increasing:
		if !batchResult.SwapPrice.GTE(currentPoolPrice) {
			panic("invariant check fails due to incorrect price direction")
		}
	case types.Decreasing:
		if !batchResult.SwapPrice.LTE(currentPoolPrice) {
			panic("invariant check fails due to incorrect price direction")
		}
	case types.Staying:
		if !batchResult.SwapPrice.Equal(currentPoolPrice) {
			panic("invariant check fails due to incorrect price direction")
		}
	}
}

// SwapMsgStatesInvariants checks swap match result states invariants.
func SwapMsgStatesInvariants(matchResultXtoY, matchResultYtoX []types.MatchResult, matchResultMap map[uint64]types.MatchResult,
	swapMsgStates []*types.SwapMsgState, XtoY, YtoX []*types.SwapMsgState) {
	if len(matchResultXtoY)+len(matchResultYtoX) != len(matchResultMap) {
		panic("invalid length of match result")
	}

	for k, v := range matchResultMap {
		if k != v.SwapMsgState.MsgIndex {
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
			if sms.MsgIndex == msgAfter.SwapMsgState.MsgIndex {
				if *(sms) != *(msgAfter.SwapMsgState) || sms != msgAfter.SwapMsgState {
					panic("batch message not matched")
				} else {
					break
				}
			} else {
				panic("fail msg pointer consistency")
			}
		}
	}
}

// SwapOrdersExecutionStateInvariants checks all executed orders have order price which is not "executable" or not "unexecutable".
func SwapOrdersExecutionStateInvariants(matchResultMap map[uint64]types.MatchResult, swapMsgStates []*types.SwapMsgState,
	batchResult types.BatchResult, denomX string) {
	for _, sms := range swapMsgStates {
		if _, ok := matchResultMap[sms.MsgIndex]; ok {
			if !sms.Executed || !sms.Succeeded {
				panic("swap msg state consistency error, matched but not succeeded")
			}

			if sms.Msg.OfferCoin.Denom == denomX {
				// buy orders having equal or higher order price than found swapPrice
				if !sms.Msg.OrderPrice.GTE(batchResult.SwapPrice) {
					panic("execution validity failed, executed but unexecutable")
				}
			} else {
				// sell orders having equal or lower order price than found swapPrice
				if !sms.Msg.OrderPrice.LTE(batchResult.SwapPrice) {
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
				if !sms.Msg.OrderPrice.LTE(batchResult.SwapPrice) {
					panic("execution validity failed, unexecuted but executable")
				}
			} else {
				// sell orders having equal or higher order price than found swapPrice
				if !sms.Msg.OrderPrice.GTE(batchResult.SwapPrice) {
					panic("execution validity failed, unexecuted but executable")
				}
			}
		}
	}
}
