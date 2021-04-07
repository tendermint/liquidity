package types

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params Params, liquidityPoolRecords []PoolRecord) *GenesisState {
	return &GenesisState{
		Params:      params,
		PoolRecords: liquidityPoolRecords,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []PoolRecord{}) // TODO: 0 or 1
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	for _, record := range data.PoolRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate Liquidity Pool Record after init or after export
func (record PoolRecord) Validate() error {
	if record.PoolBatch.DepositMsgIndex == 0 ||
		(len(record.DepositMsgStates) > 0 && record.PoolBatch.DepositMsgIndex != record.DepositMsgStates[len(record.DepositMsgStates)-1].MsgIndex+1) {
		return ErrBadBatchMsgIndex
	}
	if record.PoolBatch.WithdrawMsgIndex == 0 ||
		(len(record.WithdrawMsgStates) != 0 && record.PoolBatch.WithdrawMsgIndex != record.WithdrawMsgStates[len(record.WithdrawMsgStates)-1].MsgIndex+1) {
		return ErrBadBatchMsgIndex
	}
	if record.PoolBatch.SwapMsgIndex == 0 ||
		(len(record.SwapMsgStates) != 0 && record.PoolBatch.SwapMsgIndex != record.SwapMsgStates[len(record.SwapMsgStates)-1].MsgIndex+1) {
		return ErrBadBatchMsgIndex
	}
	// TODO: add verify of escrow amount and poolcoin amount with compare to remaining msgs
	return nil
}
