package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	//types.RegisterQueryServer(queryHelper, app.LiquidityKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestParams(t *testing.T) {
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
