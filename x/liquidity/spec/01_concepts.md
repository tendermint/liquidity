<!-- order: 1 -->

 # Concepts

## Liquidity Module

The liquidity module serves an automated market maker (AMM) style decentralized exchange on the Cosmos SDK. An AMM style exchange provides a unique coin swap model for its users, liquidity providers, and swap requestors.

AMMs are a class of decentralized exchanges that rely on mathematical formulas to set the price of a token.

### Democratized Liquidity

AMM allows liquidity providers to play market maker roles without technically sophisticated real-time orderbook management software or significant capital investments.

Use the liquidity module to deposit coins into liquidity pools, monitor asset composition changes, and accumulate fee rewards from liquidity providing.

Democratized liquidity provides activities and lowers the cost of liquidity and provides an enriched quality liquidity provided on the AMM exchange.

### Liquidity Pool

Liquidity pool is a coin reserve with two kinds of coins to provide liquidity for coin swap requests between the two coins in the liquidity pool. A liquidity pool contains two assets in a trading pair. The liquidity pool acts as the opposite party of swap requests as the role of market makers in the AMM style exchange.

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
