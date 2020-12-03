package types

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params Params, liquidityPoolRecords []LiquidityPoolRecord) *GenesisState {
	return &GenesisState{
		Params: params,
		LiquidityPoolRecords: liquidityPoolRecords,
	}
}
//// NewGenesisState is the constructor function for GenesisState
//func NewGenesisState(params Params, liquidityPools []LiquidityPool, liquidityPoolsMetaData []LiquidityPoolMetaData,
//	liquidityPoolBatches []LiquidityPoolBatch, batchPoolDepositMsgs []BatchPoolDepositMsg,
//	batchPoolWithdrawMsgs []BatchPoolWithdrawMsg, batchPoolSwapMsgs []BatchPoolSwapMsg, batchPoolSwapMsgRecords []BatchPoolSwapMsgRecord,
//	liquidityPoolRecords []QueryLiquidityPoolResponse) *GenesisState {
//
//	return &GenesisState{
//		Params: params,
//		LiquidityPools: liquidityPools,
//		LiquidityPoolsMetaData: liquidityPoolsMetaData,
//		LiquidityPoolBatches: liquidityPoolBatches,
//		BatchPoolDepositMsgs: batchPoolDepositMsgs,
//		BatchPoolWithdrawMsgs: batchPoolWithdrawMsgs,
//		BatchPoolSwapMsgs: batchPoolSwapMsgs,
//		BatchPoolSwapMsgRecords: batchPoolSwapMsgRecords,
//		LiquidityPoolRecords: liquidityPoolRecords,
//	}
//}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []LiquidityPoolRecord{})
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	// TODO: add validate only type level without keeper
	for _, record := range data.LiquidityPoolRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate Liquidity Pool Record after init or after export
func (record LiquidityPoolRecord) Validate() error {
	// TODO: add validate only type level without keeper

	if len(record.BatchPoolDepositMsgs)!=0 && record.LiquidityPoolBatch.DepositMsgIndex != record.BatchPoolDepositMsgs[len(record.BatchPoolDepositMsgs)-1].MsgIndex+1 {
		return ErrBadBatchMsgIndex
	}
	if len(record.BatchPoolWithdrawMsgs)!=0 && record.LiquidityPoolBatch.WithdrawMsgIndex != record.BatchPoolWithdrawMsgs[len(record.BatchPoolWithdrawMsgs)-1].MsgIndex {
		return ErrBadBatchMsgIndex
	}
	if len(record.BatchPoolSwapMsgs)!=0 && record.LiquidityPoolBatch.SwapMsgIndex != record.BatchPoolSwapMsgs[len(record.BatchPoolSwapMsgs)-1].MsgIndex {
		return ErrBadBatchMsgIndex
	}

	// TODO: add verify of escrow amount and poolcoin amount with compare to remaining msgs
	return nil
}

