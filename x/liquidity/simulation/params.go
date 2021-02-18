package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

const (
	keyMinInitDepositToPool   = "KeyMinInitDepositToPool"
	keyInitPoolCoinMintAmount = "KeyInitPoolCoinMintAmount"
	keySwapFeeRate            = "KeySwapFeeRate"
	keyUnitBatchSize          = "KeyUnitBatchSize"
	keyWithdrawFeeRate        = "KeyWithdrawFeeRate"
	keyMaxOrderAmountRatio    = "KeyMaxOrderAmountRatio"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMinInitDepositToPool,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinInitDepositToPool(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyInitPoolCoinMintAmount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenInitPoolCoinMintAmount(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keySwapFeeRate,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSwapFeeRate(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyUnitBatchSize,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenUnitBatchSize(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyWithdrawFeeRate,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenWithdrawFeeRate(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMaxOrderAmountRatio,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMaxOrderAmountRatio(r))
			},
		),
	}
}
