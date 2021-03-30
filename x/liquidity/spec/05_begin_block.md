<!-- order: 5 -->

# Begin-Block

## Delete pool batch messages to delete and reset states for pool batch messages

Delete `{*action}MsgState messages that have`ToBeDeleted`state and then reset states remaining`{*action}MsgState`messages for execute on`end-block` of next batch index

## Reinitialize executed pool batch to next liquidity pool batch

Reinitialization executed `PoolBatch` for to be next batch, the Reinitialization process includes the following actions

- Increase state `BatchIndex` of the batch
- Reset state `BeginHeight` as current block height
- Reset state `Executed` as `false`
