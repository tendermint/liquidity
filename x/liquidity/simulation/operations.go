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
		simAccount, _ := simtypes.RandomAcc(r, accs)
		denomA, denomB := randomDenoms(r)
		reserveCoinDenoms := []string{denomA, denomB}
		err := mintCoins(r, simAccount.Address, reserveCoinDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "unable to mint and send coins"), nil, nil
		}

		poolKey := types.GetPoolKey(reserveCoinDenoms, types.DefaultPoolTypeIndex)
		reserveAcc := types.GetPoolReserveAcc(poolKey)

		// ensure the liquidity pool doesn't exist
		_, found := k.GetLiquidityPoolByReserveAccIndex(ctx, reserveAcc)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "liquidity pool already exists"), nil, nil
		}

		balanceA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if !balanceA.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "balanceA is negative"), nil, nil
		}

		balanceB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if !balanceB.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "balanceB is negative"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateLiquidityPool, "unable to generate fees"), nil, err
		}

		poolCreator := account.GetAddress()
		depositCoinA := randomDepositCoin(r, denomA)
		depositCoinB := randomDepositCoin(r, denomB)
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

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
		if len(k.GetAllLiquidityPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "number of liquidity pools equals zero"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "unable to pick liquidity pool"), nil, nil
		}

		// mint pool denoms to the simulated account
		err := mintCoins(r, simAccount.Address, pool.ReserveCoinDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "unable to mint and send coins"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "unable to generate fees"), nil, err
		}

		depositCoinA := randomNewDepositCoin(r, pool.ReserveCoinDenoms[0])
		depositCoinB := randomNewDepositCoin(r, pool.ReserveCoinDenoms[1])
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

		msg := types.NewMsgDepositToLiquidityPool(account.GetAddress(), pool.PoolId, depositCoins)

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
		if len(k.GetAllLiquidityPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "number of liquidity pools equals zero"), nil, nil
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "unable to pick liquidity pool"), nil, nil
		}

		poolCoinDenom := pool.GetPoolCoinDenom()

		// need to retrieve random simulation account to retrieve PrivKey
		simAccount := accs[simtypes.RandIntBetween(r, 0, len(accs))]

		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "account private key is nil"), nil, nil
		}

		// mint pool coin to the simulated account
		err := mintPoolCoin(r, simAccount.Address, poolCoinDenom, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositToLiquidityPool, "unable to mint and send coins"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawFromLiquidityPool, "unable to generate fees"), nil, err
		}

		withdrawer := account.GetAddress()
		withdrawCoin := randomWithdrawCoin(r, poolCoinDenom)

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
		if len(k.GetAllLiquidityPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "number of liquidity pools equals zero"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "unable to pick liquidity pool"), nil, nil
		}

		// mint pool denoms to the simulated account
		err := mintCoins(r, simAccount.Address, pool.ReserveCoinDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "unable to mint and send coins"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "unable to generate fees"), nil, err
		}

		swapRequester := account.GetAddress()
		demandCoinDenom := pool.ReserveCoinDenoms[1]
		orderPrice := randomOrderPrice(r)

		offerAmount, err := simtypes.RandPositiveInt(r, sdk.NewInt(r.Int63()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwap, "unable to generate positive amount"), nil, err
		}

		offerCoin := sdk.NewCoin(pool.ReserveCoinDenoms[0], offerAmount)
		swapFeeRate := GenSwapFeeRate(r)

		msg := types.NewMsgSwap(swapRequester, pool.PoolId, types.DefaultSwapType,
			offerCoin, demandCoinDenom, orderPrice, swapFeeRate)

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

// mintCoins mints coins relative to the denoms and send them to the simulated account
func mintCoins(r *rand.Rand, address sdk.AccAddress, denoms []string, bk types.BankKeeper, ctx sdk.Context) error {
	// mint random amounts of denomA and denomB coins
	mintCoinA := randomMintAmount(r, denoms[0])
	mintCoinB := randomMintAmount(r, denoms[1])
	mintCoins := sdk.NewCoins(mintCoinA, mintCoinB)
	err := bk.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return err
	}

	// transfer random amounts to the simulated random account
	coinA := sdk.NewCoin(denoms[0], sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e12, 1e14))))
	coinB := sdk.NewCoin(denoms[1], sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e12, 1e14))))
	coins := sdk.NewCoins(coinA, coinB)
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, coins)
	if err != nil {
		return err
	}
	return nil
}

// mintPoolCoin mints pool coin and send random amount to the simulated account
func mintPoolCoin(r *rand.Rand, address sdk.AccAddress, poolCoinDenom string, bk types.BankKeeper, ctx sdk.Context) error {
	mintCoins := sdk.NewCoins(sdk.NewCoin(poolCoinDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e7, 1e8)))))
	err := bk.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(poolCoinDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e5, 1e6)))))
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, coins)
	if err != nil {
		return err
	}
	return nil
}

// randomDenoms randomizes denoms with a length from 4 to 6 characters
func randomDenoms(r *rand.Rand) (string, string) {
	denomA := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
	denomB := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
	return denomA, denomB
}

// randomMintAmount randomizes minting coins in a range of 1e12 to 1e15
func randomMintAmount(r *rand.Rand, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e14, 1e15))))
}

// randomLiquidity returns a random liquidity pool given access to the keeper and ctx
func randomLiquidity(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) (pool types.LiquidityPool, ok bool) {
	pools := k.GetAllLiquidityPools(ctx)
	if len(pools) == 0 {
		return types.LiquidityPool{}, false
	}

	i := r.Intn(len(pools))

	return pools[i], true
}

// randomDepositCoin randomizes deposit coin greater than DefaultMinInitDepositToPool
func randomDepositCoin(r *rand.Rand, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, int(types.DefaultMinInitDepositToPool.Int64()), 1e10))))
}

// randomNewDepositCoin randomizes new deposit coin
func randomNewDepositCoin(r *rand.Rand, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1, 1e8))))
}

// randomWithdrawCoin randomizes withdraw coin
func randomWithdrawCoin(r *rand.Rand, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1, 1e5))))
}

// randomOrderPrice randomized order price ranging from 0.01 to 1
func randomOrderPrice(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 1e2)), 2)
}
