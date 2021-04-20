package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	k.IterateAllPoolBatches(ctx, func(poolBatch types.PoolBatch) bool {
		// Delete and init next batch when not empty batch on before block
		if poolBatch.Executed {
			// On the other hand, BatchDeposit, BatchWithdraw, is all handled by the endblock if there is no error.
			// If there are BatchMsgs left, reset the Executed, Succeeded flag so that it can be executed in the next batch.
			depositMsgs := k.GetAllRemainingPoolBatchDepositMsgStates(ctx, poolBatch)
			if len(depositMsgs) > 0 {
				for _, msg := range depositMsgs {
					msg.Executed = false
					msg.Succeeded = false
				}
				k.SetPoolBatchDepositMsgStatesByPointer(ctx, poolBatch.PoolId, depositMsgs)
			}

			withdrawMsgs := k.GetAllRemainingPoolBatchWithdrawMsgStates(ctx, poolBatch)
			if len(withdrawMsgs) > 0 {
				for _, msg := range withdrawMsgs {
					msg.Executed = false
					msg.Succeeded = false
				}
				k.SetPoolBatchWithdrawMsgStatesByPointer(ctx, poolBatch.PoolId, withdrawMsgs)
			}

			height := ctx.BlockHeight()

			// reinitialize remaining batch msgs
			// In the case of BatchSwapMsgs, it is often fractional matched or has not yet expired since it has not passed ExpiryHeight.
			swapMsgs := k.GetAllRemainingPoolBatchSwapMsgStates(ctx, poolBatch)
			if len(swapMsgs) > 0 {
				for _, msg := range swapMsgs {
					if height > msg.OrderExpiryHeight {
						msg.ToBeDeleted = true
					} else {
						msg.Executed = false
						msg.Succeeded = false
					}
				}
				k.SetPoolBatchSwapMsgStatesByPointer(ctx, poolBatch.PoolId, swapMsgs)
			}

			// delete batch messages that were executed or ready to delete
			k.DeleteAllReadyPoolBatchDepositMsgStates(ctx, poolBatch)
			k.DeleteAllReadyPoolBatchWithdrawMsgStates(ctx, poolBatch)
			k.DeleteAllReadyPoolBatchSwapMsgStates(ctx, poolBatch)

			// Increase the batch index and initialize the values.
			if err := k.InitNextBatch(ctx, poolBatch); err != nil {
				panic(err)
			}
		}
		return false
	})
}

// Increase the index of the already executed batch for processing as the next batch and reinitialize the values.
func (k Keeper) InitNextBatch(ctx sdk.Context, poolBatch types.PoolBatch) error {
	if !poolBatch.Executed {
		return types.ErrBatchNotExecuted
	}

	poolBatch.Index = k.GetNextPoolBatchIndexWithUpdate(ctx, poolBatch.PoolId)
	poolBatch.BeginHeight = ctx.BlockHeight()
	poolBatch.Executed = false

	k.SetPoolBatch(ctx, poolBatch)

	return nil
}

// In case of deposit, withdraw, and swap msgs, unlike other normal tx msgs,
// collect them in the liquidity pool batch and perform an execution once at the endblock to calculate and use the universal price.
func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllPoolBatches(ctx, func(poolBatch types.PoolBatch) bool {
		params := k.GetParams(ctx)

		if !poolBatch.Executed && ctx.BlockHeight()-poolBatch.BeginHeight+1 >= int64(params.UnitBatchHeight) {
			executedMsgCount, err := k.SwapExecution(ctx, poolBatch)
			if err != nil {
				panic(err)
			}

			k.IterateAllPoolBatchDepositMsgStates(ctx, poolBatch, func(batchMsg types.DepositMsgState) bool {
				if batchMsg.Executed || batchMsg.ToBeDeleted || batchMsg.Succeeded {
					return false
				}
				executedMsgCount++
				if err := k.DepositLiquidityPool(ctx, batchMsg, poolBatch); err != nil {
					if err := k.RefundDepositLiquidityPool(ctx, batchMsg, poolBatch); err != nil {
						panic(err)
					}
				}
				return false
			})

			k.IterateAllPoolBatchWithdrawMsgStates(ctx, poolBatch, func(batchMsg types.WithdrawMsgState) bool {
				if batchMsg.Executed || batchMsg.ToBeDeleted || batchMsg.Succeeded {
					return false
				}
				executedMsgCount++
				if err := k.WithdrawLiquidityPool(ctx, batchMsg, poolBatch); err != nil {
					if err := k.RefundWithdrawLiquidityPool(ctx, batchMsg, poolBatch); err != nil {
						panic(err)
					}
				}
				return false
			})

			// set executed when something executed
			if executedMsgCount > 0 {
				poolBatch.Executed = true
				k.SetPoolBatch(ctx, poolBatch)
			}
		}
		return false
	})
}

// In order to deal with the batch at once, the coins of msgs deposited in escrow.
func (k Keeper) HoldEscrow(ctx sdk.Context, depositor sdk.AccAddress, depositCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, depositCoins); err != nil {
		return err
	}
	return nil
}

// If batch messages has expired or has not been processed, will be refunded the escrow that had been deposited through this function.
func (k Keeper) ReleaseEscrow(ctx sdk.Context, withdrawer sdk.AccAddress, withdrawCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawer, withdrawCoins); err != nil {
		return err
	}
	return nil
}

// Generate inputs and outputs to treat escrow refunds atomically.
func (k Keeper) ReleaseEscrowForMultiSend(withdrawer sdk.AccAddress, withdrawCoins sdk.Coins) (
	banktypes.Input, banktypes.Output, error) {
	var input banktypes.Input
	var output banktypes.Output

	input = banktypes.NewInput(k.accountKeeper.GetModuleAddress(types.ModuleName), withdrawCoins)
	output = banktypes.NewOutput(withdrawer, withdrawCoins)

	if err := banktypes.ValidateInputsOutputs([]banktypes.Input{input}, []banktypes.Output{output}); err != nil {
		return banktypes.Input{}, banktypes.Output{}, err
	}

	return input, output, nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositWithinBatch) (types.DepositMsgState, error) {
	if err := k.ValidateMsgDepositLiquidityPool(ctx, *msg); err != nil {
		return types.DepositMsgState{}, err
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.DepositMsgState{}, types.ErrPoolBatchNotExists
	}

	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	msgState := types.DepositMsgState{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.DepositMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetDepositor(), msg.DepositCoins); err != nil {
		return types.DepositMsgState{}, err
	}

	poolBatch.DepositMsgIndex++
	k.SetPoolBatch(ctx, poolBatch)
	k.SetPoolBatchDepositMsgState(ctx, poolBatch.PoolId, msgState)

	return msgState, nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawWithinBatch) (types.WithdrawMsgState, error) {
	if err := k.ValidateMsgWithdrawLiquidityPool(ctx, *msg); err != nil {
		return types.WithdrawMsgState{}, err
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.WithdrawMsgState{}, types.ErrPoolBatchNotExists
	}

	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.WithdrawMsgState{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.WithdrawMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetWithdrawer(), sdk.NewCoins(msg.PoolCoin)); err != nil {
		return types.WithdrawMsgState{}, err
	}

	poolBatch.WithdrawMsgIndex++
	k.SetPoolBatch(ctx, poolBatch)
	k.SetPoolBatchWithdrawMsgState(ctx, poolBatch.PoolId, batchPoolMsg)

	return batchPoolMsg, nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) SwapLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgSwapWithinBatch, OrderExpirySpanHeight int64) (*types.SwapMsgState, error) {
	if err := k.ValidateMsgSwapWithinBatch(ctx, *msg); err != nil {
		return nil, err
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.SwapMsgState{
		MsgHeight:            ctx.BlockHeight(),
		MsgIndex:             poolBatch.SwapMsgIndex,
		Executed:             false,
		Succeeded:            false,
		ToBeDeleted:          false,
		ExchangedOfferCoin:   sdk.NewCoin(msg.OfferCoin.Denom, sdk.ZeroInt()),
		RemainingOfferCoin:   msg.OfferCoin,
		ReservedOfferCoinFee: msg.OfferCoinFee,
		Msg:                  msg,
	}

	batchPoolMsg.OrderExpiryHeight = batchPoolMsg.MsgHeight + OrderExpirySpanHeight

	if err := k.HoldEscrow(ctx, msg.GetSwapRequester(), sdk.NewCoins(msg.OfferCoin.Add(msg.OfferCoinFee))); err != nil {
		return nil, err
	}

	poolBatch.SwapMsgIndex++
	k.SetPoolBatch(ctx, poolBatch)
	k.SetPoolBatchSwapMsgState(ctx, poolBatch.PoolId, batchPoolMsg)

	return &batchPoolMsg, nil
}
