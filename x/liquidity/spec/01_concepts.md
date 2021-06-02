<!-- order: 1 -->

 # Concepts

## Liquidity Module

The liquidity module is a module that can be used on any Cosmos SDK based applications. It implements decentralized exchange (DEX) that serves liquidity providing and coin swap functions. Anyone can create a liquidity pool with a pair of coins, provide liquidity by depositing reserve coins into the liquidity pool, and trade coins using the liquidity pool. All the logics are designed in a way that always protect the pool investors.

## Liquidity Pool

A liquidity pool is a coin reserve that contains two different types of coins in a trading pair. It has to be unique. A liquidity provider can be anyone (permissionless) who provides liquidity by depositing reserve coins into the pool. He/she earns the accumulated swap fees with respect to their pool share. Pool share is represented as possession of pool coins. All matchable swap requests are expected to be executed and unmatched ones are removed.
## Equivalent Swap Price Model (ESPM)

The liquidity module is a Cosmos SDK implementation of an AMM system with a novel economic model called the Equivalent Swap Price Model (ESPM). The key distinguishing feature of the ESPM model from the Constant Product Market Maker (CPMM) model (e.g. Uniswap) is the implementation of a hybrid system. This system combines an orderbook model exchange and a simple liquidity pool model by which the order book is governed by a set of order rules and execution is performed in batches. In the ESPM, the pool price is always equal to the last swap price, thus reducing opportunities for arbitrage.

The ESPM model is intended to provide protection against price volatility, transaction ordering vulnerabilities, and losses due to arbitrage, which are widely known vulnerabilities in AMMs such as Uniswap.
## Batch Execution

The liquidity module uses batch execution methodology. Making deposit, withdrawal, or swap orders are accumulated in a liquidity pool for a pre-defined period, which can be one or more blocks in length. Orders are then added to the pool and executed at the end of the batch. The size of each batch is changeable by using governance parameter called `UnitBatchSize`.

## Price Discovery

Swap prices in liquidity pools are determined by the current pool coin reserves and the requested swap amount. Arbitrageurs buy or sell coins in liquidity pools to gain instant profit that results in real-time price discovery of liquidity pools.

## Escrow Process

The liquidity module uses module account that acts as an escrow account that holds and releases the coin amount during batch execution.

## Fees
### PoolCreationFee

The liquidity module has `PoolCreationFee` parameter that is paid upon pool creation. This param is defined at genesis and it can be updated via the governance proposal. The purpose of this fee is to prevent from creating useless pools and making limitless transactions. This fee goes to community fund.
### WithdrawalFeeRate

The liquidity module has `WithdrawFeeRate` parameter that is paid upon withdrawal. This param is defined at genesis and it can be updated via the governance proposal. The purpose of this fee is to prevent from making limitless withdrawals.

### SwapFeeRate

Swap fees are paid upon swap orders. They are accumulated in the pools and are shared among the liquidity providers. The liquidity module implements half-half-fee mechanism that minimizes the impact of fee payment process. Read this issue to have more context.
## Pool Identification

The pools in the liquidity module are identified with:
### PoolName

- Concatenate two different reserve coin denoms and pool type id and forward slash `/` separator. 
  - Example: `uatom/stake/1`
### PoolReserveAccount

- `sdk.AccAddress(crypto.AddressHash([]byte(PoolName)))`
  - Example: `cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3` (`D35A0CC16EE598F90B044CE296A405BA9C381E38`)
### PoolCoinDenom

- `fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(PoolName)))`
- Use `PoolCoinDenomPrefix` for `pool`
  - Example: `poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4`



