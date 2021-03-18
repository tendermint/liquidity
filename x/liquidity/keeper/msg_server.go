package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Message server, handler for CreatePool msg
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, err := k.Keeper.CreatePool(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(pool.Id, 10)),
			sdk.NewAttribute(types.AttributeValuePoolTypeId, fmt.Sprintf("%d", msg.PoolTypeId)),
			sdk.NewAttribute(types.AttributeValuePoolName, pool.Name()),
			sdk.NewAttribute(types.AttributeValueReserveAccount, pool.ReserveAccountAddress),
			sdk.NewAttribute(types.AttributeValueDepositCoins, msg.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, pool.PoolCoinDenom),
		),
	)
	return &types.MsgCreatePoolResponse{}, nil
}

// Message server, handler for MsgDepositWithinBatch
func (k msgServer) DepositWithinBatch(goCtx context.Context, msg *types.MsgDepositWithinBatch) (*types.MsgDepositWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO: remove redundant GetPoolBatch
	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}
	batchMsg, err := k.Keeper.DepositLiquidityPoolToBatch(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueDepositCoins, batchMsg.Msg.DepositCoins.String()),
		),
	)
	return &types.MsgDepositWithinBatchResponse{}, nil
}

// Message server, handler for MsgWithdrawWithinBatch
func (k msgServer) WithdrawWithinBatch(goCtx context.Context, msg *types.MsgWithdrawWithinBatch) (*types.MsgWithdrawWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO: remove redundant GetPoolBatch
	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}
	batchMsg, err := k.Keeper.WithdrawLiquidityPoolToBatch(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, batchMsg.Msg.PoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, batchMsg.Msg.PoolCoin.Amount.String()),
		),
	)
	return &types.MsgWithdrawWithinBatchResponse{}, nil
}

// Message server, handler for MsgSwapWithinBatch
func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwapWithinBatch) (*types.MsgSwapWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if msg.OfferCoinFee.IsZero() {
		msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
	}
	// TODO: remove redundant GetPoolBatch
	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}
	batchMsg, err := k.Keeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
	if err != nil {
		return &types.MsgSwapWithinBatchResponse{}, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwapWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueSwapTypeId, strconv.FormatUint(uint64(batchMsg.Msg.SwapTypeId), 10)),
			sdk.NewAttribute(types.AttributeValueOfferCoinDenom, batchMsg.Msg.OfferCoin.Denom),
			sdk.NewAttribute(types.AttributeValueOfferCoinAmount, batchMsg.Msg.OfferCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueOfferCoinFeeAmount, batchMsg.Msg.OfferCoinFee.Amount.String()),
			sdk.NewAttribute(types.AttributeValueDemandCoinDenom, batchMsg.Msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeValueOrderPrice, batchMsg.Msg.OrderPrice.String()),
		),
	)
	return &types.MsgSwapWithinBatchResponse{}, nil
}
