<!--
order: 4
-->

# Messages

## MsgCreateLiquidityPool

```go
type MsgCreateLiquidityPool struct {
	PoolCreatorAddress  string         // account address of the origin of this message
	PoolTypeIndex       uint32         // index of the liquidity pool type of this new liquidity pool
	ReserveCoinDenoms   []string       // list of reserve coin denoms for this new liquidity pool, store in alphabetical order
	DepositCoins 	    sdk.Coins      // deposit coins for initial pool deposit into this new liquidity pool
}
```

**Validity check**
- `MsgCreateLiquidityPool` fails if
  - `PoolCreator` address does not exist
  - `PoolTypeIndex` does not exist in parameters
  - there exists duplicated `LiquidityPool` with same `PoolTypeIndex` and `ReserveCoinDenoms`
  - if one or more coins in ReserveCoinDenoms do not exist in `bank` module
  - if the balance of `PoolCreator` does not have enough amount of coins for `DepositCoins`
  - if the balance of `PoolCreator` does not have enough amount of coins for paying `LiquidityPoolCreationFee`

## MsgDepositToLiquidityPool

```go
type MsgDepositToLiquidityPool struct {
	DepositorAddress    string         // account address of the origin of this message
	PoolId              uint64         // id of the liquidity pool where this message is belong to
	DepositCoins 	    sdk.Coins      // deposit coins of this pool deposit message
}
```

**Validity check**
- `MsgDepositToLiquidityPool` failes if
  - `Depositor` address does not exist
  - `PoolId` does not exist
  - if the denoms of `DepositCoins` are not composed of `ReserveCoinDenoms` of the `LiquidityPool` with given `PoolId`
  - if the balance of `Depositor` does not have enough amount of coins for `DepositCoins`
  
## MsgWithdrawFromLiquidityPool

```go
type MsgWithdrawFromLiquidityPool struct {
	WithdrawerAddress string         // account address of the origin of this message
	PoolId            uint64         // id of the liquidity pool where this message is belong to
	PoolCoin          sdk.Coin       // pool coin sent for reserve coin withdraw
}
```

**Validity check**
- `MsgWithdrawFromLiquidityPool` failes if
  - `Withdrawer` address does not exist
  - `PoolId` does not exist
  - if the denom of `PoolCoin` are not equal to the `PoolCoinDenom` of the `LiquidityPool` with given `PoolId`
  - if the balance of `Depositor` does not have enough amount of coins for `PoolCoin`
  
## MsgSwap

```go
type MsgSwap struct {
	SwapRequesterAddress string     // account address of the origin of this message
	PoolId               uint64     // id of the liquidity pool where this message is belong to
	SwapType             uint32     // swap type of this swap message, default 1: InstantSwap, requesting instant swap
	OfferCoin            sdk.Coin   // offer coin of this swap message
	DemandCoinDenom      sdk.Coin   // denom of demand coin of this swap message
	OrderPrice           sdk.Dec    // order price of this swap message
}
```

**Validity check**
- `MsgSwap` failes if
  - `SwapRequester` address does not exist
  - `PoolId` does not exist
  - `SwapType` does not exist
  - denoms of `OfferCoin` or `DemandCoin` do not exist in `bank` module
  - if the balance of `SwapRequester` does not have enough amount of coins for `OfferCoin`
  - if `OrderPrice` <= zero
