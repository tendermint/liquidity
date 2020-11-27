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
		batches := k.GetAllLiquidityPoolBatches(ctx)
		for _, batch := range batches {
			swapMsgs := k.GetAllNotToDeleteLiquidityPoolBatchSwapMsgs(ctx, batch)
			for _, msg := range swapMsgs {
				remainingCoins = remainingCoins.Add(msg.RemainingOfferCoin)
			}
			depositMsgs := k.GetAllNotToDeleteLiquidityPoolBatchDepositMsgs(ctx, batch)
			for _, msg := range depositMsgs {
				remainingCoins = remainingCoins.Add(msg.Msg.DepositCoins...)
			}
			withdrawMsgs := k.GetAllNotToDeleteLiquidityPoolBatchWithdrawMsgs(ctx, batch)
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

//// CanWithdrawInvariant checks that current rewards can be completely withdrawn
//func CanWithdrawInvariant(k Keeper) sdk.Invariant {
//	return func(ctx sdk.Context) (string, bool) {
//
//		// cache, we don't want to write changes
//		ctx, _ = ctx.CacheContext()
//
//		var remaining sdk.DecCoins
//
//		valDelegationAddrs := make(map[string][]sdk.AccAddress)
//		for _, del := range k.stakingKeeper.GetAllSDKDelegations(ctx) {
//			valAddr := del.GetValidatorAddr().String()
//			valDelegationAddrs[valAddr] = append(valDelegationAddrs[valAddr], del.GetDelegatorAddr())
//		}
//
//		// iterate over all validators
//		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
//			_, _ = k.WithdrawValidatorCommission(ctx, val.GetOperator())
//
//			delegationAddrs, ok := valDelegationAddrs[val.GetOperator().String()]
//			if ok {
//				for _, delAddr := range delegationAddrs {
//					if _, err := k.WithdrawDelegationRewards(ctx, delAddr, val.GetOperator()); err != nil {
//						panic(err)
//					}
//				}
//			}
//
//			remaining = k.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
//			if len(remaining) > 0 && remaining[0].Amount.IsNegative() {
//				return true
//			}
//
//			return false
//		})
//
//		broken := len(remaining) > 0 && remaining[0].Amount.IsNegative()
//		return sdk.FormatInvariant(types.ModuleName, "can withdraw",
//			fmt.Sprintf("remaining coins: %v\n", remaining)), broken
//	}
//}
//
//// ReferenceCountInvariant checks that the number of historical rewards records is correct
//func ReferenceCountInvariant(k Keeper) sdk.Invariant {
//	return func(ctx sdk.Context) (string, bool) {
//
//		valCount := uint64(0)
//		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
//			valCount++
//			return false
//		})
//		dels := k.stakingKeeper.GetAllSDKDelegations(ctx)
//		slashCount := uint64(0)
//		k.IterateValidatorSlashEvents(ctx,
//			func(_ sdk.ValAddress, _ uint64, _ types.ValidatorSlashEvent) (stop bool) {
//				slashCount++
//				return false
//			})
//
//		// one record per validator (last tracked period), one record per
//		// delegation (previous period), one record per slash (previous period)
//		expected := valCount + uint64(len(dels)) + slashCount
//		count := k.GetValidatorHistoricalReferenceCount(ctx)
//		broken := count != expected
//
//		return sdk.FormatInvariant(types.ModuleName, "reference count",
//			fmt.Sprintf("expected historical reference count: %d = %v validators + %v delegations + %v slashes\n"+
//				"total validator historical reference count: %d\n",
//				expected, valCount, len(dels), slashCount, count)), broken
//	}
//}
//
//// ModuleAccountInvariant checks that the coins held by the distr ModuleAccount
//// is consistent with the sum of validator outstanding rewards
//func ModuleAccountInvariant(k Keeper) sdk.Invariant {
//	return func(ctx sdk.Context) (string, bool) {
//
//		var expectedCoins sdk.DecCoins
//		k.IterateValidatorOutstandingRewards(ctx, func(_ sdk.ValAddress, rewards types.ValidatorOutstandingRewards) (stop bool) {
//			expectedCoins = expectedCoins.Add(rewards.Rewards...)
//			return false
//		})
//
//		communityPool := k.GetFeePoolCommunityCoins(ctx)
//		expectedInt, _ := expectedCoins.Add(communityPool...).TruncateDecimal()
//
//		macc := k.GetDistributionAccount(ctx)
//		balances := k.bankKeeper.GetAllBalances(ctx, macc.GetAddress())
//
//		broken := !balances.IsEqual(expectedInt)
//		return sdk.FormatInvariant(
//			types.ModuleName, "ModuleAccount coins",
//			fmt.Sprintf("\texpected ModuleAccount coins:     %s\n"+
//				"\tdistribution ModuleAccount coins: %s\n",
//				expectedInt, balances,
//			),
//		), broken
//	}
//}
