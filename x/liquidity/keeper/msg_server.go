package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"strconv"
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

// Message server, handler for CreateLiquidityPool msg
func (k msgServer) CreateLiquidityPool(goCtx context.Context, msg *types.MsgCreateLiquidityPool) (*types.MsgCreateLiquidityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.Keeper.CreateLiquidityPool(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeCreateLiquidityPool,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.PoolCreatorAddress),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, ""),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolTypeIndex, fmt.Sprintf("%d", msg.PoolTypeIndex)),
			sdk.NewAttribute(types.AttributeValueReserveCoinDenoms, ""),
			sdk.NewAttribute(types.AttributeValueReserveAccount, ""),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, ""),
			sdk.NewAttribute(types.AttributeValueSwapFeeRate, ""),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolFeeRate, ""),
			sdk.NewAttribute(types.AttributeValueBatchSize, ""),
		),
	)
	return &types.MsgCreateLiquidityPoolResponse{}, nil
}

// Message server, handler for MsgDepositToLiquidityPool
func (k msgServer) DepositToLiquidityPool(goCtx context.Context, msg *types.MsgDepositToLiquidityPool) (*types.MsgDepositToLiquidityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	batchMsg, err := k.Keeper.DepositLiquidityPoolToBatch(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeDepositToLiquidityPoolToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueDepositCoins, batchMsg.Msg.DepositCoins.String()),
			// TODO: AttributeValueBatchIndex
			//sdk.NewAttribute(types.AttributeValueBatchIndex, ),
		),
	)
	return &types.MsgDepositToLiquidityPoolResponse{}, nil
}

// Message server, handler for MsgWithdrawFromLiquidityPool
func (k msgServer) WithdrawFromLiquidityPool(goCtx context.Context, msg *types.MsgWithdrawFromLiquidityPool) (*types.MsgWithdrawFromLiquidityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	batchMsg, err := k.Keeper.WithdrawLiquidityPoolToBatch(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeWithdrrawFromLiquidityPoolToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.WithdrawerAddress),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, batchMsg.Msg.PoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, batchMsg.Msg.PoolCoin.Amount.String()),
		),
	)
	return &types.MsgWithdrawFromLiquidityPoolResponse{}, nil
}

// Message server, handler for MsgSwap
func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.OfferCoinFee.IsZero() {
		msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin)
	}
	batchMsg, err := k.Keeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
	if err != nil {
		return &types.MsgSwapResponse{}, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			//sdk.NewAttribute(sdk.AttributeKeySender, msg.SwapRequesterAddress),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueSwapType, strconv.FormatUint(uint64(batchMsg.Msg.SwapType), 10)),
			sdk.NewAttribute(types.AttributeValueOfferCoinDenom, batchMsg.Msg.OfferCoin.Denom),
			sdk.NewAttribute(types.AttributeValueOfferCoinAmount, batchMsg.Msg.OfferCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueDemandCoinDenom, batchMsg.Msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeValueOrderPrice, batchMsg.Msg.OrderPrice.String()),
		),
	)
	return &types.MsgSwapResponse{}, nil
}
