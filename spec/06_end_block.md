<!--
order: 6
-->

# End-Block

## 1) Create New LiquidityPool

`MsgCreateLiquidityPool` is verified and executed in the end block.

After successful verification, a new `LiquidityPool` is created and the initial `DepositTokensAmount` are deposited to the `ReserveAccount` of newly created `LiquidityPool`.

## 2) Create New LiquidityPoolBatch

When there exists no `LiquidityPoolBatch` for the incoming `MsgDepositToLiquidityPool`, `MsgWithdrawFromLiquidityPool`, or `MsgSwap` of corresponding `LiquidityPool`, a new `LiquidityPoolBatch` is created.

And, `LastLiquidityPoolBatchIndex` of the corresponding `LiquidityPool` is updated to the `LiquidityPoolBatchIndex` of the newly created `LiquidityPoolBatch`.

## 3) Append Messages to LiquidityPoolBatch

After successful message verification and token escrow process, the incoming `MsgDepositToLiquidityPool`, `MsgWithdrawFromLiquidityPool`, and `MsgSwap` are appended into the current `LiquidityPoolBatch` of the corresponding `LiquidityPool`.

## 4) Execute LiquidityPoolBatch upon its Execution Heights

If current `BlockHeight` *mod* `BatchSize` of current `LiquidityPoolBatch` equals *zero*, the `LiquidityPoolBatch` is executed.