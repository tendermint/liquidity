package cli_test

//
//import (
//	"testing"
//
//	"github.com/cosmos/cosmos-sdk/testutil/network"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/client/testutil"
//	"github.com/stretchr/testify/suite"
//	testnet "github.com/cosmos/cosmos-sdk/testutil/network"
//)
//
//type IntegrationTestSuite struct {
//	suite.Suite
//
//	cfg     network.Config
//	network *network.Network
//}
//
//func (s *IntegrationTestSuite) SetupSuite() {
//	s.T().Log("setting up integration test suite")
//
//	if testing.Short() {
//		s.T().Skip("skipping test in unit-tests mode.")
//	}
//
//	cfg := testnet.DefaultConfig()
//	genesisState := cfg.GenesisState
//	cfg.NumValidators = 1
//
//
//	cfg := network.DefaultConfig()
//	cfg.NumValidators = 2
//
//	s.cfg = cfg
//	s.network = network.New(s.T(), cfg)
//
//	_, err := s.network.WaitForHeight(1)
//	s.Require().NoError(err)
//
//	unbond, err := sdk.ParseCoin("10stake")
//	s.Require().NoError(err)
//
//	val := s.network.Validators[0]
//	val2 := s.network.Validators[1]
//
//	// redelegate
//	_, err = stakingtestutil.MsgRedelegateExec(val.ClientCtx, val.Address, val.ValAddress, val2.ValAddress, unbond)
//	s.Require().NoError(err)
//	_, err = s.network.WaitForHeight(1)
//	s.Require().NoError(err)
//
//	// unbonding
//	_, err = stakingtestutil.MsgUnbondExec(val.ClientCtx, val.Address, val.ValAddress, unbond)
//	s.Require().NoError(err)
//	_, err = s.network.WaitForHeight(1)
//	s.Require().NoError(err)
//}
//
//func (s *IntegrationTestSuite) TearDownSuite() {
//	s.T().Log("tearing down integration test suite")
//	s.network.Cleanup()
//}
