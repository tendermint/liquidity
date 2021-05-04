package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SoftForkAirdrop(ctx sdk.Context, providerAddr string, omittedAddrList []string, distributionCoin sdk.Coin) error {
	cachedCtx, writeCache := ctx.CacheContext()
	providerAcc, err := sdk.AccAddressFromBech32(providerAddr)
	if err != nil {
		return err
	}
	var distributionCoins = sdk.NewCoins(distributionCoin)
	var accList []sdk.AccAddress
	for _, addr := range omittedAddrList {
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