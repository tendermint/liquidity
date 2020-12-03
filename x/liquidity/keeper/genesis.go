package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	for _, record := range genState.LiquidityPoolRecords {
		k.SetLiquidityPoolRecord(ctx, &record)
		if err := k.ValidateLiquidityPoolRecord(ctx, &record); err != nil {
			panic(err)
		}
	}
	if err := k.ValidateGenesis(ctx, genState); err != nil {
		panic(err)
	}
	// TODO: reset heights variables when init or export
}


func (k Keeper) ValidateGenesis(ctx sdk.Context, genState types.GenesisState) error {
	if err := genState.Params.Validate(); err != nil {
		return err
	}
	for _, record := range genState.LiquidityPoolRecords {
		k.SetLiquidityPoolRecord(ctx, &record)
		if err := k.ValidateLiquidityPoolRecord(ctx, &record); err != nil {
			return err
		}
	}
	return nil
}



// TODO: WIP
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	var poolRecords []types.LiquidityPoolRecord

	pools := k.GetAllLiquidityPools(ctx)
	for _, pool := range pools {
		record, found := k.GetLiquidityPoolRecord(ctx, pool)
		// TODO: verify LiquidityPoolRecord
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
