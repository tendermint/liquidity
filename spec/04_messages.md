<!--
order: 4
-->

# Messages

## MsgCreateLiquidityPool

```go
type MsgCreateLiquidityPool struct {
	PoolCreator         sdk.AccAddress // account address of the origin of this message
	PoolTypeIndex       uint32         // index of the liquidity pool type of this new liquidity pool
	ReserveTokenDenoms  []string       // list of reserve token denoms for this new liquidity pool, store alphabetical
	DepositTokensAmount sdk.Coins      // deposit token for initial pool deposit into this new liquidity pool
}
```

## MsgDepositToLiquidityPool

```go
type MsgDepositToLiquidityPool struct {
	Depositor           sdk.AccAddress // account address of the origin of this message
	PoolID              uint64         // id of the liquidity pool where this message is belong to
	DepositTokensAmount sdk.Coins      // deposit token of this pool deposit message
}
```

## MsgWithdrawFromLiquidityPool

```go
type MsgWithdrawFromLiquidityPool struct {
	Withdrawer      sdk.AccAddress // account address of the origin of this message
	PoolID          uint64         // id of the liquidity pool where this message is belong to
	PoolTokenAmount sdk.Coins      // pool token sent for reserve token withdraw
}
```

## MsgSwap

```go
type MsgSwap struct {
	SwapRequester sdk.AccAddress // account address of the origin of this message
	PoolID        uint64         // id of the liquidity pool where this message is belong to
	PoolTypeIndex uint32         // index of the liquidity pool type where this message is belong to
	SwapType      uint32         // swap type of this swap message, default 1: InstantSwap, requesting instant swap
	OfferToken    sdk.Coin       // offer token of this swap message
	DemandToken   sdk.Coin       // denom of demand token of this swap message
	OrderPrice    sdk.Dec        // order price of this swap message
}
```