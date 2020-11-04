package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

func TestLiquidityPool(t *testing.T) {
	app, ctx := createTestInput()
	lp := types.LiquidityPool{
		PoolID:            0,
		PoolTypeIndex:     0,
		ReserveCoinDenoms: []string{"a", "b"},
		ReserveAccount:    nil,
		PoolCoinDenom:     "poolCoin",
	}
	app.LiquidityKeeper.SetLiquidityPool(ctx, lp)

	lpGet, found := app.LiquidityKeeper.GetLiquidityPool(ctx, 0)
	require.True(t, found)
	require.Equal(t, lp, lpGet)
}

func TestCreateLiquidityPool(t *testing.T) {
	simapp, ctx := createTestInput()

	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	poolTypeIndex := uint32(0)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	denoms := []string{denomA, denomB}

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)

	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(0), lpList[0].PoolID)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextLiquidityPoolIDWithUpdate(ctx))
	require.Equal(t, denomA, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomB, lpList[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	err = simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
	require.Error(t, err, types.ErrPoolAlreadyExists)
}

func TestDepositLiquidityPool(t *testing.T) {
	simapp, ctx := createTestInput()

	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	poolTypeIndex := uint32(0)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	denoms := []string{denomA, denomB}

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)
	app.SaveAccount(simapp, ctx, addrs[1], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	depositA = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomA)
	depositB = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomB)
	depositBalance = sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	createMsg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)

	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, createMsg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	lp := lpList[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)

	depositMsg := types.NewMsgDepositToLiquidityPool(addrs[1], lp.PoolID, deposit)
	err = simapp.LiquidityKeeper.DepositLiquidityPool(ctx, depositMsg)
	require.NoError(t, err)

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)
	depositorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], lp.PoolCoinDenom)
	require.Equal(t, poolCoin.Sub(poolCoinBefore), depositorBalance.Amount)
}

func TestWithdrawLiquidityPool(t *testing.T) {
	simapp, ctx := createTestInput()

	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	poolTypeIndex := uint32(0)
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	denoms := []string{denomA, denomB}

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	createMsg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)

	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, createMsg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	lp := lpList[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)
	withdrawerPoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.PoolCoinDenom)

	fmt.Println(poolCoinBefore, withdrawerPoolCoinBefore.Amount)
	require.Equal(t, poolCoinBefore, withdrawerPoolCoinBefore.Amount)
	withdrawMsg := types.NewMsgWithdrawFromLiquidityPool(addrs[0], lp.PoolID, sdk.NewCoins(sdk.NewCoin(lp.PoolCoinDenom, poolCoinBefore)))
	err = simapp.LiquidityKeeper.WithdrawLiquidityPool(ctx, withdrawMsg)
	require.NoError(t, err)

	poolCoinAfter := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)
	withdrawerPoolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.PoolCoinDenom)
	require.True(t, true, poolCoinAfter.IsZero())
	require.True(t, true, withdrawerPoolCoinAfter.IsZero())
	withdrawerDenomAbalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.ReserveCoinDenoms[0])
	withdrawerDenomBbalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.ReserveCoinDenoms[1])
	require.Equal(t, deposit.AmountOf(lp.ReserveCoinDenoms[0]), withdrawerDenomAbalance.Amount)
	require.Equal(t, deposit.AmountOf(lp.ReserveCoinDenoms[1]), withdrawerDenomBbalance.Amount)

}
