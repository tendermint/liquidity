package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/client/cli"
	liquiditytestutil "github.com/tendermint/liquidity/x/liquidity/client/testutil"
	liquiditytypes "github.com/tendermint/liquidity/x/liquidity/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tmcli "github.com/tendermint/tendermint/libs/cli"
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
func (s *IntegrationTestSuite) SetupTest() {
	s.T().Log("setting up integration test suite")

	cfg := liquiditytestutil.NewConfig()
	genesisState := cfg.GenesisState
	cfg.NumValidators = 2

	var liquidtyData liquiditytypes.GenesisState
	s.Require().NoError(cfg.Codec.UnmarshalJSON(cfg.GenesisState[liquiditytypes.ModuleName], &liquidtyData))

	// TODO: any params to set for the integration tests?
	liquidtyData.Params = liquiditytypes.DefaultParams()

	liquidtyDataBz, err := cfg.Codec.MarshalJSON(&liquidtyData)
	s.Require().NoError(err)
	genesisState[liquiditytypes.ModuleName] = liquidtyDataBz
	cfg.GenesisState = genesisState

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// TearDownTest cleans up the curret test network after _each_ test.
func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestNewCreatePoolCmd() {
	simapp, ctx := app.CreateTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, liquiditytypes.DefaultParams())

	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("denomX", "denomY")

	initCoins := sdk.NewCoins(
		sdk.NewCoin(denomX, sdk.NewInt(100_000_000)),
		sdk.NewCoin(denomY, sdk.NewInt(100_000_000)),
		sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200)),
	)

	val := s.network.Validators[0]
	app.SaveAccount(simapp, ctx, val.Address, initCoins)

	defaultPoolTypeId := fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid transaction",
			[]string{
				defaultPoolTypeId,
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid transaction (more then two denoms)",
			[]string{
				defaultPoolTypeId,
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000)), sdk.NewCoin("denomZ", sdk.NewInt(100_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCreatePoolCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewDepositWithinBatchCmd()  {}
func (s *IntegrationTestSuite) TestNewWithdrawWithinBatchCmd() {}
func (s *IntegrationTestSuite) TestNewSwapWithinBatchCmd()     {}

func (s *IntegrationTestSuite) TestGetCmdQueryParams() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			`{"pool_types":[{"id":1,"name":"DefaultPoolType","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}],"min_init_deposit_amount":"1000000","init_pool_coin_mint_amount":"1000000","max_reserve_coin_amount":"0","pool_creation_fee":[{"denom":"stake","amount":"100000000"}],"swap_fee_rate":"0.003000000000000000","withdraw_fee_rate":"0.003000000000000000","max_order_amount_ratio":"0.100000000000000000","unit_batch_height":1}`,
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=text", tmcli.OutputFlag)},
			`init_pool_coin_mint_amount: "1000000"
max_order_amount_ratio: "0.100000000000000000"
max_reserve_coin_amount: "0"
min_init_deposit_amount: "1000000"
pool_creation_fee:
- amount: "100000000"
  denom: stake
pool_types:
- description: ""
  id: 1
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: DefaultPoolType
swap_fee_rate: "0.003000000000000000"
unit_batch_height: 1
withdraw_fee_rate: "0.003000000000000000"`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryParams()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPool() {
	val := s.network.Validators[0]

	// TODO: test pool from test_helpers.go

	_, err := s.network.WaitForHeight(4)
	s.Require().NoError(err)

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedOutput string
	}{
		{
			"empty pool id",
			[]string{
				fmt.Sprintf("--%s=3", flags.FlagHeight),
				"",
			},
			true,
			"",
		},
		{
			"invalid pool id",
			[]string{
				fmt.Sprintf("--%s=3", flags.FlagHeight),
				"ABC",
			},
			true,
			"",
		},
		{
			"pool doesn't exist",
			[]string{
				fmt.Sprintf("--%s=3", flags.FlagHeight),
				"2",
			},
			true,
			"",
		},
		{
			"json output",
			[]string{
				fmt.Sprintf("--%s=3", flags.FlagHeight),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				"1",
			},
			false,
			`{"commission":[{"denom":"stake","amount":"116.130000000000000000"}]}`,
		},
		// 		{
		// 			"text output",
		// 			[]string{
		// 				fmt.Sprintf("--%s=3", flags.FlagHeight),
		// 				fmt.Sprintf("--%s=text", tmcli.OutputFlag),
		// 				"1",
		// 			},
		// 			false,
		// 			`commission:
		// - amount: "116.130000000000000000"
		//   denom: stake`,
		// 		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryLiquidityPool()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				fmt.Println("tc.name: ", tc.name)
				fmt.Println("strings.TrimSpace(out.String()): ", out)
				fmt.Println("strings.TrimSpace(out.String()): ", out.String())
				fmt.Println("strings.TrimSpace(out.String()): ", strings.TrimSpace(out.String()))
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPools()        {}
func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPoolBatch()    {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsg()   {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsgs()  {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsg()  {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsgs() {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsg()      {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsgs()     {}
