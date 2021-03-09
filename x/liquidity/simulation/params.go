package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMinInitDepositToPool),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinInitDepositToPool(r).Int64())
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyInitPoolCoinMintAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenInitPoolCoinMintAmount(r).Int64())
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyReserveCoinLimitAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenReserveCoinLimitAmount(r).Int64())
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeySwapFeeRate),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSwapFeeRate(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyWithdrawFeeRate),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenWithdrawFeeRate(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxOrderAmountRatio),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMaxOrderAmountRatio(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyUnitBatchSize),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenUnitBatchSize(r))
			},
		),
	}
}
