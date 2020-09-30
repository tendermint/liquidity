package types

// -------------- Batch ----------------
// TODO: to protobuf

type LiquidityPoolBatch struct {
	BatchIndex              uint64                     // index of this batch
	PoolID                  uint64                     // id of target liquidity pool
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

// -------------- Swap ----------------

// SwapPriceFunction list
const (
	ConstantProductFunctionName       = "constant_product_function"
	OtherSwapPriceFunctionExampleName = "other_swap_price_function_example"
)

type LiquidityPoolTypeLegacy struct {
	PoolTypeIndex         uint16
	NumOfReserveTokens    uint16
	SwapPriceFunctionName string
	Description           string
}
