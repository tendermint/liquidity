<!--
order: 8
-->

# Parameters

## Parameters

The liquidity module contains the following parameters:

|Key                                 |Type                |Example                                                                                                                                             |
|------------------------------------|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
|LiquidityPoolTypes                  |[]LiquidityPoolType |[{"description":"ConstantProductLiquidityPool","num_of_reserve_tokens":2,"pool_type_index":0},"swap_price_function_name":"ConstantProductFunction"}]|
|MinInitDepositToPool                |string (sdk.Int)    |"1000000"                                                                                                                                           |
|InitPoolTokenMintAmount             |string (sdk.Int)    |"1000000"                                                                                                                                           |
|SwapFeeRate                         |string (sdk.Dec)    |"0.001000000000000000"                                                                                                                              |
|LiquidityPoolFeeRate                |string (sdk.Dec)    |"0.002000000000000000"                                                                                                                              |

## LiquidityPoolTypes

List of available LiquidityPoolType

```go
type LiquidityPoolType struct {
	PoolTypeIndex         uint32
	NumOfReserveTokens    uint32
	SwapPriceFunctionName string
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
