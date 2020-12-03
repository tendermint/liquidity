package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

// TODO: fix test for latest version of genesis state
func TestGenesisState(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	//simapp := app.Setup(false)
	//ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})
	//params := simapp.LiquidityKeeper.GetParams(ctx)

	//genesisState := types.NewGenesisState(params)
	//require.NotNil(t, genesisState)
	//require.Equal(t, params, genesisState.Params)
	//
	//genesisState = types.DefaultGenesisState()
	//require.NotNil(t, genesisState)
	//require.Equal(t, params, genesisState.Params)
	//
	//err := types.ValidateGenesis(*genesisState)
	//require.NoError(t, err)
}
