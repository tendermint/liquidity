package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestGenesisState(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	simapp := app.Setup(false)

	ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})
	params := simapp.LiquidityKeeper.GetParams(ctx)
	genesis := types.DefaultGenesisState()

	params.LiquidityPoolCreationFee = sdk.Coins{sdk.Coin{"invalid denom---", sdk.NewInt(0)}}
	err := params.Validate()
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.SwapFeeRate = sdk.NewDec(-1)
	genesisState := types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.SwapFeeRate = sdk.NewDec(2)
	genesisState = types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.InitPoolCoinMintAmount = sdk.NewInt(0)
	genesisState = types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.InitPoolCoinMintAmount = sdk.NewInt(-1)
	genesisState = types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.InitPoolCoinMintAmount = sdk.NewInt(10)
	genesisState = types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.LiquidityPoolCreationFee = sdk.Coins{sdk.Coin{"invalid denom---", sdk.NewInt(0)}}
	err = params.Validate()
	require.Error(t, err)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.MinInitDepositAmount = sdk.NewInt(0)
	err = params.Validate()
	require.Error(t, err)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.MinInitDepositAmount = sdk.NewInt(-1)
	err = params.Validate()
	require.Error(t, err)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.PoolTypes = []types.PoolType{types.DefaultPoolType, types.DefaultPoolType}
	err = params.Validate()
	require.Error(t, err)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.PoolTypes = []types.PoolType{}
	err = params.Validate()
	require.Error(t, err)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	malformedPoolType := types.DefaultPoolType
	malformedPoolType.Id = 0
	params.PoolTypes = []types.PoolType{malformedPoolType, types.DefaultPoolType}
	err = params.Validate()
	require.Error(t, err)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	genesisState = types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.NoError(t, err)

	require.NotNil(t, genesisState)
	require.Equal(t, params, genesisState.Params)

	genesisState = types.DefaultGenesisState()
	require.NotNil(t, genesisState)

	err = types.ValidateGenesis(*genesisState)
	require.NoError(t, err)
}
