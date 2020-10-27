<!--
order: 8
-->

# Parameters

## Parameters

The liquidity module contains the following parameters:

| Key                      | Type                | Example                                                      |
| ------------------------ | ------------------- | ------------------------------------------------------------ |
| LiquidityPoolTypes       | []LiquidityPoolType | [{"pool_type_index":0,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}] |
| MinInitDepositToPool     | string (sdk.Int)    | "1000000"                                                    |
| InitPoolCoinMintAmount   | string (sdk.Int)    | "1000000"                                                    |
| SwapFeeRate              | string (sdk.Dec)    | "0.003000000000000000"                                       |
| LiquidityPoolCreationFee | sdk.Coins           | [{"denom":"uatom","amount":"100000000"}]                     |
| UnitBatchSize            | uint32              | 1                                                            |

## LiquidityPoolTypes

List of available LiquidityPoolType

```go
type LiquidityPoolType struct {
	PoolTypeIndex         uint32
	Name		          string
	MinReserveCoinNum     uint32
	MaxReserveCoinNum     uint32
	Description           string
}
```

## MinInitDepositToPool

Minimum number of coins to be deposited to the liquidity pool upon pool creation

## InitPoolCoinMintAmount

Initial mint amount of pool coin upon pool creation

## SwapFeeRate

Swap fee rate for every executed swap

## LiquidityPoolCreationFee

Fee paid for new LiquidityPool creation to prevent spamming

## UnitBatchSize

The smallest unit batch size for every liquidity pool
