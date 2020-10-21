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
		batch := types.NewLiquidityPoolBatch(liquidityPool.PoolID, k.GetNextBatchIndex(ctx, liquidityPool.PoolID))
		k.SetLiquidityPoolBatch(ctx, batch)
		return false
	})
}

func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if liquidityPoolBatch.ExecutionStatus {
			// TODO: ordering policy
			k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolDepositMsg) bool {
				if err := k.DepositLiquidityPool(ctx, batchMsg.Msg); err != nil {

				}
				return false
			})
			k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolWithdrawMsg) bool {
				k.WithdrawLiquidityPool(ctx, batchMsg.Msg)
				return false
			})
			k.IterateAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolSwapMsg) bool {
				k.SwapExecution(ctx, batchMsg.Msg)
				return false
			})
			liquidityPoolBatch.ExecutionStatus = true
			k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
		}
		return false
	})

	//k.IterateAllLiquidityPools(ctx, func(liquidityPool types.LiquidityPool) bool {
	//	return false
	//})
}

