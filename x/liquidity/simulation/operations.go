package simulation

import (
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	liquidityparams "github.com/tendermint/liquidity/app/params"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

var (
	// create simulated accounts due to gas usage overflow issue.
	// Read this issue: https://github.com/tendermint/liquidity/issues/349
	randomAccounts []simtypes.Account
)

// Simulation operation weights constants.
const (
	OpWeightMsgCreatePool          = "op_weight_msg_create_pool"
	OpWeightMsgDepositWithinBatch  = "op_weight_msg_deposit_to_pool"
	OpWeightMsgWithdrawWithinBatch = "op_weight_msg_withdraw_from_pool"
	OpWeightMsgSwapWithinBatch     = "op_weight_msg_swap"
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreatePool int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = liquidityparams.DefaultWeightMsgCreatePool
		},
	)

	var weightMsgDepositWithinBatch int
	appParams.GetOrGenerate(cdc, OpWeightMsgDepositWithinBatch, &weightMsgDepositWithinBatch, nil,
		func(_ *rand.Rand) {
			weightMsgDepositWithinBatch = liquidityparams.DefaultWeightMsgDepositWithinBatch
		},
	)

	var weightMsgMsgWithdrawWithinBatch int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawWithinBatch, &weightMsgMsgWithdrawWithinBatch, nil,
		func(_ *rand.Rand) {
			weightMsgMsgWithdrawWithinBatch = liquidityparams.DefaultWeightMsgWithdrawWithinBatch
		},
	)

	var weightMsgSwapWithinBatch int
	appParams.GetOrGenerate(cdc, OpWeightMsgSwapWithinBatch, &weightMsgSwapWithinBatch, nil,
		func(_ *rand.Rand) {
			weightMsgSwapWithinBatch = liquidityparams.DefaultWeightMsgSwapWithinBatch
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreatePool,
			SimulateMsgCreatePool(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDepositWithinBatch,
			SimulateMsgDepositWithinBatch(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgMsgWithdrawWithinBatch,
			SimulateMsgWithdrawWithinBatch(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSwapWithinBatch,
			SimulateMsgSwapWithinBatch(ak, bk, k),
		),
	}
}

// SimulateMsgCreatePool generates a MsgCreatePool with random values
// nolint: interfacer
func SimulateMsgCreatePool(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var simAccount simtypes.Account
		randomAccounts = simtypes.RandomAccounts(r, 3)
		simAccount = randomAccounts[r.Intn(3)]

		params := k.GetParams(ctx)
		params.MaxReserveCoinAmount = GenMaxReserveCoinAmount(r)
		k.SetParams(ctx, params)

		// get randomized two denoms to create liquidity pool
		var mintingDenoms []string
		denomA, denomB := randomDenoms(r)
		reserveCoinDenoms := []string{denomA, denomB}
		mintingDenoms = append(mintingDenoms, reserveCoinDenoms...)

		// simAccount should have some fees to pay for transaction and pool creation fee
		var feeDenoms []string
		for _, fee := range params.PoolCreationFee {
			feeDenoms = append(feeDenoms, fee.GetDenom())
		}
		mintingDenoms = append(mintingDenoms, feeDenoms...)

		// mint coins of randomized and fee denoms
		err := mintCoins(r, simAccount, mintingDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to mint and send coins"), nil, err
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		poolName := types.PoolName(reserveCoinDenoms, types.DefaultPoolTypeID)
		reserveAcc := types.GetPoolReserveAcc(poolName)

		// ensure the liquidity pool doesn't exist
		_, found := k.GetPoolByReserveAccIndex(ctx, reserveAcc)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "liquidity pool already exists"), nil, nil
		}

		balanceA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if balanceA.IsNegative() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "balanceA is negative"), nil, nil
		}

		balanceB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if balanceB.IsNegative() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "balanceB is negative"), nil, nil
		}

		poolCreator := account.GetAddress()
		depositCoinA := randomDepositCoin(r, params.MinInitDepositAmount, denomA)
		depositCoinB := randomDepositCoin(r, params.MinInitDepositAmount, denomB)
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

		// it will fail if the total reserve coin amount after the deposit is larger than the parameter
		err = types.ValidateReserveCoinLimit(params.MaxReserveCoinAmount, depositCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "can not exceed reserve coin limit amount"), nil, nil
		}

		msg := types.NewMsgCreatePool(poolCreator, types.DefaultPoolTypeID, depositCoins)

		fees, err := randomFees(r, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to generate fees"), nil, err
		}

		txGen := liquidityparams.MakeTestEncodingConfig().TxConfig
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

// SimulateMsgDepositWithinBatch  generates a MsgDepositWithinBatch  with random values
// nolint: interfacer
func SimulateMsgDepositWithinBatch(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(k.GetAllPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "number of liquidity pools equals zero"), nil, nil
		}

		// use pre-created simulated accounts if randomAccounts are empty for the unit test
		if len(randomAccounts) == 0 {
			randomAccounts = accs
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "unable to pick liquidity pool"), nil, nil
		}

		reserveCoinDenomA := pool.ReserveCoinDenoms[0]
		reserveCoinDenomB := pool.ReserveCoinDenoms[1]

		// select random simulated account and mint reserve coins
		// note that select the simulated account that has some balances of reserve coin denoms result in
		// many failed transactions due to random accounts change after a creating pool.
		simAccount := randomAccounts[r.Intn(len(randomAccounts))]
		err := mintCoins(r, simAccount, []string{reserveCoinDenomA, reserveCoinDenomB}, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "unable to mint and send coins"), nil, err
		}

		params := k.GetParams(ctx)
		params.MaxReserveCoinAmount = GenMaxReserveCoinAmount(r)
		k.SetParams(ctx, params)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		depositor := account.GetAddress()
		depositCoinA := randomDepositCoin(r, params.MinInitDepositAmount, reserveCoinDenomA)
		depositCoinB := randomDepositCoin(r, params.MinInitDepositAmount, reserveCoinDenomB)
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

		reserveCoins := k.GetReserveCoins(ctx, pool)

		// it will fail if the total reserve coin amount after the deposit is larger than the parameter
		err = types.ValidateReserveCoinLimit(params.MaxReserveCoinAmount, reserveCoins.Add(depositCoinA, depositCoinB))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "can not exceed reserve coin limit amount"), nil, nil
		}

		fees, err := randomFees(r, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgDepositWithinBatch(depositor, pool.Id, depositCoins)

		txGen := liquidityparams.MakeTestEncodingConfig().TxConfig
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

// SimulateMsgWithdrawWithinBatch generates a MsgWithdrawWithinBatch with random values.
// nolint: interfacer
func SimulateMsgWithdrawWithinBatch(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(k.GetAllPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "number of liquidity pools equals zero"), nil, nil
		}

		// use pre-created simulated accounts if randomAccounts are empty for the unit test
		if len(randomAccounts) == 0 {
			randomAccounts = accs
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "unable to pick liquidity pool"), nil, nil
		}

		poolCoinDenom := pool.GetPoolCoinDenom()

		// select random simulated account and mint reserve coins
		// note that select the simulated account that has some balance of pool coin denom result in
		// many failed transactions due to random accounts change after a creating pool.
		simAccount := randomAccounts[r.Intn(len(randomAccounts))]
		err := mintCoins(r, simAccount, []string{poolCoinDenom}, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "unable to mint and send coins"), nil, err
		}

		// if simAccount.PrivKey == nil, then no account has pool coin denom balanace
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "account private key is nil"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		balance := bk.GetBalance(ctx, simAccount.Address, poolCoinDenom)
		withdrawer := account.GetAddress()
		withdrawCoin := randomWithdrawCoin(r, poolCoinDenom, balance.Amount)

		msg := types.NewMsgWithdrawWithinBatch(withdrawer, pool.Id, withdrawCoin)

		fees, err := randomFees(r, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "unable to generate fees"), nil, err
		}

		txGen := liquidityparams.MakeTestEncodingConfig().TxConfig
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

// SimulateMsgSwapWithinBatch generates a MsgSwapWithinBatch with random values
// nolint: interfacer
func SimulateMsgSwapWithinBatch(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(k.GetAllPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwapWithinBatch, "number of liquidity pools equals zero"), nil, nil
		}

		// use pre-created simulated accounts if randomAccounts are empty for the unit test
		if len(randomAccounts) == 0 {
			randomAccounts = accs
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwapWithinBatch, "unable to pick liquidity pool"), nil, nil
		}

		reserveCoinDenomA := pool.ReserveCoinDenoms[0]
		reserveCoinDenomB := pool.ReserveCoinDenoms[1]

		// select random simulated account and mint reserve coins
		// note that select the simulated account that has some balances of reserve coin denoms result in
		// many failed transactions due to random accounts change after a creating pool.
		simAccount := randomAccounts[r.Intn(len(randomAccounts))]
		err := mintCoins(r, simAccount, []string{reserveCoinDenomA, reserveCoinDenomB}, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "unable to mint and send coins"), nil, err
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		swapRequester := account.GetAddress()
		offerCoin := randomOfferCoin(r, k, ctx, pool, pool.ReserveCoinDenoms[0])
		demandCoinDenom := pool.ReserveCoinDenoms[1]
		orderPrice := randomOrderPrice(r)
		swapFeeRate := GenSwapFeeRate(r)

		// set randomly generated swap fee rate in params to prevent from miscalculation
		params := k.GetParams(ctx)
		params.SwapFeeRate = swapFeeRate
		k.SetParams(ctx, params)

		msg := types.NewMsgSwapWithinBatch(swapRequester, pool.Id, types.DefaultSwapTypeID, offerCoin, demandCoinDenom, orderPrice, swapFeeRate)

		fees, err := randomFees(r, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to generate fees"), nil, err
		}

		txGen := liquidityparams.MakeTestEncodingConfig().TxConfig
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

// mintCoins mints coins relative to the number of denoms and send them to the simulated account
func mintCoins(r *rand.Rand, acc simtypes.Account, denoms []string, bk types.BankKeeper, ctx sdk.Context) error {
	var mintCoins, sendCoins sdk.Coins
	for _, denom := range denoms {
		mintAmt := sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e15, 1e16)))
		sendAmt := sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e13, 1e14)))
		mintCoins = mintCoins.Add(sdk.NewCoin(denom, mintAmt))
		sendCoins = sendCoins.Add(sdk.NewCoin(denom, sendAmt))
	}

	feeCoin := int64(simtypes.RandIntBetween(r, 1e13, 1e14))
	mintCoins = mintCoins.Add(sdk.NewInt64Coin(sdk.DefaultBondDenom, feeCoin))
	sendCoins = sendCoins.Add(sdk.NewInt64Coin(sdk.DefaultBondDenom, feeCoin))

	err := bk.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return err
	}

	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc.Address, sendCoins)
	if err != nil {
		return err
	}

	return nil
}

// randomLiquidity returns random liquidity pool with given access to the keeper and ctx.
func randomLiquidity(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) (pool types.Pool, ok bool) {
	pools := k.GetAllPools(ctx)
	if len(pools) == 0 {
		return types.Pool{}, false
	}

	i := r.Intn(len(pools))

	return pools[i], true
}

// randomDepositCoin returns deposit amount between minInitDepositAmount+1 and 1e9.
func randomDepositCoin(r *rand.Rand, minInitDepositAmount sdk.Int, denom string) sdk.Coin {
	amount := int64(simtypes.RandIntBetween(r, int(minInitDepositAmount.Int64()+1), 1e8))
	return sdk.NewInt64Coin(denom, amount)
}

// randomWithdrawCoin returns random withdraw amount between 1 and the account's current balance divide by 10.
func randomWithdrawCoin(r *rand.Rand, denom string, balance sdk.Int) sdk.Coin {
	// prevent panic from RandIntBetween
	if balance.Quo(sdk.NewInt(10)).Int64() <= 1 {
		return sdk.NewInt64Coin(denom, 1)
	}

	amount := int64(simtypes.RandIntBetween(r, 1, int(balance.Quo(sdk.NewInt(10)).Int64())))
	return sdk.NewInt64Coin(denom, amount)
}

// randomOfferCoin returns random offer amount of coin.
func randomOfferCoin(r *rand.Rand, k keeper.Keeper, ctx sdk.Context, pool types.Pool, denom string) sdk.Coin {
	params := k.GetParams(ctx)
	reserveCoinAmt := k.GetReserveCoins(ctx, pool).AmountOf(denom)
	maximumOrderableAmt := reserveCoinAmt.ToDec().Mul(params.MaxOrderAmountRatio).TruncateInt()
	amt := int64(simtypes.RandIntBetween(r, 1, int(maximumOrderableAmt.Int64())))
	return sdk.NewInt64Coin(denom, amt)
}

// randomOrderPrice returns random order price amount between 0.01 to 1.
func randomOrderPrice(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 1e2)), 2)
}

// randomFees returns a random amount of bond denom fee and
// if the user doesn't have enough funds for paying fees, it returns empty coins.
func randomFees(r *rand.Rand, spendableCoins sdk.Coins) (sdk.Coins, error) {
	if spendableCoins.Empty() {
		return nil, nil
	}

	if spendableCoins.AmountOf(sdk.DefaultBondDenom).Equal(sdk.ZeroInt()) {
		return nil, nil
	}

	amt, err := simtypes.RandPositiveInt(r, spendableCoins.AmountOf(sdk.DefaultBondDenom))
	if err != nil {
		return nil, err
	}

	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amt))

	return fees, nil
}

// randomDenoms returns two random denoms with random string length anywhere from 4 to 6.
func randomDenoms(r *rand.Rand) (string, string) {
	denomA := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
	denomB := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
	denomA, denomB = types.AlphabeticalDenomPair(strings.ToLower(denomA), strings.ToLower(denomB))
	return denomA, denomB
}
