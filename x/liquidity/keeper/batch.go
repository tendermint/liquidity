package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	// Delete already executed batches
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if liquidityPoolBatch.ExecutionStatus {
			k.DeleteLiquidityPoolBatch(ctx, liquidityPoolBatch)
		}
		return false
	})

	// Init empty batch
	k.IterateAllLiquidityPools(ctx, func(liquidityPool types.LiquidityPool) bool {
		batch := types.NewLiquidityPoolBatch(liquidityPool.PoolId, k.GetNextBatchIndexWithUpdate(ctx, liquidityPool.PoolId))
		k.SetLiquidityPoolBatch(ctx, batch)
		return false
	})
}

func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if liquidityPoolBatch.ExecutionStatus {
			if err := k.SwapExecution(ctx, liquidityPoolBatch); err != nil {
				// TODO: WIP
			}
			k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolDepositMsg) bool {
				if err := k.DepositLiquidityPool(ctx, batchMsg.Msg); err != nil {
					// TODO: err handling
				}
				return false
			})
			k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolWithdrawMsg) bool {
				if err := k.WithdrawLiquidityPool(ctx, batchMsg.Msg); err != nil {
					// TODO: err handling
				}
				return false
			})
			liquidityPoolBatch.ExecutionStatus = true
			k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
		}
		return false
	})
}
