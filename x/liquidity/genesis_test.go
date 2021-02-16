package liquidity_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestGenesisState(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
	simapp := app.Setup(false)

	ctx := simapp.BaseApp.NewContext(false, tmproto.Header{})
	genesis := types.DefaultGenesisState()

	liquidity.InitGenesis(ctx, simapp.LiquidityKeeper, *genesis)

	genesisExported := liquidity.ExportGenesis(ctx, simapp.LiquidityKeeper)

	fmt.Println(genesis)
	fmt.Println(genesisExported)

	require.Equal(t, genesis, genesisExported)
}
