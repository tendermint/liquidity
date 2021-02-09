<!--
order: 7
-->

# Events

## Handlers

### MsgCreateLiquidityPool

| Type                  | Attribute Key             | Attribute Value          |
| --------------------- | ------------------------- | ------------------------ |
| create_liquidity_pool | liquidity_pool_id         | {liquidityPoolId}        |
| create_liquidity_pool | liquidity_pool_type_index | {liquidityPoolTypeIndex} |
| create_liquidity_pool | liquidity_pool_key        |{AttributeValueLiquidityPoolKey}|
| create_liquidity_pool | reserve_account           | {reserveAccountAddress}  |
| create_liquidity_pool | deposit_coins             | {depositCoins}           |
| create_liquidity_pool | pool_coin_denom           | {poolCoinDenom}          |
| message               | module                    | liquidity                |
| message               | action                    | create_liquidity_pool    |
| message               | sender                    | {senderAddress}          |

### MsgDepositToLiquidityPool

| Type                               | Attribute Key     | Attribute Value           |
| ---------------------------------- | ----------------- | ------------------------- |
| deposit_to_liquidity_pool_to_batch | liquidity_pool_id | {liquidityPoolId}         |
| deposit_to_liquidity_pool_to_batch | batch_index       | {batchIndex}              |
| deposit_to_liquidity_pool_to_batch | msg_index         | {depositMsgIndex}         |
| deposit_to_liquidity_pool_to_batch | deposit_coins     | {depositCoins}            |
| message                            | module            | liquidity                 |
| message                            | action            | deposit_to_liquidity_pool |
| message                            | sender            | {senderAddress}           |

### MsgWithdrawFromLiquidityPool

| Type                                  | Attribute Key     | Attribute Value              |
| ------------------------------------- | ----------------- | ---------------------------- |
| withdraw_from_liquidity_pool_to_batch | liquidity_pool_id | {liquidityPoolId}            |
| withdraw_from_liquidity_pool_to_batch | batch_index       | {batchIndex}                 |
| withdraw_from_liquidity_pool_to_batch | msg_index         | {withdrawMsgIndex}           |
| withdraw_from_liquidity_pool_to_batch | pool_coin_denom   | {poolCoinDenom}              |
| withdraw_from_liquidity_pool_to_batch | pool_coin_amount  | {poolCoinAmount}             |
| message                               | module            | liquidity                    |
| message                               | action            | withdraw_from_liquidity_pool |
| message                               | sender            | {senderAddress}              |

### MsgSwap

| Type          | Attribute Key     | Attribute Value   |
| ------------- | ----------------- | ----------------- |
| swap_to_batch | liquidity_pool_id | {liquidityPoolId} |
| swap_to_batch | batch_index       | {batchIndex}      |
| swap_to_batch | msg_index         | {swapMsgIndex}    |
| swap_to_batch | swap_type         | {swapType}        |
| swap_to_batch | offer_coin_denom  | {offerCoinDenom}  |
| swap_to_batch | offer_coin_amount | {offerCoinAmount} |
| swap_to_batch | demand_coin_denom | {demandCoinDenom} |
| swap_to_batch | order_price       | {orderPrice}      |
| message       | module            | liquidity         |
| message       | action            | swap              |
| message       | sender            | {senderAddress}   |

## EndBlocker

### Batch Result for MsgDepositToLiquidityPool

| Type                      | Attribute Key     | Attribute Value    |
| ------------------------- | ----------------- | ------------------ |
| deposit_to_liquidity_pool | liquidity_pool_id | {liquidityPoolId}  |
| deposit_to_liquidity_pool | batch_index       | {batchIndex}       |
| deposit_to_liquidity_pool | msg_index         | {depositMsgIndex}  |
| deposit_to_liquidity_pool | depositor         | {depositorAddress} |
| deposit_to_liquidity_pool | accepted_coins    | {acceptedCoins}    |
| deposit_to_liquidity_pool | refunded_coins    | {refundedCoins}    |
| deposit_to_liquidity_pool | pool_coin_denom   | {poolCoinDenom}     |
| deposit_to_liquidity_pool | pool_coin_amount  | {poolCoinAmount}    |
| deposit_to_liquidity_pool | success           | {success}          |

### Batch Result for MsgWithdrawFromLiquidityPool

| Type                         | Attribute Key     | Attribute Value     |
| ---------------------------- | ----------------- | ------------------- |
| withdraw_from_liquidity_pool | liquidity_pool_id | {liquidityPoolId}   |
| withdraw_from_liquidity_pool | batch_index       | {batchIndex}        |
| withdraw_from_liquidity_pool | msg_index         | {withdrawMsgIndex}  |
| withdraw_from_liquidity_pool | withdrawer        | {withdrawerAddress} |
| withdraw_from_liquidity_pool | pool_coin_denom   | {poolCoinDenom}     |
| withdraw_from_liquidity_pool | pool_coin_amount  | {poolCoinAmount}    |
| withdraw_from_liquidity_pool | withdraw_coins    | {withdrawCoins}     |
| withdraw_from_liquidity_pool | success           | {success}           |

### Batch Result for MsgSwap

| Type            | Attribute Key               | Attribute Value            |
| --------------- | --------------------------- | -------------------------- |
| swap_transacted | liquidity_pool_id           | {liquidityPoolId}          |
| swap_transacted | batch_index                 | {batchIndex}               |
| swap_transacted | msg_index                   | {swapMsgIndex}             |
| swap_transacted | swap_requester              | {swapRequesterAddress}     |
| swap_transacted | swap_type                   | {swapType}                 |
| swap_transacted | offer_coin_denom            | {offerCoinDenom}           |
| swap_transacted | offer_coin_amount           | {offerCoinAmount}          |
| swap_transacted | order_price                 | {orderPrice}               |
| swap_transacted | swap_price                  | {swapPrice}                |
| swap_transacted | transacted_coin_amount      | {transactedCoinAmount}     |
| swap_transacted | remaining_offer_coin_amount | {remainingOfferCoinAmount} |
| swap_transacted | exchanged_offer_coin_amount | {exchangedOfferCoinAmount} |
| swap_transacted | offer_coin_fee_amount       | {offerCoinFeeAmount}       |
| swap_transacted | offer_coin_fee_reserve_amount   | {offerCoinFeeReserveAmount}    |
| swap_transacted | order_expiry_height         | {orderExpiryHeight}        |
| swap_transacted | success                     | {success}                  |

### Cancel Result for MsgSwap on Batch

| Type        | Attribute Key               | Attribute Value            |
| ----------- | --------------------------- | -------------------------- |
| swap_cancel | liquidity_pool_id           | {liquidityPoolId}          |
| swap_cancel | batch_index                 | {batchIndex}               |
| swap_cancel | msg_index                   | {swapMsgIndex}             |
| swap_cancel | swap_requester              | {swapRequesterAddress}     |
| swap_cancel | swap_type                   | {swapType}                 |
| swap_cancel | offer_coin_denom            | {offerCoinDenom}           |
| swap_cancel | offer_coin_amount           | {offerCoinAmount}          |
| swap_cancel | offer_coin_fee_amount       | {offerCoinFeeAmount}       |
| swap_cancel | offer_coin_fee_reserve_amount   | {offerCoinFeeReserveAmount}    |
| swap_cancel | order_price                 | {orderPrice}               |
| swap_cancel | swap_price                  | {swapPrice}                |
| swap_cancel | cancelled_coin_amount       | {cancelledOfferCoinAmount} |
| swap_cancel | remaining_offer_coin_amount | {remainingOfferCoinAmount} |
| swap_cancel | order_expiry_height         | {orderExpiryHeight}        |
| swap_cancel | success                     | {success}                  |
