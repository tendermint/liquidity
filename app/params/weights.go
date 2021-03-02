package params

// Default simulation operation weights for messages and gov proposals
const (
	DefaultWeightMsgSend                        int = 100
	DefaultWeightMsgMultiSend                   int = 10
	DefaultWeightMsgSetWithdrawAddress          int = 50
	DefaultWeightMsgWithdrawDelegationReward    int = 50
	DefaultWeightMsgWithdrawValidatorCommission int = 50
	DefaultWeightMsgFundCommunityPool           int = 50
	DefaultWeightMsgDeposit                     int = 100
	DefaultWeightMsgVote                        int = 67
	DefaultWeightMsgUnjail                      int = 100
	DefaultWeightMsgCreateValidator             int = 100
	DefaultWeightMsgEditValidator               int = 5
	DefaultWeightMsgDelegate                    int = 100
	DefaultWeightMsgUndelegate                  int = 100
	DefaultWeightMsgBeginRedelegate             int = 100

	DefaultWeightCommunitySpendProposal int = 5
	DefaultWeightTextProposal           int = 5
	DefaultWeightParamChangeProposal    int = 5

	// liquidity module simulation operation weights for messages
	DefaultWeightMsgCreateLiquidityPool       int = 100
	DefaultWeightMsgDepositToLiquidityPool    int = 50
	DefaultWeightMsgWithdrawFromLiquidityPool int = 30
	DefaultWeightMsgSwap                      int = 0
	// DefaultWeightMsgCreateLiquidityPool       int = 5
	// DefaultWeightMsgDepositToLiquidityPool    int = 10
	// DefaultWeightMsgWithdrawFromLiquidityPool int = 10
	// DefaultWeightMsgSwap                      int = 90
)
