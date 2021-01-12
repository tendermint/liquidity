<!--
order: 8
-->

# Parameters

## Parameters

The liquidity module contains the following parameters:

| Key                      | Type                | Example                                                      |
| ------------------------ | ------------------- | ------------------------------------------------------------ |
| LiquidityPoolTypes       | []LiquidityPoolType | [{"pool_type_index":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}] |
| MinInitDepositToPool     | string (sdk.Int)    | "1000000"                                                    |
| InitPoolCoinMintAmount   | string (sdk.Int)    | "1000000"                                                    |
| SwapFeeRate              | string (sdk.Dec)    | "0.003000000000000000"                                       |
| LiquidityPoolCreationFee | sdk.Coins           | [{"denom":"uatom","amount":"100000000"}]                     |

## LiquidityPoolTypes

List of available LiquidityPoolType

```go
type LiquidityPoolType struct {
	PoolTypeIndex         uint32
	Name                  string
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

Swap fee rate for every executed swap, when Swap request Reserved half of Swap fee as OfferCoinFee
and remaining half of fee as `ExchangedCoinFee` is collected when batch is executed,   

## LiquidityPoolCreationFee

Fee paid for new LiquidityPool creation to prevent spamming

# Constant Variables

| Key                 | Type   | Constant Value |
| ------------------- | ------ | -------------- |
| UnitBatchSize       | uint32 | 1              |
| CancelOrderLifeSpan | int64  | 0              |
| MinReserveCoinNum   | uint32 | 2              |
| MaxReserveCoinNum   | uint32 | 2              |

## UnitBatchSize

The smallest unit batch size for every liquidity pool

## CancelOrderLifeSpan

The life span of swap orders in block heights

## MinReserveCoinNum, MaxReserveCoinNum

min, max number of reserveCoins for LiquidityPoolType on this spec