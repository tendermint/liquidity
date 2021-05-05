package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) SoftForkAirdrop(ctx sdk.Context, providerAddr string, targetAddrs []string, distributionCoin sdk.Coin) error {
	cachedCtx, writeCache := ctx.CacheContext()
	providerAcc, err := sdk.AccAddressFromBech32(providerAddr)
	if err != nil {
		return err
	}
	var distributionCoins = sdk.NewCoins(distributionCoin)
	var accList []sdk.AccAddress
	for _, addr := range targetAddrs {
		acc, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return err
		}
		accList = append(accList, acc)
	}
	providerSpendableCoins := k.bankKeeper.SpendableCoins(cachedCtx, providerAcc)
	if !providerSpendableCoins.IsAllGTE(sdk.Coins{sdk.NewCoin(distributionCoin.Denom, distributionCoin.Amount.MulRaw(int64(len(accList))))}) {
		return fmt.Errorf("insufficient balances of provider account for softfork distribution")
	}
	for _, acc := range accList {
		err := k.bankKeeper.SendCoins(cachedCtx, providerAcc, acc, distributionCoins)
		if err != nil {
			return err
		}
	}
	// Write ctx only when it's done without errors.
	writeCache()
	return  nil
}

func (k Keeper) SoftForkAirdropMultiCoins(ctx sdk.Context, providerAddr string, airdropPairs []types.AirdropPair) error {
	cachedCtx, writeCache := ctx.CacheContext()
	providerAcc, err := sdk.AccAddressFromBech32(providerAddr)
	if err != nil {
		return err
	}
	totalDistributionCoins := sdk.NewCoins()
	for i, pair := range airdropPairs {
		if pair.TargetAcc == nil && pair.TargetAddress != "" {
			targetAcc, err := sdk.AccAddressFromBech32(pair.TargetAddress)
			if err != nil {
				return err
			}
			airdropPairs[i].TargetAcc = targetAcc
		}
		if err := pair.DistributionCoins.Validate(); err != nil {
			return err
		}
		totalDistributionCoins = totalDistributionCoins.Add(pair.DistributionCoins...)
	}
	providerSpendableCoins := k.bankKeeper.SpendableCoins(cachedCtx, providerAcc)
	if !providerSpendableCoins.IsAllGTE(totalDistributionCoins) {
		return fmt.Errorf("insufficient balances of provider account for softfork distribution")
	}
	for _, pair := range airdropPairs {
		fmt.Println(pair)
		err := k.bankKeeper.SendCoins(cachedCtx, providerAcc, pair.TargetAcc, pair.DistributionCoins)
		if err != nil {
			return err
		}
	}
	// Write ctx only when it's done without errors.
	writeCache()
	return  nil
}