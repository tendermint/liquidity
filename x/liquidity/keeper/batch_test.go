package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

const (
	DefaultPoolTypeId = uint32(1)
	DenomX            = "denomX"
	DenomY            = "denomY"
	DenomA            = "denomA"
	DenomB            = "denomB"
)

func TestCreateDepositWithdrawLiquidityPoolToBatch(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	denomA, denomB := types.AlphabeticalDenomPair(DenomA, DenomB)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

	A := sdk.NewInt(1000000000000)
	B := sdk.NewInt(1000000000000)
	depositAB := sdk.NewCoins(sdk.NewCoin(denomA, A), sdk.NewCoin(denomB, B))

	// set accounts for creator, depositor, withdrawer, balance for deposit
	addrs := app.AddTestAddrs(simapp, ctx, 4, params.LiquidityPoolCreationFee)

	app.SaveAccount(simapp, ctx, addrs[0], deposit.Add(depositAB...)) // pool creator
	depositX := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositY := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositX, depositY)
	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], DenomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalanceAB := sdk.NewCoins(depositA, depositB)
	require.Equal(t, deposit, depositBalance)
	require.Equal(t, depositAB, depositBalanceAB)
	feePoolAcc := simapp.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)
	feePoolBalance := simapp.BankKeeper.GetAllBalances(ctx, feePoolAcc)

	// Success case, create Liquidity pool
	poolTypeId := types.DefaultPoolTypeId
	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// Verify PoolCreationFee pay successfully
	feePoolBalance = feePoolBalance.Add(params.LiquidityPoolCreationFee...)
	require.Equal(t, params.LiquidityPoolCreationFee, feePoolBalance)

	// Fail case, reset deposit balance for pool already exists case
	app.SaveAccount(simapp, ctx, addrs[0], deposit)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Equal(t, types.ErrPoolAlreadyExists, err)

	// reset deposit balance without LiquidityPoolCreationFee of pool creator
	// Fail case, insufficient balances for pool creation fee case
	msgAB := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeId, depositBalanceAB)
	app.SaveAccount(simapp, ctx, addrs[0], depositAB)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msgAB)
	require.Equal(t, types.ErrInsufficientPoolCreationFee, err)

	// Success case, create another pool
	msgAB = types.NewMsgCreateLiquidityPool(addrs[0], poolTypeId, depositBalanceAB)
	app.SaveAccount(simapp, ctx, addrs[0], depositAB.Add(params.LiquidityPoolCreationFee...))
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msgAB)
	require.NoError(t, err)

	// Verify PoolCreationFee pay successfully
	feePoolBalance = simapp.BankKeeper.GetAllBalances(ctx, feePoolAcc)
	require.Equal(t, params.LiquidityPoolCreationFee.Add(params.LiquidityPoolCreationFee...), feePoolBalance)

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolId := lpList[0].Id
	require.Equal(t, 2, len(lpList))
	//require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// begin block, init
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// set pool depositor account
	app.SaveAccount(simapp, ctx, addrs[1], deposit) // pool creator
	depositX = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	depositY = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	depositBalance = sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	depositMsg := types.NewMsgDepositToLiquidityPool(addrs[1], poolId, depositBalance)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// check escrow balance of module account
	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleAccEscrowAmtX := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomX)
	moduleAccEscrowAmtY := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomY)
	require.Equal(t, depositX, moduleAccEscrowAmtX)
	require.Equal(t, depositY, moduleAccEscrowAmtY)

	// endblock
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// verify minted pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	depositorPoolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	require.NotEqual(t, sdk.ZeroInt(), depositBalance)
	require.Equal(t, poolCoin, depositorPoolCoinBalance.Amount.Add(creatorBalance.Amount))

	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, bool)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 1)
	require.True(t, msgs[0].Executed)
	require.True(t, msgs[0].Succeeded)
	require.True(t, msgs[0].ToBeDeleted)
	require.Equal(t, uint64(1), batch.Index)

	// error balance after endblock
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	// msg deleted
	_, found := simapp.LiquidityKeeper.GetPoolBatchDepositMsgState(ctx, poolId, msgs[0].MsgIndex)
	require.False(t, found)

	msgs = simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 0)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(1), batch.Index)

	// withdraw
	withdrawerBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	withdrawerBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, lpList[0].PoolCoinDenom)
	require.Equal(t, sdk.ZeroInt(), moduleAccEscrowAmtPool.Amount)
	withdrawMsg := types.NewMsgWithdrawFromLiquidityPool(addrs[1], poolId, withdrawerBalancePoolCoinBefore)
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
	require.NoError(t, err)

	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount.Add(depositorPoolCoinBalance.Amount))

	// check escrow balance of module account
	moduleAccEscrowAmtPool = simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, lpList[0].PoolCoinDenom)
	require.Equal(t, withdrawerBalancePoolCoinBefore, moduleAccEscrowAmtPool)

	// endblock
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// verify burned pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	require.Equal(t, depositX.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceX.Amount)
	require.Equal(t, depositY.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	batch, bool = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, bool)
	withdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Len(t, withdrawMsgs, 1)
	require.True(t, withdrawMsgs[0].Executed)
	require.True(t, withdrawMsgs[0].Succeeded)
	require.True(t, withdrawMsgs[0].ToBeDeleted)
	require.Equal(t, uint64(1), batch.Index)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// msg deleted
	withdrawMsgs = simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Len(t, withdrawMsgs, 0)
	_, found = simapp.LiquidityKeeper.GetPoolBatchWithdrawMsgState(ctx, poolId, 0)
	require.False(t, found)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(2), batch.Index)
	require.False(t, batch.Executed)
}

func TestCreateDepositWithdrawLiquidityPoolToBatch2(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	deposit2 := sdk.NewCoins(sdk.NewCoin(denomX, X.QuoRaw(2)), sdk.NewCoin(denomY, Y.QuoRaw(2)))

	// set accounts for creator, depositor, withdrawer, balance for deposit
	params := simapp.LiquidityKeeper.GetParams(ctx)
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)
	app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
	depositX := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositY := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeId := types.DefaultPoolTypeId
	msg := types.NewMsgCreateLiquidityPool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolId := lpList[0].Id
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, lpList[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// begin block, init
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// set pool depositor account
	app.SaveAccount(simapp, ctx, addrs[1], deposit2) // pool creator
	depositX = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	depositY = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	depositBalance = sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit2, depositBalance)

	depositMsg := types.NewMsgDepositToLiquidityPool(addrs[1], poolId, depositBalance)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// check escrow balance of module account
	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleAccEscrowAmtX := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomX)
	moduleAccEscrowAmtY := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomY)
	require.Equal(t, depositX, moduleAccEscrowAmtX)
	require.Equal(t, depositY, moduleAccEscrowAmtY)

	// endblock
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// verify minted pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	depositorPoolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	require.NotEqual(t, sdk.ZeroInt(), depositBalance)
	require.Equal(t, poolCoin, depositorPoolCoinBalance.Amount.Add(creatorBalance.Amount))

	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, bool)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 1)
	require.True(t, msgs[0].Executed)
	require.True(t, msgs[0].Succeeded)
	require.True(t, msgs[0].ToBeDeleted)
	require.Equal(t, uint64(1), batch.Index)

	// error balance after endblock
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	// msg deleted
	_, found := simapp.LiquidityKeeper.GetPoolBatchDepositMsgState(ctx, poolId, msgs[0].MsgIndex)
	require.False(t, found)

	msgs = simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 0)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(1), batch.Index)

	// withdraw
	withdrawerBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	withdrawerBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, lpList[0].PoolCoinDenom)
	require.Equal(t, sdk.ZeroInt(), moduleAccEscrowAmtPool.Amount)
	withdrawMsg := types.NewMsgWithdrawFromLiquidityPool(addrs[1], poolId, withdrawerBalancePoolCoinBefore)
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
	require.NoError(t, err)

	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount.Add(depositorPoolCoinBalance.Amount))

	// check escrow balance of module account
	moduleAccEscrowAmtPool = simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, lpList[0].PoolCoinDenom)
	require.Equal(t, withdrawerBalancePoolCoinBefore, moduleAccEscrowAmtPool)

	// endblock
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// verify burned pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin = simapp.BankKeeper.GetBalance(ctx, addrs[1], lpList[0].PoolCoinDenom)
	require.Equal(t, depositX.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceX.Amount)
	require.Equal(t, depositY.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	batch, bool = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, bool)
	withdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Len(t, withdrawMsgs, 1)
	require.True(t, withdrawMsgs[0].Executed)
	require.True(t, withdrawMsgs[0].Succeeded)
	require.True(t, withdrawMsgs[0].ToBeDeleted)
	require.Equal(t, uint64(1), batch.Index)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// msg deleted
	withdrawMsgs = simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Len(t, withdrawMsgs, 0)
	_, found = simapp.LiquidityKeeper.GetPoolBatchWithdrawMsgState(ctx, poolId, 0)
	require.False(t, found)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(2), batch.Index)
	require.False(t, batch.Executed)
}

func TestLiquidityScenario(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)
	//denoms := []string{denomX, denomY}

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	poolId2 := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, "testDenom", addrs[0])
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, uint64(2), poolId2)

	// begin block, init
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:10], poolId, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	_, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)

	// msg deleted
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 0)

	//balance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	//balance = simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	//require.Len(t, balance)

	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500000), addrs[1:10], poolId, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// msg deleted
	withdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Len(t, withdrawMsgs, 0)
	_, found = simapp.LiquidityKeeper.GetPoolBatchWithdrawMsgState(ctx, poolId, 0)
	require.False(t, found)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(2), batch.Index)
	require.False(t, batch.Executed)
}

func TestLiquidityScenario2(t *testing.T) {
	simapp, ctx := createTestInput()
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

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	price, _ := sdk.NewDecFromStr("1.1")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000))}
	orderPriceList := []sdk.Dec{price}
	orderAddrList := addrs[1:2]
	batchMsgs, batch := app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	fmt.Println(batchMsgs, batch)
	batchMsgs, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	batchMsgs, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	fmt.Println(batchMsgs, batch)
	batchMsgs, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, true)
	fmt.Println(batchMsgs, batch)
}

func TestLiquidityScenario3(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(5000), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(5000), addrs[2:3], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[2:3], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[2:3], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// refund Deposit scenario
func TestLiquidityScenario4(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:2], poolId, false)
	balanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.ZeroInt(), balanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), balanceY.Amount)
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	simapp.LiquidityKeeper.DeletePool(ctx, pool)
	pool, found = simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.False(t, found)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	balanceXrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, X, balanceXrefunded.Amount)
	require.Equal(t, Y, balanceYrefunded.Amount)
	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// refund Withdraw scenario, TODO: fix
func TestLiquidityScenario5(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	poolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	app.TestWithdrawPool(t, simapp, ctx, poolCoin.Amount, addrs[0:1], poolId, false)

	poolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.ZeroInt(), poolCoinAfter.Amount)

	PoolCoinDenom := pool.PoolCoinDenom
	simapp.LiquidityKeeper.DeletePool(ctx, pool)
	pool, found = simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.False(t, found)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[0], PoolCoinDenom)
	require.Equal(t, poolCoin.Amount, poolCoinRefunded.Amount)
	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// Verification of results when deposit 30A and 20B is put into a reserve pool with 100A and 200B deposited
// state : 100A, 200B, 10PoolCoin(total supply)
// deposit 30A, 20B ->
// - 10A, 20B
// - 1 PoolCoin received
// - 20A refunded
func TestLiquidityScenario6(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	poolCoins := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	app.TestDepositPool(t, simapp, ctx, sdk.NewInt(30000000), sdk.NewInt(20000000), addrs[1:2], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(100000), poolCoinBalance.Amount)
	require.Equal(t, poolCoins.QuoRaw(10), poolCoinBalance.Amount)

	balanceXrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.NewInt(20000000), balanceXrefunded.Amount)
	require.Equal(t, sdk.ZeroInt(), balanceYrefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// Verification of results when deposit 10A and 20B is put into a reserve pool with 100A and 200B deposited
// state : 100A, 200B, 10PoolCoin(total supply)
// deposit 10A, 30B ->
// - 10A, 20B
// - 1 PoolCoin received
// - 10B refunded
func TestLiquidityScenario7(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	poolCoins := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	app.TestDepositPool(t, simapp, ctx, sdk.NewInt(10000000), sdk.NewInt(30000000), addrs[1:2], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(100000), poolCoinBalance.Amount)
	require.Equal(t, poolCoins.QuoRaw(10), poolCoinBalance.Amount)

	balanceXrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYrefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.ZeroInt(), balanceXrefunded.Amount)
	require.Equal(t, sdk.NewInt(10000000), balanceYrefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// Verification of results when withdraw 1 PoolCoin from a reserve pool with 100A and 200B deposited
// state : 100A, 200B, 10PoolCoin(total supply)
// withdraw 1 PoolCoin ->
// - 1 PoolCoin burned
// - 10A, 20B received
func TestLiquidityScenario8(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	poolCoins := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(1000000), poolCoins)
	require.Equal(t, sdk.NewInt(1000000), poolCoinBalance.Amount)
	app.TestWithdrawPool(t, simapp, ctx, poolCoins.QuoRaw(10), addrs[0:1], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoins = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	poolCoinBalance = simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(900000), poolCoins)
	require.Equal(t, sdk.NewInt(900000), poolCoinBalance.Amount)
	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// Test UnitBatchSize when over 1
func TestLiquidityUnitBatchSize(t *testing.T) {
	simapp, ctx := createTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)
	params.UnitBatchSize = 2
	simapp.LiquidityKeeper.SetParams(ctx, params)

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	poolCoins := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(1000000), poolCoins)
	require.Equal(t, sdk.NewInt(1000000), poolCoinBalance.Amount)
	app.TestWithdrawPool(t, simapp, ctx, poolCoins.QuoRaw(10), addrs[0:1], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// batch not executed, 1 >= 2(UnitBatchSize)
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, pool.Id)
	require.True(t, found)
	require.False(t, batch.Executed)
	batchWithdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Equal(t, 1, len(batchWithdrawMsgs))

	poolCoins = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	poolCoinBalance = simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(1000000), poolCoins)
	require.Equal(t, sdk.NewInt(900000), poolCoinBalance.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	batchWithdrawMsgs = simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Equal(t, 1, len(batchWithdrawMsgs))
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// batch executed, 2 >= 2(UnitBatchSize)
	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, pool.Id)
	require.True(t, found)
	require.True(t, batch.Executed)
	batchWithdrawMsgs = simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Equal(t, 1, len(batchWithdrawMsgs))

	poolCoins = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	poolCoinBalance = simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(900000), poolCoins)
	require.Equal(t, sdk.NewInt(900000), poolCoinBalance.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// batch msg deleted after batch execution
	batchWithdrawMsgs = simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Equal(t, 0, len(batchWithdrawMsgs))
}

func TestInitNextBatch(t *testing.T) {
	simapp, ctx := createTestInput()
	pool := types.Pool{
		Id:                    1,
		TypeId:                1,
		ReserveCoinDenoms:     nil,
		ReserveAccountAddress: "",
		PoolCoinDenom:         "",
	}
	simapp.LiquidityKeeper.SetPool(ctx, pool)

	batch := types.NewPoolBatch(pool.Id, 1)

	simapp.LiquidityKeeper.SetPoolBatch(ctx, batch)
	simapp.LiquidityKeeper.SetPoolBatchIndex(ctx, batch.PoolId, batch.Index)
	err := simapp.LiquidityKeeper.InitNextBatch(ctx, batch)
	require.Error(t, err)

	batch.Executed = true
	simapp.LiquidityKeeper.SetPoolBatch(ctx, batch)

	err = simapp.LiquidityKeeper.InitNextBatch(ctx, batch)
	require.NoError(t, err)

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.False(t, batch.Executed)
	require.Equal(t, uint64(2), batch.Index)

}

func TestDeleteAndInitPoolBatchDeposit(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)

	depositsAll := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Equal(t, 6, len(depositsAll))
	depositsAll[0].Executed = true
	depositsAll[0].ToBeDeleted = false
	simapp.LiquidityKeeper.SetPoolBatchDepositMsgStates(ctx, poolId, depositsAll)
	depositsRemaining := simapp.LiquidityKeeper.GetAllRemainingPoolBatchDepositMsgStates(ctx, batch)
	batch.Executed = true
	simapp.LiquidityKeeper.SetPoolBatch(ctx, batch)
	simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)
	depositsAfter := simapp.LiquidityKeeper.GetAllRemainingPoolBatchDepositMsgStates(ctx, batch)

	fmt.Println(depositsAll)
	fmt.Println(depositsRemaining)
	fmt.Println(depositsAfter)
	require.Equal(t, 1, len(depositsRemaining))
	require.Equal(t, 0, len(depositsAfter))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}
func TestDeleteAndInitPoolBatchWithdraw(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(5000), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[1:2], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(5000), addrs[2:3], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500), addrs[2:3], poolId, false)
	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(50), addrs[2:3], poolId, false)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)

	withdrawsAll := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Equal(t, 6, len(withdrawsAll))
	withdrawsAll[0].Executed = true
	withdrawsAll[0].ToBeDeleted = false
	simapp.LiquidityKeeper.SetPoolBatchWithdrawMsgStates(ctx, poolId, withdrawsAll)
	withdrawsRemaining := simapp.LiquidityKeeper.GetAllRemainingPoolBatchWithdrawMsgStates(ctx, batch)
	batch.Executed = true
	simapp.LiquidityKeeper.SetPoolBatch(ctx, batch)
	simapp.LiquidityKeeper.DeleteAndInitPoolBatch(ctx)
	withdrawsAfter := simapp.LiquidityKeeper.GetAllRemainingPoolBatchWithdrawMsgStates(ctx, batch)
	require.Equal(t, 1, len(withdrawsRemaining))
	require.Equal(t, 0, len(withdrawsAfter))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}
