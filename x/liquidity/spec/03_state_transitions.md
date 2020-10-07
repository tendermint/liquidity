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

- excel simulation

    - [https://docs.google.com/spreadsheets/d/1yBhDF1DU0b_3ykuLmlvKtdrYKq4F-sg2cVf588TE-ZE/edit#gid=0](https://docs.google.com/spreadsheets/d/1yBhDF1DU0b_3ykuLmlvKtdrYKq4F-sg2cVf588TE-ZE/edit#gid=0)
- process

    1) swap price delta

    - definitions
        - all swap orders are seen as buy/sell limit orders from X coin to Y coin
            - swap order sending X coin to demand Y coin : buy order (of Y coin)
            - swap order sending Y coin to demand X coin : sell order (of Y coin)
            - order price = unit price of Y coin in X coin
        - S = sum of sell order amount with order price equal or lower than current swap price
        - B = sum of buy order amount with order price equal or higher than current swap price
        - NX = number of X coin in the liquidity pool
        - NY = number of X coin in the liquidity pool
        - P(t) = latest swap price from pool coin ratio = NX / NY
        - SwapPrice(t+1) = swap price for this batch ( to find! )
            - P(t) is not equal to SwapPrice(t) !
            - P(t+1) is not equal to SwapPrice(t+1) !
    - swap price delta
        - *if* S ≥ B *then* P(t+1) - P(t) ≤ 0 : price is non-increasing
        - *if* S < B *then* P(t+1) - P(t) ≥ 0 : price is non-decreasing

    2) simulate batch for all order prices of swap requests in the batch ( for price non-decreasing case )

    (step1) finding adjusted price based on constant product equation

    - definitions
        - SimP_i = order price of i-th swap request = the swap price for this simulation
            - SimP_i ≥ P(t) : price non-decreasing case only
                - ignore SimP_i with SimP_i < P(t)
        - SX_i = sum of buy order amount with order price equal or higher than SimP_i, in X coin, which sends X coin and demands Y coin
            - self swap : swap requests which can be matchable without utilizing pool liquidity
        - SY_i = sum of sell order amount with order price equal or lower than SimP_i, in Y coin, which sends Y coin and demands X coin
    - calculation process
        - find AdjP_i for each simulation
            - constant product equation
                - NX*NY = (NX + SX_i - AdjP_i*SY_i) * (NY + SY_i - AdjP_i*SX_i)
                    - *if* SY_i == 0 or SX_i == 0 : above equation is linear equation → unique solution for AdjP_i
                    - *if* SY_i > 0 and SX_i > 0 : above equation is quadratic equation → two solutions can be found for AdjP_i
                        - choose AdjP_i which is nearer to P(t) (less price impact)
            - range criteria for AdjP_i
                - range criteria : AdjP_i should be located at first left or first right of SimP_i
                    - MAX_j(SimP_j | SimP_j < SimP_i) < AdjP_i < MIN_j(SimP_j | SimP_j > SimP_i)
                    - so that the AdjP_i possesses same SX_i and SY_i as SimP_i does
                        - adjustment available only inside the territory of SimP_i
                    - if above inequality does not hold, AdjP_i = SimP_i (fail to adjust price)

    (step2) actual swap simulation

    - definitions
        - PY_i = available pool liquidity amount in Y coin, to be provided for matching, based on constant product equation
        - TY_i = available swap/pool amounts in Y coin, to be provided for matching
        - MX_i = total matched X coin amount by self-swap or pool-swap
        - MSX_i = self matched X coin amount without utilizing pool liquidity
        - MPX_i = pool matched X coin amount via pool liquidity
        - CPEDev_i = deviation of constant product value from NX*NY to the pool status after simulated swap
    - calculation process
        - calculate PY_i
            - constant product equation : NX*NY = (NX + PY_i*AdjP_i)*(NY - PY_i)
            - we can derive PY_i because other variables are known
            - this amount of liquidity provided by the pool can be seen as a limit order from the pool with order price AdjP_i
        - calculate TY_i = SY_i + PY_i
        - calculate MX_i = MIN(SX_i, AdjP_i*TY_i)
        - calculate MSX_i = MIN(AdjP_i*SY_i, MX_i)
        - calculate MPX_i = MIN(MX_i-MSX_i, AdjP_i*PY_i)
        - calculate CPEDev_i = | NX*NY - (NX + MPX_i)*(NY - MPX_i/AdjP_i) |
        - finding optimized swap price from simulations
            - CPEDev_i should be zero : satisfying constant product equation
            - maximize MX_i : maximum swap amount for coin X
                - when there exists multiple simulation with maximum MX : choose one with minimal price impact ( |AdjP_i-P(t)| )
            - the chosen AdjP_max is assigned as SwapPrice(t+1)
            - the chosen simulation result is chosen to become the actual batch execution result

    3) fee payment

    - TBD
