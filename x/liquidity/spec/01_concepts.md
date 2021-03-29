<!-- order: 1 -->

 # Concepts

## Liquidity Module

The Liquidity Module lets you implement a decentralized token exchange (DEX) on any Cosmos SDK-based network. Any user can create a liquidity pool with a pair of tokens, provide liquidity by depositing reserve tokens into the liquidity pool, and trade tokens using the liquidity pool.

### How the Liquidity Module works

These features of the Liquidity Module create incentives to transfer tokens:

- Combines a traditional order book-based exchange system with a Uniswap-like AMM (Automated Market Maker) mechanism. This hybrid system deepens liquidity for the token swap marketplace.

- Executes batch-style swaps that minimize front-running risk and sub-second latency competition, thereby protecting ordinary traders.

  - The order book accumulates incoming limit orders into a batch.
  - The Liquidity Module matches accumulated limit orders and orders from the liquidity pool at an equivalent swap price at each batch execution height.
  - All limit orders in a batch are treated equally and executed at the same swap price.

### Democratized Liquidity

Democratized liquidity lowers the cost of liquidity and provides an enriched quality liquidity provided on the AMM exchange.

A liquidity pool is a collection of funds locked in a smart contract. Use the Liquidity Module to deposit coins into liquidity pools, monitor asset composition changes, and accumulate fee rewards from liquidity providing.

AMM allows liquidity providers to play market maker roles without investing in technically sophisticated real-time orderbook management software or making significant capital investments.

### Liquidity Pool

A liquidity pool contains two assets in a trading pair. A liquidity pool is a coin reserve with a pair of tokens to provide liquidity for coin swap requests between the two coins. The liquidity pool acts as the opposite party of swap requests as the role of market makers in the AMM style exchange.

Liquidity providers deposit the two kinds of coins into the liquidity pool and then share the accumulated swap fee with respect to their pool share. Pool share is represented as possession of pool coins.

### Coin Swap

You can request coin swap to a liquidity pool on an AMM style exchange without using a classic order book mechanism.

1. The requested coin swap is executed with a swap price that is calculated from the given swap price function.

2. The current other coin swap requests and the current liquidity pool coin reserve status.

3. Swap orders are executed only when execution swap price is equal to or greater than the submitted order price of the swap order.

### Price Discovery

Coin swap prices in liquidity pools are determined by the current liquidity pool coin reserves and current requested swap amount. Arbitrageurs constantly buy or sell coins in liquidity pools to gain instant profit that results in real-time price discovery of liquidity pools.

### Escrow Process

For swap orders and deposit and withdrawal messages, the module withdraws the escrow amount of coins from the users' balance.

### Swap Fees

Coin swap requestors pay swap fees to liquidity pools. Swap fees are accumulated in the liquidity pools so that the pool coin owners accumulate profit.

### Batches and Swap Executions

Coin swaps are executed for every batch. A batch is composed of one or more consecutive blocks. The size of each batch can be decided by governance parameters and the algorithm in the liquidity module.

### Pool Identification

#### PoolName

- `reserveCoinDenoms1/reserveCoinDenoms2/poolTypeId`

- string join with reserve coin denoms and `poolTypeId`

- Forward slash `/` separator

- Example: `denomX/denomY/1`

#### PoolReserveAcc

- `sdk.AccAddress(crypto.AddressHash([]byte(PoolName)))`

- Example: `cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3` (`D35A0CC16EE598F90B044CE296A405BA9C381E38`)

#### PoolCoinDenom

- `fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(PoolName)))`

- PoolCoinDenomPrefix: `pool`

- Example: `poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4`
