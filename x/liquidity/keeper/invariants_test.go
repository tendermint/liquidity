package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

func TestLiquidityPoolsEscrowAmountInvariant(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, true)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, true)

	invariant := keeper.AllInvariants(simapp.LiquidityKeeper)
	msg, broken := invariant(ctx)
	require.False(t, broken)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	msg, broken = invariant(ctx)
	require.False(t, broken)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	msg, broken = invariant(ctx)
	require.False(t, broken)

	price, _ := sdk.NewDecFromStr("1.1")
	priceY, _ := sdk.NewDecFromStr("1.2")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000))}
	offerCoinListY := []sdk.Coin{sdk.NewCoin(denomY, sdk.NewInt(5000))}
	orderPriceList := []sdk.Dec{price}
	orderPriceListY := []sdk.Dec{priceY}
	orderAddrList := addrs[1:2]
	orderAddrListY := addrs[2:3]
	_, batch := app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	_, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	_, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	_, batch = app.TestSwapPool(t, simapp, ctx, offerCoinListY, orderPriceListY, orderAddrListY, poolId, false)

	msg, broken = invariant(ctx)
	require.False(t, broken)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	msg, broken = invariant(ctx)
	require.False(t, broken)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	msg, broken = invariant(ctx)
	require.False(t, broken)

	batchEscrowAcc := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	escrowAmt := simapp.BankKeeper.GetAllBalances(ctx, batchEscrowAcc)
	simapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addrs[0],
		sdk.NewCoins(sdk.NewCoin(offerCoinList[0].Denom, offerCoinList[0].Amount.QuoRaw(2))))
	escrowAmt = simapp.BankKeeper.GetAllBalances(ctx, batchEscrowAcc)

	msg, broken = invariant(ctx)
	fmt.Println(msg, escrowAmt, batch)
	require.True(t, broken)
}
