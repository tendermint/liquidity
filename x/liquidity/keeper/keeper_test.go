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

	app   *app.LiquidityApp
	ctx   sdk.Context
	addrs []sdk.AccAddress
	//vals        []stakingtypes.Validator
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

	//_, validators, vals := createValidators(suite.T(), ctx, app, []int64{9, 8, 7})
	//for i, addr := range validators {
	//	addr := sdk.AccAddress(addr)
	//	app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(addr, pubkeys[i], uint64(i), 0))
	//}
	//suite.vals = vals

	types.RegisterQueryServer(queryHelper, app.LiquidityKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestParams(t *testing.T) {
	//app := simapp.Setup(false)
	//ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	//
	//expParams := types.DefaultParams()
	//
	////check that the empty keeper loads the default
	//resParams := app.StakingKeeper.GetParams(ctx)
	//require.True(t, expParams.Equal(resParams))
	//
	////modify a params, save, and retrieve
	//expParams.MaxValidators = 777
	//app.StakingKeeper.SetParams(ctx, expParams)
	//resParams = app.StakingKeeper.GetParams(ctx)
	//require.True(t, expParams.Equal(resParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
