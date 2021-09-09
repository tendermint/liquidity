# Contents

1. Introduction

2. Order Book

2.1. Context

2.1.1. Restricted Freedom of Traders on DEXs

2.1.2 Context of Advanced Consensus Mechanism: Proof of Stake

2.1.3. AMM on Order Book

2.1.4. Multiple Liquidity Pools for Each Coin Pair

2.1.5. Tick System

2.1.6. Limit, Market, Modify, Cancel Orders

2.2. Order Matching Algorithm

2.2.1. Matching Rules

2.2.2. Matching Process

3. Ranged Liquidity

3.1. Introduction

3.2. Pre-Calculation

3.2.1. Design

3.2.2. Finding "a" and "b"

3.3. State Changes

3.3.1. Creating a New Pool

3.3.2. Deposit/Withdraw

3.3.3. Swap

3.4. Allocation of Pool Liquidity on Order Book

4. Fees

4.1. Introduction

4.2. New Swap Fee Design

4.2.1. Pool Swap Fee

4.2.2. Swap Tax

4.2.3. Virtual Swap Fee Accumulation

5. Dynamic Farming

# 1. Introduction

Gravity DEX is a Cosmos-SDK module on Cosmos Hub which serves the utility of a decentralized exchange. Tendermint team introduce the Gravity DEX V2.0, with expanded freedom for traders to participate on the market, and customized pool creation for improved capital efficiency of liquidity providers. Improved efficiency of the platform is provided by introducing two major new features to the Gravity DEX, "Order Book" and "Ranged Liquidity".

# 2. Order Book

## 2.1. Context

### 2.1.1. Restricted Freedom of Traders on DEXs

By introducing AMM(Automated Market Maker), blockchain industry has innovated how assets are traded among market participants, achieving better financial efficiency and utility. However, because most of the innovation happened in Ethereum ecosystem, its one of the biggest problem, high gas cost, forced developers to give up the most wisely used utility of all existing exchange, the order book. In Ethereum network, it costs too much gas to store and update order book, therefore resulting in removing of such useful utility.

Without order book, traders have very restricted freedom on how to participate in the trading activity because it does not allow trades among traders, but only allow trades with the liquidity pools. This monopolized trading connectivity significantly reduces traders' rights to freely bid and ask among themselves, hence greatly degrades traders' utility and liquidity efficiency of DEXs.

### 2.1.2 Context of Advanced Consensus Mechanism: Proof of Stake

PoS(Proof of Stake) consensus algorithm introduced new ways to create and manage decentralized ledger without high electricity cost for miners and high gas cost for traders. Therefore, DEXs on PoS or dPoS(delegated Proof of Stake) now do not have a reason to not utilize one of the most efficient utility of any kind of exchange, the order book.

### 2.1.3. AMM on Order Book

We introduce a decentralized exchange system where trades are matched on the order book, and liquidity pools provide liquidity on the order book by distributing their liquidity to each price tick of the order book. Liquidity pools locate limit orders on each tick with order amount which is calculated from its AMM equations. This process allows us to transform pool liquidity into executable limit orders, so that order matching can be processed by the widely used order book matching algorithm.

### 2.1.4. Multiple Liquidity Pools for Each Coin Pair

Because of the adoption of ranged liquidity feature on Gravity DEX V2.0, we now can have multiple liquidity pools participating into the marketplace for each coin pair. These liquidity pools will provide liquidity into the corresponding marketplace by transforming the liquidity into split limit orders calculated from AMM equations.

### 2.1.5. Tick System

We introduce tick system in Gravity DEX V2.0, alongside with enabling order book feature. This is a natural consequence because most exchanges with order book have its own tick system. Because we don't define base currency for the exchange, such as USD in Nasdaq, we end up with two possible order book existence (A/B and B/A) for each coin pair A and B.

We expect that the market participants will gather around one order book with more conventional price definition, therefore one specific order book will be mostly used. Because two order books exist for each coin pair, liquidity pool creators should select which order book the liquidity pool is defined to belong to.

### 2.1.6. Limit, Market, Modify, Cancel Orders

For order book utilization, we introduce four most basic order types, limit/market/modify/cancel orders. 

**Limit Order**

- buy/sell order with specific order price
- limit order can be executed at the order price or better price (lower for buy, higher for sell)
- limit orders are automatically cancelled after predefined number of blocks

**Market Order**

- buy/sell orders without specific order price
- a buy market order is equivalent to a buy limit order with the highest possible order price
- a sell market order is equivalent to a sell limit order with the lowest possible order price

**Modify Order**

- existing order in order book can be modified by modify order
- order price and order amount (can be only decreased) can be modified
- if the target order is partially executed at the same block, then rest of the order will be modified
- if the target order is fully executed at the same block, then modify order will fail

**Cancel Order**

- existing order in order book can be cancelled by cancel order
- if the target order is partially executed at the same block, then rest of the order will be cancelled
- if the target order is fully executed at the same block, then cancel order will fail

## 2.2. Order Matching Algorithm

### 2.2.1. Matching Rules

**1) matching priority** : matching priority is decided only by the order price of each order

- buy orders : higher order price has priority
- sell orders : lower order price has priority

**2) matching price** :

- matching phase 1 : last price
- matching phase 2 :
    - price increasing case : each matching sell order price
    - price decreasing case : each matching buy order price

**3) partial matching**

- every order with same order price is matched with equal percentage

### 2.2.2. Matching Process

**1) get last price**

- last price : the last matching price in the last batch
- if does not exist → last price = initial pool price

**2) matching phase 1**

- definitions
    - matchable buy order amount (MBX) : buys with order price higher than last price
    - matchable sell order amount (MSY) : sells with order price lower than last price
    - last price buy order amount (LPBX) : buys with order price equal to last price
    - last price sell order amount (LPSY) : sells with order price equal to last price
    - last price (LP) : last price
- price direction
    - increasing : if MBX > (MSY+LPSY)*LP
    - decreasing : if MSY*LP > MBX+LPBX
    - stay : otherwise
- matching
    - priority : higher buy orders matched with lower sell orders
    - matching price : matching price is at the last price
    - matching condition : matching happens until there exists no more ...
        - buy order with order price equal or higher than last price
        - sell order with order price equal or lower than last price

**3) matching phase 2**

- price increasing case
    - matching
        - priority : higher buy orders matched with lower sell orders
        - matching price : matching price is at each sell order price
        - matching conditions
            - it is matchable when buy order price is equal or higher than sell order price
            - matching happens until the highest buy order price is lower than the lowest sell order price
- price decreasing case
    - matching
        - priority : higher buy orders matched with lower sell orders
        - matching price : matching price is at each buy order price
        - matching conditions
            - it is matchable when sell order price is equal or lower than buy order price
            - matching happens until the highest buy order price is lower than the lowest sell order price
- price stay case
    - no matching for phase 2

# 3. Ranged Liquidity

## 3.1. Introduction

Ranged Liquidity is a new way to provide liquidity for automated market making. Liquidity is provided only within the predefined range of price. Ranged liquidity can be utilized for leveraging the LP positions, and also for stable pools with significantly improved capital efficiency.

## 3.2. Pre-Calculation

### 3.2.1. Design

- shifting constant product curve to left(a) & downward(b)
    - constant product : <img src="https://render.githubusercontent.com/render/math?math=(X%2Ba)*(Y%2Bb) = k"> —— (1)
    - pool price : <img src="https://render.githubusercontent.com/render/math?math=P = (X%2Ba)/(Y%2Bb)"> —— (2)
- ranged pool definition
    - coin pair : X and Y
    - price range : range of the pool price P
        - L : maximum pool price
        - M : minimum pool price

### 3.2.2. Finding "a" and "b"

- swap dx→Y : all Y is traded from dx swap
    - constant product : <img src="https://render.githubusercontent.com/render/math?math=k = (X%2Ba)*(Y%2Bb) = (X%2Bdx%2Ba)*(b)"> → <img src="https://render.githubusercontent.com/render/math?math=dx = (a*Y%2BX*Y)/b">
    - swap price : <img src="https://render.githubusercontent.com/render/math?math=dx/Y = ((a*Y%2BX*Y)/b)/Y = L"> —— (3)
- swap dy→X : all X is traded from dy swap
    - constant product : <img src="https://render.githubusercontent.com/render/math?math=k = (X%2Ba)*(Y%2Bb) = (a)*(Y%2Bdy%2Bb)"> → <img src="https://render.githubusercontent.com/render/math?math=dy = (b*X%2BX*Y)/a">
    - swap price : <img src="https://render.githubusercontent.com/render/math?math=X/dy = X/((b*X%2BX*Y)/a) = M"> —— (4)
- from (3) and (4)
    - <img src="https://render.githubusercontent.com/render/math?math=a = (M*X%2BL*M*Y)/(L-M)"> —— (5)
    - <img src="https://render.githubusercontent.com/render/math?math=b = (X%2BM*Y)/(L-M)"> —— (6)

## 3.3. State Changes

### 3.3.1. Creating a New Pool

- creating a pool with below parameters
    - input parameters
        - X : amount of X token
        - P : initial pool price
        - L : maximum price of the price range of this pool
        - M : minimum price of the price range of this pool
    - target to calculate
        - Y : amount of Y token
            - Y can be derived from above input parameters
        - a, b
- finding calculation targets
    - from (2)
        - <img src="https://render.githubusercontent.com/render/math?math=P = (X%2Ba)/(Y%2Bb) = (X%2B(M*X%2BL*M*Y)/(L-M))/(Y%2B(X%2BM*Y)/(L-M))"> —— (7)
    - solve (7) for Y
        - <img src="https://render.githubusercontent.com/render/math?math=Y = X*(L-P)/(L*P-L*M)"> —— (8)
        - <img src="https://render.githubusercontent.com/render/math?math=X = Y*(L*P-L*M)/(L-P)">
    - combining (5), (6) and (8)
        - <img src="https://render.githubusercontent.com/render/math?math=a = M*X/(P-M)">
        - <img src="https://render.githubusercontent.com/render/math?math=b = P*X/(L*(P-M)) = P*Y/(L-P)">

### 3.3.2. **Deposit/Withdraw**

- deposit/withdraw with the current price P(X,Y)
    - (dx, dy) is deposit amount
        - dx, dy can be negative (withdrawal)
    - price should not change after deposit/withdraw:
        - <img src="https://render.githubusercontent.com/render/math?math=P(X,Y) = P(X%2Bdx,Y%2Bdy)"> —— (9)
- a and b are depend on pool's status
    - <img src="https://render.githubusercontent.com/render/math?math=(X%2Ba)/(Y%2Bb) = (X%2Bdx%2Ba')/(Y%2Bdy%2Bb')">
        - <img src="https://render.githubusercontent.com/render/math?math=(X%2Ba)/(Y%2Bb) = (X%2BM*X/(P-M))/(Y%2BP*X/(L*(P-M)))">
        - <img src="https://render.githubusercontent.com/render/math?math=(X%2Bdx%2Ba')/(Y%2Bdy%2Bb') = (X%2Bdx%2BM*(X%2Bdx)/(P-M))/(Y%2Bdy%2BP*(X%2Bdx)/(L*(P-M)))">
        - <img src="https://render.githubusercontent.com/render/math?math=dx/dy = X/Y">

### 3.3.3. Swap

- price increasing case : swap dx→dy at swap price P
    - finding dy when P is given ( usable Y liquidity at given price P )
        - constant product : <img src="https://render.githubusercontent.com/render/math?math=(X%2Ba)*(Y%2Bb) = (X%2Ba%2BP*dy)*(Y%2Bb-dy)"> —— (10)
        - solving (10) for dy yields
            - <img src="https://render.githubusercontent.com/render/math?math=dy = (Y%2Bb)-(X%2Ba)/P">

- price decreasing case : swap dy→dx at swap price P
    - finding dx when P is given ( usable X liquidity at given price P )
        - constant product : <img src="https://render.githubusercontent.com/render/math?math=(X%2Ba)*(Y%2Bb) = (X%2Ba-dx)*(Y%2Bb%2Bdx/P)"> —— (11)
        - solving (11) for dx yields
            - <img src="https://render.githubusercontent.com/render/math?math=dx = (X%2Ba)-P*(Y%2Bb)">

## 3.4. Allocation of Pool Liquidity on Order Book

- From 3.3.3, it can be inferred that pool can provide Y liquidity for the price which is higher than current pool price and X liquidity for the price lower than current pool price.
    - calculating liquidity of each tick provided by pool
        - Y liquidity of tick price P (price of tick is higher than current pool price)
            - total Y liquidity provided by pool from current price to P (From 3.3.3)
                - <img src="https://render.githubusercontent.com/render/math?math=dy = (Y%2Bb)-(X%2Ba)/P"> —— (12)
            - total Y liquidity provided by pool from current price to P'(price 1 tick lower than P)
                - <img src="https://render.githubusercontent.com/render/math?math=dy' = (Y%2Bb)-(X%2Ba)/P'"> —— (13)
            - subtracting (13) from (12)
                - Y liquidity of tick price P = <img src="https://render.githubusercontent.com/render/math?math=(X%2Ba)/P' - (X%2Ba)/P">
            - If P approaches L, from (5), (6) and (12), dy approaches Y. It means that the pool provides Y liquidity only until predetermined maximum pool price.
        - X liquidity of tick price P (price of tick is lower than current pool price)
            - total X liquidity provided by pool from current price to P (From 3.3.3)
                - <img src="https://render.githubusercontent.com/render/math?math=dx = (X%2Ba)-P*(Y%2Bb)"> —— (14)

            - total Y liquidity provided by pool from current price to P'(price 1 tick higher than P)
                - <img src="https://render.githubusercontent.com/render/math?math=dx' = (X%2Ba)-P'*(Y%2Bb)"> —— (15)
            - subtracting (15) from (14)
                - X liquidity of tick price P = <img src="https://render.githubusercontent.com/render/math?math=P'*(Y%2Bb) - P*(Y%2Bb) = (P'-P)*(Y%2Bb)">
            - If P approaches M, from (5), (6) and (14), dx approaches X. It means that the pool provides X liquidity only until predetermined minimum pool price.

# 4. Fees

## 4.1. Introduction

In Gravity DEX V2.0, with multiple pools with different swap fee rates, combined by order book existence, the swap fee structure should be re-designed to possess economically efficient and accurate financial incentive mechanism for every participants.

## 4.2. New Swap Fee Design

### 4.2.1. Pool Swap Fee

**Pool swap fee : just an adjustment of pool's quoting prices**

- variables
    - <img src="https://render.githubusercontent.com/render/math?math=QuotePrice(i)"> : pool's quoted price at each price tick in order book
    - <img src="https://render.githubusercontent.com/render/math?math=Direction(i)">
        - +1, if pool is selling
        - -1, if pool is buying
    - <img src="https://render.githubusercontent.com/render/math?math=FeeRate"> : swap fee rate of the pool
- effective swap price calculation
    - <img src="https://render.githubusercontent.com/render/math?math=EffectiveSwapPrice(i) = QuotePrice(i)*(1%2BDirection(i)*FeeRate)">

**Adjusted quote price**

- include swap fee inside effective quote price
    - <img src="https://render.githubusercontent.com/render/math?math=AdjQuotePrice(i) = QuotePrice(i)*(1%2BDirection(i)*FeeRate)">
    - <img src="https://render.githubusercontent.com/render/math?math=EffectiveSwapPrice(i) = AdjQuotePrice(i)">
- migrate pool swap fee into adjusted quote price
    - pool swap fee is migrated to adjusted quote price
    - concept of pool swap fee disappears : pool does not need to collect swap fee anymore
    - fee rate concept should be replaced by "QuoteSpread"
        - <img src="https://render.githubusercontent.com/render/math?math=QuoteSpread = (HBP-LAP)/((HBP %2B LAP)/2)">
        - HBP : highest bid price from the pool quotes
        - LAP : lowest ask price from the pool quotes
        - QuoteSpread should be divided by 2 to replace FeeRate

### 4.2.2. Swap Tax

**Introduction**

- swap tax should be paid by every executed order on Gravity DEX
- default swap tax rate : 0%

**Swap Tax Cut Rate**

- some orders can have swap tax cut
    - MakerOrderTaxCutRate
        - all non-pool maker orders get tax cut with this rate
        - default value : 100%
        - maker order definition
            - when price went up from last order book price : selling orders are maker orders
            - when price went down from last order book price : buying orders are maker orders
            - when price stays at last order book price : all orders are maker orders
    - PoolOrderTaxCutRate
        - all orders from pools get tax cut with this rate
        - default value : 100%

**Swap Tax Accumulation**

- swap tax is accumulated in a Liquidity module account
- usage of the swap tax is still not decided yet, and will be decided by a governance procedure

### 4.2.3. Virtual Swap Fee Accumulation

**Reasoning**

- because pool swap fee is retired in Gravity DEX V2.0, it is difficult for liquidity providers to measure past profit from liquidity providing
- hence, we accumulate every profit by liquidity providing from quote spread for every pool, so that liquidity providers can be informed about the profit amount accumulated

**Calculation**

- <img src="https://render.githubusercontent.com/render/math?math=VirtualSwapFee = |\frac{AdjustedPrice}{QuotePrice}-1|*SwapAmount">

**Use for APY Calculation**

- using accumulated virtual swap fee for specific range of time, APY(from quote spread) can be calculated as below

    <img src="https://render.githubusercontent.com/render/math?math=APY = \frac{AccVirtualSwapFee}{PoolReserve}*\frac{365}{AccumulationPeriod}">

- note that this APY is not the total APY because the pool might receive farming rewards in future
