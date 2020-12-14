package liquidity_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := app.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abcitypes.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)
	params := app.LiquidityKeeper.GetParams(ctx)
	require.NotNil(t, params)
}
