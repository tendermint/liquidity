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
	invariantCheckFlag = true // TODO: better way to handle below invariant checks?
)

// MintingPoolCoinsInvariant checks the correct ratio of minting amount of pool coins.
func MintingPoolCoinsInvariant(poolCoinTotalSupply, mintPoolCoin, depositCoinA, depositCoinB, lastReserveCoinA, lastReserveCoinB sdk.Dec) {
	poolCoinRatio := mintPoolCoin.Quo(poolCoinTotalSupply)
	depositCoinARatio := depositCoinA.Quo(lastReserveCoinA)
	depositCoinBRatio := depositCoinB.Quo(lastReserveCoinB)

	// TODO: handle case when someone sends coins to escrow module account

	// TODO: double check if this is intended result
	// there may be decimal error differences which should be smaller than 1
	poolCoinRatio = poolCoinRatio.Add(sdk.NewDec(1.0))

	// NewPoolCoinAmount / LastPoolCoinSupply = DepositCoinA / LastReserveCoinA
	// NewPoolCoinAmount / LastPoolCoinSupply = DepositCoinB / LastReserveCoinB
	if !poolCoinRatio.GTE(depositCoinARatio) || !poolCoinRatio.GTE(depositCoinBRatio) {
		panic("invariant check fails due to incorrect ratio of pool coins")
	}
}

// DepositReserveCoinsInvariant checks the after deposit amounts.
func DepositReserveCoinsInvariant(depositCoinA, depositCoinB, lastReserveCoinA, lastReserveCoinB, afterReserveCoinA, afterReserveCoinB sdk.Dec) {
	// Debugging... TestLiquidityScenario2
	// fmt.Println("afterReserveCoinA: ", afterReserveCoinA)
	// fmt.Println("afterReserveCoinB: ", afterReserveCoinB)
	// fmt.Println("lastReserveCoinA.Add(depositCoinADec): ", lastReserveCoinA.Add(depositCoinA))
	// fmt.Println("lastReserveCoinB.Add(depositCoinBDec): ", lastReserveCoinB.Add(depositCoinB))

	// AfterDepositReserveCoinA = LastReserveCoinA + DepositCoinA
	// AfterDepositReserveCoinB = LastReserveCoinB + DepositCoinB
	// if !afterReserveCoinA.Equal(lastReserveCoinA.Add(depositCoinA)) ||
	// 	!afterReserveCoinB.Equal(lastReserveCoinB.Add(depositCoinB)) {
	// 	panic("invariant check fails due to incorrect deposit amounts")
	// }
}

// DepositRatioInvariant checks the correct ratio of deposit coin amounts.
func DepositRatioInvariant(depositCoinA, depositCoinB, lastReserveCoinRatio sdk.Dec) {
	// depositCoinRatio := depositCoinA.Quo(depositCoinB)

	// // DepositCoinA / DepositCoinB = LastReserveCoinA / LastReserveCoinB
	// if !depositCoinRatio.Equal(lastReserveCoinRatio) {
	// 	panic("invariant check fails due to incorrect deposit ratio")
	// }
}

// ImmutablePoolPriceAfterDepositInvariant checks immutable pool price after depositing coins.
func ImmutablePoolPriceAfterDepositInvariant(lastReserveCoinRatio, afterReserveCoinRatio sdk.Dec) {
	// LastReserveCoinA / LastReserveCoinB = AfterDepositReserveCoinA / AfterDepositReserveCoinB
	// if !lastReserveCoinRatio.Equal(afterReserveCoinRatio) {
	// 	panic("invariant check fails due to incorrect pool price ratio")
	// }
}

// BurningPoolCoinsInvariant checks the correct burning amount of pool coins.
func BurningPoolCoinsInvariant(burnedPoolCoin, withdrawCoinA, withdrawCoinB, reserveCoinA, reserveCoinB, lastPoolCoinSupply, withdrawProportion sdk.Dec) {
	burningPoolCoinRatio := burnedPoolCoin.Quo(lastPoolCoinSupply)
	withdrawCoinARatio := withdrawCoinA.Add(withdrawProportion).Quo(reserveCoinA)
	withdrawCoinBRatio := withdrawCoinB.Add(withdrawProportion).Quo(reserveCoinB)

	// TODO: double check if this is intended result
	// there may be decimal error differences which should be smaller than 1
	burningPoolCoinRatio = burningPoolCoinRatio.Add(sdk.NewDec(1))

	// BurnedPoolCoinAmount / LastPoolCoinSupply = (WithdrawCoinA+WithdrawFeeCoinA) / LastReserveCoinA
	// BurnedPoolCoinAmount / LastPoolCoinSupply = (WithdrawCoinB+WithdrawFeeCoinB) / LastReserveCoinB
	if !burningPoolCoinRatio.GTE(withdrawCoinARatio) || !burningPoolCoinRatio.GTE(withdrawCoinBRatio) {
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

// WithdrawRatioInvariant checks the correct ratio of withdraw coin amounts.
func WithdrawRatioInvariant(withdrawCoinA, withdrawCoinB, reserveCoinA, reserveCoinB sdk.Dec) {
	// withdrawCoinRatio := withdrawCoinA.Quo(withdrawCoinB)
	// reserveCoinRatio := reserveCoinA.Quo(reserveCoinB)

	// WithdrawCoinA / WithdrawCoinB = LastReserveCoinA / LastReserveCoinB
	// if !withdrawCoinRatio.Equal(reserveCoinRatio) {
	// 	panic("invariant check fails due to incorrect ratio of withdraw coin amounts")
	// }
}

// ImmutablePoolPriceAfterWithdrawInvariant checks the immutable pool price after withdrawing coins.
func ImmutablePoolPriceAfterWithdrawInvariant(reserveCoinA, reserveCoinB, afterReserveCoinA, afterReserveCoinB sdk.Dec) {
	// reserveCoinRatio := reserveCoinA.Quo(reserveCoinB)
	// afterReserveCoinRatio := afterReserveCoinA.Quo(afterReserveCoinB)

	// LastReserveCoinA / LastReserveCoinB = AfterWithdrawReserveCoinA / AfterWithdrawReserveCoinB
	// if !reserveCoinRatio.Equal(afterReserveCoinRatio) {
	// panic("invariant check fails due to incorrect pool price ratio")
	// }
}

// SwapPriceInvariants checks the calculated swap price is increased, decreased, or equal from the last pool price.
func SwapPriceInvariants(XtoY, YtoX []*types.SwapMsgState, matchResultXtoY, matchResultYtoX []types.MatchResult,
	fractionalCntX, fractionalCntY int, poolXDelta, poolYDelta, poolXDelta2, poolYDelta2, decimalErrorX, decimalErrorY sdk.Dec, result types.BatchResult) {
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

	invariantCheckX = invariantCheckX.Add(poolXDelta)
	invariantCheckY = invariantCheckY.Add(poolYDelta)

	// print the invariant check and validity with swap, match result
	if invariantCheckX.IsZero() && invariantCheckY.IsZero() {
	} else {
		panic(invariantCheckX)
	}

	if !poolXDelta.Add(decimalErrorX).Equal(poolXDelta2) || !poolYDelta.Add(decimalErrorY).Equal(poolYDelta2) {
		panic(poolXDelta)
	}

	validitySwapPrice := types.CheckSwapPrice(matchResultXtoY, matchResultYtoX, result.SwapPrice)
	if !validitySwapPrice {
		panic("validitySwapPrice")
	}
}

// OrdersWithNotExecutedStateInvariants checks all executed orders have order price which is not "executable" or not "unexecutable".
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
