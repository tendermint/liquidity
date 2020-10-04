<!--
order: 1
-->

# Concepts


## The Liquidity module on the Cosmos-SDK

The liquidity module serves AMM style decentralized exchange on the Cosmos-SDK. AMM style exchange provides unique token swap model for its users, liquidity providers and swap requestors.

### Democratized Liquidity Providing

AMM allows liquidity providers to play market maker roles without technically sophisticated real-time orderbook management and significant capital requirement. The liquidity provides only need to deposit tokens into liquidity pools, and monitor asset composition changes and accumulated fee rewards from liquidity providing.

It results in democratized liquidity providing activities, hence lowering the cost of liquidity and more enriched quality liquidity provided on the AMM exchange.

### Liquidity Pool

Liquidity pool is a token reserve with two kinds of tokens to provide liquidity for token swap requests between the two tokens in the liquidity pool. The liquidity pool acts as the opposite party of swap requests as the role of market makers in the AMM style exchange.

Liquidity providers deposit the two kinds of tokens into the liquidity pool, and share swap fee accumulated in the liquidity pool with respect to their pool share, which is represented as possession of pool tokens.

### Token Swap

Users can request token swap to a liquidity pool on an AMM style exchange without interacting with constantly changing orderbooks. The requested token swap is executed with a swap price calculated from given swap price function, the current other swap requests and the current liquidity pool token reserve status.

### Price Discovery

Token swap prices in liquidity pools are determined by the current liquidity pool token reserves and current requested swap amount. Arbitrageurs constantly buy or sell tokens in liquidity pools to gain instant profit which results in real-time price discovery of liquidity pools.

### Swap Fees

Token swap requestors pay swap fees to liquidity pools, which are accumulated in the liquidity pools so that ultimately the pool token owners will accumulate profit from them.

### Batches and Swap Executions

Token swaps are executed for every batch, which is composed of one or more consecutive blocks. The size of each batch can be decided by governance parameters and the algorithm in the liquidity module.

