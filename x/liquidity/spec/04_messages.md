<!-- order: 4 -->

 # Messages

Messages (Msg) trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps Liquidity Module messages from transactions.

All Liquidity Module messages require a corresponding handler that performs validation logic. See [State Transitions](./03_state_transitions).

## MsgCreatePool

This message is sent when a liquidity pool is created.

```go
type MsgCreatePool struct {
    PoolCreatorAddress  string         // account address of the origin of this message
    PoolTypeId          uint32         // id of the liquidity pool type of this new liquidity pool
    DepositCoins         sdk.Coins      // deposit coins for initial pool deposit into this new liquidity pool
}
```

### Validity checks

The MsgCreatePool message performs these validity checks:

- `MsgCreatePool` fails if

  - `PoolCreator` address does not exist
  - `PoolTypeId` does not exist in parameters
  - A duplicate `LiquidityPool` with same `PoolTypeId` and Reserve Coin Denoms exists
  - One or more coins in ReserveCoinDenoms do not exist in `bank` module
  - The balance of `PoolCreator` does not have enough amount of coins for `DepositCoins`
  - The balance of `PoolCreator` does not have enough amount of coins for paying `PoolCreationFee`

## MsgDepositWithinBatch

This message is sent when a deposit in a liquidity pool batch occurs.

```go
type MsgDepositWithinBatch struct {
    DepositorAddress    string         // account address of the origin of this message
    PoolId              uint64         // id of the liquidity pool where this message is belong to
    DepositCoins         sdk.Coins      // deposit coins of this pool deposit message
}
```

### Validity checks

The MsgDepositWithinBatch message performs these validity checks:

- `MsgDepositWithinBatch` fails if

  - `Depositor` address does not exist
  - `PoolId` does not exist
  - The denoms of `DepositCoins` are not composed of `ReserveCoinDenoms` of the `LiquidityPool` with given `PoolId`
  - The balance of `Depositor` does not have enough coins for `DepositCoins`

## MsgWithdrawWithinBatch

This message is sent when a withdrawal from a liquidity pool batch occurs.

```go
type MsgWithdrawWithinBatch struct {
    WithdrawerAddress string         // account address of the origin of this message
    PoolId            uint64         // id of the liquidity pool where this message is belong to
    PoolCoin          sdk.Coin       // pool coin sent for reserve coin withdraw
}
```

### Validity checks

The MsgWithdrawWithinBatch message performs these validity checks:

- `MsgWithdrawWithinBatch` fails if

  - `Withdrawer` address does not exist
  - `PoolId` does not exist
  - The denom of `PoolCoin` are not equal to the `PoolCoinDenom` of the `LiquidityPool` with given `PoolId`
  - The balance of `Depositor` does not have enough coins for `PoolCoin`

## MsgSwapWithinBatch

This message is sent when coins are swapped between liquidity pools.

```go
type MsgSwapWithinBatch struct {
    SwapRequesterAddress string     // account address of the origin of this message
    PoolId               uint64     // id of the liquidity pool where this message is belong to
    SwapTypeId           uint32     // swap type id of this swap message, default 1: InstantSwap, requesting instant swap
    OfferCoin            sdk.Coin   // offer coin of this swap message
    DemandCoinDenom      string     // denom of demand coin of this swap message
    OfferCoinFee         sdk.Coin   // offer coin fee for pay fees in half offer coin
    OrderPrice           sdk.Dec    // limit order price for the order, the price is the exchange ratio of X/Y where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically
}
```

### Validity checks

The MsgWithdrawWithinBatch message performs these validity checks:

- `MsgSwapWithinBatch` fails if

  - `SwapRequester` address does not exist
  - `PoolId` does not exist
  - `SwapTypeId` does not exist
  - Denoms of `OfferCoin` or `DemandCoin` do not exist in `bank` module
  - The balance of `SwapRequester` does not have enough amount of coins for `OfferCoin`
  - `OrderPrice` <= zero
  - `OfferCoinFee` Equal `OfferCoin` _`params.SwapFeeRate`_ `0.5` with truncating Int
  - Has sufficient balance `OfferCoinFee` to reserve offer coin fee
