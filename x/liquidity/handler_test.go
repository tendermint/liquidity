package liquidity_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestBadMsg(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	handler := liquidity.NewHandler(simapp.LiquidityKeeper)
	_, err := handler(ctx, nil)
	require.ErrorIs(t, err, sdkerrors.ErrUnknownRequest)
	_, err = handler(ctx, banktypes.NewMsgMultiSend(nil, nil))
	require.ErrorIs(t, err, sdkerrors.ErrUnknownRequest)
}

func TestMsgServerCreatePool(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	handler := liquidity.NewHandler(simapp.LiquidityKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(1), lpList[0].Id)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextPoolId(ctx)-1)
	require.Equal(t, denomA, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomB, lpList[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Error(t, err, types.ErrPoolAlreadyExists)
}

func TestMsgServerDepositLiquidityPool(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := app.AddTestAddrs(simapp, ctx, 4, params.LiquidityPoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

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

	createMsg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createMsg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	lp := lpList[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)

	depositMsg := types.NewMsgDepositWithinBatch(addrs[1], lp.Id, deposit)

	handler := liquidity.NewHandler(simapp.LiquidityKeeper)
	_, err = handler(ctx, depositMsg)
	require.NoError(t, err)

	poolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, depositMsg.PoolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	err = simapp.LiquidityKeeper.DepositLiquidityPool(ctx, msgs[0], poolBatch)
	require.NoError(t, err)

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)
	depositorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[1], lp.PoolCoinDenom)
	require.Equal(t, poolCoin.Sub(poolCoinBefore), depositorBalance.Amount)
}

func TestMsgServerWithdrawLiquidityPool(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	createMsg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createMsg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	lp := lpList[0]

	poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)
	withdrawerPoolCoinBefore := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.PoolCoinDenom)

	fmt.Println(poolCoinBefore, withdrawerPoolCoinBefore.Amount)
	require.Equal(t, poolCoinBefore, withdrawerPoolCoinBefore.Amount)
	withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[0], lp.Id, sdk.NewCoin(lp.PoolCoinDenom, poolCoinBefore))

	handler := liquidity.NewHandler(simapp.LiquidityKeeper)
	_, err = handler(ctx, withdrawMsg)
	require.NoError(t, err)

	poolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, withdrawMsg.PoolId)
	require.True(t, found)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, poolBatch)
	require.Equal(t, 1, len(msgs))

	err = simapp.LiquidityKeeper.WithdrawLiquidityPool(ctx, msgs[0], poolBatch)
	require.NoError(t, err)

	poolCoinAfter := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lp)
	withdrawerPoolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.PoolCoinDenom)
	require.True(t, true, poolCoinAfter.IsZero())
	require.True(t, true, withdrawerPoolCoinAfter.IsZero())
	withdrawerDenomAbalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.ReserveCoinDenoms[0])
	withdrawerDenomBbalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lp.ReserveCoinDenoms[1])
	require.Equal(t, deposit.AmountOf(lp.ReserveCoinDenoms[0]).ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerDenomAbalance.Amount)
	require.Equal(t, deposit.AmountOf(lp.ReserveCoinDenoms[1]).ToDec().Mul(sdk.OneDec().Sub(params.WithdrawFeeRate)).TruncateInt(), withdrawerDenomBbalance.Amount)

}

func TestMsgServerGetLiquidityPoolMetadata(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := app.AddTestAddrs(simapp, ctx, 3, params.LiquidityPoolCreationFee)

	denomA := "uETH"
	denomB := "uUSD"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	app.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)

	handler := liquidity.NewHandler(simapp.LiquidityKeeper)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	lpList := simapp.LiquidityKeeper.GetAllPools(ctx)
	require.Equal(t, 1, len(lpList))
	require.Equal(t, uint64(1), lpList[0].Id)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextPoolId(ctx)-1)
	require.Equal(t, denomA, lpList[0].ReserveCoinDenoms[0])
	require.Equal(t, denomB, lpList[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], lpList[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Error(t, err, types.ErrPoolAlreadyExists)

	metaData := simapp.LiquidityKeeper.GetPoolMetaData(ctx, lpList[0])
	require.Equal(t, lpList[0].Id, metaData.PoolId)

	reserveCoin := simapp.LiquidityKeeper.GetReserveCoins(ctx, lpList[0])
	require.Equal(t, reserveCoin, metaData.ReserveCoins)
	require.Equal(t, msg.DepositCoins, metaData.ReserveCoins)

	totalSupply := sdk.NewCoin(lpList[0].PoolCoinDenom, simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, lpList[0]))
	require.Equal(t, totalSupply, metaData.PoolCoinTotalSupply)
	require.Equal(t, creatorBalance, metaData.PoolCoinTotalSupply)
}

func TestMsgServerSwap(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)
	// init test app and context

	// define test denom X, Y for Liquidity Pool
	denomX, denomY := types.AlphabeticalDenomPair("denomX", "denomY")
	X := params.MinInitDepositAmount
	Y := params.MinInitDepositAmount

	// init addresses for the test
	addrs := app.AddTestAddrs(simapp, ctx, 20, params.LiquidityPoolCreationFee)

	// Create pool
	// The create pool msg is not run in batch, but is processed immediately.
	poolId := app.TestCreatePool(t, simapp, ctx, X, Y, denomX, denomY, addrs[0])

	// In case of deposit, withdraw, and swap msg, unlike other normal tx msgs,
	// collect them in the batch and perform an execution at once at the endblock.

	// add a deposit to pool and run batch execution on endblock
	app.TestDepositPool(t, simapp, ctx, X, Y, addrs[1:2], poolId, true)

	// next block, reinitialize batch and increase batchIndex at beginBlocker,
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// Create swap msg for test purposes and put it in the batch.
	price, _ := sdk.NewDecFromStr("1.1")
	priceY, _ := sdk.NewDecFromStr("1.2")
	offerCoinList := []sdk.Coin{sdk.NewCoin(denomX, sdk.NewInt(10000)),
		sdk.NewCoin(denomX, sdk.NewInt(10000)),
		sdk.NewCoin(denomX, sdk.NewInt(10000))}
	offerCoinListY := []sdk.Coin{sdk.NewCoin(denomY, sdk.NewInt(5000))}
	orderPriceList := []sdk.Dec{price, price, price}
	orderPriceListY := []sdk.Dec{priceY}
	orderAddrList := addrs[1:4]
	orderAddrListY := addrs[5:6]

	msg1 := app.GetSwapMsg(t, simapp, ctx, offerCoinList, orderPriceList, orderAddrList, poolId)
	msg4 := app.GetSwapMsg(t, simapp, ctx, offerCoinListY, orderPriceListY, orderAddrListY, poolId)

	handler := liquidity.NewHandler(simapp.LiquidityKeeper)
	_, err := handler(ctx, msg1[0])
	require.NoError(t, err)
	_, err = handler(ctx, msg1[1])
	require.NoError(t, err)
	_, err = handler(ctx, msg1[2])
	require.NoError(t, err)
	_, err = handler(ctx, msg4[0])
	require.NoError(t, err)
	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
	notProcessedMsgs := simapp.LiquidityKeeper.GetAllNotProcessedPoolBatchSwapMsgStates(ctx, batch)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchSwapMsgStatesAsPointer(ctx, batch)
	require.Equal(t, 4, len(msgs))
	require.Equal(t, 4, len(notProcessedMsgs))
}
