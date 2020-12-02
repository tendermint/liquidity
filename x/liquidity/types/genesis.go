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



	//// todo: it could be delete when if LiquidityPoolRecords is working
	//for _, i := range data.LiquidityPools{
	//	return i.Validate()
	//}
	//
	//// todo: it could be delete when if LiquidityPoolRecords is working
	//for _, i := range data.LiquidityPoolsMetaData{
	//	if err := i.ReserveCoins.Validate(); err != nil {
	//		return err
	//	}
	//	return i.Validate()
	//}
	//// todo: it could be delete when if LiquidityPoolRecords is working
	//for _, i := range data.LiquidityPoolBatches{
	//	if err := i.Validate(); err != nil {
	//		return err
	//	}
	//}
	//for _, i := range data.BatchPoolDepositMsgs{
	//	// TODO: check states, end reset if need
	//	if err := i.Msg.ValidateBasic(); err != nil {
	//		return err
	//	}
	//}
	//for _, i := range data.BatchPoolWithdrawMsgs{
	//	if err := i.Msg.ValidateBasic(); err != nil {
	//		// TODO: check states, end reset if need
	//		return err
	//	}
	//}
	//for _, i := range data.BatchPoolSwapMsgs{
	//	if err := i.Msg.ValidateBasic(); err != nil {
	//		// TODO: check states, end reset if need
	//		return err
	//	}
	//}
	//for _, i := range data.BatchPoolSwapMsgRecords{
	//	if err := i.Msg.ValidateBasic(); err != nil {
	//		// TODO: check not pointer
	//		return err
	//	}
	//}
	//for _, i := range data.LiquidityPoolRecords{
	//	if err := i.LiquidityPool.Validate(); err != nil {
	//		return err
	//	}
	//	//if err := i.LiquidityPoolBatch.Validate(); err != nil {
	//	//	return err
	//	//}
	//	if err := i.LiquidityPoolMetaData.ReserveCoins.Validate(); err != nil {
	//		return err
	//	}
	//}

	return nil
	// TODO: validate
}


