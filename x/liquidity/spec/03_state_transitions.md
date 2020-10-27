<!--
order: 3
-->

# State Transitions

## Coin Escrow for Liquidity Module Messages

Three messages on the liquidity module need prior coin escrow before confirmation, which causes state transition on `Bank` module. Below lists are describing coin escrow processes for each given message type.

### MsgDepositToLiquidityPool

To deposit coins into existing `LiquidityPool`, the depositor needs to escrow `DepositCoins` into `LiquidityModuleEscrowAccount`.

### MsgWithdrawFromLiquidityPool

To withdraw coins from `LiquidityPool`, the withdrawer needs to escrow `PoolCoin` into `LiquidityModuleEscrowAccount`.

### MsgSwap

To request coin swap, swap requestor needs to escrow `OfferCoin` into `LiquidityModuleEscrowAccount`.

## LiquidityPoolBatch Execution

Batch execution causes state transitions on `Bank` module. Below categories describes state transition executed by each process in `LiquidityPoolBatch` execution.

### Coin Swap

After successful coin swap, coins accumulated in `LiquidityModuleEscrowAccount` for coin swaps are sent to other swap requestors(self-swap) or to the `LiquidityPool`(pool-swap). Also fees are sent to the `LiquidityPool`.

### LiquidityPool Deposit and Withdraw

For deposit, after successful deposit, escrowed coins are sent to the `ReserveAccount` of targeted `LiquidityPool`, and new pool coins are minted and sent to the depositor.

For withdrawal, after successful withdraw, escrowed pool coins are burnt, and corresponding amount of reserve coins are sent to the withdrawer from the `LiquidityPool`.

### Pseudo Algorithm for LiquidityPoolBatch Execution

**1) Swap Price Calculation**

**Finding price direction**

- Variables
    - `X` : Reserve of X coin, `Y` : Reserve of Y coin (before this batch execution)
    - `PoolSwapPrice` = `X`/`Y`
    - `XOverLastPrice` : amount of orders which swap X for Y with order price higher than last `PoolSwapPrice`
    - `XAtLastPrice` : amount of orders which swap X for Y with order price equal to last `PoolSwapPrice`
    - `YUnderLastPrice` : amount of orders which swap Y for X with order price lower than last `PoolSwapPrice`
    - `YAtLastPrice` : amount of orders which swap Y for X with order price equal to last `PoolSwapPrice`
- **Increase** : swap price is increased from last `PoolSwapPrice`
    - `XOverLastPrice` > (`YUnderLastPrice`+`YAtLastPrice`)*`PoolSwapPrice`
- **Decrease** : swap price is decreased from last `PoolSwapPrice`
    - `YUnderLastPrice` > (`XOverLastPrice`+`XAtLastPrice`)/`PoolSwapPrice`
- **Stay** : swap price is not changed from last `PoolSwapPrice`
    - when both above inequalities do not hold

**Increase case**

- Iteration
    - `orderPrice(i)` : Iterate order prices of all swap orders from **low to high**
    - `EX(i)` : All executable orders which swap X for Y with order price equal or higher than this `orderPrice(i)`
    - `EY(i)` : All executable orders which swap Y for X with order price equal or lower than this `orderPrice(i)`
    - `PoolY`(Y coins provided by the liquidity pool) = `Y` - `X`/`orderPrice(i)`
- Find the `orderPrice(k)` where below value has the first negative number
    - `EX(k)` - `EY(k)`*`orderPrice(k)` - `PoolY`*`orderPrice(k)`
- **ExactMatch** : swapPrice is located between `orderPrice(k-1)` and `orderPrice(k)` ?
    - `swapPrice` = (`X` + `EX(k)`)/(`Y` + `EY(k-1)`)
        - `orderPrice(k-1)` < `swapPrice` < `orderPrice(k)`
    - `adjPoolY` = (`Y`*`EX(k)` - `X`*`EY(k-1)`)/(`X` + `EX(k)`)
        - `adjPoolY` >= 0
    - If both conditions are met, `swapPrice` is the swap price for this batch
        - Amount of X coins matched from swap orders = `EX(k)`
        - Amount of Y coins matched from swap orders = `EY(k-1)`
        - Amount of Y coins provided from liquidity pool = `adjPoolY`
        - Three parts are perfectly matched without fractional match
    - If one of above conditions doesn’t hold, go to next step : FractionalMatch
- **FractionalMatch** : `swapPrice` = `orderPrice(k-1)`
    - Amount of X coins matched from swap orders :
        - `FracEX` = min(`EX(k-1)`, `EY(k-1)`*`swapPrice`+`PoolY`*`swapPrice`)
    - Amount of Y coins matched from swap orders = `EY(k-1)`
    - Amount of Y coins provided from liquidity pool = `PoolY`
    - Swap orders which swap X for Y are fractionally matched
        - `FractionalRatio` = `FracEX` / `EX(k-1)`

**Decrease case**

- Iteration
    - `orderPrice(i)` : Iterate order prices of all swap orders from high to low
    - `EX(i)` : All executable orders which swap X for Y with order price equal or higher than this `orderPrice(i)`
    - `EY(i)` : All executable orders which swap Y for X with order price equal or lower than this `orderPrice(i)`
    - `PoolX`(X coins provided by the liquidity pool) = `X` - `Y`*`orderPrice(i)`
- Find the `orderPrice(k)` where below value has the first negative number
    - `EY(k)` - `EX(k)`/`orderPrice(k)` - `PoolX`/`orderPrice(k)`
- **ExactMatch** : `swapPrice` is located between `orderPrice(k-1)` and `orderPrice(k)` ?
    - `swapPrice` = (`X` + `EX(k-1)`)/(`Y` + `EY(k)`)
        - `orderPrice(k)` < `swapPrice` < `orderPrice(k-1)`
    - `adjPoolX` = (`X`*`EY(k)` - `Y`*`EX(k-1)`)/(`Y` + `EY(k)`)
        - `adjPoolX` >= 0
    - If both conditions are met, swapPrice is the swap price for this batch
        - Amount of X coins matched from swap orders = `EX(k-1)`
        - Amount of Y coins matched from swap orders = `EY(k)`
        - Amount of X coins provided from liquidity pool = `adjPoolX`
        - Three parts are perfectly matched without fractional match
    - If one of above conditions doesn’t hold, go to next step : FractionalMatch
- **FractionalMatch** : `swapPrice` = `orderPrice(k-1)`
    - Amount of Y coins matched from swap orders :
        - `FracEY` = min(`EY(k-1)`, `EX(k-1)`/`swapPrice`+`PoolX`/`swapPrice`)
    - Amount of X coins matched from swap orders = `EX(k-1)`
    - Amount of X coins provided from liquidity pool = `PoolX`
    - Swap orders which swap Y for X are fractionally matched
        - `FractionalRatio` = `FracEY` / `EY(k-1)`

**Stay case**

- Variables
    - `swapPrice` = last `PoolSwapPrice`
    - `EX` : All executable orders which swap X for Y with order price equal or higher than last `PoolSwapPrice`
    - `EY` : All executable orders which swap Y for X with order price equal or lower than last `PoolSwapPrice`
- **ExactMatch** : If `EX` == `EY`*`swapPrice`
    - Amount of X coins matched from swap orders = `EX`
    - Amount of Y coins matched from swap orders = `EY`
    - All two parts are perfectly matched without fractional match
- **FractionalMatch**
    - If `EX` > `EY`*`swapPrice` : Residual X order amount remains
        - Amount of X coins matched from swap orders = `EY`*`swapPrice`
        - Amount of Y coins matched from swap orders = `EY`
    - If `EY` > `EX`/`swapPrice` : Residual Y order amount remains
        - Amount of X coins matched from swap orders = `EX`
        - Amount of Y coins matched from swap orders = `EX`/`swapPrice`

**2) Swap Fee Payment**

- Swap fees are calculated after above calculation process
- Swap fees are proportional to the coins received from matched swap orders
    - `SwapFee` = `ReceivedMatchedCoin` * `SwapFeeRate`
- Swap fees are sent to the liquidity pool
