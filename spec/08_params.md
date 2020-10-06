<!--
order: 8
-->

# Parameters

## Parameters

The liquidity module contains the following parameters:

|Key                                 |Type                |Example                                                                                                                                             |
|------------------------------------|--------------------|----------------------------------------------------------------------------------------------------------------|
|LiquidityPoolTypes                  |[]LiquidityPoolType |[{"pool_type_index":0,</br>"name":"ConstantProductLiquidityPool",</br>"range_of_reserve_coin_num":[2,2],</br>"description":""}]|
|MinInitDepositToPool                |string (sdk.Int)    |"1000000"|
|InitPoolTokenMintAmount             |string (sdk.Int)    |"1000000"|
|SwapFeeRate                         |string (sdk.Dec)    |"0.001000000000000000"|
|LiquidityPoolFeeRate                |string (sdk.Dec)    |"0.002000000000000000"|
|LiquidityPoolCreationFee            |sdk.Coin            |100000000uatom|
|UnitBatchSize  	             |string (sdk.Int)    |"1"|

## LiquidityPoolTypes

List of available LiquidityPoolType

```go
type LiquidityPoolType struct {
	PoolTypeIndex         uint32
	Name		      string
	RangeOfReserveCoinNum []uint32
	Description           string
}
```

## MinInitDepositToPool

Minimum number of tokens to be deposited to the liquidity pool upon pool creation

## InitPoolTokenMintAmount

Initial mint amount of pool token upon pool creation

## SwapFeeRate

Swap fee rate for every executed swap

## LiquidityPoolFeeRate

Liquidity pool fee rate only for swaps consumed pool liquidity

## LiquidityPoolCreationFee

Fee paid for new LiquidityPool creation to prevent spamming

## UnitBatchSize

The smallest unit batch size for every liquidity pool
