package keeper

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (k Keeper) MakeQueryLiquidityPoolResponse(ctx sdk.Context, pool types.LiquidityPool) (*types.QueryLiquidityPoolResponse, error) {
	batch, found := k.GetLiquidityPoolBatch(ctx, pool.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	return &types.QueryLiquidityPoolResponse{LiquidityPool: pool,
		LiquidityPoolMetaData: k.GetPoolMetaData(ctx, pool),
		LiquidityPoolBatch:    batch}, nil
}

func (k Keeper) MakeQueryLiquidityPoolsResponse(ctx sdk.Context, pools types.LiquidityPools) (*[]types.QueryLiquidityPoolResponse, error) {
	resp := make([]types.QueryLiquidityPoolResponse, len(pools))
	for i, pool := range pools {
		batch, found := k.GetLiquidityPoolBatch(ctx, pool.PoolId)
		if !found {
			return nil, types.ErrPoolBatchNotExists
		}
		meta := k.GetPoolMetaData(ctx, pool)
		res := types.QueryLiquidityPoolResponse{
			LiquidityPool:pool,
			LiquidityPoolMetaData:meta,
			LiquidityPoolBatch:batch,
		}
		resp[i] = res
	}

	return &resp, nil
}

func (k Keeper) LiquidityPool(c context.Context, req *types.QueryLiquidityPoolRequest) (*types.QueryLiquidityPoolResponse, error) {
	empty := &types.QueryLiquidityPoolRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetLiquidityPool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool %d doesn't exist", req.PoolId)
	}
	return k.MakeQueryLiquidityPoolResponse(ctx, pool)
}

func (k Keeper) LiquidityPools(c context.Context, req *types.QueryLiquidityPoolsRequest) (*types.QueryLiquidityPoolsResponse, error) {
	empty := &types.QueryLiquidityPoolsRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	poolStore := prefix.NewStore(store, types.LiquidityPoolKeyPrefix)
	var pools types.LiquidityPools

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		pool, err := types.UnmarshalLiquidityPool(k.cdc, value)
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

	return &types.QueryLiquidityPoolsResponse{*response, pageRes}, nil
}

func (k Keeper) LiquidityPoolsBatch(c context.Context, req *types.QueryLiquidityPoolsBatchRequest) (*types.QueryLiquidityPoolsBatchResponse, error) {
	empty := &types.QueryLiquidityPoolsBatchRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	batchStore := prefix.NewStore(store, types.LiquidityPoolBatchKeyPrefix)
	var response []types.QueryLiquidityPoolBatchResponse

	pageRes, err := query.Paginate(batchStore, req.Pagination, func(key []byte, value []byte) error {
		batch, err := types.UnmarshalLiquidityPoolBatch(k.cdc, value)
		if err != nil {
			return err
		}
		res := &types.QueryLiquidityPoolBatchResponse{
			LiquidityPoolBatch: batch,
		}
		response = append(response, *res)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLiquidityPoolsBatchResponse{response, pageRes}, nil
}

func (k Keeper) LiquidityPoolBatch(c context.Context, req *types.QueryLiquidityPoolBatchRequest) (*types.QueryLiquidityPoolBatchResponse, error) {
	empty := &types.QueryLiquidityPoolBatchRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	batch, found := k.GetLiquidityPoolBatch(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool batch %d doesn't exist", req.PoolId)
	}
	return &types.QueryLiquidityPoolBatchResponse{LiquidityPoolBatch: batch}, nil
}

func (k Keeper) PoolBatchSwapMsgs(c context.Context, req *types.QueryPoolBatchSwapMsgsRequest) (*types.QueryPoolBatchSwapMsgsResponse, error) {
	empty := &types.QueryPoolBatchSwapMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.LiquidityPoolBatchSwapMsgIndexKeyPrefix)
	var msgs []types.BatchPoolSwapMsg

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalBatchPoolSwapMsg(k.cdc, value)
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
		SwapMsgs:   msgs,
		Pagination: pageRes,
	}, nil
	return nil, nil
}

func (k Keeper) PoolBatchDepositMsgs(c context.Context, req *types.QueryPoolBatchDepositMsgsRequest) (*types.QueryPoolBatchDepositMsgsResponse, error) {
	empty := &types.QueryPoolBatchDepositMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.LiquidityPoolBatchDepositMsgIndexKeyPrefix)
	var msgs []types.BatchPoolDepositMsg

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalBatchPoolDepositMsg(k.cdc, value)
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
		DepositMsgs: msgs,
		Pagination:  pageRes,
	}, nil
}

func (k Keeper) PoolBatchWithdrawMsgs(c context.Context, req *types.QueryPoolBatchWithdrawMsgsRequest) (*types.QueryPoolBatchWithdrawMsgsResponse, error) {
	empty := &types.QueryPoolBatchWithdrawMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.LiquidityPoolBatchWithdrawMsgIndexKeyPrefix)
	var msgs []types.BatchPoolWithdrawMsg

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalBatchPoolWithdrawMsg(k.cdc, value)
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
		WithdrawMsgs: msgs,
		Pagination:   pageRes,
	}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
