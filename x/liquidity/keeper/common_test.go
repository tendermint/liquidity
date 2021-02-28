package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/liquidity/app"
	lapp "github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

// createTestInput Returns a simapp with custom LiquidityKeeper
// to avoid messing with the hooks.
func createTestInput() (*lapp.LiquidityApp, sdk.Context) {
	return lapp.CreateTestInput()
}

var (
	valTokens = sdk.TokensFromConsensusPower(42)
	//TestProposal        = types.NewTextProposal("Test", "description")
	TestDescription     = stakingtypes.NewDescription("T", "E", "S", "T", "Z")
	TestCommissionRates = stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

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
