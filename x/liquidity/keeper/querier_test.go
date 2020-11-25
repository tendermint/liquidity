package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"strings"
	"testing"
)

const custom = "custom"

func getQueriedLiquidityPool(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, poolId uint64) (types.LiquidityPool, error) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryLiquidityPool}, "/"),
		Data: cdc.MustMarshalJSON(types.QueryLiquidityPoolParams{PoolId: poolId}),
	}

	pool := types.LiquidityPool{}
	bz, err := querier(ctx, []string{types.QueryLiquidityPool}, query)
	if err != nil {
		return pool, err
	}
	require.Nil(t, cdc.UnmarshalJSON(bz, &pool))
	return pool, nil
}

func TestNewQuerier(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	simapp := app.Setup(false)
	ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})
	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	querier := keeper.NewQuerier(simapp.LiquidityKeeper, cdc)

	poolId := testCreatePool(t, simapp, ctx, X, Y, DenomX, DenomY, addrs[0])
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryLiquidityPool}, "/"),
		Data: cdc.MustMarshalJSON(types.QueryLiquidityPoolParams{PoolId: poolId}),
	}
	queryFailCase := abci.RequestQuery{
		Path: strings.Join([]string{"failCustom", "failRoute", "failQuery"}, "/"),
		Data: cdc.MustMarshalJSON(types.LiquidityPool{}),
	}
	pool := types.LiquidityPool{}
	bz, err := querier(ctx, []string{types.QueryLiquidityPool}, query)
	require.NoError(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &pool))

	bz, err = querier(ctx, []string{"fail"}, queryFailCase)
	require.Error(t, err)
	require.Error(t, cdc.UnmarshalJSON(bz, &pool))
}

func TestQueries(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)

	simapp := app.Setup(false)
	ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	querier := keeper.NewQuerier(simapp.LiquidityKeeper, cdc)

	poolId := testCreatePool(t, simapp, ctx, X, Y, DenomX, DenomY, addrs[0])
	require.Equal(t, uint64(1), poolId)
	poolRes, err := getQueriedLiquidityPool(t, ctx, cdc, querier, poolId)
	require.NoError(t, err)
	require.Equal(t, poolId, poolRes.PoolId)
	require.Equal(t, DefaultPoolTypeIndex, poolRes.PoolTypeIndex)
	require.Equal(t, []string{DenomX, DenomY}, poolRes.ReserveCoinDenoms)
	require.NotNil(t, poolRes.PoolCoinDenom)
	require.NotNil(t, poolRes.ReserveAccountAddress)

	poolResEmpty, err := getQueriedLiquidityPool(t, ctx, cdc, querier, uint64(2))
	require.Error(t, err)
	require.Equal(t, uint64(0), poolResEmpty.PoolId)
}
