package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestGenesis(t *testing.T) {
	simapp, ctx := app.CreateTestInput()

	lk := simapp.LiquidityKeeper

	// default genesis state
	genState := types.DefaultGenesisState()
	require.Equal(t, sdk.NewDecWithPrec(3, 3), genState.Params.SwapFeeRate)

	// change swap fee rate
	params := lk.GetParams(ctx)
	params.SwapFeeRate = sdk.NewDecWithPrec(5, 3)

	// set params
	lk.SetParams(ctx, params)

	newGenState := lk.ExportGenesis(ctx)
	require.Equal(t, sdk.NewDecWithPrec(5, 3), newGenState.Params.SwapFeeRate)

	fmt.Println("newGenState: ", newGenState)
}

func TestGenesisState(t *testing.T) {
	simapp, ctx := app.CreateTestInput()

	params := simapp.LiquidityKeeper.GetParams(ctx)
	paramsDefault := simapp.LiquidityKeeper.GetParams(ctx)
	genesis := types.DefaultGenesisState()

	params.PoolCreationFee = sdk.Coins{sdk.Coin{Denom: "invalid denom---", Amount: sdk.NewInt(0)}}
	require.Error(t, params.Validate())

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.SwapFeeRate = sdk.NewDec(-1)
	genesisState := types.NewGenesisState(params, genesis.PoolRecords)
	require.Error(t, types.ValidateGenesis(*genesisState))

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	X := sdk.NewInt(100_000_000)
	Y := sdk.NewInt(200_000_000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10_000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	poolCoins := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	app.TestDepositPool(t, simapp, ctx, sdk.NewInt(30_000_000), sdk.NewInt(20_000_000), addrs[1:2], poolId, false)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinBalanceCreator := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(100_000), poolCoinBalance.Amount)
	require.Equal(t, poolCoins.QuoRaw(10), poolCoinBalance.Amount)

	balanceXRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.NewInt(20000000), balanceXRefunded.Amount)
	require.Equal(t, sdk.ZeroInt(), balanceYRefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// validate pool records
	newGenesis := simapp.LiquidityKeeper.ExportGenesis(ctx)
	genesisState = types.NewGenesisState(paramsDefault, newGenesis.PoolRecords)
	require.NoError(t, types.ValidateGenesis(*genesisState))

	pool.TypeId = 5
	simapp.LiquidityKeeper.SetPool(ctx, pool)
	newGenesisBrokenPool := simapp.LiquidityKeeper.ExportGenesis(ctx)
	require.NoError(t, types.ValidateGenesis(*newGenesisBrokenPool))
	require.Equal(t, 1, len(newGenesisBrokenPool.PoolRecords))

	err := simapp.LiquidityKeeper.ValidatePoolRecord(ctx, newGenesisBrokenPool.PoolRecords[0])
	require.Error(t, err)

	// not initialized genState of other module (auth, bank, ... ) only liquidity module
	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	require.Equal(t, 2, len(reserveCoins))
	simapp2 := app.Setup(false)
	ctx2 := simapp2.BaseApp.NewContext(false, tmproto.Header{})
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
	require.NoError(t, simapp2.BankKeeper.SetBalances(ctx2, pool.GetReserveAccount(), reserveCoins))
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
	require.NoError(t, simapp2.BankKeeper.SetBalances(ctx2, addrs[0], sdk.Coins{poolCoinBalanceCreator}))
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
	require.NoError(t, simapp2.BankKeeper.SetBalances(ctx2, addrs[1], sdk.Coins{poolCoinBalance}))
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
}
