package keeper_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestSimulationSwapExecutionFindEdgeCase(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		r := rand.New(rand.NewSource(seed))

		simapp, ctx := createTestInput()
		params := simapp.LiquidityKeeper.GetParams(ctx)

		// define test denom X, Y for Liquidity Pool
		denomX := "denomX"
		denomY := "denomY"
		denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

		// get random X, Y amount for create pool
		param := simapp.LiquidityKeeper.GetParams(ctx)
		X, Y := app.GetRandPoolAmt(r, param.MinInitDepositAmount)
		deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

		// set pool creator account, balance for deposit
		addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)
		app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
		depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
		depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
		depositBalance := sdk.NewCoins(depositA, depositB)
		require.Equal(t, deposit, depositBalance)

		// create Liquidity pool
		poolTypeId := types.DefaultPoolTypeId
		msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
		_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
		require.NoError(t, err)

		for i := 0; i < 20; i++ {
			ctx = ctx.WithBlockHeight(int64(i))
			testSwapEdgeCases(t, r, simapp, ctx, X, Y, depositBalance, addrs)
		}
	}
}

func TestSwapExecution(t *testing.T) {
	// TODO: to with simulation, invariants, ransim
	for seed := int64(0); seed < 50; seed++ {
		s := rand.NewSource(seed)
		r := rand.New(s)
		simapp, ctx := createTestInput()
		simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
		params := simapp.LiquidityKeeper.GetParams(ctx)

		// define test denom X, Y for Liquidity Pool
		denomX := "denomX"
		denomY := "denomY"
		denomX, denomY = types.AlphabeticalDenomPair(denomX, denomY)

		// get random X, Y amount for create pool
		X, Y := app.GetRandPoolAmt(r, params.MinInitDepositAmount)
		deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

		// set pool creator account, balance for deposit
		addrs := app.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)
		app.SaveAccount(simapp, ctx, addrs[0], deposit) // pool creator
		depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomX)
		depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomY)
		depositBalance := sdk.NewCoins(depositA, depositB)
		require.Equal(t, deposit, depositBalance)

		// create Liquidity pool
		poolTypeId := types.DefaultPoolTypeId
		msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
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

		var XtoY []*types.MsgSwapWithinBatch // buying Y from X
		var YtoX []*types.MsgSwapWithinBatch // selling Y for X

		// make random orders, set buyer, seller accounts for the orders
		XtoY, YtoX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 250, 250)
		buyerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
		sellerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

		for i, msg := range XtoY {
			app.SaveAccountWithFee(simapp, ctx, buyerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = buyerAddrs[i].String()
			msg.PoolId = poolId
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
		for i, msg := range YtoX {
			app.SaveAccountWithFee(simapp, ctx, sellerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = sellerAddrs[i].String()
			msg.PoolId = poolId
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}

		// begin block, delete and init pool batch
		liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

		// handle msgs, set order msgs to batch
		for _, msg := range XtoY {
			_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
			require.NoError(t, err)
		}
		for _, msg := range YtoX {
			_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, 0)
			require.NoError(t, err)
		}

		// verify pool batch
		liquidityPoolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
		require.True(t, found)
		require.NotNil(t, liquidityPoolBatch)

		// end block, swap execution
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	}
}

func testSwapEdgeCases(t *testing.T, r *rand.Rand, simapp *app.LiquidityApp, ctx sdk.Context, X, Y sdk.Int, depositBalance sdk.Coins, addrs []sdk.AccAddress) {
	//simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	denomX := depositBalance[0].Denom
	denomY := depositBalance[1].Denom

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

	var XtoY []*types.MsgSwapWithinBatch // buying Y from X
	var YtoX []*types.MsgSwapWithinBatch // selling Y for X

	batch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)

	remainingSwapMsgs := simapp.LiquidityKeeper.GetAllNotProcessedPoolBatchSwapMsgStates(ctx, batch)
	if ctx.BlockHeight() == 0 || len(remainingSwapMsgs) == 0 {
		// make random orders, set buyer, seller accounts for the orders
		XtoY, YtoX = app.GetRandomSizeOrders(denomX, denomY, X, Y, r, 100, 100)
		buyerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(XtoY), sdk.NewInt(0))
		sellerAddrs := app.AddTestAddrsIncremental(simapp, ctx, len(YtoX), sdk.NewInt(0))

		for i, msg := range XtoY {
			app.SaveAccountWithFee(simapp, ctx, buyerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = buyerAddrs[i].String()
			msg.PoolId = poolId
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
		for i, msg := range YtoX {
			app.SaveAccountWithFee(simapp, ctx, sellerAddrs[i], sdk.NewCoins(msg.OfferCoin), msg.OfferCoin)
			msg.SwapRequesterAddress = sellerAddrs[i].String()
			msg.PoolId = poolId
			msg.OfferCoinFee = types.GetOfferCoinFee(msg.OfferCoin, params.SwapFeeRate)
		}
	}

	// begin block, delete and init pool batch
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	// handle msgs, set order msgs to batch
	for _, msg := range XtoY {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, int64(r.Intn(4)))
		require.NoError(t, err)
	}
	for _, msg := range YtoX {
		_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, msg, int64(r.Intn(4)))
		require.NoError(t, err)
	}

	// verify pool batch
	liquidityPoolBatch, found := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, found)
	require.NotNil(t, liquidityPoolBatch)

	// end block, swap execution
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}

func TestGetRandomOrders(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		r := rand.New(rand.NewSource(seed))

		// get random X, Y amount for create pool
		X, Y := app.GetRandPoolAmt(r, types.DefaultMinInitDepositAmount)
		XtoY, YtoX := app.GetRandomSizeOrders("denomX", "denomY", X, Y, r, 50, 50)
		_, _ = XtoY, YtoX // TODO: test against these
		require.Equal(t, X.ToDec().MulInt64(2).TruncateInt(), X.MulRaw(2))
		require.Equal(t, X.ToDec().MulInt64(2), X.MulRaw(2).ToDec())
		require.NotEqual(t, X.ToDec().MulInt64(2).TruncateInt(), X.MulRaw(2).ToDec())
	}
}

func TestBadSwapExecution(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	simapp, ctx := app.CreateTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)
	denomX, denomY := types.AlphabeticalDenomPair("denomX", "denomY")

	// add pool creator account
	X, Y := app.GetRandPoolAmt(r, params.MinInitDepositAmount)
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	creatorAddr := app.AddRandomTestAddr(simapp, ctx, deposit.Add(params.PoolCreationFee...))
	balanceX := simapp.BankKeeper.GetBalance(ctx, creatorAddr, denomX)
	balanceY := simapp.BankKeeper.GetBalance(ctx, creatorAddr, denomY)
	creatorBalance := sdk.NewCoins(balanceX, balanceY)
	require.Equal(t, deposit, creatorBalance)

	// create pool
	createPoolMsg := types.NewMsgCreatePool(creatorAddr, types.DefaultPoolTypeId, creatorBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, createPoolMsg)
	require.NoError(t, err)

	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

	offerCoin := sdk.NewCoin(denomX, sdk.NewInt(10000))
	offerCoinFee := types.GetOfferCoinFee(offerCoin, params.SwapFeeRate)
	testAddr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins(offerCoin.Add(offerCoinFee)))

	currentPrice := X.ToDec().Quo(Y.ToDec())
	swapMsg := types.NewMsgSwapWithinBatch(testAddr, 0, types.DefaultSwapTypeId, offerCoin, denomY, currentPrice, params.SwapFeeRate)
	_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, swapMsg, 0)
	require.ErrorIs(t, err, types.ErrPoolNotExists)

	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
}

func TestBalancesAfterSwap(t *testing.T) {
	for price := int64(9800); price < 10000; price++ {
		simapp, ctx := app.CreateTestInput()
		params := simapp.LiquidityKeeper.GetParams(ctx)
		denomX, denomY := types.AlphabeticalDenomPair("denomx", "denomy")
		X, Y := sdk.NewInt(100_000_000), sdk.NewInt(100_000_000)

		creatorCoins := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
		creatorAddr := app.AddRandomTestAddr(simapp, ctx, creatorCoins.Add(params.PoolCreationFee...))

		orderPrice := sdk.NewDecWithPrec(price, 4)
		aliceCoin := sdk.NewCoin(denomY, sdk.NewInt(10_000_000))
		aliceAddr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins(aliceCoin))

		pool, err := simapp.LiquidityKeeper.CreatePool(ctx, types.NewMsgCreatePool(creatorAddr, types.DefaultPoolTypeId, creatorCoins))
		require.NoError(t, err)

		liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

		offerAmt := aliceCoin.Amount.ToDec().Quo(sdk.MustNewDecFromStr("1.0015")).TruncateInt()
		offerCoin := sdk.NewCoin(denomY, offerAmt)

		_, err = simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, types.NewMsgSwapWithinBatch(
			aliceAddr, pool.Id, types.DefaultSwapTypeId, offerCoin, denomX, orderPrice, params.SwapFeeRate), 0)
		require.NoError(t, err)

		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		deltaX := simapp.BankKeeper.GetBalance(ctx, aliceAddr, denomX).Amount
		deltaY := simapp.BankKeeper.GetBalance(ctx, aliceAddr, denomY).Amount.Sub(aliceCoin.Amount)
		require.Truef(t, !deltaX.IsNegative(), "deltaX should not be negative: %s", deltaX)
		require.Truef(t, deltaY.IsNegative(), "deltaY should be negative: %s", deltaY)

		deltaXWithoutFee := deltaX.ToDec().Quo(sdk.MustNewDecFromStr("0.9985"))
		deltaYWithoutFee := deltaY.ToDec().Quo(sdk.MustNewDecFromStr("1.0015"))
		effectivePrice := deltaXWithoutFee.Quo(deltaYWithoutFee.Neg())
		priceDiffRatio := orderPrice.Sub(effectivePrice).Abs().Quo(orderPrice)
		require.Truef(t, priceDiffRatio.LT(sdk.MustNewDecFromStr("0.01")), "effectivePrice differs too much from orderPrice")
	}
}

func TestRefundEscrow(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		r := rand.New(rand.NewSource(seed))

		X := sdk.NewInt(1_000_000)
		Y := app.GetRandRange(r, 10_000_000_000_000_000, 1_000_000_000_000_000_000)

		simapp, ctx := createTestInput()
		params := simapp.LiquidityKeeper.GetParams(ctx)

		addr := app.AddRandomTestAddr(simapp, ctx, sdk.NewCoins())

		pool, err := createPool(simapp, ctx, X, Y, DenomX, DenomY)
		require.NoError(t, err)

		for i := 0; i < 100; i++ {
			poolBalances := simapp.BankKeeper.GetAllBalances(ctx, pool.GetReserveAccount())
			RX := poolBalances.AmountOf(DenomX)
			RY := poolBalances.AmountOf(DenomY)
			poolPrice := RX.ToDec().Quo(RY.ToDec())

			offerAmt := RY.ToDec().Mul(sdk.NewDecFromIntWithPrec(app.GetRandRange(r, 1, 100_000_000_000_000_000), sdk.Precision))    // RY * (0, 0.1)
			offerAmtWithFee := offerAmt.Quo(sdk.OneDec().Add(params.SwapFeeRate.QuoInt64(2))).TruncateInt()                          // offerAmt / (1 + swapFeeRate/2)
			orderPrice := poolPrice.Mul(sdk.NewDecFromIntWithPrec(app.GetRandRange(r, 1, 1_000_000_000_000_000_000), sdk.Precision)) // poolPrice * (0, 1)

			require.NoError(t, simapp.BankKeeper.SetBalance(ctx, addr, sdk.NewCoin(DenomY, offerAmt.Ceil().TruncateInt())))

			liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)

			_, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, types.NewMsgSwapWithinBatch(
				addr, pool.Id, types.DefaultSwapTypeId, sdk.NewCoin(DenomY, offerAmtWithFee), DenomX, orderPrice, params.SwapFeeRate), 0)
			require.NoError(t, err)

			liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
		}

		require.True(t, simapp.BankKeeper.GetAllBalances(ctx, simapp.AccountKeeper.GetModuleAddress(types.ModuleName)).IsZero(), "there must be no remaining coin escrow")
	}
}

func createPool(simapp *app.LiquidityApp, ctx sdk.Context, X, Y sdk.Int, denomX, denomY string) (types.Pool, error) {
	params := simapp.LiquidityKeeper.GetParams(ctx)

	coins := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	addr := app.AddRandomTestAddr(simapp, ctx, coins.Add(params.PoolCreationFee...))

	return simapp.LiquidityKeeper.CreatePool(ctx, types.NewMsgCreatePool(addr, types.DefaultPoolTypeId, coins))
}
