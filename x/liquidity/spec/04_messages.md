<!--
order: 4
-->

# Messages

## MsgCreatePool

```go
type MsgCreatePool struct {
	PoolCreatorAddress  string         // account address of the origin of this message
	PoolTypeId          uint32         // id of the liquidity pool type of this new liquidity pool
	DepositCoins 	    sdk.Coins      // deposit coins for initial pool deposit into this new liquidity pool
}
```

**Validity check**

- `MsgCreatePool` fails if
  - `PoolCreator` address does not exist
  - `PoolTypeId` does not exist in parameters
  - there exists duplicated `LiquidityPool` with same `PoolTypeId` and Reserve Coin Denoms
  - if one or more coins in ReserveCoinDenoms do not exist in `bank` module
  - if the balance of `PoolCreator` does not have enough amount of coins for `DepositCoins`
  - if the balance of `PoolCreator` does not have enough amount of coins for paying `PoolCreationFee`

## MsgDepositWithinBatch

```go
type MsgDepositWithinBatch struct {
	DepositorAddress    string         // account address of the origin of this message
	PoolId              uint64         // id of the liquidity pool where this message is belong to
	DepositCoins 	    sdk.Coins      // deposit coins of this pool deposit message
}
```

**Validity check**

- `MsgDepositWithinBatch` failes if
  - `Depositor` address does not exist
  - `PoolId` does not exist
  - if the denoms of `DepositCoins` are not composed of `ReserveCoinDenoms` of the `LiquidityPool` with given `PoolId`
  - if the balance of `Depositor` does not have enough amount of coins for `DepositCoins`

## MsgWithdrawWithinBatch

```go
type MsgWithdrawWithinBatch struct {
	WithdrawerAddress string         // account address of the origin of this message
	PoolId            uint64         // id of the liquidity pool where this message is belong to
	PoolCoin          sdk.Coin       // pool coin sent for reserve coin withdraw
}
```

**Validity check**

- `MsgWithdrawWithinBatch` failes if
  - `Withdrawer` address does not exist
  - `PoolId` does not exist
  - if the denom of `PoolCoin` are not equal to the `PoolCoinDenom` of the `LiquidityPool` with given `PoolId`
  - if the balance of `Depositor` does not have enough amount of coins for `PoolCoin`

## MsgSwapWithinBatch

```go
type MsgSwapWithinBatch struct {
	SwapRequesterAddress string     // account address of the origin of this message
	PoolId               uint64     // id of the liquidity pool where this message is belong to
	SwapTypeId           uint32     // swap type id of this swap message, default 1: InstantSwap, requesting instant swap
	OfferCoin            sdk.Coin   // offer coin of this swap message
	DemandCoinDenom      string     // denom of demand coin of this swap message
	OfferCoinFee         sdk.Coin   // offer coin fee for pay fees in half offer coin
	OrderPrice           sdk.Dec    // order price of this swap message
}
```

**Validity check**

- `MsgSwapWithinBatch` failes if
  - `SwapRequester` address does not exist
  - `PoolId` does not exist
  - `SwapTypeId` does not exist
  - denoms of `OfferCoin` or `DemandCoin` do not exist in `bank` module
  - if the balance of `SwapRequester` does not have enough amount of coins for `OfferCoin`
  - if `OrderPrice` <= zero
  - if `OfferCoinFee` Equal `OfferCoin` _ `params.SwapFeeRate` _ `0.5` with truncating Int
  - if has sufficient balance `OfferCoinFee` to reserve offer coin fee.
