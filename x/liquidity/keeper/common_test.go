package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/liquidity/app"
	lapp "github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

// createTestInput Returns a simapp with custom LiquidityKeeper
// to avoid messing with the hooks.
func createTestInput() (*lapp.LiquidityApp, sdk.Context) {
	return lapp.CreateTestInput()
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

// setTotalSupply provides the total supply based on accAmt * totalAccounts.
func setTotalSupply(app *lapp.LiquidityApp, ctx sdk.Context, accAmt sdk.Int, totalAccounts int) {
	totalSupply := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(totalAccounts))))
	prevSupply := app.BankKeeper.GetSupply(ctx)
	app.BankKeeper.SetSupply(ctx, banktypes.NewSupply(prevSupply.GetTotal().Add(totalSupply...)))
}

// saveAccount saves the provided account into the simapp with balance based on initCoins.
func saveAccount(app *lapp.LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := app.BankKeeper.AddCoins(ctx, addr, initCoins)
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

func addTestAddrs(app *lapp.LiquidityApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	setTotalSupply(app, ctx, accAmt, accNum)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		saveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(app *simapp.SimApp, ctx sdk.Context, numAddrs int) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := simapp.AddTestAddrsIncremental(app, ctx, numAddrs, sdk.NewInt(10000))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}

var (
	valTokens = sdk.TokensFromConsensusPower(42)
	//TestProposal        = types.NewTextProposal("Test", "description")
	TestDescription     = stakingtypes.NewDescription("T", "E", "S", "T", "Z")
	TestCommissionRates = stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

var (
	pubkeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}
)

func createValidators(t *testing.T, ctx sdk.Context, app *lapp.LiquidityApp, powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []stakingtypes.Validator) {
	addrs := lapp.AddTestAddrsIncremental(app, ctx, 5, sdk.NewInt(30000000))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	pks := simapp.CreateTestPubKeys(5)

	appCodec, _ := simapp.MakeCodecs()
	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		app.GetKey(stakingtypes.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	val1, _ := stakingtypes.NewValidator(valAddrs[0], pks[0], stakingtypes.Description{})
	val2, _ := stakingtypes.NewValidator(valAddrs[1], pks[1], stakingtypes.Description{})
	val3, _ := stakingtypes.NewValidator(valAddrs[2], pks[2], stakingtypes.Description{})
	vals := []stakingtypes.Validator{val1, val2, val3}

	app.StakingKeeper.SetValidator(ctx, val1)
	app.StakingKeeper.SetValidator(ctx, val2)
	app.StakingKeeper.SetValidator(ctx, val3)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val1)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val2)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val3)
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val1)
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val2)
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val3)

	_, _ = app.StakingKeeper.Delegate(ctx, addrs[0], sdk.TokensFromConsensusPower(powers[0]), stakingtypes.Unbonded, val1, true)
	_, _ = app.StakingKeeper.Delegate(ctx, addrs[1], sdk.TokensFromConsensusPower(powers[1]), stakingtypes.Unbonded, val2, true)
	_, _ = app.StakingKeeper.Delegate(ctx, addrs[2], sdk.TokensFromConsensusPower(powers[2]), stakingtypes.Unbonded, val3, true)

	_ = staking.EndBlocker(ctx, app.StakingKeeper)

	return addrs, valAddrs, vals
}

func createLiquidity(t *testing.T, ctx sdk.Context, simapp *lapp.LiquidityApp) (
	[]sdk.AccAddress, []types.LiquidityPool, []types.LiquidityPoolBatch,
	[]types.BatchPoolDepositMsg, []types.BatchPoolWithdrawMsg) {
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	denomA, denomB := types.AlphabeticalDenomPair("denomA", "denomB")

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)
	A := sdk.NewInt(500000000)
	B := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := lapp.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	lapp.TestCreatePool(t, simapp, ctx, A, B, denomA, denomB, addrs[1])

	lapp.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	price, _ := sdk.NewDecFromStr("1.1")
	priceY, _ := sdk.NewDecFromStr("1.2")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000))}
	offerCoinListY := []sdk.Coin{sdk.NewCoin(denomY, sdk.NewInt(5000))}
	orderPriceList := []sdk.Dec{price}
	orderPriceListY := []sdk.Dec{priceY}
	orderAddrList := addrs[1:2]
	orderAddrListY := addrs[2:3]

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	lapp.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	lapp.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	lapp.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[1:2], poolId, false)
	lapp.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[1:2], poolId, false)
	lapp.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[2:3], poolId, false)
	lapp.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[2:3], poolId, false)

	lapp.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	lapp.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	lapp.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	lapp.TestSwapPool(t, simapp, ctx, offerCoinListY, orderPriceListY, orderAddrListY, poolId, false)

	pools := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	batches := simapp.LiquidityKeeper.GetAllLiquidityPoolBatches(ctx)
	depositMsgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchDepositMsgs(ctx, batches[0])
	withdrawMsgs := simapp.LiquidityKeeper.GetAllLiquidityPoolBatchWithdrawMsgs(ctx, batches[0])
	return addrs, pools, batches, depositMsgs, withdrawMsgs
}
