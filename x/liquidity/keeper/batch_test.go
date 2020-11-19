package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"math/rand"
	"testing"
	"time"
)

func TestDepositLiquidityPoolToBatch(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX := "denomX"
	denomY := "denomY"
	denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)
	denoms := []string{denomX, denomY}

	// get random X, Y amount for create pool
	X, Y := getRandPoolAmt(r)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	fmt.Println("-------------------------------------------------------")
	fmt.Println("X/Y", X.ToDec().Quo(Y.ToDec()), "X", X, "Y", Y)

	// set pool creator account, balance for deposit
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 3, sdk.NewInt(10000))
	app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeIndex := uint32(0)
	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeIndex, denoms, depositBalance)
	err := simapp.LiquidityKeeper.CreateLiquidityPool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllLiquidityPools(ctx)
	poolId := lpList[0].PoolId
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(0), poolId)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)


	// begin block, init
	// TODO: beginBlocker or direct function
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	//simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)

	// set pool depositor account
	app.SaveAccount(simapp, ctx, addrs[1], deposit) // pool creator
	depositA = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	depositB = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	depositBalance = sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	depositMsg := types.NewMsgDepositToLiquidityPool(addrs[1], poolId, depositBalance)
	err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)
	// TODO: escrow balance verity

	depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// TODO: endBlocker or direct function
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	//simapp.LiquidityKeeper.ExecutePoolBatch(ctx)

	// verify minted pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	depositorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	require.NotEqual(t, sdk.ZeroInt(), depositBalance)
	require.Equal(t, poolCoin, depositorBalance.Amount.Add(creatorBalance.Amount))
}
