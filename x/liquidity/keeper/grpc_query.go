package keeper

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// LiquidityPool queries a liquidity pool with the given pool id.
func (k Querier) LiquidityPool(c context.Context, req *types.QueryLiquidityPoolRequest) (*types.QueryLiquidityPoolResponse, error) {
	empty := &types.QueryLiquidityPoolRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool %d doesn't exist", req.PoolId)
	}

	return k.MakeQueryLiquidityPoolResponse(ctx, pool)
}

// LiquidityPoolBatch queries a liquidity pool batch with the given pool id.
func (k Querier) LiquidityPoolBatch(c context.Context, req *types.QueryLiquidityPoolBatchRequest) (*types.QueryLiquidityPoolBatchResponse, error) {
	empty := &types.QueryLiquidityPoolBatchRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	batch, found := k.GetPoolBatch(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool batch %d doesn't exist", req.PoolId)
	}

	return &types.QueryLiquidityPoolBatchResponse{
		Batch: batch,
	}, nil
}

// Pools queries all liquidity pools currently existed with each liquidity pool with batch and metadata.
func (k Querier) LiquidityPools(c context.Context, req *types.QueryLiquidityPoolsRequest) (*types.QueryLiquidityPoolsResponse, error) {
	empty := &types.QueryLiquidityPoolsRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	poolStore := prefix.NewStore(store, types.PoolKeyPrefix)

	var pools types.Pools

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		pool, err := types.UnmarshalPool(k.cdc, value)
		if err != nil {
			return err
		}
		pools = append(pools, pool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := k.MakeQueryLiquidityPoolsResponse(ctx, pools)

	return &types.QueryLiquidityPoolsResponse{
		Pools:      *response,
		Pagination: pageRes,
	}, nil
}

// LiquidityPoolsBatches queries all liquidity pools batch.
func (k Querier) LiquidityPoolsBatches(c context.Context, req *types.QueryLiquidityPoolsBatchesRequest) (*types.QueryLiquidityPoolsBatchesResponse, error) {
	empty := &types.QueryLiquidityPoolsBatchesRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	batchStore := prefix.NewStore(store, types.PoolBatchKeyPrefix)
	var batches []types.PoolBatch

	pageRes, err := query.Paginate(batchStore, req.Pagination, func(key []byte, value []byte) error {
		batch, err := types.UnmarshalPoolBatch(k.cdc, value)
		if err != nil {
			return err
		}
		batches = append(batches, batch)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLiquidityPoolsBatchesResponse{
		Batches:    batches,
		Pagination: pageRes,
	}, nil
}

// PoolBatchSwapMsg queries the pool batch swap message with the message index of the liquidity pool.
func (k Querier) PoolBatchSwapMsg(c context.Context, req *types.QueryPoolBatchSwapMsgRequest) (*types.QueryPoolBatchSwapMsgResponse, error) {
	empty := &types.QueryPoolBatchSwapMsgRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	msg, found := k.GetPoolBatchSwapMsgState(ctx, req.PoolId, req.MsgIndex)
	if !found {
		return nil, status.Errorf(codes.NotFound, "the msg given msg_index %d doesn't exist or deleted", req.MsgIndex)
	}

	return &types.QueryPoolBatchSwapMsgResponse{
		Swaps: msg,
	}, nil
}

// PoolBatchSwapMsgs queries all pool batch swap messages of the liquidity pool.
func (k Querier) PoolBatchSwapMsgs(c context.Context, req *types.QueryPoolBatchSwapMsgsRequest) (*types.QueryPoolBatchSwapMsgsResponse, error) {
	empty := &types.QueryPoolBatchSwapMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.GetPoolBatchSwapMsgStatesPrefix(req.PoolId))

	var msgs []types.SwapMsgState

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalSwapMsgState(k.cdc, value)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolBatchSwapMsgsResponse{
		Swaps:      msgs,
		Pagination: pageRes,
	}, nil
}

// PoolBatchDepositMsg queries the pool batch deposit message with the msg_index of the liquidity pool.
func (k Querier) PoolBatchDepositMsg(c context.Context, req *types.QueryPoolBatchDepositMsgRequest) (*types.QueryPoolBatchDepositMsgResponse, error) {
	empty := &types.QueryPoolBatchDepositMsgRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	msg, found := k.GetPoolBatchDepositMsgState(ctx, req.PoolId, req.MsgIndex)
	if !found {
		return nil, status.Errorf(codes.NotFound, "the msg given msg_index %d doesn't exist or deleted", req.MsgIndex)
	}

	return &types.QueryPoolBatchDepositMsgResponse{
		Deposits: msg,
	}, nil
}

// PoolBatchDepositMsgs queries all pool batch deposit messages of the liquidity pool.
func (k Querier) PoolBatchDepositMsgs(c context.Context, req *types.QueryPoolBatchDepositMsgsRequest) (*types.QueryPoolBatchDepositMsgsResponse, error) {
	empty := &types.QueryPoolBatchDepositMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.GetPoolBatchDepositMsgStatesPrefix(req.PoolId))
	var msgs []types.DepositMsgState

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalDepositMsgState(k.cdc, value)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolBatchDepositMsgsResponse{
		Deposits:   msgs,
		Pagination: pageRes,
	}, nil
}

// PoolBatchWithdrawMsg queries the pool batch withdraw message with the msg_index of the liquidity pool.
func (k Querier) PoolBatchWithdrawMsg(c context.Context, req *types.QueryPoolBatchWithdrawMsgRequest) (*types.QueryPoolBatchWithdrawMsgResponse, error) {
	empty := &types.QueryPoolBatchWithdrawMsgRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	msg, found := k.GetPoolBatchWithdrawMsgState(ctx, req.PoolId, req.MsgIndex)
	if !found {
		return nil, status.Errorf(codes.NotFound, "the msg given msg_index %d doesn't exist or deleted", req.MsgIndex)
	}

	return &types.QueryPoolBatchWithdrawMsgResponse{
		Withdraws: msg,
	}, nil
}

// PoolBatchWithdrawMsgs queries all pool batch withdraw messages of the liquidity pool.
func (k Querier) PoolBatchWithdrawMsgs(c context.Context, req *types.QueryPoolBatchWithdrawMsgsRequest) (*types.QueryPoolBatchWithdrawMsgsResponse, error) {
	empty := &types.QueryPoolBatchWithdrawMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.GetPoolBatchWithdrawMsgsPrefix(req.PoolId))
	var msgs []types.WithdrawMsgState

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalWithdrawMsgState(k.cdc, value)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolBatchWithdrawMsgsResponse{
		Withdraws:  msgs,
		Pagination: pageRes,
	}, nil
}

// Params queries params of liquidity module.
func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// MakeQueryLiquidityPoolResponse wraps MakeQueryLiquidityPoolResponse.
func (k Querier) MakeQueryLiquidityPoolResponse(ctx sdk.Context, pool types.Pool) (*types.QueryLiquidityPoolResponse, error) {
	batch, found := k.GetPoolBatch(ctx, pool.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	return &types.QueryLiquidityPoolResponse{
		PoolId:                pool.PoolId,
		PoolTypeIndex:         pool.PoolTypeIndex,
		ReserveCoinDenoms:     pool.ReserveCoinDenoms,
		ReserveAccountAddress: pool.ReserveAccountAddress,
		PoolCoinDenom:         pool.PoolCoinDenom,
		Metadata:              k.GetPoolMetaData(ctx, pool),
		Batch:                 batch,
	}, nil
}

// MakeQueryLiquidityPoolsResponse wraps a list of QueryLiquidityPoolResponses.
func (k Querier) MakeQueryLiquidityPoolsResponse(ctx sdk.Context, pools types.Pools) (*[]types.QueryLiquidityPoolResponse, error) {
	resp := make([]types.QueryLiquidityPoolResponse, len(pools))
	for i, pool := range pools {
		batch, found := k.GetPoolBatch(ctx, pool.PoolId)
		if !found {
			return nil, types.ErrPoolBatchNotExists
		}

		meta := k.GetPoolMetaData(ctx, pool)

		res := types.QueryLiquidityPoolResponse{
			PoolId:                pool.PoolId,
			PoolTypeIndex:         pool.PoolTypeIndex,
			ReserveCoinDenoms:     pool.ReserveCoinDenoms,
			ReserveAccountAddress: pool.ReserveAccountAddress,
			PoolCoinDenom:         pool.PoolCoinDenom,
			Metadata:              meta,
			Batch:                 batch,
		}

		resp[i] = res
	}

	return &resp, nil
}
