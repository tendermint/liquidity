package types_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestParams(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	simapp := app.Setup(false)
	ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})
	params := types.DefaultParams()
	currentParams := simapp.LiquidityKeeper.GetParams(ctx)
	require.Equal(t, params, currentParams)

	paramsNew := types.NewParams(params.PoolTypes, params.MinInitDepositToPool, params.InitPoolCoinMintAmount,
		params.ReserveCoinLimitAmount, params.LiquidityPoolCreationFee, params.SwapFeeRate, params.WithdrawFeeRate,
		params.MaxOrderAmountRatio, params.UnitBatchSize)
	require.NotNil(t, paramsNew)
	require.Equal(t, params, paramsNew)

	res := types.ParamKeyTable()
	require.IsType(t, paramstypes.KeyTable{}, res)

	resPair := params.ParamSetPairs()
	require.IsType(t, paramstypes.ParamSetPairs{}, resPair)

	genesisStr := `pool_types:
- id: 1
  name: DefaultPoolType
  min_reserve_coin_num: 2
  max_reserve_coin_num: 2
  description: ""
min_init_deposit_to_pool: "1000000"
init_pool_coin_mint_amount: "1000000"
reserve_coin_limit_amount: "0"
liquidity_pool_creation_fee:
- denom: stake
  amount: "100000000"
swap_fee_rate: "0.003000000000000000"
withdraw_fee_rate: "0.003000000000000000"
max_order_amount_ratio: "0.100000000000000000"
unit_batch_size: 1
`
	fmt.Println(params.String())
	require.Equal(t, genesisStr, params.String())
	require.NoError(t, params.Validate())

	params = types.DefaultParams()
	params.PoolTypes = nil
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
