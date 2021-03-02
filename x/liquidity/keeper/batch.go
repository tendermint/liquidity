package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		// Delete and init next batch when not empty batch on before block
		if liquidityPoolBatch.Executed {

			// On the other hand, BatchDeposit, BatchWithdraw, is all handled by the endblock if there is no error.
			// If there are BatchMsgs left, reset the Executed, Succeeded flag so that it can be executed in the next batch.
			depositMsgs := k.GetAllRemainingLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch)
			if len(depositMsgs) > 0 {
				for _, msg := range depositMsgs {
					msg.Executed = false
					msg.Succeeded = false
				}
				k.SetLiquidityPoolBatchDepositMsgsByPointer(ctx, liquidityPoolBatch.PoolId, depositMsgs)
			}

			withdrawMsgs := k.GetAllRemainingLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch)
			if len(withdrawMsgs) > 0 {
				for _, msg := range withdrawMsgs {
					msg.Executed = false
					msg.Succeeded = false
				}
				k.SetLiquidityPoolBatchWithdrawMsgsByPointer(ctx, liquidityPoolBatch.PoolId, withdrawMsgs)
			}

			height := ctx.BlockHeight()
			// reinitialize remaining batch msgs
			// In the case of BatchSwapMsgs, it is often fractional matched or has not yet expired since it has not passed ExpiryHeight.
			swapMsgs := k.GetAllRemainingLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
			if len(swapMsgs) > 0 {
				for _, msg := range swapMsgs {
					if height > msg.OrderExpiryHeight {
						msg.ToBeDeleted = true
					} else {
						msg.Executed = false
						msg.Succeeded = false
					}
				}
				k.SetLiquidityPoolBatchSwapMsgPointers(ctx, liquidityPoolBatch.PoolId, swapMsgs)
			}

			// delete batch messages that were executed or ready to delete
			k.DeleteAllReadyLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllReadyLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllReadyLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)

			// Increase the batch index and initialize the values.
			k.InitNextBatch(ctx, liquidityPoolBatch)
		}
		return false
	})
}

// Increase the index of the already executed batch for processing as the next batch and reinitialize the values.
func (k Keeper) InitNextBatch(ctx sdk.Context, liquidityPoolBatch types.LiquidityPoolBatch) error {
	if !liquidityPoolBatch.Executed {
		return types.ErrBatchNotExecuted
	}
	liquidityPoolBatch.BatchIndex = k.GetNextBatchIndexWithUpdate(ctx, liquidityPoolBatch.PoolId)
	liquidityPoolBatch.BeginHeight = ctx.BlockHeight()
	liquidityPoolBatch.Executed = false
	k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
	return nil
}

// In case of deposit, withdraw, and swap msgs, unlike other normal tx msgs,
// collect them in the liquidity pool batch and perform an execution once at the endblock to calculate and use the universal price.
func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		params := k.GetParams(ctx)
		if !liquidityPoolBatch.Executed && ctx.BlockHeight()-liquidityPoolBatch.BeginHeight+1 >= int64(params.UnitBatchSize) {
			executedMsgCount, err := k.SwapExecution(ctx, liquidityPoolBatch)
			if err != nil {
				panic(err)
			}
			k.IterateAllLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolDepositMsg) bool {
				executedMsgCount++
				if err := k.DepositLiquidityPool(ctx, batchMsg, liquidityPoolBatch); err != nil {
					k.RefundDepositLiquidityPool(ctx, batchMsg, liquidityPoolBatch)
				}
				return false
			})
			k.IterateAllLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch, func(batchMsg types.BatchPoolWithdrawMsg) bool {
				executedMsgCount++
				if err := k.WithdrawLiquidityPool(ctx, batchMsg, liquidityPoolBatch); err != nil {
					k.RefundWithdrawLiquidityPool(ctx, batchMsg, liquidityPoolBatch)
				}
				return false
			})
			// set executed when something executed
			if executedMsgCount > 0 {
				liquidityPoolBatch.Executed = true
				k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
			}
		}
		return false
	})
}

// In order to deal with the batch at once, the coins of msgs deposited in escrow.
func (k Keeper) HoldEscrow(ctx sdk.Context, depositor sdk.AccAddress, depositCoins sdk.Coins) error {
	for _, dc := range depositCoins {
		balance := k.bankKeeper.GetBalance(ctx, depositor, dc.Denom)
		fmt.Println("balance: ", balance) // 2591591574toeA
		fmt.Println("depositCoin: ", dc)  // 8834780187toeA
		if balance.Amount.LT(dc.Amount) {
			continue
		}
	}

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
func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) (types.BatchPoolDepositMsg, error) {
	if err := k.ValidateMsgDepositLiquidityPool(ctx, *msg); err != nil {
		return types.BatchPoolDepositMsg{}, err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.BatchPoolDepositMsg{}, types.ErrPoolBatchNotExists
	}
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolDepositMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.DepositMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetDepositor(), msg.DepositCoins); err != nil {
		return types.BatchPoolDepositMsg{}, err
	}

	poolBatch.DepositMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchDepositMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return batchPoolMsg, nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) (types.BatchPoolWithdrawMsg, error) {
	if err := k.ValidateMsgWithdrawLiquidityPool(ctx, *msg); err != nil {
		return types.BatchPoolWithdrawMsg{}, err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.BatchPoolWithdrawMsg{}, types.ErrPoolBatchNotExists
	}
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolWithdrawMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.WithdrawMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetWithdrawer(), sdk.NewCoins(msg.PoolCoin)); err != nil {
		return types.BatchPoolWithdrawMsg{}, err
	}

	poolBatch.WithdrawMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return batchPoolMsg, nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) SwapLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgSwap, OrderExpirySpanHeight int64) (*types.BatchPoolSwapMsg, error) {
	if err := k.ValidateMsgSwap(ctx, *msg); err != nil {
		return nil, err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolSwapMsg{
		MsgHeight:           ctx.BlockHeight(),
		MsgIndex:            poolBatch.SwapMsgIndex,
		Executed:            false,
		Succeeded:           false,
		ToBeDeleted:         false,
		ExchangedOfferCoin:  sdk.NewCoin(msg.OfferCoin.Denom, sdk.ZeroInt()),
		RemainingOfferCoin:  msg.OfferCoin,
		OfferCoinFeeReserve: msg.OfferCoinFee,
		Msg:                 msg,
	}
	// TODO: add logic if OrderExpiryHeight==0, pass on batch logic
	batchPoolMsg.OrderExpiryHeight = batchPoolMsg.MsgHeight + OrderExpirySpanHeight

	if err := k.HoldEscrow(ctx, msg.GetSwapRequester(), sdk.NewCoins(msg.OfferCoin)); err != nil {
		return nil, err
	}

	if err := k.HoldEscrow(ctx, msg.GetSwapRequester(), sdk.NewCoins(msg.OfferCoinFee)); err != nil {
		return nil, err
	}

	poolBatch.SwapMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchSwapMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return &batchPoolMsg, nil
}
