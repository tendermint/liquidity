<!--
order: 6
-->

# End-Block

## 1) Append messages to LiquidityPoolBatch

After successful message verification and coin escrow process, the incoming `MsgDepositToLiquidityPool`, `MsgWithdrawFromLiquidityPool`, and `MsgSwap` are appended into the current `LiquidityPoolBatch` of the corresponding `LiquidityPool`.

## 2) Execute LiquidityPoolBatch upon its execution heights

If current `BlockHeight` *mod* `BatchSize` equals *zero*, the `LiquidityPoolBatch` is executed.