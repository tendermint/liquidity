package types

const (
	// QueryLiquidityPool liquidity query endpoint supported by the liquidity querier
	QueryLiquidityPool = "pool"
)

// QueryLiquidityPoolParams is the query parameters for 'custom/liquidity'
type QueryLiquidityPoolParams struct {
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`
}
