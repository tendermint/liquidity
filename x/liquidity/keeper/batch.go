package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// TODO: testcodes
func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	// Delete already executed batches
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if liquidityPoolBatch.Executed {
			k.DeleteAllReadyLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllReadyLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllReadyLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)

			// TODO: clean, check order expiry height delete for swap

			k.InitNextBatch(ctx, liquidityPoolBatch)
		}
		return false
	})
}

func (k Keeper) InitNextBatch(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	if !liquidityPoolBatch.Executed {
		return types.ErrBatchNotExecuted
	}
	liquidityPoolBatch.BatchIndex = k.GetNextBatchIndexWithUpdate(ctx, liquidityPoolBatch.PoolId)
	liquidityPoolBatch.BeginHeight = ctx.BlockHeight()
	k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
	return nil
}

func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if !liquidityPoolBatch.Executed {
			if err := k.SwapExecution(ctx, liquidityPoolBatch); err != nil {
				// TODO: WIP
			}
			k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolDepositMsg) bool {
				if err := k.DepositLiquidityPool(ctx, batchMsg); err != nil {
					k.RefundDepositLiquidityPool(ctx, batchMsg)
				}
				return false
			})
			k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolWithdrawMsg) bool {
				if err := k.WithdrawLiquidityPool(ctx, batchMsg); err != nil {
					k.RefundWithdrawLiquidityPool(ctx, batchMsg)
				}
				return false
			})
			liquidityPoolBatch.Executed = true
			k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
		}
		return false
	})
}

func (k Keeper) HoldEscrow(ctx sdk.Context, depositor sdk.AccAddress, depositCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, depositCoins); err != nil {
		return err
	}
	return nil
}

func (k Keeper) ReleaseEscrow(ctx sdk.Context, withdrawer sdk.AccAddress, withdrawCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawer, withdrawCoins); err != nil {
		return err
	}
	return nil
}

func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	if err := k.ValidateMsgDepositLiquidityPool(ctx, *msg); err != nil {
		return err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before add to batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolDepositMsg{
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetDepositor(), msg.DepositCoins); err != nil {
		return err
	}

	poolBatch.DepositMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchDepositMsg(ctx, poolBatch.PoolId, poolBatch.DepositMsgIndex, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return nil
}

func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	if err := k.ValidateMsgWithdrawLiquidityPool(ctx, *msg); err != nil {
		return err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before add to batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolWithdrawMsg{
		MsgHeight: ctx.BlockHeight(),
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetWithdrawer(), sdk.NewCoins(msg.PoolCoin)); err != nil {
		return err
	}

	poolBatch.WithdrawMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, poolBatch.PoolId, poolBatch.WithdrawMsgIndex, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return nil
}

func (k Keeper) SwapLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgSwap) error {
	if err := k.ValidateMsgSwap(ctx, *msg); err != nil {
		return err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	// TODO: add validate msg before add to batch
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	poolBatch.SwapMsgIndex += 1
	batchPoolMsg := types.BatchPoolSwapMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.SwapMsgIndex,
		Msg:       msg,
	}
	batchPoolMsg.OrderExpiryHeight = batchPoolMsg.MsgHeight + types.CancelOrderLifeSpan

	if err := k.HoldEscrow(ctx, msg.GetSwapRequester(), sdk.NewCoins(msg.OfferCoin)); err != nil {
		return err
	}

	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchSwapMsg(ctx, poolBatch.PoolId, poolBatch.SwapMsgIndex, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return nil
}
