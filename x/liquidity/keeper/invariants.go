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
func MintingPoolCoinsInvariant() {
	// NewPoolTokenAmount / LastPoolTokenSupply = DepositTokenA / LastReserveTokenA
	// NewPoolTokenAmount / LastPoolTokenSupply = DepositTokenB / LastReserveTokenB
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

// OrdersWithNotExecutedStateInvariants checks all executed orders have order price which is not "unexecutable"
func OrdersWithNotExecutedStateInvariants() {
}

// OrdersWithExecutedStateInvariants checks all unexecuted orders have order price which is not "executable"
func OrdersWithExecutedStateInvariants() {
}
