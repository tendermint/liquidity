<!--
order: 2
-->

# State

## Pool

`Pool` stores static information of a liquidity pool

```go
type Pool struct {
	PoolId                 uint64         // index of this liquidity pool
	PoolTypeIndex          uint32         // pool type of this liquidity pool
	ReserveCoinDenoms      []string       // list of reserve coin denoms for this liquidity pool
	ReserveAccountAddress  string         // reserve account address for this liquidity pool to store reserve coins
	PoolCoinDenom          string         // denom of pool coin for this liquidity pool
}
```

Pool: `0x11 | PoolId -> amino(Pool)`

PoolByReserveAccIndex: `0x12 | ReserveAcc -> nil`

GlobalLiquidityPoolIdKey: `[]byte("globalLiquidityPoolId")`

ModuleName, RouterKey, StoreKey, QuerierRoute: `liquidity`

PoolCoinDenomPrefix: `pool`

## PoolBatch

```go
type PoolBatch struct {
	PoolId           uint64  // id of target liquidity pool
	BatchIndex       uint64  // index of this batch
	BeginHeight      uint64  // height where this batch is begun
	DepositMsgIndex  uint64  // last index of BatchPoolDepositMsgs
	WithdrawMsgIndex uint64  // last index of BatchPoolWithdrawMsgs
	SwapMsgIndex     uint64  // last index of BatchPoolSwapMsgs
	Executed         bool    // true if executed, false if not executed yet
}


type DepositMsgState struct {
	MsgHeight  uint64 // height where this message is appended to the batch
	MsgIndex   uint64 // index of this deposit message in this liquidity pool
	Executed   bool   // true if executed on this batch, false if not executed yet
	Succeeded  bool   // true if executed successfully on this batch, false if failed
	ToBeDelete bool   // true if ready to be deleted on kvstore, false if not ready to be deleted
	Msg        MsgDepositToLiquidityPool
}

type WithdrawMsgState struct {
	MsgHeight  uint64 // height where this message is appended to the batch
	MsgIndex   uint64 // index of this withdraw message in this liquidity pool
	Executed   bool   // true if executed on this batch, false if not executed yet
	Succeeded  bool   // true if executed successfully on this batch, false if failed
	ToBeDelete bool   // true if ready to be deleted on kvstore, false if not ready to be deleted
	Msg        MsgWithdrawFromLiquidityPool
}

type SwapMsgState struct {
	MsgHeight          uint64 // height where this message is appended to the batch
	MsgIndex           uint64 // index of this swap message in this liquidity pool
	Executed           bool   // true if executed on this batch, false if not executed yet
	Succeeded          bool   // true if executed successfully on this batch, false if failed
	ToBeDelete         bool   // true if ready to be deleted on kvstore, false if not ready to be deleted
	OrderExpiryHeight  int64  // swap orders are cancelled when current height is equal or higher than ExpiryHeight
	ExchangedOfferCoin sdk.Coin // offer coin exchanged until now
	RemainingOfferCoin sdk.Coin // offer coin currently remaining to be exchanged
	Msg                MsgSwap
}

```

PoolBatchIndex: `0x21 | PoolId -> amino(int64)`

PoolBatch: `0x22 | PoolId -> amino(PoolBatch)`

PoolBatchDepositMsgStateIndex: `0x31 | PoolId -> nil`

PoolBatchDepositMsgStates: `0x31 | PoolId | MsgIndex -> amino(DepositMsgState)`

PoolBatchWithdrawMsgStateIndex: `0x32 | PoolId -> nil`

PoolBatchWithdrawMsgStates: `0x32 | PoolId | MsgIndex -> amino(WithdrawMsgState)`

PoolBatchSwapMsgStateIndex: `0x33 | PoolId -> nil`

PoolBatchSwapMsgStates: `0x33 | PoolId | MsgIndex -> amino(SwapMsgState)`
