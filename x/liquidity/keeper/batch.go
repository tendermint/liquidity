package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	// Delete already executed batches
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if liquidityPoolBatch.Executed {
			// TODO: verify clean, check order expiry height delete for swap
			depositMsgs := k.GetAllRemainingLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch)
			if len(depositMsgs) > 0 {
				for _, msg := range depositMsgs {
					msg.Executed = false
					msg.Succeed = false
				}
				k.SetLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch.PoolId, depositMsgs)
			}
			// TODO: verify set

			withdrawMsgs := k.GetAllRemainingLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch)
			if len(withdrawMsgs) > 0 {
				for _, msg := range withdrawMsgs {
					msg.Executed = false
					msg.Succeed = false
				}
				k.SetLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch.PoolId, withdrawMsgs)
			}

			swapMsgs := k.GetAllRemainingLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)
			if len(swapMsgs) > 0 {
				for _, msg := range swapMsgs {
					msg.Executed = false
					msg.Succeed = false
				}
				k.SetLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch.PoolId, swapMsgs)
			}

			// TODO: optimizing iteration
			k.DeleteAllReadyLiquidityPoolBatchDepositMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllReadyLiquidityPoolBatchWithdrawMsgs(ctx, liquidityPoolBatch)
			k.DeleteAllReadyLiquidityPoolBatchSwapMsgs(ctx, liquidityPoolBatch)

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
	liquidityPoolBatch.Executed = false
	k.SetLiquidityPoolBatch(ctx, liquidityPoolBatch)
	return nil
}

func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllLiquidityPoolBatches(ctx, func(liquidityPoolBatch types.LiquidityPoolBatch) bool {
		if !liquidityPoolBatch.Executed {
			if liquidityPoolBatch.Executed {
				return false
			}
			if err := k.SwapExecution(ctx, liquidityPoolBatch); err != nil {
				panic(err)
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

func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	if err := k.ValidateMsgDepositLiquidityPool(ctx, *msg); err != nil {
		return err
	}
	poolBatch, found := k.GetLiquidityPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
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
		return err
	}

	poolBatch.DepositMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchDepositMsg(ctx, poolBatch.PoolId, batchPoolMsg)
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
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolWithdrawMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.WithdrawMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetWithdrawer(), sdk.NewCoins(msg.PoolCoin)); err != nil {
		return err
	}

	poolBatch.WithdrawMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchWithdrawMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return nil
}

func (k Keeper) SwapLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgSwap) (*types.BatchPoolSwapMsg, error) {
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
		MsgHeight:          ctx.BlockHeight(),
		MsgIndex:           poolBatch.SwapMsgIndex,
		Executed:           false,
		Succeed:            false,
		ToDelete:           false,
		ExchangedOfferCoin: sdk.NewCoin(msg.OfferCoin.Denom, sdk.ZeroInt()),
		RemainingOfferCoin: msg.OfferCoin,
		Msg:                msg,
	}
	// TODO: add logic if OrderExpiryHeight==0, pass on batch logic
	batchPoolMsg.OrderExpiryHeight = batchPoolMsg.MsgHeight + types.CancelOrderLifeSpan

	if err := k.HoldEscrow(ctx, msg.GetSwapRequester(), sdk.NewCoins(msg.OfferCoin)); err != nil {
		return nil, err
	}

	poolBatch.SwapMsgIndex += 1
	k.SetLiquidityPoolBatch(ctx, poolBatch)
	k.SetLiquidityPoolBatchSwapMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return &batchPoolMsg, nil
}
