<!--
order: 1
-->

# Concepts


## The Liquidity module on the Cosmos-SDK

The liquidity module serves AMM style decentralized exchange on the Cosmos-SDK. AMM style exchange provides unique coin swap model for its users, liquidity providers and swap requestors.

### Democratized Liquidity Providing

AMM allows liquidity providers to play market maker roles without technically sophisticated real-time orderbook management and significant capital requirement. The liquidity provides only need to deposit coins into liquidity pools, and monitor asset composition changes and accumulated fee rewards from liquidity providing.

It results in democratized liquidity providing activities, hence lowering the cost of liquidity and more enriched quality liquidity provided on the AMM exchange.

### Liquidity Pool

Liquidity pool is a coin reserve with two kinds of coins to provide liquidity for coin swap requests between the two coins in the liquidity pool. The liquidity pool acts as the opposite party of swap requests as the role of market makers in the AMM style exchange.

Liquidity providers deposit the two kinds of coins into the liquidity pool, and share swap fee accumulated in the liquidity pool with respect to their pool share, which is represented as possession of pool coins.

### Coin Swap

Users can request coin swap to a liquidity pool on an AMM style exchange without interacting with constantly changing orderbooks. The requested coin swap is executed with a swap price calculated from given swap price function, the current other swap requests and the current liquidity pool coin reserve status.

### Price Discovery

Coin swap prices in liquidity pools are determined by the current liquidity pool coin reserves and current requested swap amount. Arbitrageurs constantly buy or sell coins in liquidity pools to gain instant profit which results in real-time price discovery of liquidity pools.

### Swap Fees

Coin swap requestors pay swap fees to liquidity pools, which are accumulated in the liquidity pools so that ultimately the pool coin owners will accumulate profit from them.

### Batches and Swap Executions

Coin swaps are executed for every batch, which is composed of one or more consecutive blocks. The size of each batch can be decided by governance parameters and the algorithm in the liquidity module.

