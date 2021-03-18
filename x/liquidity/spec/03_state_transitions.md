<!--
order: 3
-->

# State Transitions

## Coin Escrow for Liquidity Module Messages

Three messages on the liquidity module need prior coin escrow before confirmation, which causes state transition on `Bank` module. Below lists are describing coin escrow processes for each given message type.

### MsgDepositWithinBatch

To deposit coins into existing `Pool`, the depositor needs to escrow `DepositCoins` into `LiquidityModuleEscrowAccount`.

### MsgWithdrawWithinBatch

To withdraw coins from `Pool`, the withdrawer needs to escrow `PoolCoin` into `LiquidityModuleEscrowAccount`.

### MsgSwapWithinBatch

To request coin swap, swap requestor needs to escrow `OfferCoin` into `LiquidityModuleEscrowAccount`.

## LiquidityPoolBatch Execution

Batch execution causes state transitions on `Bank` module. Below categories describes state transition executed by each process in `PoolBatch` execution.

### Coin Swap

After successful coin swap, coins accumulated in `LiquidityModuleEscrowAccount` for coin swaps are sent to other swap requestors(self-swap) or to the `Pool`(pool-swap). Also fees are sent to the `Pool`.

### LiquidityPool Deposit and Withdraw

For deposit, after successful deposit, escrowed coins are sent to the `ReserveAccount` of targeted `Pool`, and new pool coins are minted and sent to the depositor.

For withdrawal, after successful withdraw, escrowed pool coins are burnt, and corresponding amount of reserve coins are sent to the withdrawer from the `Pool`.

### Pseudo Algorithm for LiquidityPoolBatch Execution

simulation script (in python) : [https://github.com/b-harvest/Liquidity-Module-For-the-Hub/blob/master/pseudo-batch-execution-logic/batch.py](https://github.com/b-harvest/Liquidity-Module-For-the-Hub/blob/master/pseudo-batch-execution-logic/batch.py)

**1) Swap Price Calculation**

**Finding price direction**

- Variables

  - `X` : Reserve of X coin, `Y` : Reserve of Y coin (before this batch execution)
  - `PoolPrice` = `X`/`Y`
  - `XOverLastPrice` : amount of orders which swap X for Y with order price higher than last `PoolPrice`
  - `XAtLastPrice` : amount of orders which swap X for Y with order price equal to last `PoolPrice`
  - `YUnderLastPrice` : amount of orders which swap Y for X with order price lower than last `PoolPrice`
  - `YAtLastPrice` : amount of orders which swap Y for X with order price equal to last `PoolPrice`

- **Increase** : swap price is increased from last `PoolPrice`

  - `XOverLastPrice` > (`YUnderLastPrice`+`YAtLastPrice`)\*`PoolPrice`

- **Decrease** : swap price is decreased from last `PoolPrice`

  - `YUnderLastPrice` > (`XOverLastPrice`+`XAtLastPrice`)/`PoolPrice`

- **Stay** : swap price is not changed from last `PoolPrice`
  - when both above inequalities do not hold

**Stay case**

- Variables
  - `swapPrice` = last `PoolPrice`
  - `EX` : All executable orders which swap X for Y with order price equal or higher than last `PoolPrice`
  - `EY` : All executable orders which swap Y for X with order price equal or lower than last `PoolPrice`
- **ExactMatch** : If `EX` == `EY`\*`swapPrice`
  - Amount of X coins matched from swap orders = `EX`
  - Amount of Y coins matched from swap orders = `EY`
- **FractionalMatch**
  - If `EX` > `EY`\*`swapPrice` : Residual X order amount remains
    - Amount of X coins matched from swap orders = `EY`\*`swapPrice`
    - Amount of Y coins matched from swap orders = `EY`
  - If `EY` > `EX`/`swapPrice` : Residual Y order amount remains
    - Amount of X coins matched from swap orders = `EX`
    - Amount of Y coins matched from swap orders = `EX`/`swapPrice`

**Increase case**

- Iteration : iterate `orderPrice(i)` of all swap orders from low to high

  - variables
    - `EX(i)` : Sum of all order amount of swap orders which swap X for Y with order price equal or higher than this `orderPrice(i)`
    - `EY(i)` : Sum of all order amount of swap orders which swap Y for X with order price equal or lower than this `orderPrice(i)`
  - ExactMatch : SwapPrice is found between two orderPrices
    - `swapPrice(i)` = (`X` + 2*`EX(i)`)/(`Y` + 2*`EY(i-1)`)
      - condition1) `orderPrice(i-1)` < `swapPrice(i)` < `orderPrice(i)`
    - `PoolY(i)` = (`swapPrice(i)`_`Y` - `X`) / (2_`swapPrice(i)`)
      - condition2) `PoolY(i)` >= 0
    - If both above conditions are met, `swapPrice` is the swap price for this iteration
      - Amount of X coins matched = `EX(i)`
    - If one of above conditions doesn’t hold, go to FractionalMatch
  - FractionalMatch : SwapPrice is found at an orderPrice
    - `swapPrice(i)` = `orderPrice(i)`
    - `PoolY(i)` = (`swapPrice(i)`_`Y` - `X`) / (2_`swapPrice(i)`)
    - Amount of X coins matched :
      - `EX(i)` ← min[ `EX(i)`, (`EY(i)`+`PoolY(i)`)*`swapPrice(i)` ]

- Find optimized swapPrice :
  - Find `swapPrice(k)` which has the largest amount of X coins matched
    - this is our optimized swap price
    - corresponding swap result variables
      - `swapPrice(k)`, `EX(k)`, `EY(k)`, `PoolY(k)`

**Decrease case**

- Iteration : iterate `orderPrice(i)` of all swap orders from high to low

  - variables
    - `EX(i)` : Sum of all order amount of swap orders which swap X for Y with order price equal or higher than this `orderPrice(i)`
    - `EY(i)` : Sum of all order amount of swap orders which swap Y for X with order price equal or lower than this `orderPrice(i)`
  - ExactMatch : SwapPrice is found between two orderPrices
    - `swapPrice(i)` = (`X` + 2*`EX(i)`)/(`Y` + 2*`EY(i-1)`)
      - condition1) `orderPrice(i)` < `swapPrice(i)` < `orderPrice(i-1)`
    - `PoolX(i)` = (`X` - `swapPrice(i)`\*`Y`)/2
      - condition2) `PoolX(i)` >= 0
    - If both above conditions are met, `swapPrice` is the swap price for this iteration
      - Amount of Y coins matched = `EY(i)`
    - If one of above conditions doesn’t hold, go to FractionalMatch
  - FractionalMatch : SwapPrice is found at an orderPrice
    - `swapPrice(i)` = `orderPrice(i)`
    - `PoolX(i)` = (`X` - `swapPrice(i)`\*`Y`)/2
    - Amount of Y coins matched :
      - `EY(i)` ← min[ `EY(i)`, (`EX(i)`+`PoolX(i)`)/`swapPrice(i)` ]

- Find optimized swapPrice
  - Find `swapPrice(k)` which has the largest amount of Y coins matched
    - this is our optimized swap price
    - corresponding swap result variables
      - `swapPrice(k)`, `EX(k)`, `EY(k)`, `PoolX(k)`

**Calculate matching result**

- for swap orders from X to Y

  - Iteration : iterate `orderPrice(i)` of swap orders from X to Y (high to low)
    - sort by order price (high to low), sum all order amount with each `orderPrice(i)`
    - if `EX(i)` ≤ `EX(k)`
      - `fractionalRatio` = 1
    - if `EX(i)` > `EX(k)`
      - `fractionalRatio(i)` = (`EX(k)` - `EX(i-1)`) / (`EX(i)` - `EX(i-1)`)
      - break the iteration
    - matching amount for swap orders with this `orderPrice(i)` :
      - `matchingAmt` = `offerAmt` \* `fractionalRatio(i)`

- for swap orders from Y to X
  - Iteration : iterate `orderPrice(i)` of swap orders from Y to X (low to high)
    - sort by order price (low to high), sum all order amount with each `orderPrice(i)`
    - if `EY(i)` ≤ `EY(k)`
      - `fractionalRatio` = 1
    - if `EY(i)` > `EY(k)`
      - `fractionalRatio(i)` = (`EY(k)` - `EY(i-1)`) / (`EY(i)` - `EY(i-1)`)
      - break the iteration
    - matching amount for swap orders with this `orderPrice(i)` :
      - `matchingAmt` = `offerAmt` \* `fractionalRatio(i)`

**2) Swap Fee Payment**

- Swap fees are calculated after above calculation process
- Swap fees are proportional to the coins received from matched swap orders
  - `SwapFee` = `ReceivedMatchedCoin` \* `SwapFeeRate`
- Swap fees are sent to the liquidity pool

**3) Cancel unexecuted swap orders with expired CancelHeight**

After execution of `PoolBatch`, all remaining swap orders with `CancelHeight` equal or higher than current height are cancelled.

**4) Refund escrowed coins**

Refund escrowed coins for cancelled swap order and failed create pool, deposit, withdraw messages.
