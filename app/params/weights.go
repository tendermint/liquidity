package params

const (
	// liquidity module simulation operation weights for messages
	DefaultWeightMsgCreateLiquidityPool       int = 5
	DefaultWeightMsgDepositToLiquidityPool    int = 10
	DefaultWeightMsgWithdrawFromLiquidityPool int = 50
	DefaultWeightMsgSwap                      int = 85
)
