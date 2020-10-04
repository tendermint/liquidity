<!--
order: 7
-->

# Events


## Handlers

### MsgCreateLiquidityPool

|Type                 |Attribute Key            |Attribute Value      |
|---------------------|-------------------------|---------------------|
|create_liquidity_pool|liquidity_pool_id        |                     |
|create_liquidity_pool|liquidity_pool_type_index|                     |
|create_liquidity_pool|reserve_token_denoms     |                     |
|create_liquidity_pool|reserve_account          |                     |
|create_liquidity_pool|pool_token_denom         |                     |
|create_liquidity_pool|swap_fee_rate            |                     |
|create_liquidity_pool|liquidity_pool_fee_rate  |                     |
|create_liquidity_pool|batch_size               |                     |
|message              |module                   |liquidity            |
|message              |action                   |create_liquidity_pool|
|message              |sender                   |{senderAddress}      |


### MsgDepositToLiquidityPool

|Type                              |Attribute Key|Attribute Value          |
|----------------------------------|-------------|-------------------------|
|deposit_to_liquidity_pool_to_batch|batch_id     |                         |
|message                           |module       |liquidity                |
|message                           |action       |deposit_to_liquidity_pool|
|message                           |sender       |{senderAddress}          |

### MsgWithdrawFromLiquidityPool

|Type                                 |Attribute Key|Attribute Value             |
|-------------------------------------|-------------|----------------------------|
|withdraw_from_liquidity_pool_to_batch|batch_id     |                            |
|message                              |module       |liquidity                   |
|message                              |action       |withdraw_from_liquidity_pool|
|message                              |sender       |{senderAddress}             |

### MsgSwap

|Type         |Attribute Key|Attribute Value|
|-------------|-------------|---------------|
|swap_to_batch|batch_id     |               |
|message      |module       |liquidity      |
|message      |action       |swap           |
|message      |sender       |{senderAddress}|

## EndBlocker

### Batch Result for MsgDepositToLiquidityPool

| Type                      | Attribute Key         | Attribute Value |
| ------------------------- | --------------------- | --------------- |
| deposit_to_liquidity_pool | tx_hash               |                 |
| deposit_to_liquidity_pool | depositor             |                 |
| deposit_to_liquidity_pool | liquidity_pool_id     |                 |
| deposit_to_liquidity_pool | accepted_token_amount |                 |
| deposit_to_liquidity_pool | refunded_token_amount |                 |
| deposit_to_liquidity_pool | success               |                 |

### Batch Result for MsgWithdrawFromLiquidityPool

| Type                         | Attribute Key         | Attribute Value |
| ---------------------------- | --------------------- | --------------- |
| withdraw_from_liquidity_pool | tx_hash               |                 |
| withdraw_from_liquidity_pool | withdrawer            |                 |
| withdraw_from_liquidity_pool | liquidity_pool_id     |                 |
| withdraw_from_liquidity_pool | pool_token_amount     |                 |
| withdraw_from_liquidity_pool | withdraw_token_amount |                 |
| withdraw_from_liquidity_pool | success               |                 |

### Batch Result for MsgSwap

| Type | Attribute Key           | Attribute Value |
| ---- | ----------------------- | --------------- |
| swap | tx_hash                 |                 |
| swap | swap_requester          |                 |
| swap | liquidity_pool_id       |                 |
| swap | swap_type               |                 |
| swap | accepted_offer_token    |                 |
| swap | refunded_offer_token    |                 |
| swap | received_demand_token   |                 |
| swap | swap_price              |                 |
| swap | paid_swap_fee           |                 |
| swap | paid_liquidity_pool_fee |                 |
| swap | success                 |                 |
