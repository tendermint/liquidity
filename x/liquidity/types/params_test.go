package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestParams(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	simapp := app.Setup(false)
	ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})
	params := types.DefaultParams()
	require.Equal(t, params, simapp.LiquidityKeeper.GetParams(ctx))

	paramsNew := types.NewParams(params.LiquidityPoolTypes, params.MinInitDepositToPool, params.InitPoolCoinMintAmount,
		params.SwapFeeRate, params.LiquidityPoolCreationFee)
	require.NotNil(t, paramsNew)
	require.Equal(t, params, paramsNew)

	res := types.ParamKeyTable()
	require.IsType(t, paramtypes.KeyTable{}, res)

	resPair := params.ParamSetPairs()
	require.IsType(t, paramtypes.ParamSetPairs{}, resPair)

	genesisStr := `liquidity_pool_types:
- pool_type_index: 1
  name: DefaultPoolType
  min_reserve_coin_num: 2
  max_reserve_coin_num: 2
  description: ""
min_init_deposit_to_pool: "1000000"
init_pool_coin_mint_amount: "1000000"
swap_fee_rate: "0.003000000000000000"
liquidity_pool_creation_fee:
- denom: uatom
  amount: "100000000"
`

	require.Equal(t, genesisStr, params.String())
	require.NoError(t, params.Validate())

	params = types.DefaultParams()
	params.LiquidityPoolTypes = nil
	require.Error(t, params.Validate())

	params = types.DefaultParams()
	dec, _ := sdk.NewDecFromStr("2.0")
	params.SwapFeeRate = dec
	require.Error(t, params.Validate())
	dec, _ = sdk.NewDecFromStr("-0.5")
	params.SwapFeeRate = dec
	require.Error(t, params.Validate())

	params = types.DefaultParams()
	params.LiquidityPoolCreationFee = sdk.NewCoins()
	require.Error(t, params.Validate())

	params = types.DefaultParams()
	params.InitPoolCoinMintAmount = sdk.ZeroInt()
	require.Error(t, params.Validate())

	params = types.DefaultParams()
	params.MinInitDepositToPool = sdk.ZeroInt()
	require.Error(t, params.Validate())

}
