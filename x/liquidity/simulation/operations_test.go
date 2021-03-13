package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/tendermint/liquidity/app"
	simappparams "github.com/tendermint/liquidity/app/params"
	"github.com/tendermint/liquidity/x/liquidity/simulation"
	"github.com/tendermint/liquidity/x/liquidity/types"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// TestWeightedOperations tests the weights of the operations.
func TestWeightedOperations(t *testing.T) {
	app, ctx := createTestApp(false)

	ctx.WithChainID("test-chain")

	cdc := app.AppCodec()
	appParams := make(simtypes.AppParams)

	weightesOps := simulation.WeightedOperations(appParams, cdc, app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)

	s := rand.NewSource(1)
	r := rand.New(s)
	accs := simtypes.RandomAccounts(r, 3)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{simappparams.DefaultWeightMsgCreateLiquidityPool, types.ModuleName, types.TypeMsgCreateLiquidityPool},
		{simappparams.DefaultWeightMsgDepositToLiquidityPool, types.ModuleName, types.TypeMsgDepositToLiquidityPool},
		{simappparams.DefaultWeightMsgWithdrawFromLiquidityPool, types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool},
		{simappparams.DefaultWeightMsgSwap, types.ModuleName, types.TypeMsgSwap},
	}

	for i, w := range weightesOps {
		operationMsg, _, _ := w.Op()(r, app.BaseApp, ctx, accs, ctx.ChainID())
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		require.Equal(t, expected[i].weight, w.Weight(), "weight should be the same")
		require.Equal(t, expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		require.Equal(t, expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

// TestSimulateMsgCreateLiquidityPool tests the normal scenario of a valid message of type TypeMsgCreateLiquidityPool.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateLiquidityPool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup randomly generated liquidity pool creation fees
	feeCoins := simulation.GenLiquidityPoolCreationFee(r)
	params := app.LiquidityKeeper.GetParams(ctx)
	params.LiquidityPoolCreationFee = feeCoins
	app.LiquidityKeeper.SetParams(ctx, params)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgCreateLiquidityPool(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgCreateLiquidityPool
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.GetPoolCreator().String())
	require.Equal(t, types.DefaultPoolTypeIndex, msg.PoolTypeIndex)
	require.Equal(t, "171625357wLfFy,279341739zDmT", msg.DepositCoins.String())
	require.Equal(t, types.TypeMsgCreateLiquidityPool, msg.Type())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgDepositToLiquidityPool tests the normal scenario of a valid message of type TypeMsgDepositToLiquidityPool.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgDepositToLiquidityPool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup random liquidity pools
	setupLiquidityPools(t, r, app, ctx, accounts)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgDepositToLiquidityPool(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgDepositToLiquidityPool
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "cosmos1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7u4x9a0", msg.GetDepositor().String())
	require.Equal(t, "160538706Qfyze,478362889VIkPZ", msg.DepositCoins.String())
	require.Equal(t, types.TypeMsgDepositToLiquidityPool, msg.Type())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgWithdrawFromLiquidityPool tests the normal scenario of a valid message of type TypeMsgWithdrawFromLiquidityPool.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgWithdrawFromLiquidityPool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup random liquidity pools
	setupLiquidityPools(t, r, app, ctx, accounts)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgWithdrawFromLiquidityPool(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgWithdrawFromLiquidityPool
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "cosmos1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7u4x9a0", msg.GetWithdrawer().String())
	require.Equal(t, "70867pool/2D59CF15954FA399BBEA5EE6A2E73D09BC39FC8720F2E922AC17C9AC06758EA8", msg.PoolCoin.String())
	require.Equal(t, types.TypeMsgWithdrawFromLiquidityPool, msg.Type())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgSwap tests the normal scenario of a valid message of type TypeMsgSwap.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgSwap(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// setup random liquidity pools
	setupLiquidityPools(t, r, app, ctx, accounts)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgSwap(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgSwap
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.GetSwapRequester().String())
	require.Equal(t, "960168fGaE", msg.OfferCoin.String())
	require.Equal(t, "jXUlr", msg.DemandCoinDenom)
	require.Equal(t, types.TypeMsgSwap, msg.Type())
	require.Len(t, futureOperations, 0)
}

// returns context and an app
func createTestApp(isCheckTx bool) (*app.LiquidityApp, sdk.Context) {
	app := app.Setup(false)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *app.LiquidityApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := sdk.TokensFromConsensusPower(1e6)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		err := app.BankKeeper.SetBalances(ctx, account.Address, initCoins)
		require.NoError(t, err)
	}

	return accounts
}

func setupLiquidityPools(t *testing.T, r *rand.Rand, app *app.LiquidityApp, ctx sdk.Context, accounts []simtypes.Account) {
	params := app.StakingKeeper.GetParams(ctx)

	for _, account := range accounts {
		// random denom with a length from 4 to 6 characters
		denomA := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
		denomB := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
		denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

		// random fees
		fees := sdk.NewCoin(params.GetBondDenom(), sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e10, 1e12))))

		// mint random amounts of denomA and denomB coins
		mintCoinA := sdk.NewCoin(denomA, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e14, 1e15))))
		mintCoinB := sdk.NewCoin(denomB, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e14, 1e15))))
		mintCoins := sdk.NewCoins(mintCoinA, mintCoinB, fees)
		err := app.BankKeeper.MintCoins(ctx, types.ModuleName, mintCoins)
		require.NoError(t, err)

		// transfer random amounts to the simulated random account
		coinA := sdk.NewCoin(denomA, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e12, 1e14))))
		coinB := sdk.NewCoin(denomB, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e12, 1e14))))
		coins := sdk.NewCoins(coinA, coinB, fees)
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, account.Address, coins)
		require.NoError(t, err)

		// create liquidity pool with random deposit amounts
		account := app.AccountKeeper.GetAccount(ctx, account.Address)
		depositCoinA := sdk.NewCoin(denomA, sdk.NewInt(int64(simtypes.RandIntBetween(r, int(types.DefaultMinInitDepositToPool.Int64()), 1e8))))
		depositCoinB := sdk.NewCoin(denomB, sdk.NewInt(int64(simtypes.RandIntBetween(r, int(types.DefaultMinInitDepositToPool.Int64()), 1e8))))
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

		createPoolMsg := types.NewMsgCreateLiquidityPool(account.GetAddress(), types.DefaultPoolTypeIndex, depositCoins)

		_, err = app.LiquidityKeeper.CreatePool(ctx, createPoolMsg)
		require.NoError(t, err)
	}
}
