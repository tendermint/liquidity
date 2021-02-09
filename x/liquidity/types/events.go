// nolint
package types

// liquidity module event types, it will be improve the completeness of Milestone 2.
const (
	EventTypeCreateLiquidityPool              = "create_liquidity_pool"
	EventTypeDepositToLiquidityPoolToBatch    = "deposit_to_liquidity_pool_to_batch"
	EventTypeWithdrawFromLiquidityPoolToBatch = "withdraw_from_liquidity_pool_to_batch"
	EventTypeSwapToBatch                      = "swap_to_batch"
	EventTypeDepositToLiquidityPool           = "deposit_to_liquidity_pool"
	EventTypeWithdrawFromLiquidityPool        = "withdraw_from_liquidity_pool"
	EventTypeSwap                             = "swap"

	AttributeValueLiquidityPoolId        = "liquidity_pool_id"
	AttributeValueLiquidityPoolTypeIndex = "liquidity_pool_type_index"
	AttributeValueLiquidityPoolFeeRate   = "liquidity_pool_fee_rate"
	//AttributeValueSwapPriceFunction      = "swap_price_function"
	AttributeValueLiquidityPoolKey  = "liquidity_pool_key"
	AttributeValueReserveCoinDenoms = "reserve_coin_denoms"
	AttributeValueReserveAccount    = "reserve_account"
	AttributeValuePoolCoinDenom     = "pool_coin_denom"
	AttributeValuePoolCoinAmount    = "pool_coin_amount"
	AttributeValueSwapFeeRate       = "swap_fee_rate"
	AttributeValueBatchSize         = "batch_size"
	AttributeValueBatchIndex        = "batch_index"
	AttributeValueMsgIndex          = "msg_index"

	AttributeValueDepositCoins = "deposit_coins"

	AttributeValueOfferCoinDenom     = "offer_coin_denom"
	AttributeValueOfferCoinAmount    = "offer_coin_amount"
	AttributeValueOfferCoinFeeAmount = "offer_coin_fee_amount"
	AttributeValueDemandCoinDenom    = "demand_coin_denom"
	AttributeValueOrderPrice         = "order_price"

	AttributeValueDepositor            = "depositor"
	AttributeValueAcceptedCoinAmount   = "accepted_coin_amount"
	AttributeValueRefundedCoinAmount   = "refunded_coin_amount"
	AttributeValueSuccess              = "success"
	AttributeValueWithdrawer           = "withdrawer"
	AttributeValuePoolCoin             = "pool_coin"
	AttributeValueWithdrawCoinAmount   = "withdraw_coin_amount"
	AttributeValueSwapRequester        = "swap_requester"
	AttributeValueSwapType             = "swap_type"
	AttributeValueAcceptedOfferCoin    = "accepted_offer_coin"
	AttributeValueRefundedOfferCoin    = "refunded_offer_coin"
	AttributeValueReceivedDemandCoin   = "received_demand_coin"
	AttributeValueSwapPrice            = "swap_price"
	AttributeValuePaidSwapFee          = "paid_swap_fee"
	AttributeValuePaidLiquidityPoolFee = "paid_liquidity_pool_fee"

	AttributeValueCategory = ModuleName
)
