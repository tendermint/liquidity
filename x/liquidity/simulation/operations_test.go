package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/simulation"
	"github.com/tendermint/liquidity/x/liquidity/types"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// TestSimulateMsgCreateLiquidityPool tests the normal scenario of a valid message of type TypeMsgCreateLiquidityPool.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateLiquidityPool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 1)

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
	require.Equal(t, sdk.NewInt(1000000), msg.DepositCoins.AmountOf("denomA"))
	require.Equal(t, sdk.NewInt(1000000), msg.DepositCoins.AmountOf("denomB"))
	require.Equal(t, types.TypeMsgCreateLiquidityPool, msg.Type())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgDepositToLiquidityPool tests the normal scenario of a valid message of type TypeMsgDepositToLiquidityPool.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgDepositToLiquidityPool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgDepositToLiquidityPool(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgDepositToLiquidityPool
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, uint64(1), msg.PoolId)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.GetDepositor().String())
	require.Equal(t, sdk.NewInt(5000), msg.DepositCoins.AmountOf("denomA"))
	require.Equal(t, sdk.NewInt(5000), msg.DepositCoins.AmountOf("denomB"))
	require.Equal(t, types.TypeMsgDepositToLiquidityPool, msg.Type())
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgWithdrawFromLiquidityPool tests the normal scenario of a valid message of type TypeMsgWithdrawFromLiquidityPool.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgWithdrawFromLiquidityPool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 1)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgWithdrawFromLiquidityPool(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgWithdrawFromLiquidityPool
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, uint64(1), msg.PoolId)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.GetWithdrawer().String())
	require.Equal(t, sdk.NewInt(5000), msg.PoolCoin.Amount)
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

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgSwap(app.AccountKeeper, app.BankKeeper, app.LiquidityKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgSwap
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, uint64(1), msg.PoolId)
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.GetSwapRequester().String())
	require.Equal(t, sdk.NewInt(5000), msg.OfferCoin.Amount)
	require.Equal(t, "denomB", msg.DemandCoinDenom)
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
