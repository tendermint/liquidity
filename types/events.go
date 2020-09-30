// nolint
package types

// liquidity module event types
const (
	EventTypeCreateLiquidityPool               = "create_liquidity_pool"
	EventTypeDepositToLiquidityPoolToBatch     = "deposit_to_liquidity_pool_to_batch"
	EventTypeWithdrrawFromLiquidityPoolToBatch = "withdraw_from_liquidity_pool_to_batch"
	EventTypeSwapToBatch                       = "swap_to_batch"
	EventTypeDepositToLiquidityPool            = "deposit_to_liquidity_pool"
	EventTypeWithdrrawFromLiquidityPool        = "withdraw_from_liquidity_pool"
	EventTypeSwap                              = "swap"

	AttributeValueLiquidityPoolID        = "liquidity_pool_id"
	AttributeValueLiquidityPoolTypeIndex = "liquidity_pool_type_index"
	AttributeValueLiquidityPoolFeeRate   = "liquidity_pool_fee_rate"
	//AttributeValueSwapPriceFunction      = "swap_price_function"
	AttributeValueReserveTokenDenoms = "reserve_token_denoms"
	AttributeValueReserveAccount     = "reserve_account"
	AttributeValuePoolTokenDenom     = "pool_token_denom"
	AttributeValueSwapFeeRate        = "swap_fee_rate"
	AttributeValueBatchSize          = "batch_size"
	AttributeValueBatchID            = "batch_id"
	AttributeValueTxHash             = "tx_hash"

	AttributeValueDepositor            = "depositor"
	AttributeValueAcceptedTokenAmount  = "accepted_token_amount"
	AttributeValueRefundedTokenAmount  = "refunded_token_amount"
	AttributeValueSuccess              = "success"
	AttributeValueWithdrawer           = "withdrawer"
	AttributeValuePoolTokenAmount      = "pool_token_amount"
	AttributeValueWithdrawTokenAmount  = "withdraw_token_amount"
	AttributeValueSwapRequester        = "swap_requester"
	AttributeValueSwapType             = "swap_type"
	AttributeValueAcceptedOfferToken   = "accepted_offer_token"
	AttributeValueRefundedOfferToken   = "refunded_offer_token"
	AttributeValueReceivedDemandToken  = "received_demand_token"
	AttributeValueSwapPrice            = "swap_price"
	AttributeValuePaidSwapFee          = "paid_swap_fee"
	AttributeValuePaidLiquidityPoolFee = "paid_liquidity_pool_fee"

	AttributeValueCategory = ModuleName
)
