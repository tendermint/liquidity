package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	lapp "github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestLiquidityPool(t *testing.T) {
	app, ctx := createTestInput()
	lp := types.Pool{
		Id:                    0,
		TypeId:                0,
		ReserveCoinDenoms:     []string{"a", "b"},
		ReserveAccountAddress: "",
		PoolCoinDenom:         "poolCoin",
	}
	app.LiquidityKeeper.SetPool(ctx, lp)

	lpGet, found := app.LiquidityKeeper.GetPool(ctx, 0)
	require.True(t, found)
	require.Equal(t, lp, lpGet)
}

func TestCreatePool(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	invalidMsg := types.NewMsgCreatePool(addrs[0], 0, depositBalance)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, invalidMsg)
	require.Error(t, err, types.ErrBadPoolTypeId)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	require.Equal(t, 1, len(pools))
	require.Equal(t, uint64(1), pools[0].Id)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextPoolId(ctx)-1)
	require.Equal(t, denomA, pools[0].ReserveCoinDenoms[0])
	require.Equal(t, denomB, pools[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Error(t, err, types.ErrPoolAlreadyExists)
}

func TestPoolCreationFee(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	// Set PoolCreationFee for fail (insufficient balances for pool creation fee)
	params.PoolCreationFee = depositBalance
	simapp.LiquidityKeeper.SetParams(ctx, params)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Equal(t, types.ErrInsufficientPoolCreationFee, err)

	// Set PoolCreationFee for success
	params.PoolCreationFee = types.DefaultPoolCreationFee
	simapp.LiquidityKeeper.SetParams(ctx, params)
	feePoolAcc := simapp.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)
	feePoolBalance := simapp.BankKeeper.GetAllBalances(ctx, feePoolAcc)
	msg = types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// Verify PoolCreationFee pay successfully
	feePoolBalance = feePoolBalance.Add(params.PoolCreationFee...)
	require.Equal(t, params.PoolCreationFee, feePoolBalance)
	require.Equal(t, feePoolBalance, simapp.BankKeeper.GetAllBalances(ctx, feePoolAcc))
}

func TestDepositLiquidityPool(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 4, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)
	lapp.SaveAccount(simapp, ctx, addrs[1], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	depositA = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomA)
	depositB = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomB)
	depositBalance = sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	createMsg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createMsg)
	require.NoError(t, err)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	pool := pools[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

	depositMsg := types.NewMsgDepositWithinBatch(addrs[1], pool.Id, deposit)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	poolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, depositMsg.PoolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	err = simapp.LiquidityKeeper.DepositLiquidityPool(ctx, msgs[0], poolBatch)
	require.NoError(t, err)

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	depositorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, poolCoin.Sub(poolCoinBefore), depositorBalance.Amount)
}

func TestReserveCoinLimit(t *testing.T) {
	simapp, ctx := createTestInput()
	params := types.DefaultParams()
	params.MaxReserveCoinAmount = sdk.NewInt(1000000000000)
	simapp.LiquidityKeeper.SetParams(ctx, params)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, params.MaxReserveCoinAmount), sdk.NewCoin(denomB, sdk.NewInt(1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Equal(t, types.ErrExceededReserveCoinLimit, err)

	params.MaxReserveCoinAmount = sdk.ZeroInt()
	simapp.LiquidityKeeper.SetParams(ctx, params)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	params.MaxReserveCoinAmount = sdk.NewInt(1000000000000)
	simapp.LiquidityKeeper.SetParams(ctx, params)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	pool := pools[0]

	deposit = sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(1000000)), sdk.NewCoin(denomB, sdk.NewInt(1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[1], deposit)
	depositMsg := types.NewMsgDepositWithinBatch(addrs[1], pool.Id, deposit)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.Equal(t, types.ErrExceededReserveCoinLimit, err)

	params.MaxReserveCoinAmount = sdk.ZeroInt()
	simapp.LiquidityKeeper.SetParams(ctx, params)

	depositMsg = types.NewMsgDepositWithinBatch(addrs[1], pool.Id, deposit)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	poolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, depositMsg.PoolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	simapp.LiquidityKeeper.ExecutePoolBatch(ctx)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)
	lapp.SaveAccount(simapp, ctx, addrs[1], deposit)
	depositMsg = types.NewMsgDepositWithinBatch(addrs[1], pool.Id, deposit)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	params.MaxReserveCoinAmount = sdk.NewInt(1000000000000)
	simapp.LiquidityKeeper.SetParams(ctx, params)

	poolBatch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, depositMsg.PoolId)
	require.True(t, found)
	msgs = simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	err = simapp.LiquidityKeeper.DepositLiquidityPool(ctx, msgs[0], poolBatch)
	require.Equal(t, types.ErrExceededReserveCoinLimit, err)
}

func TestWithdrawLiquidityPool(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	createMsg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createMsg)
	require.NoError(t, err)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	pool := pools[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	withdrawerPoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)

	require.Equal(t, poolCoinBefore, withdrawerPoolCoinBefore.Amount)
	withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[0], pool.Id, sdk.NewCoin(pool.PoolCoinDenom, poolCoinBefore))

	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
	require.NoError(t, err)

	poolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, withdrawMsg.PoolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	err = simapp.LiquidityKeeper.WithdrawLiquidityPool(ctx, msgs[0], poolBatch)
	require.NoError(t, err)

	poolCoinAfter := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	withdrawerPoolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.True(t, true, poolCoinAfter.IsZero())
	require.True(t, true, withdrawerPoolCoinAfter.IsZero())
	withdrawerDenomABalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.ReserveCoinDenoms[0])
	withdrawerDenomBBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.ReserveCoinDenoms[1])
	require.Equal(t, deposit.AmountOf(pool.ReserveCoinDenoms[0]).ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerDenomABalance.Amount)
	require.Equal(t, deposit.AmountOf(pool.ReserveCoinDenoms[1]).ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerDenomBBalance.Amount)
}

func TestReinitializePool(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)
	params.WithdrawFeeRate = sdk.ZeroDec()
	simapp.LiquidityKeeper.SetParams(ctx, params)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(100*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	createMsg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createMsg)
	require.NoError(t, err)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	pool := pools[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	withdrawerPoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)

	reserveCoins := simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	require.True(t, reserveCoins.IsEqual(deposit))

	require.Equal(t, poolCoinBefore, withdrawerPoolCoinBefore.Amount)
	withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[0], pool.Id, sdk.NewCoin(pool.PoolCoinDenom, poolCoinBefore))

	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
	require.NoError(t, err)

	poolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, withdrawMsg.PoolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	err = simapp.LiquidityKeeper.WithdrawLiquidityPool(ctx, msgs[0], poolBatch)
	require.NoError(t, err)

	poolCoinAfter := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	withdrawerPoolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.True(t, true, poolCoinAfter.IsZero())
	require.True(t, true, withdrawerPoolCoinAfter.IsZero())
	withdrawerDenomABalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.ReserveCoinDenoms[0])
	withdrawerDenomBBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.ReserveCoinDenoms[1])
	require.Equal(t, deposit.AmountOf(pool.ReserveCoinDenoms[0]), withdrawerDenomABalance.Amount)
	require.Equal(t, deposit.AmountOf(pool.ReserveCoinDenoms[1]), withdrawerDenomBBalance.Amount)

	reserveCoins = simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	require.True(t, reserveCoins.IsZero())

	depositMsg := types.NewMsgDepositWithinBatch(addrs[0], pool.Id, deposit)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	depositMsgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, poolBatch)
	require.Equal(t, 1, len(depositMsgs))

	err = simapp.LiquidityKeeper.DepositLiquidityPool(ctx, depositMsgs[0], poolBatch)
	require.NoError(t, err)

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	depositorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, poolCoin, depositorBalance.Amount)

	reserveCoins = simapp.LiquidityKeeper.GetReserveCoins(ctx, pool)
	require.True(t, reserveCoins.IsEqual(deposit))
}

func TestGetLiquidityPoolMetadata(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	require.Equal(t, 1, len(pools))
	require.Equal(t, uint64(1), pools[0].Id)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextPoolId(ctx)-1)
	require.Equal(t, denomA, pools[0].ReserveCoinDenoms[0])
	require.Equal(t, denomB, pools[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Error(t, err, types.ErrPoolAlreadyExists)

	metaData := simapp.LiquidityKeeper.GetPoolMetaData(ctx, pools[0])
	require.Equal(t, pools[0].Id, metaData.PoolId)

	reserveCoin := simapp.LiquidityKeeper.GetReserveCoins(ctx, pools[0])
	require.Equal(t, reserveCoin, metaData.ReserveCoins)
	require.Equal(t, msg.DepositCoins, metaData.ReserveCoins)

	totalSupply := sdk.NewCoin(pools[0].PoolCoinDenom, simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0]))
	require.Equal(t, totalSupply, metaData.PoolCoinTotalSupply)
	require.Equal(t, creatorBalance, metaData.PoolCoinTotalSupply)
}

func TestIsPoolCoinDenom(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "denomA"
	denomB := "denomB"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)
	getPool, found := simapp.LiquidityKeeper.GetPool(ctx, pool.Id)
	require.True(t, found)
	require.Equal(t, pool, getPool)

	require.Equal(t, "denomA/denomB/1", pool.Name())
	poolCoinDenom := types.GetPoolCoinDenom(pool.Name())
	require.Equal(t, pool.PoolCoinDenom, poolCoinDenom)
	require.True(t, simapp.LiquidityKeeper.IsPoolCoinDenom(ctx, pool.PoolCoinDenom))
	require.False(t, simapp.LiquidityKeeper.IsPoolCoinDenom(ctx, pool.Name()))
}

func TestGetPoolByReserveAccIndex(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	fmt.Println(pool)
	poolStored, found := simapp.LiquidityKeeper.GetPool(ctx, pool.Id)
	require.True(t, found)
	require.Equal(t, pool, poolStored)
	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	require.Equal(t, pool, pools[0])

	poolByReserveAcc, found := simapp.LiquidityKeeper.GetPoolByReserveAccIndex(ctx, pool.GetReserveAccount())
	require.True(t, found)
	require.Equal(t, pool, poolByReserveAcc)

	poolCoinDenom := types.GetPoolCoinDenom(pool.Name())
	require.Equal(t, pool.PoolCoinDenom, poolCoinDenom)
	require.True(t, simapp.LiquidityKeeper.IsPoolCoinDenom(ctx, pool.PoolCoinDenom))
	require.False(t, simapp.LiquidityKeeper.IsPoolCoinDenom(ctx, pool.Name()))
	//SetPoolByReserveAccIndex
}
