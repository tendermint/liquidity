<!-- order: 5 -->

 # Begin-Block

Begin block operations for the Liquidity Module.

## Delete pool batch messages and reset states for pool batch messages

Delete `{*action}MsgState` messages that have `ToBeDeleted` state and then reset states for the remaining `{*action}MsgState` messages to execute on `end-block` of next batch index

## Reinitialize executed pool batch to next liquidity pool batch

Reinitialize executed `PoolBatch` for the next batch. The Reinitialization process includes the following actions:

- Increase state `BatchIndex` of the batch
- Reset state `BeginHeight` as current block height
- Reset state `Executed` as `false`
