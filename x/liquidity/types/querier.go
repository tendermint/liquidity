package types

const (
	// QueryLiquidityPool liquidity query endpoint supported by the liquidity querier
	QueryLiquidityPool = "liquidity_pool"
)

// QueryLiquidityPoolParams is the query parameters for 'custom/liquidity'
type QueryLiquidityPoolParams struct {
	PoolId uint64 `json:"pool_id" yaml:"pool_id"`
}


func NewQueryLiquidityPoolParams(poolId uint64) QueryLiquidityPoolParams {
	return QueryLiquidityPoolParams{
		PoolId:poolId,
	}
}