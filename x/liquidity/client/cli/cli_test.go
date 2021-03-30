package cli_test

import (
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	testnet "github.com/cosmos/cosmos-sdk/testutil/network"
	liquiditytypes "github.com/tendermint/liquidity/x/liquidity/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

// SetupTest creates a new network for _each_ integration test. We create a new
// network for each test because there are some state modifications that are
// needed to be made in order to make useful queries. However, we don't want
// these state changes to be present in other tests.
func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := testnet.DefaultConfig()
	genesisState := cfg.GenesisState
	cfg.NumValidators = 2

	var liquidtyData liquiditytypes.GenesisState
	s.Require().NoError(cfg.Codec.UnmarshalJSON(genesisState[liquiditytypes.ModuleName], &liquidtyData))

	// TODO: any params to set for the integration tests?
	// liquidtyData.Params.? =

	// liquidtyDataBz, err := cfg.Codec.MarshalJSON(&liquidtyData)
	// s.Require().NoError(err)
	// genesisState[minttypes.ModuleName] = liquidtyDataBz
	// cfg.GenesisState = genesisState

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// TearDownTest cleans up the curret test network after _each_ test.
func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGetCmdQueryParams()                {}
func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPool()         {}
func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPools()        {}
func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPoolBatch()    {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsgs()  {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsg()   {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsgs() {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsg()  {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsgs()     {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsg()      {}

func (s *IntegrationTestSuite) NewCreatePoolCmd() {}

func (s *IntegrationTestSuite) NewDepositWithinBatchCmd() {}

func (s *IntegrationTestSuite) NewWithdrawWithinBatchCmd() {}

func (s *IntegrationTestSuite) NewSwapWithinBatchCmd() {}
