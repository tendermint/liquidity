package keeper_test

import (
	gocontext "context"
	"fmt"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// TODO: after rebase latest stable sdk 0.40.0
func (suite *KeeperTestSuite) TestGRPCLiquidityPool() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	pool, found := app.LiquidityKeeper.GetLiquidityPool(ctx, suite.pools[0].PoolId)
	suite.True(found)

	var req *types.QueryLiquidityPoolRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryLiquidityPoolRequest{}  // TODO: empty request be 0, to 1
			},
			false,
		},
		{"valid request",
			func() {
				req = &types.QueryLiquidityPoolRequest{PoolId: suite.pools[0].PoolId}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.LiquidityPool(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(pool.Equal(&res.LiquidityPool))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}