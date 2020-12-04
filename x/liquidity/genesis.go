package liquidity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// InitGenesis new liquidity genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	// validate logic on module.go/InitGenesis
	for _, record := range data.LiquidityPoolRecords {
		keeper.SetLiquidityPoolRecord(ctx, &record)
	}
	keeper.InitGenesis(ctx, data)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return keeper.ExportGenesis(ctx)
}
