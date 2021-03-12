package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	params := simapp.LiquidityKeeper.GetParams(ctx)
	paramsDefault := simapp.LiquidityKeeper.GetParams(ctx)
	genesis := types.DefaultGenesisState()

	params.LiquidityPoolCreationFee = sdk.Coins{sdk.Coin{"invalid denom---", sdk.NewInt(0)}}
	err := params.Validate()
	require.Error(t, err)

	params = simapp.LiquidityKeeper.GetParams(ctx)
	params.SwapFeeRate = sdk.NewDec(-1)
	genesisState := types.NewGenesisState(params, genesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.Error(t, err)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	poolCoins := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	app.TestDepositPool(t, simapp, ctx, sdk.NewInt(30000000), sdk.NewInt(20000000), addrs[1:2], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinBalanceCreator := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(100000), poolCoinBalance.Amount)
	require.Equal(t, poolCoins.QuoRaw(10), poolCoinBalance.Amount)

	balanceXrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.NewInt(20000000), balanceXrefunded.Amount)
	require.Equal(t, sdk.ZeroInt(), balanceYrefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// validate pool records

	newGenesis := simapp.LiquidityKeeper.ExportGenesis(ctx)
	genesisState = types.NewGenesisState(paramsDefault, newGenesis.PoolRecords)
	err = types.ValidateGenesis(*genesisState)
	require.NoError(t, err)

	pool.PoolTypeIndex = 5
	simapp.LiquidityKeeper.SetPool(ctx, pool)
	newGenesisBrokenPool := simapp.LiquidityKeeper.ExportGenesis(ctx)
	err = types.ValidateGenesis(*newGenesisBrokenPool)
	require.NoError(t, err)
	require.Equal(t, 1, len(newGenesisBrokenPool.PoolRecords))
	err = simapp.LiquidityKeeper.ValidatePoolRecord(ctx, &newGenesisBrokenPool.PoolRecords[0])
	require.Error(t, err)

	// not initialized genState of other module (auth, bank, ... ) only liquidity module
	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	require.Equal(t, 2, len(reserveCoins))
	simapp2 := app.Setup(false)
	ctx2 := simapp2.BaseApp.NewContext(false, tmproto.Header{})
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
	simapp2.BankKeeper.SetBalances(ctx2, pool.GetReserveAccount(), reserveCoins)
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
	simapp2.BankKeeper.SetBalances(ctx2, addrs[0], sdk.Coins{poolCoinBalanceCreator})
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})
	simapp2.BankKeeper.SetBalances(ctx2, addrs[1], sdk.Coins{poolCoinBalance})
	require.Panics(t, func() {
		simapp2.LiquidityKeeper.InitGenesis(ctx2, *newGenesis)
	})

}
