<!-- order: 1 -->

 # Concepts

## Liquidity Module

The liquidity module is a module that can be used in Cosmos SDK applications. It implements decentralized token exchange (DEX) that serves liquidity providing and coin swap functions. Anyone can create a liquidity pool with a pair of coins, provide liquidity by depositing reserve coins into the liquidity pool, and trade coins using the liquidity pool. All the logics are designed in a way that always protect the pool investors.

### How the Liquidity Module democratizes liquidity

The liquidity module democratizes liquidity because anyone can provide liquidity by depositting reserve coins into the liquidity pool and earn fees.

The liquidity module has the following characteristics: 

- Combines a traditional order book-based exchange system with a Uniswap-like AMM mechanism. This hybrid system deepens liquidity for the token swap marketplace.

- Executes batch-style swaps that minimize front-running risk and sub-second latency competition, thereby protecting retail traders.

  - The order book accumulates incoming limit orders into a batch.
  - The liquidity module matches the accumulated limit orders and orders from the liquidity pool at an equivalent swap price at each batch execution height.
  - All limit orders in a batch are treated equally and executed at the same swap price.

## Features of the Liquidity Module

The liquidity module has the following features:
### Liquidity Pool

A liquidity pool is a coin reserve that contains two different types of coins in a trading pair. It has to be unique.

A liquidity provider can be anyone (permissionless) who provides liquidity by depositing reserve coins into the pool. He/she earns the accumulated swap fees with respect to their pool share. Pool share is represented as possession of pool coins. All matchable swap requests are expected to be executed and unmatched ones are removed.

### Pool Creation Fee

The liquidity module has `PoolCreationFee` parameter that is paid upon pool creation. This param is defined at genesis and it can be updated via the governance proposal. The purpose of this fee is to prevent from creating useless pools and making limitless transactions. This fee goes to community fund. 

### Withdrawal Fee Rate

The liquidity module has `WithdrawFeeRate` parameter that is paid upon withdrawal. This param is defined at genesis and it can be updated via the governance proposal. The purpose of this fee is to prevent from making limitless withdrawals.

### Swap Fee Rate

Swap fees are paid upon swap orders. They are accumulated in the pools and are shared among the liquidity providers. The liquidity module implements half-half-fee mechanism that minimizes the impact of fee payment process. Read this [issue](https://github.com/tendermint/liquidity/issues/41) to have more context.
### Swap

You can trade coins using the liquidity pool by making swap requests.Under the hood, coin swaps occur with universal swap ratio for all swap requests.

- The requested coin swap is executed with a swap price that is calculated from the given swap price.

- The current other coin swap requests and the current liquidity pool coin reserve status.

- Swap orders are executed only when the execution swap price is equal to or greater than the submitted order price of the swap order.

### Price Discovery

Swap prices in a liquidity pool are determined by the current liquidity pool coin reserves and requested swap amount.

Arbitrageurs buy or sell coins in a liquidity pool to gain instant profit that results in real-time price discovery of the liquidity pool.
### Escrow Process

The liquidity module uses module account that acts as an escrow account that holds and releases the coin amount during batch execution when there is deposit, withdrawal, or swap order. 
### Batches and Swap Executions

Coin swaps are executed for every batch. A batch is composed of one or more consecutive blocks. The size of each batch can be decided by governance parameters and the algorithm in the liquidity module.

### Pool Identification

The pools in the liquidity module are identified with:

#### Pool Name

- `reserveCoinDenoms1/reserveCoinDenoms2/poolTypeId`

- string join with reserve coin denoms and `poolTypeId`

- Forward slash `/` separator

- Example: `denomX/denomY/1`

#### Pool Reserve Account

- `sdk.AccAddress(crypto.AddressHash([]byte(PoolName)))`

- Example: `cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3` (`D35A0CC16EE598F90B044CE296A405BA9C381E38`)

#### Pool Coin Denom

- `fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(PoolName)))`

- PoolCoinDenomPrefix: `pool`

- Example: `poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4`
