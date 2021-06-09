package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestWithdrawRatioInvariant(t *testing.T) {
	require.NotPanics(t, func() {
		keeper.WithdrawAmountInvariant(sdk.NewInt(1), sdk.NewInt(1), sdk.NewInt(2), sdk.NewInt(3), sdk.NewInt(1), sdk.NewInt(2), types.DefaultParams().WithdrawFeeRate)
	})
	require.Panics(t, func() {
		keeper.WithdrawAmountInvariant(sdk.NewInt(1), sdk.NewInt(1), sdk.NewInt(2), sdk.NewInt(5), sdk.NewInt(1), sdk.NewInt(2), types.DefaultParams().WithdrawFeeRate)
	})
}

func TestLiquidityPoolsEscrowAmountInvariant(t *testing.T) {
	simapp, ctx := app.CreateTestInput()

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolID := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolID, true)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolID, true)

	invariant := keeper.AllInvariants(simapp.LiquidityKeeper)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	_, broken = invariant(ctx)
	require.False(t, broken)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	_, broken = invariant(ctx)
	require.False(t, broken)

	price, _ := sdk.NewDecFromStr("1.1")
	priceY, _ := sdk.NewDecFromStr("1.2")
	xOfferCoins := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000))}
	yOfferCoins := []sdk.Coin{sdk.NewCoin(denomY, sdk.NewInt(5000))}
	xOrderPrices := []sdk.Dec{price}
	yOrderPrices := []sdk.Dec{priceY}
	xOrderAddrs := addrs[1:2]
	yOrderAddrs := addrs[2:3]
	app.TestSwapPool(t, simapp, ctx, xOfferCoins, xOrderPrices, xOrderAddrs, poolID, false)
	app.TestSwapPool(t, simapp, ctx, xOfferCoins, xOrderPrices, xOrderAddrs, poolID, false)
	app.TestSwapPool(t, simapp, ctx, xOfferCoins, xOrderPrices, xOrderAddrs, poolID, false)
	app.TestSwapPool(t, simapp, ctx, yOfferCoins, yOrderPrices, yOrderAddrs, poolID, false)

	_, broken = invariant(ctx)
	require.False(t, broken)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	_, broken = invariant(ctx)
	require.False(t, broken)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	_, broken = invariant(ctx)
	require.False(t, broken)

	batchEscrowAcc := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	escrowAmt := simapp.BankKeeper.GetAllBalances(ctx, batchEscrowAcc)
	require.NotEmpty(t, escrowAmt)
	err := simapp.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addrs[0],
		sdk.NewCoins(sdk.NewCoin(xOfferCoins[0].Denom, xOfferCoins[0].Amount.QuoRaw(2))))
	require.NoError(t, err)
	escrowAmt = simapp.BankKeeper.GetAllBalances(ctx, batchEscrowAcc)

	msg, broken := invariant(ctx)
	require.True(t, broken)
	require.Equal(t, "liquidity: batch escrow amount invariant broken invariant\nbatch escrow amount LT batch remaining amount\n", msg)
}
