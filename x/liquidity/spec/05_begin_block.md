<!--
order: 5
-->

# Begin-Block

## 1) Delete pool batch messages to delete, And reset statuses for pool batch messages  

Delete `LiquidityPoolBatch{*action}Msgs` which have `ToDelete` status and then reset statuses remaining `LiquidityPoolBatch{*action}Msgs`, for execute on `end-block` of next batch index 

## 2) Reinitialization executed pool batch to next liquidity pool batch 

Reinitialization executed `LiquidityPoolBatch` for to be next batch, the Reinitialization process includes the following actions
- Increase status `BatchIndex` of the batch
- Reset status `BeginHeight` as current block height
- Reset status `Executed` as `false`
