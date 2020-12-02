package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	app          *app.LiquidityApp
	ctx          sdk.Context
	addrs        []sdk.AccAddress
	pools        []types.LiquidityPool
	batches      []types.LiquidityPoolBatch
	depositMsgs  []types.BatchPoolDepositMsg
	withdrawMsgs []types.BatchPoolWithdrawMsg
	swapMsgs     []types.BatchPoolSwapMsg
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
