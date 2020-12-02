package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// NewQuerier creates a querier for liquidity REST endpoints
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		fmt.Println(path, path, req)
		switch path[0] {
		case types.QueryLiquidityPool:
			return queryLiquidityPool(ctx, path[1:], req, k, legacyQuerierCdc)
		default:
			fmt.Println("querier defalt case")
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path of liquidity module: %s", path[0])
		}
	}
}

func queryLiquidityPool(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryLiquidityPoolParams

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	liquidityPool, found := k.GetLiquidityPool(ctx, params.PoolId)
	if !found {
		return nil, types.ErrPoolNotExists
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, liquidityPool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
