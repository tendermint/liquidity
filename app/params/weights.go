package params

const (
	// liquidity module simulation operation weights for messages
	DefaultWeightMsgCreatePool          int = 5
	DefaultWeightMsgDepositWithinBatch  int = 10
	DefaultWeightMsgWithdrawWithinBatch int = 50
	DefaultWeightMsgSwapWithinBatch     int = 85
)
