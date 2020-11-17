<!--
order: 6
-->

# End-Block

## 1) Create new LiquidityPoolBatch

When there exists no `LiquidityPoolBatch` for the incoming `MsgDepositToLiquidityPool`, `MsgWithdrawFromLiquidityPool`, or `MsgSwap` of corresponding `LiquidityPool`, a new `LiquidityPoolBatch` is created.

## 2) Append messages to LiquidityPoolBatch

After successful message verification and coin escrow process, the incoming `MsgDepositToLiquidityPool`, `MsgWithdrawFromLiquidityPool`, and `MsgSwap` are appended into the current `LiquidityPoolBatch` of the corresponding `LiquidityPool`.

## 3) Execute LiquidityPoolBatch upon its execution heights

If current `BlockHeight` *mod* `BatchSize` equals *zero*, the `LiquidityPoolBatch` is executed.