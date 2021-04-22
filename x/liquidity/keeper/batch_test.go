package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

const (
	DenomX = "denomX"
	DenomY = "denomY"
	DenomA = "denomA"
	DenomB = "denomB"
)

func TestBadDeposit(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)

	depositCoins := sdk.NewCoins(sdk.NewCoin(DenomX, params.MinInitDepositAmount), sdk.NewCoin(DenomY, params.MinInitDepositAmount))
	depositorAddr := app.AddRandomTestAddr(simapp, ctx, depositCoins.Add(params.PoolCreationFee...))

	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, &types.MsgCreatePool{
		PoolCreatorAddress: depositorAddr.String(),
		PoolTypeId:         types.DefaultPoolTypeId,
		DepositCoins:       depositCoins,
	})
	require.NoError(t, err)

	// deposit with empty message
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, &types.MsgDepositWithinBatch{})
	require.Error(t, err)

	// deposit coins more than it has
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, &types.MsgDepositWithinBatch{
		DepositorAddress: depositorAddr.String(),
		PoolId:           pool.Id,
		DepositCoins:     sdk.NewCoins(sdk.NewCoin(DenomX, sdk.OneInt()), sdk.NewCoin(DenomY, sdk.OneInt())),
	})
	require.Error(t, err)

	// forcefully delete current pool batch
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, pool.Id)
	require.True(t, found)
	simapp.LiquidityKeeper.DeletePoolBatch(ctx, batch)
	// deposit coins when there is no pool batch
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, &types.MsgDepositWithinBatch{
		DepositorAddress: depositorAddr.String(),
		PoolId:           pool.Id,
		DepositCoins:     sdk.NewCoins(sdk.NewCoin(DenomX, sdk.OneInt()), sdk.NewCoin(DenomY, sdk.OneInt())),
	})
	require.ErrorIs(t, err, types.ErrPoolBatchNotExists)
}

func TestBadWithdraw(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)

	depositCoins := sdk.NewCoins(sdk.NewCoin(DenomX, params.MinInitDepositAmount), sdk.NewCoin(DenomY, params.MinInitDepositAmount))
	depositorAddr := app.AddRandomTestAddr(simapp, ctx, depositCoins.Add(params.PoolCreationFee...))

	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, &types.MsgCreatePool{
		PoolCreatorAddress: depositorAddr.String(),
		PoolTypeId:         types.DefaultPoolTypeId,
		DepositCoins:       depositCoins,
	})
	require.NoError(t, err)

	// withdraw with empty message
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, &types.MsgWithdrawWithinBatch{})
	require.Error(t, err)

	balance := simapp.BankKeeper.GetBalance(ctx, depositorAddr, pool.PoolCoinDenom)

	// mint extra pool coins to test if below fails
	require.NoError(t, simapp.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, sdk.NewInt(1000)))))
	// withdraw pool coins more than it has
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, &types.MsgWithdrawWithinBatch{
		WithdrawerAddress: depositorAddr.String(),
		PoolId:            pool.Id,
		PoolCoin:          balance.Add(sdk.NewCoin(pool.PoolCoinDenom, sdk.OneInt())),
	})
	require.Error(t, err)

	// forcefully delete current pool batch
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, pool.Id)
	require.True(t, found)
	simapp.LiquidityKeeper.DeletePoolBatch(ctx, batch)
	// withdraw pool coins when there is no pool batch
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, &types.MsgWithdrawWithinBatch{
		WithdrawerAddress: depositorAddr.String(),
		PoolId:            pool.Id,
		PoolCoin:          sdk.NewCoin(pool.PoolCoinDenom, sdk.OneInt()),
	})
	require.ErrorIs(t, err, types.ErrPoolBatchNotExists)
}

func TestBadSwap(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)

	depositCoins := sdk.NewCoins(sdk.NewCoin(DenomX, params.MinInitDepositAmount), sdk.NewCoin(DenomY, params.MinInitDepositAmount))
	depositorAddr := app.AddRandomTestAddr(simapp, ctx, depositCoins.Add(params.PoolCreationFee...))

	pool, err := simapp.LiquidityKeeper.CreatePool(ctx, &types.MsgCreatePool{
		PoolCreatorAddress: depositorAddr.String(),
		PoolTypeId:         types.DefaultPoolTypeId,
		DepositCoins:       depositCoins,
	})
	require.NoError(t, err)

	// swap with empty message
	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, &types.MsgSwapWithinBatch{}, 0)
	require.Error(t, err)

	orderPrice := sdk.OneDec()

	// swap coin more than it has
	offerCoin := sdk.NewCoin(DenomX, sdk.NewInt(100000))
	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, &types.MsgSwapWithinBatch{
		SwapRequesterAddress: depositorAddr.String(),
		PoolId:               pool.Id,
		SwapTypeId:           types.DefaultSwapTypeId,
		OfferCoin:            offerCoin,
		DemandCoinDenom:      DenomY,
		OfferCoinFee:         types.GetOfferCoinFee(offerCoin, params.SwapFeeRate),
		OrderPrice:           orderPrice,
	}, 0)
	require.Error(t, err)

	// when swap fails, user's balance should never change
	require.NoError(t, simapp.BankKeeper.SetBalance(ctx, depositorAddr, offerCoin))
	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, &types.MsgSwapWithinBatch{
		SwapRequesterAddress: depositorAddr.String(),
		PoolId:               pool.Id,
		SwapTypeId:           types.DefaultSwapTypeId,
		OfferCoin:            offerCoin,
		DemandCoinDenom:      DenomY,
		OfferCoinFee:         types.GetOfferCoinFee(offerCoin, params.SwapFeeRate),
		OrderPrice:           orderPrice,
	}, 0)
	require.Error(t, err)
	balance := simapp.BankKeeper.GetBalance(ctx, depositorAddr, DenomX)
	require.Equal(t, offerCoin, balance)

	// forcefully delete current pool batch
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, pool.Id)
	require.True(t, found)
	simapp.LiquidityKeeper.DeletePoolBatch(ctx, batch)
	// swap coin when there is no pool batch
	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, &types.MsgSwapWithinBatch{
		SwapRequesterAddress: depositorAddr.String(),
		PoolId:               pool.Id,
		SwapTypeId:           types.DefaultSwapTypeId,
		OfferCoin:            offerCoin,
		DemandCoinDenom:      DenomY,
		OfferCoinFee:         types.GetOfferCoinFee(offerCoin, params.SwapFeeRate),
		OrderPrice:           orderPrice,
	}, 0)
	require.ErrorIs(t, err, types.ErrPoolBatchNotExists)
}

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
	addrs := app.AddTestAddrs(simapp, ctx, 4, params.PoolCreationFee)

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
	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance, nil)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// Verify PoolCreationFee pay successfully
	feePoolBalance = feePoolBalance.Add(params.PoolCreationFee...)
	require.Equal(t, params.PoolCreationFee, feePoolBalance)

	// Fail case, reset deposit balance for pool already exists case
	app.SaveAccount(simapp, ctx, addrs[0], deposit)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Equal(t, types.ErrPoolAlreadyExists, err)

	// reset deposit balance without PoolCreationFee of pool creator
	// Fail case, insufficient balances for pool creation fee case
	msgAB := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalanceAB, nil)
	app.SaveAccount(simapp, ctx, addrs[0], depositAB)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msgAB)
	require.Equal(t, types.ErrInsufficientPoolCreationFee, err)

	// Success case, create another pool
	msgAB = types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalanceAB, nil)
	app.SaveAccount(simapp, ctx, addrs[0], depositAB.Add(params.PoolCreationFee...))
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msgAB)
	require.NoError(t, err)

	// Verify PoolCreationFee pay successfully
	feePoolBalance = simapp.BankKeeper.GetAllBalances(ctx, feePoolAcc)
	require.Equal(t, params.PoolCreationFee.Add(params.PoolCreationFee...), feePoolBalance)

	// verify created liquidity pool
	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolId := pools[0].Id
	require.Equal(t, 2, len(pools))
	//require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, pools[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, pools[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// begin block, init
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// set pool depositor account
	app.SaveAccount(simapp, ctx, addrs[1], deposit) // pool creator
	depositX = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	depositY = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	depositBalance = sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	depositMsg := types.NewMsgDepositWithinBatch(addrs[1], poolId, depositBalance)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
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
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	depositorPoolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	require.NotEqual(t, sdk.ZeroInt(), depositBalance)
	require.Equal(t, poolCoin, depositorPoolCoinBalance.Amount.Add(creatorBalance.Amount))

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 1)
	require.True(t, msgs[0].Executed)
	require.True(t, msgs[0].Succeeded)
	require.True(t, msgs[0].ToBeDeleted)
	require.Equal(t, uint64(1), batch.Index)

	// error balance after endblock
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	// msg deleted
	_, found = simapp.LiquidityKeeper.GetPoolBatchDepositMsgState(ctx, poolId, msgs[0].MsgIndex)
	require.False(t, found)

	msgs = simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 0)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(1), batch.Index)

	// withdraw
	withdrawerBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	withdrawerBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pools[0].PoolCoinDenom)
	require.Equal(t, sdk.ZeroInt(), moduleAccEscrowAmtPool.Amount)
	withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[1], poolId, withdrawerBalancePoolCoinBefore)
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
	require.NoError(t, err)

	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount.Add(depositorPoolCoinBalance.Amount))

	// check escrow balance of module account
	moduleAccEscrowAmtPool = simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pools[0].PoolCoinDenom)
	require.Equal(t, withdrawerBalancePoolCoinBefore, moduleAccEscrowAmtPool)

	// endblock
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// verify burned pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	require.Equal(t, depositX.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceX.Amount)
	require.Equal(t, depositY.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
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
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)
	app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
	depositX := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
	depositY := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
	depositBalance := sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeId := types.DefaultPoolTypeId
	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance, nil)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	poolId := pools[0].Id
	require.Equal(t, 1, len(pools))
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, denomX, pools[0].ReserveCoinDenoms[0])
	require.Equal(t, denomY, pools[0].ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	// begin block, init
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// set pool depositor account
	app.SaveAccount(simapp, ctx, addrs[1], deposit2) // pool creator
	depositX = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	depositY = simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	depositBalance = sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit2, depositBalance)

	depositMsg := types.NewMsgDepositWithinBatch(addrs[1], poolId, depositBalance)
	_, err = simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
	require.NoError(t, err)

	depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
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
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	depositorPoolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	require.NotEqual(t, sdk.ZeroInt(), depositBalance)
	require.Equal(t, poolCoin, depositorPoolCoinBalance.Amount.Add(creatorBalance.Amount))

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 1)
	require.True(t, msgs[0].Executed)
	require.True(t, msgs[0].Succeeded)
	require.True(t, msgs[0].ToBeDeleted)
	require.Equal(t, uint64(1), batch.Index)

	// error balance after endblock
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	depositorBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	depositorBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	require.Equal(t, sdk.ZeroInt(), depositorBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), depositorBalanceY.Amount)
	require.Equal(t, denomX, depositorBalanceX.Denom)
	require.Equal(t, denomY, depositorBalanceY.Denom)
	// msg deleted
	_, found = simapp.LiquidityKeeper.GetPoolBatchDepositMsgState(ctx, poolId, msgs[0].MsgIndex)
	require.False(t, found)

	msgs = simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 0)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(1), batch.Index)

	// withdraw
	withdrawerBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	withdrawerBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pools[0].PoolCoinDenom)
	require.Equal(t, sdk.ZeroInt(), moduleAccEscrowAmtPool.Amount)
	withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[1], poolId, withdrawerBalancePoolCoinBefore)
	_, err = simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
	require.NoError(t, err)

	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount.Add(depositorPoolCoinBalance.Amount))

	// check escrow balance of module account
	moduleAccEscrowAmtPool = simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pools[0].PoolCoinDenom)
	require.Equal(t, withdrawerBalancePoolCoinBefore, moduleAccEscrowAmtPool)

	// endblock
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// verify burned pool coin
	poolCoin = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	withdrawerBalanceX = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[0])
	withdrawerBalanceY = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].ReserveCoinDenoms[1])
	withdrawerBalancePoolCoin = simapp.BankKeeper.GetBalance(ctx, addrs[1], pools[0].PoolCoinDenom)
	require.Equal(t, depositX.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceX.Amount)
	require.Equal(t, depositY.Amount.ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerBalanceY.Amount)
	require.Equal(t, sdk.ZeroInt(), withdrawerBalancePoolCoin.Amount)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
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

// This scenario tests simple create pool, deposit to the pool, and withdraw from the pool
func TestLiquidityScenario(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create two pools with the different denomY of 1000X and 1000Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	poolId2 := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, "testDenom", addrs[0])
	require.Equal(t, uint64(1), poolId)
	require.Equal(t, uint64(2), poolId2)

	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:10], poolId, true)

	// next block starts
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	_, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)

	// deposit message is expected to be handled
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
	require.Len(t, msgs, 0)

	app.TestWithdrawPool(t, simapp, ctx, sdk.NewInt(500000), addrs[1:10], poolId, true)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// withdraw message is expected to be handled
	withdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
	require.Len(t, withdrawMsgs, 0)

	_, found = simapp.LiquidityKeeper.GetPoolBatchWithdrawMsgState(ctx, poolId, 0)
	require.False(t, found)

	batch, found = simapp.LiquidityKeeper.GetPoolBatch(ctx, batch.PoolId)
	require.True(t, found)
	require.Equal(t, uint64(2), batch.Index)
	require.False(t, batch.Executed)
}

// This scenario tests create pool, deposit to the pool, and swap coins
func TestLiquidityScenario2(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(1000000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 1000X and 1000Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])
	require.Equal(t, uint64(1), poolId)

	// make two different deposits to the same pool with different amounts of coins
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
	batchMsgs, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	batchMsgs, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, false)
	batchMsgs, batch = app.TestSwapPool(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId, true)
	require.NotNil(t, batchMsgs)
	require.NotNil(t, batch)
}

// This scenario tests to executed accumulated deposit and withdraw pool batches
func TestLiquidityScenario3(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 1000X and 500Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// make 6 different deposits to the same pool with different amounts of coins
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X.QuoRaw(10), Y, addrs[1:2], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)
	app.TestDepositPool(t, simapp, ctx, X, Y.QuoRaw(10), addrs[2:3], poolId, false)

	// execute accumulated deposit batches
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

	// execute accumulated withdraw batches
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// This scenario tests deposit refund scenario
func TestLiquidityScenario4(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 1000X and 500Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// make deposit to the pool with 1000X and 500Y coins
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:2], poolId, false)

	// balance should be zero since accounts' balances are expected to be in an escrow account
	balanceX := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceY := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.ZeroInt(), balanceX.Amount)
	require.Equal(t, sdk.ZeroInt(), balanceY.Amount)

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	// delete previously created pool
	simapp.LiquidityKeeper.DeletePool(ctx, pool)

	pool, found = simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.False(t, found)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	balanceXRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, X, balanceXRefunded.Amount)
	require.Equal(t, Y, balanceYRefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// This scenario tests refund withdraw scenario
func TestLiquidityScenario5(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(1000000000)
	Y := sdk.NewInt(500000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 1000X and 500Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	poolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)

	// withdraw all pool coin from the pool
	app.TestWithdrawPool(t, simapp, ctx, poolCoin.Amount, addrs[0:1], poolId, false)

	poolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.ZeroInt(), poolCoinAfter.Amount)

	// save pool coin denom before deleting the pool
	poolCoinDenom := pool.PoolCoinDenom

	// delete the pool
	simapp.LiquidityKeeper.DeletePool(ctx, pool)

	pool, found = simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.False(t, found)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	// pool coin should be refunded since the pool is deleted before executing pool batch
	poolCoinRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[0], poolCoinDenom)
	require.Equal(t, poolCoin.Amount, poolCoinRefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// This scenario tests pool coin and refunded amounts after depositing X and Y amounts of coins
// - 100X and 200Y in reserve pool
// - deposit 30X and 20Y coins
// - test how many pool coins to receive
// - test how many X or Y coins to be refunded
func TestLiquidityScenario6(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 100X and 200Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	poolCoinTotalSupply := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

	// use the other account to make deposit to the pool with 30X and 20Y coins
	app.TestDepositPool(t, simapp, ctx, sdk.NewInt(30000000), sdk.NewInt(20000000), addrs[1:2], poolId, false)

	// execute pool batch
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(100000), poolCoinBalance.Amount)
	require.Equal(t, poolCoinTotalSupply.QuoRaw(10), poolCoinBalance.Amount)

	balanceXRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.NewInt(20000000), balanceXRefunded.Amount)
	require.Equal(t, sdk.ZeroInt(), balanceYRefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// This scenario is similar with scenario6
// Depositing different amounts will result in different amount of refunded amounts
// - 100X and 200Y in reserve pool
// - deposit 10X and 30Y coins
// - test how many pool coins to receive
// - test how many X or Y coins to be refunded
func TestLiquidityScenario7(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 100X and 200Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	poolCoinTotalSupply := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

	// use the other account to make deposit to the pool with 10X and 30Y coins
	app.TestDepositPool(t, simapp, ctx, sdk.NewInt(10000000), sdk.NewInt(30000000), addrs[1:2], poolId, false)

	// execute pool batch
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(100000), poolCoinBalance.Amount)
	require.Equal(t, poolCoinTotalSupply.QuoRaw(10), poolCoinBalance.Amount)

	balanceXRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomX)
	balanceYRefunded := simapp.BankKeeper.GetBalance(ctx, addrs[1], denomY)
	require.Equal(t, sdk.ZeroInt(), balanceXRefunded.Amount)
	require.Equal(t, sdk.NewInt(10000000), balanceYRefunded.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// This scenario tests to withdraw amounts from reserve pool to see the impacts of how pool coin and account's balance.
// - 100X and 200Y in reserve pool
// - withdraw 10th of pool coin total supply
// - test pool coin total supply
// - test account's coin balance
func TestLiquidityScenario8(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())

	denomX, denomY := types.AlphabeticalDenomPair(DenomX, DenomY)

	X := sdk.NewInt(100000000)
	Y := sdk.NewInt(200000000)

	// create 20 random accounts with an initial balance of 0.01
	addrs := app.AddTestAddrsIncremental(simapp, ctx, 20, sdk.NewInt(10000))

	// create pool with 100X and 200Y coins
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	poolCoinTotalSupply := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

	poolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(1000000), poolCoinTotalSupply)
	require.Equal(t, sdk.NewInt(1000000), poolCoinBalance.Amount)

	// withdraw 10th of poolCoinTotalSupply from the pool
	app.TestWithdrawPool(t, simapp, ctx, poolCoinTotalSupply.QuoRaw(10), addrs[0:1], poolId, false)

	// execute pool batch
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

	poolCoinTotalSupply = simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

	poolCoinBalance = simapp.BankKeeper.GetBalance(ctx, addrs[0], pool.PoolCoinDenom)
	require.Equal(t, sdk.NewInt(900000), poolCoinTotalSupply)
	require.Equal(t, sdk.NewInt(900000), poolCoinBalance.Amount)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
}

// Test UnitBatchHeight when over 1
func TestLiquidityUnitBatchHeight(t *testing.T) {
	simapp, ctx := createTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)
	params.UnitBatchHeight = 2
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

	// batch not executed, 1 >= 2(UnitBatchHeight)
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

	// batch executed, 2 >= 2(UnitBatchHeight)
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
