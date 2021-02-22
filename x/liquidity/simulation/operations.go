package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	simappparams "github.com/tendermint/liquidity/app/params"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateLiquidityPool       = "op_weight_msg_create_liquidity_pool"
	OpWeightMsgDepositToLiquidityPool    = "op_weight_msg_deposit_to_liquidity_pool"
	OpWeightMsgWithdrawFromLiquidityPool = "op_weight_msg_withdraw_from_liquidity_pool"
	OpWeightMsgSwap                      = "op_weight_msg_swap"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreateLiquidityPool int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateLiquidityPool, &weightMsgCreateLiquidityPool, nil,
		func(_ *rand.Rand) {
			weightMsgCreateLiquidityPool = simappparams.DefaultWeightMsgCreateLiquidityPool
		},
	)

	var weightMsgDepositToLiquidityPool int
	appParams.GetOrGenerate(cdc, OpWeightMsgDepositToLiquidityPool, &weightMsgDepositToLiquidityPool, nil,
		func(_ *rand.Rand) {
			weightMsgDepositToLiquidityPool = simappparams.DefaultWeightMsgDepositToLiquidityPool
		},
	)

	var weightMsgMsgWithdrawFromLiquidityPool int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawFromLiquidityPool, &weightMsgMsgWithdrawFromLiquidityPool, nil,
		func(_ *rand.Rand) {
			weightMsgMsgWithdrawFromLiquidityPool = simappparams.DefaultWeightMsgWithdrawFromLiquidityPool
		},
	)

	var weightMsgSwap int
	appParams.GetOrGenerate(cdc, OpWeightMsgSwap, &weightMsgSwap, nil,
		func(_ *rand.Rand) {
			weightMsgSwap = simappparams.DefaultWeightMsgSwap
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateLiquidityPool,
			SimulateMsgCreateLiquidityPool(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDepositToLiquidityPool,
			SimulateMsgDepositToLiquidityPool(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgMsgWithdrawFromLiquidityPool,
			SimulateMsgWithdrawFromLiquidityPool(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSwap,
			SimulateMsgSwap(ak, bk, k),
		),
	}
}

// SimulateMsgCreateLiquidityPool generates a MsgCreateLiquidityPool with random values
// nolint: interfacer
func SimulateMsgCreateLiquidityPool(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		/*
			1. Mint new coins and send those  minted coins to simulated random account
			2. Create new liquidity pool
		*/

		simAccount, denomA, denomB := randomAccountWithNewCoins(ctx, r, accs, bk)

		amountA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if !amountA.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "amountA is negative"), nil, nil
		}

		amountB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if !amountB.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "amountB is negative"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "unable to generate fees"), nil, err
		}

		poolCreator := account.GetAddress()
		depositCoins := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(1000000)), sdk.NewCoin(denomB, sdk.NewInt(1000000)))

		msg := types.NewMsgCreateLiquidityPool(poolCreator, types.DefaultPoolTypeIndex, depositCoins)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDepositToLiquidityPool  generates a MsgDepositToLiquidityPool  with random values
// nolint: interfacer
func SimulateMsgDepositToLiquidityPool(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		/*
			1. Mint new coins and send those minted coins to simulated random account
			2. Create new liquidity pool
			3. Deposit coins to liquidity pool
		*/

		simAccount, denomA, denomB := randomAccountWithNewCoins(ctx, r, accs, bk)

		amountA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if !amountA.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "amountA is negative"), nil, nil
		}

		amountB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if !amountB.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "amountB is negative"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "unable to generate fees"), nil, err
		}

		depositCoins := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(1000000)), sdk.NewCoin(denomB, sdk.NewInt(1000000)))

		createPoolMsg := types.NewMsgCreateLiquidityPool(account.GetAddress(), types.DefaultPoolTypeIndex, depositCoins)
		err, _ = k.CreateLiquidityPool(ctx, createPoolMsg)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "unable to create liquidity pool"), nil, err
		}

		poolId := uint64(1)

		pool, found := k.GetLiquidityPool(ctx, poolId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "unable to get liquidity pool"), nil, err
		}

		depositor := account.GetAddress()
		newDepositCoins := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(5000)), sdk.NewCoin(denomB, sdk.NewInt(5000)))

		msg := types.NewMsgDepositToLiquidityPool(depositor, pool.PoolId, newDepositCoins)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawFromLiquidityPool generates a MsgWithdrawFromLiquidityPool with random values
// nolint: interfacer
func SimulateMsgWithdrawFromLiquidityPool(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		/*
			1. Mint new coins and send those minted coins to simulated random account
			2. Create new liquidity pool
			3. Withdraw some coins from liquidity pool
		*/

		simAccount, denomA, denomB := randomAccountWithNewCoins(ctx, r, accs, bk)

		amountA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if !amountA.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "amountA is negative"), nil, nil
		}

		amountB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if !amountB.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "amountB is negative"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "unable to generate fees"), nil, err
		}

		depositCoins := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(1000000)), sdk.NewCoin(denomB, sdk.NewInt(1000000)))

		createPoolMsg := types.NewMsgCreateLiquidityPool(account.GetAddress(), types.DefaultPoolTypeIndex, depositCoins)
		err, _ = k.CreateLiquidityPool(ctx, createPoolMsg)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "unable to create liquidity pool"), nil, err
		}

		poolId := uint64(1)

		pool, found := k.GetLiquidityPool(ctx, poolId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "unable to get liquidity pool"), nil, err
		}

		withdrawer := account.GetAddress()
		withdrawCoin := sdk.NewCoin(pool.PoolCoinDenom, sdk.NewInt(5000))

		msg := types.NewMsgWithdrawFromLiquidityPool(withdrawer, pool.PoolId, withdrawCoin)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgSwap generates a MsgSwap with random values
// nolint: interfacer
func SimulateMsgSwap(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		/*
			1. Mint new coins and send those minted coins to simulated random account
			2. Create new liquidity pool
			3. Swap coinA to coinB
		*/

		simAccount, denomA, denomB := randomAccountWithNewCoins(ctx, r, accs, bk)

		amountA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if !amountA.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "amountA is negative"), nil, nil
		}

		amountB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if !amountB.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "amountB is negative"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "unable to generate fees"), nil, err
		}

		depositCoins := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(1000000)), sdk.NewCoin(denomB, sdk.NewInt(1000000)))

		createPoolMsg := types.NewMsgCreateLiquidityPool(account.GetAddress(), types.DefaultPoolTypeIndex, depositCoins)
		err, _ = k.CreateLiquidityPool(ctx, createPoolMsg)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "unable to create liquidity pool"), nil, err
		}

		params := k.GetParams(ctx)

		swapRequester := account.GetAddress()
		poolId := uint64(1)
		swapType := types.DefaultSwapType
		offerCoin := sdk.NewCoin(denomA, sdk.NewInt(5000))
		demandCoinDenom := denomB
		orderPrice := sdk.NewDec(150)
		swapFeeRate := params.SwapFeeRate

		msg := types.NewMsgSwap(swapRequester, poolId, swapType, offerCoin, demandCoinDenom, orderPrice, swapFeeRate)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// randomAccountWithNewCoins creates simulated random account and have new minted coins
func randomAccountWithNewCoins(ctx sdk.Context, r *rand.Rand, accs []simtypes.Account, bk types.BankKeeper) (simtypes.Account, string, string) {
	denomA := "denomA"
	denomB := "denomB"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	coinA := sdk.NewCoin(denomA, sdk.NewInt(1e10))                                                          // amount of coinA for the simulated random account
	coinB := sdk.NewCoin(denomB, sdk.NewInt(1e10))                                                          // amount of coinB for the simulated random account
	mintCoins := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(1e12)), sdk.NewCoin(denomB, sdk.NewInt(1e12))) // mint amounts for denomA and denomB coins

	simAccount, _ := simtypes.RandomAcc(r, accs)

	bk.MintCoins(ctx, types.ModuleName, mintCoins)
	bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, simAccount.Address, sdk.NewCoins(coinA, coinB))

	return simAccount, denomA, denomB
}
