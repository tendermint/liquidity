package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// Simulation parameter constants
const (
	LiquidityPoolTypes       = "liquidity_pool_types"
	MinInitDepositToPool     = "min_init_deposit_to_pool"
	InitPoolCoinMintAmount   = "init_pool_coin_mint_amount"
	ReserveCoinLimitAmount   = "reserve_coin_limit_amount"
	LiquidityPoolCreationFee = "liquidity_pool_creation_fee"
	SwapFeeRate              = "swap_fee_rate"
	WithdrawFeeRate          = "withdraw_fee_rate"
	MaxOrderAmountRatio      = "max_order_amount_ratio"
	UnitBatchSize            = "unit_batch_size"
)

// GenLiquidityPoolTypes return default PoolType temporarily, It will be randomized in the liquidity v2
func GenLiquidityPoolTypes(r *rand.Rand) (liquidityPoolTypes []types.PoolType) {
	liquidityPoolTypes = append(liquidityPoolTypes, types.DefaultPoolType)
	return liquidityPoolTypes
}

// GenMinInitDepositToPool randomized MinInitDepositToPool
func GenMinInitDepositToPool(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(simulation.RandIntBetween(r, int(types.DefaultMinInitDepositToPool.Int64()), 1e7)))
}

// GenInitPoolCoinMintAmount randomized InitPoolCoinMintAmount
func GenInitPoolCoinMintAmount(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(simulation.RandIntBetween(r, int(types.DefaultInitPoolCoinMintAmount.Int64()), 1e8)))
}

// GenReserveCoinLimitAmount randomized ReserveCoinLimitAmount
func GenReserveCoinLimitAmount(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(simulation.RandIntBetween(r, int(types.DefaultReserveCoinLimitAmount.Int64()), 1e13)))
}

// GenLiquidityPoolCreationFee randomized LiquidityPoolCreationFee
// list of 1 to 4 coins with an amount greater than 1
func GenLiquidityPoolCreationFee(r *rand.Rand) sdk.Coins {
	var coins sdk.Coins
	var denoms []string

	count := simulation.RandIntBetween(r, 1, 4)
	for i := 0; i < count; i++ {
		randomDenom := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 4, 6))
		denoms = append(denoms, randomDenom)
	}

	sortedDenoms := types.SortDenoms(denoms)

	for i := 0; i < count; i++ {
		randomCoin := sdk.NewCoin(sortedDenoms[i], sdk.NewInt(int64(simulation.RandIntBetween(r, 1e6, 1e7))))
		coins = append(coins, randomCoin)
	}

	return coins
}

// GenSwapFeeRate randomized SwapFeeRate ranging from 0.00001 to 1
func GenSwapFeeRate(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 1, 1e5)), 5)
}

// GenWithdrawFeeRate randomized WithdrawFeeRate ranging from 0.00001 to 1
func GenWithdrawFeeRate(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 1, 1e5)), 5)
}

// GenMaxOrderAmountRatio randomized MaxOrderAmountRatio ranging from 0.00001 to 1
func GenMaxOrderAmountRatio(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 1, 1e5)), 5)
}

// GenUnitBatchSize randomized UnitBatchSize ranging from 1 to 20
func GenUnitBatchSize(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, int(types.DefaultUnitBatchSize), 20))
}

// RandomizedGenState generates a random GenesisState for liquidity
func RandomizedGenState(simState *module.SimulationState) {
	var liquidityPoolTypes []types.PoolType
	simState.AppParams.GetOrGenerate(
		simState.Cdc, LiquidityPoolTypes, &liquidityPoolTypes, simState.Rand,
		func(r *rand.Rand) { liquidityPoolTypes = GenLiquidityPoolTypes(r) },
	)

	var minInitDepositToPool sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinInitDepositToPool, &minInitDepositToPool, simState.Rand,
		func(r *rand.Rand) { minInitDepositToPool = GenMinInitDepositToPool(r) },
	)

	var initPoolCoinMintAmount sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InitPoolCoinMintAmount, &initPoolCoinMintAmount, simState.Rand,
		func(r *rand.Rand) { initPoolCoinMintAmount = GenInitPoolCoinMintAmount(r) },
	)

	var reserveCoinLimitAmount sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ReserveCoinLimitAmount, &reserveCoinLimitAmount, simState.Rand,
		func(r *rand.Rand) { reserveCoinLimitAmount = GenReserveCoinLimitAmount(r) },
	)

	var liquidityPoolCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, LiquidityPoolCreationFee, &liquidityPoolCreationFee, simState.Rand,
		func(r *rand.Rand) { liquidityPoolCreationFee = GenLiquidityPoolCreationFee(r) },
	)

	var swapFeeRate sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SwapFeeRate, &swapFeeRate, simState.Rand,
		func(r *rand.Rand) { swapFeeRate = GenSwapFeeRate(r) },
	)

	var withdrawFeeRate sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, WithdrawFeeRate, &withdrawFeeRate, simState.Rand,
		func(r *rand.Rand) { withdrawFeeRate = GenWithdrawFeeRate(r) },
	)

	var maxOrderAmountRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxOrderAmountRatio, &maxOrderAmountRatio, simState.Rand,
		func(r *rand.Rand) { maxOrderAmountRatio = GenMaxOrderAmountRatio(r) },
	)

	var unitBatchSize uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, UnitBatchSize, &unitBatchSize, simState.Rand,
		func(r *rand.Rand) { unitBatchSize = GenUnitBatchSize(r) },
	)

	liquidityGenesis := types.GenesisState{
		Params: types.Params{
			PoolTypes:                liquidityPoolTypes,
			MinInitDepositToPool:     minInitDepositToPool,
			InitPoolCoinMintAmount:   initPoolCoinMintAmount,
			ReserveCoinLimitAmount:   reserveCoinLimitAmount,
			LiquidityPoolCreationFee: liquidityPoolCreationFee,
			SwapFeeRate:              swapFeeRate,
			WithdrawFeeRate:          withdrawFeeRate,
			MaxOrderAmountRatio:      maxOrderAmountRatio,
			UnitBatchSize:            unitBatchSize,
		},
		PoolRecords: []types.PoolRecord{},
	}

	bz, _ := json.MarshalIndent(&liquidityGenesis, "", " ")
	fmt.Printf("Selected randomly generated liquidity parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&liquidityGenesis)
}
