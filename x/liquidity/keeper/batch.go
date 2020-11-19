package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	// Delete already executed batches
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if liquidityPoolBatch.ExecutionStatus {
			// TODO: remove all msgs
			k.DeleteAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
			// TODO: remove after endblock? direct delete on fail for deposit, withdraw
			// TODO: clean, check span height delete for swap
			k.DeleteLiquidityPoolBatch(ctx, liquidityPoolBatch)
			// TODO: init next Batch only for executed, no error
		}
		return false
	})

	// Init empty batch
	// TODO: init only after executed
	//k.IterateAllLiquidityPools(ctx, func(liquidityPool types.LiquidityPool) bool {
	//	batch := types.NewLiquidityPoolBatch(liquidityPool.PoolId, k.GetNextBatchIndexWithUpdate(ctx, liquidityPool.PoolId))
	//	k.SetLiquidityPoolBatch(ctx, batch)
	//	return false
	//})
}

func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if !liquidityPoolBatch.ExecutionStatus {
			if err := k.SwapExecution(ctx, liquidityPoolBatch); err != nil {
				// TODO: WIP
			}
			k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolDepositMsg) bool {
				if err := k.DepositLiquidityPool(ctx, batchMsg.Msg); err != nil {
					k.RefundDepositLiquidityPool(ctx, batchMsg)
				}
				// TODO: remove executed msg?
				return false
			})
			k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolWithdrawMsg) bool {
				if err := k.WithdrawLiquidityPool(ctx, batchMsg.Msg); err != nil {
					k.RefundWithdrawLiquidityPool(ctx, batchMsg)
				}
				// TODO: remove executed msg?
				return false
			})
			liquidityPoolBatch.ExecutionStatus = true
			k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
		}
		return false
	})
}

func (k Keeper) DepositEscrow(ctx sdk.Context, depositor sdk.AccAddress, depositCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, depositCoins); err != nil {
		return err
	}
	return nil
}

func (k Keeper) WithdrawEscrow(ctx sdk.Context, withdrawer sdk.AccAddress, withdrawCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawer, withdrawCoins); err != nil {
		return err
	}
	return nil
}

func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before executed on batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolDepositMsg{
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}

	// TODO: escrow
	if err := k.DepositEscrow(ctx, msg.Depositor, msg.DepositCoins); err != nil {
		return err
	}

	poolBatch.DepositMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchDepositMsg(ctx, poolBatch, poolBatch.DepositMsgIndex, batchPoolMsg)
	return nil
}

func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before executed on batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolWithdrawMsg{
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}

	// TODO: escrow
	if err := k.DepositEscrow(ctx, msg.Withdrawer, msg.PoolCoin); err != nil {
		return err
	}

	poolBatch.WithdrawMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, poolBatch, poolBatch.WithdrawMsgIndex, batchPoolMsg)
	return nil
}

func (k Keeper) SwapLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgSwap) error {
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before executed on batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	poolBatch.SwapMsgIndex += 1
	batchPoolMsg := types.BatchPoolSwapMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.SwapMsgIndex,
		Msg:       msg,
	}
	batchPoolMsg.CancelHeight = batchPoolMsg.MsgHeight + types.CancelOrderLifeSpan
	// TODO: escrow
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchSwapMsg(ctx, poolBatch, poolBatch.SwapMsgIndex, batchPoolMsg)
	return nil
}
