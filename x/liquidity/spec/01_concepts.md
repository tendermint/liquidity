<!-- order: 1 -->

 # Concepts

## Liquidity Module

The Liquidity Module lets you implement a decentralized token exchange (DEX) on any Cosmos SDK-based network. Any user can create a liquidity pool with a pair of tokens, provide liquidity by depositing reserve tokens into the liquidity pool, and trade tokens using the liquidity pool.

<!-- use case 1 deposit 2 withdraw 3 swap transaction liquidity = how easily can I do a transaction? some shares are liquid, and easy to sell (google shares for example) however, some shares are not so liquid and then selling these shares is not so easy liquidity offers a service to find out how much something is worth and then being able to trade it without losing value ethereum has decentralized exchanges that sit on the chain (multiple versions) that allow you to swap tokens. (most famous is Uniswap that uses a curve to work out the price of a swap) Now Cosmos Network doesn't have this feature yet. The only way to exchange ATOM is to go through centralized exchange like coinbase. with a decentralized exchange, we have a pool of assets, then there are multiple versions of that pool and the people who put money into the pool do the swap for you order book is traditional finance = match buyers and sellers together through that process, you can find price discovery combine order book with uniswap model = discover buyer and price order book = has all the orders (buys, sells; they are both swaps), knows more about prices liquidity provider = doesn't know about prices liquidity pool = the thing holding the coins withdraw = removing coins from liquidity pool module = works out the price at time of batch execution liquidity provider = pool coin owners -->

 order book (do the thing), then AMM dips into liquidity pool to make up the rest of the orders.

swap price - comes from order book pool price - from AMM

The Liquidity Module provides an AMM

complete backend utility flow of simple AMM (Automated Market Makers) functionality for liquidity providers and swap requestors.

### How the Liquidity Module democratizes liquidity

The Liquidity Module democratizes liquidity because anyone can deposit coins in the liquidity pool and earn fees.

These features of the Liquidity Module create incentives to transfer tokens:

- Combines a traditional order book-based exchange system with a Uniswap-like AMM mechanism. This hybrid system deepens liquidity for the token swap marketplace.

- Executes batch-style swaps that minimize front-running risk and sub-second latency competition, thereby protecting ordinary traders.

  - The order book accumulates incoming limit orders into a batch.
  - The Liquidity Module matches accumulated limit orders and orders from the liquidity pool at an equivalent swap price at each batch execution height.
  - All limit orders in a batch are treated equally and executed at the same swap price.

### Liquidity Pool

A liquidity pool is a coin reserve that contains two kinds of coins in a trading pair.

A liquidity provider (entity or person) deposits the coins into the liquidity pool and then shares the accumulated swap fees with respect to their pool share. Pool share is represented as possession of pool coins.

### Coin Swap

With the Liquidity Module, you can request a coin swap in a liquidity pool.

1. The requested coin swap is executed with a swap price that is calculated from the given swap price.

2. The current other coin swap requests and the current liquidity pool coin reserve status.

3. Swap orders are executed only when the execution swap price is equal to or greater than the submitted order price of the swap order.

### Price Discovery

Coin swap prices in liquidity pools are determined by the current liquidity pool coin reserves and current requested swap amount.

Arbitrageurs buy or sell coins in liquidity pools to gain instant profit that results in real-time price discovery of liquidity pools.

### Escrow Process

For swap orders and deposit and withdrawal transactions, the escrow amount of coins is withdrawn from your balance.

### Swap Fees

Swap fees are paid to the liquidity pools. Swap fees are accumulated in the liquidity pools so that the liquidity provider accumulates profit.

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
