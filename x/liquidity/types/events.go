// nolint
package types

// liquidity module event types, it will be improve the completeness of Milestone 2.
const (
	EventTypeCreatePool          = TypeMsgCreatePool
	EventTypeDepositWithinBatch  = TypeMsgDepositWithinBatch
	EventTypeWithdrawWithinBatch = TypeMsgWithdrawWithinBatch
	EventTypeSwapWithinBatch     = TypeMsgSwapWithinBatch
	EventTypeDepositToPool       = "deposit_to_pool"
	EventTypeWithdrawFromPool    = "withdraw_from_pool"
	EventTypeSwapTransacted      = "swap_transacted"

	AttributeValuePoolId         = "pool_id"
	AttributeValuePoolTypeId     = "pool_type_id"
	AttributeValuePoolName       = "pool_name"
	AttributeValueReserveAccount = "reserve_account"
	AttributeValuePoolCoinDenom  = "pool_coin_denom"
	AttributeValuePoolCoinAmount = "pool_coin_amount"
	AttributeValueBatchIndex     = "batch_index"
	AttributeValueMsgIndex       = "msg_index"

	AttributeValueDepositCoins = "deposit_coins"

	AttributeValueOfferCoinDenom     = "offer_coin_denom"
	AttributeValueOfferCoinAmount    = "offer_coin_amount"
	AttributeValueOfferCoinFeeAmount = "offer_coin_fee_amount"
	AttributeValueDemandCoinDenom    = "demand_coin_denom"
	AttributeValueOrderPrice         = "order_price"

	AttributeValueDepositor     = "depositor"
	AttributeValueRefundedCoins = "refunded_coins"
	AttributeValueAcceptedCoins = "accepted_coins"
	AttributeValueSuccess       = "success"
	AttributeValueWithdrawer    = "withdrawer"
	AttributeValueWithdrawCoins = "withdraw_coins"
	AttributeValueSwapRequester = "swap_requester"
	AttributeValueSwapTypeId    = "swap_type_id"
	AttributeValueSwapPrice     = "swap_price"

	AttributeValueTransactedCoinAmount       = "transacted_coin_amount"
	AttributeValueRemainingOfferCoinAmount   = "remaining_offer_coin_amount"
	AttributeValueExchangedOfferCoinAmount   = "exchanged_offer_coin_amount"
	AttributeValueReservedOfferCoinFeeAmount = "reserved_offer_coin_fee_amount"
	AttributeValueOrderExpiryHeight          = "order_expiry_height"

	AttributeValueCategory = ModuleName

	Success = "success"
	Failure = "failure"
)
