package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

}

// TODO: WIP
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	var poolRecords []types.LiquidityPoolRecord

	pools := k.GetAllLiquidityPools(ctx)
	for _, pool := range pools {
		record, found := k.GetLiquidityPoolRecord(ctx, pool)
		if found {
			poolRecords = append(poolRecords, *record)
		}
	}
	return types.NewGenesisState(params, poolRecords)

	//
	//k.GetAllLiquidityPoolBatches()
	//
	//k.IterateAllLiquidityPools()
	//
	//
	//// each pool?
	//k.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, pool)
	//k.GetPoolMetaData(ctx, pool)
	//
	//k.Iterate
	//
	//return types.NewGenesisState(params, )
	////return types.DefaultGenesisState()
}
