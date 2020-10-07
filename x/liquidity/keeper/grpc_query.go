package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) LiquidityPool(c context.Context, req *types.QueryLiquidityPoolRequest) (*types.QueryLiquidityPoolResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetLiquidityPool(ctx, req.PoolID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool %d doesn't exist", req.PoolID)
	}

	return &types.QueryLiquidityPoolResponse{LiquidityPool: pool}, nil
}
