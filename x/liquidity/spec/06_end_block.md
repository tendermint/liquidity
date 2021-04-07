<!--
order: 6
-->

# Before-End-Block

## 1) Append messages to LiquidityPoolBatch

After successful message verification and coin `escrow` process, the incoming `MsgDepositWithinBatch`, `MsgWithdrawWithinBatch`, and `MsgSwapWithinBatch` are appended into the current `PoolBatch` of the corresponding `Pool`

# End-Block

## 1) Execute LiquidityPoolBatch upon its execution heights

If there are `{*action}MsgState`s that is not yet executed in the `PoolBatch` for each `Pool`, the `PoolBatch` is executed. It could contains `DepositLiquidityPool`, `WithdrawLiquidityPool` and `SwapExecution` process.

## 1-A) Transact and Refund for each message

Transactions are made through `escrow`, and refunds are made for cancellations, partial cancellations, expiration, and failed messages.

## 1-B) Set states for each message according to the results

After `1-A`, Update the state of each `{*action}MsgState` according to the results. Even If the message is completed or expired, Set the `ToBeDeleted` state value as true instead of deleting the message directly from the `end-block` and then delete the messages which have `ToBeDeleted` state from the begin-block in the next block, so that each message with result state in the block can be stored to kvstore, for the past messages with result state can be searched when kvstore is not pruning.
