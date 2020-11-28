<!--
order: 6
-->

# Before-End-Block

## 1) Append messages to LiquidityPoolBatch

After successful message verification and coin `escrow` process, the incoming `MsgDepositToLiquidityPool`, `MsgWithdrawFromLiquidityPool`, and `MsgSwap` are appended into the current `LiquidityPoolBatch` of the corresponding `LiquidityPool`

# End-Block

## 1) Execute LiquidityPoolBatch upon its execution heights

If there are `LiquidityPoolBatch{*action}Msgs` that is not yet executed in the `LiquidityPoolBatch` for each `LiquidityPool`, the `LiquidityPoolBatch` is executed. It could contains `DepositLiquidityPool`, `WithdrawLiquidityPool` and `SwapExecution` process.

## 1-A) Transact and Refund for each message

Transactions are made through `escrow`, and refunds are made for cancellations, partial cancellations, expiration, and failed messages.

## 1-B) Set Statuses for each message according to the results

After `1-A`, Update the status of each `LiquidityPoolBatch{*action}Msg` according to the results. Even If the message is completed or expired, Set the `ToDelete` status value as true instead of deleting the message directly from the `end-block` and then delete the messages which have `ToDelete` status from the begin-block in the next block, so that each message with result status in the block can be stored to kvstore, for the past messages with result status can be searched when kvstore is not pruning.
 