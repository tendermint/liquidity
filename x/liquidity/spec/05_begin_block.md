<!--
order: 5
-->

# Begin-Block

## 1) Delete pool batch messages to delete, And reset states for pool batch messages

Delete `{*action}MsgState`s which have `ToBeDeleted` state and then reset states remaining `{*action}MsgState`s, for execute on `end-block` of next batch index

## 2) Reinitialization executed pool batch to next liquidity pool batch

Reinitialization executed `PoolBatch` for to be next batch, the Reinitialization process includes the following actions

- Increase state `BatchIndex` of the batch
- Reset state `BeginHeight` as current block height
- Reset state `Executed` as `false`
