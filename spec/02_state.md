<!--
order: 2
-->

# State

## LiquidityPool

`LiquidityPool` stores static information of a liquidity pool

```go
type LiquidityPool struct {
	PoolID             uint64         // index of this liquidity pool
	PoolTypeIndex      uint32         // pool type of this liquidity pool
	ReserveCoinDenoms  []string       // list of reserve coin denoms for this liquidity pool
	ReserveAccount     sdk.AccAddress // module account address for this liquidity pool to store reserve coins
	PoolCoinDenom      string         // denom of pool coin for this liquidity pool
}
```

LiquidityPool: `0x11 | LiquidityPoolID -> amino(LiquidityPool)`

LiquidityPoolByReserveAccIndex: `0x12 | ReserveAcc -> nil`


## LiquidityPoolBatch

```go
type LiquidityPoolBatch struct {
	PoolID                  uint64                     // id of target liquidity pool
	BatchIndex              uint64                     // index of this batch
	BeginHeight             uint64                     // height where this batch is begun
	SwapMessageList         []BatchSwapMessage         // list of swap messages stored in this batch
	PoolDepositMessageList  []BatchPoolDepositMessage  // list of pool deposit messages stored in this batch
	PoolWithdrawMessageList []BatchPoolWithdrawMessage // list of pool withdraw messages stored in this batch
	ExecutionStatus         bool                       // true if executed, false if not executed yet
}

type BatchSwapMessage struct {
	TxHash    string // tx hash for the original MsgSwap
	MsgHeight uint64 // height where this message is appended to the batch
	Msg       MsgSwap
}

type BatchPoolDepositMessage struct {
	TxHash    string // tx hash for the original MsgDepositToLiquidityPool
	MsgHeight uint64 // height where this message is appended to the batch
	Msg       MsgDepositToLiquidityPool
}

type BatchPoolWithdrawMessage struct {
	TxHash    string // tx hash for the original MsgWithdrawFromLiquidityPool
	MsgHeight uint64 // height where this message is appended to the batch
	Msg       MsgWithdrawFromLiquidityPool
}
```

LiquidityPoolBatchIndex: `0x21 | PoolID -> amino(int64)`

LiquidityPoolBatch: `0x22 | PoolID | BatchIndex -> amino(LiquidityPoolBatch)`
