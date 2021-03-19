package simulation

import (
	"fmt"
	"math/rand"

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

// Simulation operation weights constants
const (
	OpWeightMsgCreatePool          = "op_weight_msg_create_pool"
	OpWeightMsgDepositWithinBatch  = "op_weight_msg_deposit_to_pool"
	OpWeightMsgWithdrawWithinBatch = "op_weight_msg_withdraw_from_pool"
	OpWeightMsgSwapWithinBatch     = "op_weight_msg_swap"
)

// WeightedOperations returns all the operations from the module with their respective weights
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
		/*
			1. Get randomly created fees denoms from genesis
			2. Mint those coins and send them to the simulated account
			2. Check if the same liquidity pool already exists and balances of both denoms
			3. Create new liquidity pool with random deposit amount of coins
		*/
		params := k.GetParams(ctx)
		params.MaxReserveCoinAmount = GenMaxReserveCoinAmount(r)
		k.SetParams(ctx, params)

		// simAccount should have some fees to pay when creating liquidity pool
		var feeDenoms []string
		for _, fee := range params.PoolCreationFee {
			feeDenoms = append(feeDenoms, fee.GetDenom())
		}

		if len(feeDenoms) < 2 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "at least 2 coin denoms required"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		// mint randomly generated fee coins to the simulated account for the use of liquidity creation fee
		err := mintCoins(r, simAccount.Address, feeDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to mint and send coins"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := randomFees(r, ctx, k, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to generate fees"), nil, err
		}

		denomA, denomB := randomDenoms(r)
		reserveCoinDenoms := []string{denomA, denomB}

		// mint new random 2 coins to create new liquidity pool
		err = mintCoins(r, simAccount.Address, reserveCoinDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "unable to mint and send coins"), nil, nil
		}

		poolName := types.PoolName(reserveCoinDenoms, types.DefaultPoolTypeId)
		reserveAcc := types.GetPoolReserveAcc(poolName)

		// ensure the liquidity pool doesn't exist
		_, found := k.GetPoolByReserveAccIndex(ctx, reserveAcc)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "liquidity pool already exists"), nil, nil
		}

		balanceA := bk.GetBalance(ctx, simAccount.Address, denomA).Amount
		if !balanceA.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "balanceA is negative"), nil, nil
		}

		balanceB := bk.GetBalance(ctx, simAccount.Address, denomB).Amount
		if !balanceB.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "balanceB is negative"), nil, nil
		}

		poolCreator := account.GetAddress()
		depositCoinA := randomDepositCoin(r, params.MinInitDepositAmount, denomA)
		depositCoinB := randomDepositCoin(r, params.MinInitDepositAmount, denomB)
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

		// it will fail if the total reserve coin amount after the deposit is larger than the parameter
		err = types.ValidateReserveCoinLimit(params.MaxReserveCoinAmount, depositCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "can not exceed reserve coin limit amount"), nil, nil
		}

		msg := types.NewMsgCreatePool(poolCreator, types.DefaultPoolTypeId, depositCoins)

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
		/*
			1. Check if there's any liquidity pool created; if there isn't, then return NoOpMsg
			2. Get random liquidity pool and mint those coins to the simulated account
			3. Deposit random amount of coins the to liquidity pool
		*/

		if len(k.GetAllPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "number of liquidity pools equals zero"), nil, nil
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "unable to pick liquidity pool"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := randomFees(r, ctx, k, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "unable to generate fees"), nil, err
		}

		// mint pool denoms to the simulated account
		err = mintCoins(r, simAccount.Address, pool.ReserveCoinDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "unable to mint and send coins"), nil, nil
		}

		params := k.GetParams(ctx)
		params.MaxReserveCoinAmount = GenMaxReserveCoinAmount(r)
		k.SetParams(ctx, params)

		depositor := account.GetAddress()
		depositCoinA := randomDepositCoin(r, params.MinInitDepositAmount, pool.ReserveCoinDenoms[0])
		depositCoinB := randomDepositCoin(r, params.MinInitDepositAmount, pool.ReserveCoinDenoms[1])
		depositCoins := sdk.NewCoins(depositCoinA, depositCoinB)

		reserveCoins := k.GetReserveCoins(ctx, pool)

		// it will fail if the total reserve coin amount after the deposit is larger than the parameter
		err = types.ValidateReserveCoinLimit(params.MaxReserveCoinAmount, reserveCoins.Add(depositCoinA, depositCoinB))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositWithinBatch, "can not exceed reserve coin limit amount"), nil, nil
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

// SimulateMsgWithdrawWithinBatch generates a MsgWithdrawWithinBatch with random values
// nolint: interfacer
func SimulateMsgWithdrawWithinBatch(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		/*
			1. Check if there's any liquidity pool created; if there isn't, then return NoOpMsg
			2. Get any available simulated account and check if it has pool coin to withdraw from the pool
			2. Get random liquidity pool and mint pool coin (LP token) to the simulated account
			3. Withdraw random amounts from the liquidity pool
		*/

		if len(k.GetAllPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "number of liquidity pools equals zero"), nil, nil
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "unable to pick liquidity pool"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		// if simAccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "account private key is nil"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := randomFees(r, ctx, k, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "unable to generate fees"), nil, err
		}

		poolCoinDenom := pool.GetPoolCoinDenom()

		// make sure simAccount have pool coin balance
		balance := bk.GetBalance(ctx, simAccount.Address, poolCoinDenom)
		if !balance.IsPositive() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawWithinBatch, "account balance is negative"), nil, nil
		}

		withdrawer := account.GetAddress()
		withdrawCoin := randomWithdrawCoin(r, poolCoinDenom, balance.Amount)

		msg := types.NewMsgWithdrawWithinBatch(withdrawer, pool.Id, withdrawCoin)

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
		/*
			1. Check if there's any liquidity pool created; if there isn't, then return NoOpMsg
			2. Get random liquidity pool and mint those coins to the simulated account
			3. Swap random amount of denomA with denomB
		*/

		if len(k.GetAllPools(ctx)) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwapWithinBatch, "number of liquidity pools equals zero"), nil, nil
		}

		pool, ok := randomLiquidity(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwapWithinBatch, "unable to pick liquidity pool"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		// mint pool denoms to the simulated account
		err := mintCoins(r, simAccount.Address, pool.ReserveCoinDenoms, bk, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwapWithinBatch, "unable to mint and send coins"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := randomFees(r, ctx, k, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSwapWithinBatch, "unable to generate fees"), nil, err
		}

		swapRequester := account.GetAddress()
		offerCoin := randomOfferCoin(r, k, ctx, pool, pool.ReserveCoinDenoms[0])
		demandCoinDenom := pool.ReserveCoinDenoms[1]
		orderPrice := randomOrderPrice(r)
		swapFeeRate := GenSwapFeeRate(r)

		// set randomly generated swap fee rate in params to prevent from miscalculation
		params := k.GetParams(ctx)
		params.SwapFeeRate = swapFeeRate
		k.SetParams(ctx, params)

		msg := types.NewMsgSwapWithinBatch(swapRequester, pool.Id, types.DefaultSwapTypeId, offerCoin, demandCoinDenom, orderPrice, swapFeeRate)

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
func mintCoins(r *rand.Rand, address sdk.AccAddress, denoms []string, bk types.BankKeeper, ctx sdk.Context) error {
	var mintCoins, sendCoins sdk.Coins
	for _, denom := range denoms {
		mintCoins = append(mintCoins, sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e13, 1e14)))))
		sendCoins = append(sendCoins, sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1e11, 1e12)))))
	}

	err := bk.MintCoins(ctx, types.ModuleName, mintCoins)
	if err != nil {
		return err
	}

	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sendCoins)
	if err != nil {
		return err
	}

	return nil
}

// randomDenoms returns two random denoms with random string length anywhere from 4 to 6
func randomDenoms(r *rand.Rand) (string, string) {
	denomA := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
	denomB := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 4, 6))
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)
	return denomA, denomB
}

// randomDepositCoin returns deposit amount between minInitDepositAmount+1 and 1e9
func randomDepositCoin(r *rand.Rand, minInitDepositAmount sdk.Int, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, int(minInitDepositAmount.Int64()+1), 1e9))))
}

// randomLiquidity returns random liquidity pool with given access to the keeper and ctx
func randomLiquidity(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) (pool types.Pool, ok bool) {
	pools := k.GetAllPools(ctx)
	if len(pools) == 0 {
		return types.Pool{}, false
	}

	i := r.Intn(len(pools))

	return pools[i], true
}

// randomWithdrawCoin returns random withdraw amount between 1 and the account's current balance divide by 10
func randomWithdrawCoin(r *rand.Rand, denom string, balance sdk.Int) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1, int(balance.Quo(sdk.NewInt(10)).Int64())))))
}

// randomOrderPrice returns random order price amount between 0.01 to 1
func randomOrderPrice(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 1e2)), 2)
}

// randomOfferCoin returns random offer amount of coin
func randomOfferCoin(r *rand.Rand, k keeper.Keeper, ctx sdk.Context, pool types.Pool, denom string) sdk.Coin {
	params := k.GetParams(ctx)

	// prevent from "can not exceed max order ratio of reserve coins that can be ordered at a order" error
	reserveCoinAmt := k.GetReserveCoins(ctx, pool).AmountOf(denom)
	maximumOrderableAmt := reserveCoinAmt.ToDec().Mul(params.MaxOrderAmountRatio).TruncateInt()

	return sdk.NewCoin(denom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1, int(maximumOrderableAmt.Int64())))))
}

// randomFees returns a random fee by selecting a random coin denomination except pool coin and
// amount from the account's available balance. If the user doesn't have enough
// funds for paying fees, it returns empty coins.
func randomFees(r *rand.Rand, ctx sdk.Context, keeper keeper.Keeper, spendableCoins sdk.Coins) (sdk.Coins, error) {
	if spendableCoins.Empty() {
		return nil, nil
	}

	perm := r.Perm(len(spendableCoins))
	var randCoin sdk.Coin
	for _, index := range perm {
		if keeper.IsPoolCoinDenom(ctx, spendableCoins[index].Denom) {
			continue
		}
		randCoin = spendableCoins[index]
		if !randCoin.Amount.IsZero() {
			break
		}
	}

	if randCoin.Amount.IsZero() {
		return nil, fmt.Errorf("no coins found for random fees")
	}

	amt, err := simtypes.RandPositiveInt(r, randCoin.Amount)
	if err != nil {
		return nil, err
	}

	// Create a random fee and verify the fees are within the account's spendable
	// balance.
	fees := sdk.NewCoins(sdk.NewCoin(randCoin.Denom, amt))

	return fees, nil
}
