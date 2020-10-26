package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// TODO: WIP
func (k Keeper) SwapExecution(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	pool, found := k.GetLiquidityPool(ctx, liquidityPoolBatch.PoolID)
	if !found {
		return types.ErrPoolNotExists
	}
	//totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	var sumOfBuy sdk.Coin
	var sumOfSell sdk.Coin
	swapMsgs := k.GetAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
	for _, m := range swapMsgs {
		if m.Msg.OfferCoin.Denom == sumOfSell.Denom {
			sumOfSell = sumOfSell.Add(m.Msg.OfferCoin)
		} else if m.Msg.DemandCoin.Denom == sumOfBuy.Denom {
			sumOfBuy = sumOfBuy.Add(m.Msg.DemandCoin)
		} else {
			return types.ErrInvalidDenom
		}
	}
	pool.GetReserveAccount()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
	return nil
}
