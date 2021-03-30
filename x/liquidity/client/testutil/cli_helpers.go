package testutil

import (
	"fmt"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/app/params"
	liquiditycli "github.com/tendermint/liquidity/x/liquidity/client/cli"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	dbm "github.com/tendermint/tm-db"
)

// NewConfig returns config that defines the necessary configuration
// used to bootstrap and start an in-process local testing network.
func NewConfig() network.Config {
	encCfg := app.MakeEncodingConfig()
	cfg := network.DefaultConfig()
	cfg.AppConstructor = NewAppConstructor(encCfg)                // the ABCI application constructor
	cfg.GenesisState = app.ModuleBasics.DefaultGenesis(cfg.Codec) // liquidity genesis state to provide
	return cfg
}

// NewAppConstructor returns a new liquidity app AppConstructor.
func NewAppConstructor(encodingCfg params.EncodingConfig) network.AppConstructor {
	return func(val network.Validator) servertypes.Application {
		return app.NewLiquidityApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			encodingCfg,
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}

func MsgCreatePoolExec(clientCtx client.Context, from, to, amount fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{from.String(), to.String(), amount.String()}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewCreatePoolCmd(), args)
}
