package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	lapp "github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app          *lapp.LiquidityApp
	ctx          sdk.Context
	addrs        []sdk.AccAddress
	pools        []types.Pool
	batches      []types.PoolBatch
	depositMsgs  []types.DepositMsgState
	withdrawMsgs []types.WithdrawMsgState
	swapMsgs     []types.SwapMsgState
	queryClient  types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	app, ctx := createTestInput()

	querier := keeper.Querier{Keeper: app.LiquidityKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)

	suite.addrs, suite.pools, suite.batches, suite.depositMsgs, suite.withdrawMsgs = createLiquidity(suite.T(), ctx, app)

	suite.ctx = ctx
	suite.app = app

	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestCircuitBreaker(t *testing.T) {
	app, ctx := createTestInput()

	circuitBreaker := app.LiquidityKeeper.GetCircuitBreaker(ctx)
	require.Equal(t, false, circuitBreaker.Enabled)

	app.LiquidityKeeper.SetCircuitBreakerEnabled(ctx, true)

	circuitBreaker = app.LiquidityKeeper.GetCircuitBreaker(ctx)
	require.Equal(t, true, circuitBreaker.Enabled)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
