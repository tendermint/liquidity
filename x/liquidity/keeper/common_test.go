package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// createTestInput Returns a simapp with custom LiquidityKeeper
// to avoid messing with the hooks.
func createTestInput() (*app.LiquidityApp, sdk.Context) {
	return app.CreateTestInput()
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func AddCoins(app *app.LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}
	return app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

// setTotalSupply provides the total supply based on accAmt * totalAccounts.
func setTotalSupply(app *app.LiquidityApp, ctx sdk.Context, accAmt sdk.Int, totalAccounts int) {
	prevSupply := app.BankKeeper.GetSupply(ctx, app.StakingKeeper.BondDenom(ctx))
	diff := accAmt.MulRaw(int64(totalAccounts)).Sub(prevSupply.Amount)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), diff)))
}

// saveAccount saves the provided account into the simapp with balance based on initCoins.
func saveAccount(app *app.LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := AddCoins(app, ctx, addr, initCoins)
	if err != nil {
		panic(err)
	}
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

func addTestAddrs(app *app.LiquidityApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	setTotalSupply(app, ctx, accAmt, accNum)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		saveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func createTestPool(X, Y sdk.Coin) (*app.LiquidityApp, sdk.Context, types.Pool, sdk.AccAddress, error) {
	simapp, ctx := createTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)

	depositCoins := sdk.NewCoins(X, Y)
	creatorAddr := app.AddRandomTestAddr(simapp, ctx, depositCoins.Add(params.PoolCreationFee...))

	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, types.NewMsgCreatePool(creatorAddr, types.DefaultPoolTypeId, depositCoins))
	if err != nil {
		return nil, sdk.Context{}, types.Pool{}, sdk.AccAddress{}, err
	}

	return simapp, ctx, pool, creatorAddr, nil
}