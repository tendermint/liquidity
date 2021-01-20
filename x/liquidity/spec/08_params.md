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
| LiquidityPoolCreationFee | sdk.Coins           | [{"denom":"stake","amount":"100000000"}]                     |
| LiquidityMsgFee          | sdk.Coins           | [{"denom":"stake","amount":"50000"}]                         |
| SwapFeeRate              | string (sdk.Dec)    | "0.003000000000000000"                                       |
| WithdrawFeeRate          | string (sdk.Dec)    | "0.003000000000000000"                                       |
| MaxOrderAmountRatio      | string (sdk.Dec)    | "0.100000000000000000"                                       |
| UnitBatchSize            | uint32              | 1                                                            |

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

## LiquidityPoolCreationFee

Fee paid for new LiquidityPool creation to prevent spamming

## LiquidityMsgFee

Fees paid for each Liquidity Batch messages(swap, withdraw, deposit) distributed to proposers, validators and delegators

## SwapFeeRate

Swap fee rate for every executed swap, when Swap request Reserved half of Swap fee as OfferCoinFee
and remaining half of fee as `ExchangedCoinFee` is collected when batch is executed,   

## WithdrawFeeRate

Reserve token withdrawal with less proportion by `WithdrawFeeRate` to prevent attack vectors from repeated deposit/withdraw

## MaxOrderAmountRatio

Maximum ratio of reserve coins that can be ordered at a swap order

## UnitBatchSize

The smallest unit batch size for every liquidity pool

# Constant Variables

| Key                 | Type   | Constant Value |
| ------------------- | ------ | -------------- |
| CancelOrderLifeSpan | int64  | 0              |
| MinReserveCoinNum   | uint32 | 2              |
| MaxReserveCoinNum   | uint32 | 2              |

## CancelOrderLifeSpan

The life span of swap orders in block heights

## MinReserveCoinNum, MaxReserveCoinNum

min, max number of reserveCoins for LiquidityPoolType on this spec