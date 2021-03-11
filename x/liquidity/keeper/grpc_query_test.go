package keeper_test

import (
	gocontext "context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

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
				req = &types.QueryLiquidityPoolRequest{}
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

func (suite *KeeperTestSuite) TestGRPCQueryLiquidityPools() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	pools := app.LiquidityKeeper.GetAllLiquidityPools(ctx)

	var req *types.QueryLiquidityPoolsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numPools int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryLiquidityPoolsRequest{
					Pagination: &query.PageRequest{}}
			},
			true,
			2,
			false,
		},
		{"valid request",
			func() {
				req = &types.QueryLiquidityPoolsRequest{
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true}}
			},
			true,
			1,
			true,
		},
		{"valid request",
			func() {
				req = &types.QueryLiquidityPoolsRequest{
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true}}
			},
			true,
			2,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			resp, err := queryClient.LiquidityPools(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(resp)
				suite.Equal(tc.numPools, len(resp.Pools))
				suite.Equal(uint64(len(pools)), resp.Pagination.Total)

				if tc.hasNext {
					suite.NotNil(resp.Pagination.NextKey)
				} else {
					suite.Nil(resp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCLiquidityPoolBatch() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	batch, found := app.LiquidityKeeper.GetLiquidityPoolBatch(ctx, suite.pools[0].PoolId)
	suite.True(found)

	var req *types.QueryLiquidityPoolBatchRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryLiquidityPoolBatchRequest{}
			},
			false,
		},
		{"valid request",
			func() {
				req = &types.QueryLiquidityPoolBatchRequest{PoolId: suite.pools[0].PoolId}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.LiquidityPoolBatch(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(batch.Equal(&res.Batch))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryBatchDepositMsgs() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	msgs := app.LiquidityKeeper.GetAllLiquidityPoolBatchDepositMsgs(ctx, suite.batches[0])

	var req *types.QueryPoolBatchDepositMsgsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numMsgs  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryPoolBatchDepositMsgsRequest{}
			},
			false,
			0,
			false,
		},
		{"returns all the pool batch deposit Msgs",
			func() {
				req = &types.QueryPoolBatchDepositMsgsRequest{
					PoolId: suite.batches[0].PoolId,
				}
			},
			true,
			len(msgs),
			false,
		},
		{"valid request",
			func() {
				req = &types.QueryPoolBatchDepositMsgsRequest{
					PoolId:     suite.batches[0].PoolId,
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true}}
			},
			true,
			1,
			true,
		},
		{"valid request",
			func() {
				req = &types.QueryPoolBatchDepositMsgsRequest{
					PoolId:     suite.batches[0].PoolId,
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true}}
			},
			true,
			len(msgs),
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			resp, err := queryClient.PoolBatchDepositMsgs(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(resp)
				suite.Equal(tc.numMsgs, len(resp.Deposits))
				suite.Equal(uint64(len(msgs)), resp.Pagination.Total)

				if tc.hasNext {
					suite.NotNil(resp.Pagination.NextKey)
				} else {
					suite.Nil(resp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryBatchWithdrawMsgs() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	msgs := app.LiquidityKeeper.GetAllLiquidityPoolBatchWithdrawMsgs(ctx, suite.batches[0])

	var req *types.QueryPoolBatchWithdrawMsgsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numMsgs  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryPoolBatchWithdrawMsgsRequest{}
			},
			false,
			0,
			false,
		},
		{"returns all the pool batch withdraw Msgs",
			func() {
				req = &types.QueryPoolBatchWithdrawMsgsRequest{
					PoolId: suite.batches[0].PoolId,
				}
			},
			true,
			len(msgs),
			false,
		},
		{"valid request",
			func() {
				req = &types.QueryPoolBatchWithdrawMsgsRequest{
					PoolId:     suite.batches[0].PoolId,
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true}}
			},
			true,
			1,
			true,
		},
		{"valid request",
			func() {
				req = &types.QueryPoolBatchWithdrawMsgsRequest{
					PoolId:     suite.batches[0].PoolId,
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true}}
			},
			true,
			len(msgs),
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			resp, err := queryClient.PoolBatchWithdrawMsgs(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(resp)
				suite.Equal(tc.numMsgs, len(resp.Withdraws))
				suite.Equal(uint64(len(msgs)), resp.Pagination.Total)

				if tc.hasNext {
					suite.NotNil(resp.Pagination.NextKey)
				} else {
					suite.Nil(resp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryBatchSwapMsgs() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	msgs := app.LiquidityKeeper.GetAllLiquidityPoolBatchSwapMsgsAsPointer(ctx, suite.batches[0])

	var req *types.QueryPoolBatchSwapMsgsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numMsgs  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryPoolBatchSwapMsgsRequest{}
			},
			false,
			0,
			false,
		},
		{"returns all the pool batch swap Msgs",
			func() {
				req = &types.QueryPoolBatchSwapMsgsRequest{
					PoolId: suite.batches[0].PoolId,
				}
			},
			true,
			len(msgs),
			false,
		},
		{"valid request",
			func() {
				req = &types.QueryPoolBatchSwapMsgsRequest{
					PoolId:     suite.batches[0].PoolId,
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true}}
			},
			true,
			1,
			true,
		},
		{"valid request",
			func() {
				req = &types.QueryPoolBatchSwapMsgsRequest{
					PoolId:     suite.batches[0].PoolId,
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true}}
			},
			true,
			len(msgs),
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			resp, err := queryClient.PoolBatchSwapMsgs(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(resp)
				suite.Equal(tc.numMsgs, len(resp.Swaps))
				suite.Equal(uint64(len(msgs)), resp.Pagination.Total)

				if tc.hasNext {
					suite.NotNil(resp.Pagination.NextKey)
				} else {
					suite.Nil(resp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
