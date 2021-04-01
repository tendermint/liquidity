package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

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

	liquidtyData.Params = liquiditytypes.DefaultParams()

	liquidtyDataBz, err := cfg.Codec.MarshalJSON(&liquidtyData)
	s.Require().NoError(err)

	genesisState[liquiditytypes.ModuleName] = liquidtyDataBz
	cfg.GenesisState = genesisState

	cfg.AccountTokens = sdk.NewInt(100_000_000_000) // node%dtoken denom
	cfg.StakingTokens = sdk.NewInt(100_000_000_000) // stake denom

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
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

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
				fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid pool type id",
			[]string{
				"pooltypeidisnotnumber",
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"invalid number of denoms",
			[]string{
				fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000)), sdk.NewCoin("denomZ", sdk.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"pool type id not available",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000)), sdk.NewCoin("denomZ", sdk.NewInt(100_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"deposit coin less than minimum deposit amount",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(1_000)), sdk.NewCoin(denomY, sdk.NewInt(1_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
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

func (s *IntegrationTestSuite) TestNewDepositWithinBatchCmd() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

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
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000_000)), sdk.NewCoin(denomY, sdk.NewInt(10_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid pool id",
			[]string{
				"invalidpoolid",
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(1_000_000)), sdk.NewCoin(denomY, sdk.NewInt(1_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"invalid number of denoms",
			[]string{
				fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(1_000_000)), sdk.NewCoin(denomY, sdk.NewInt(1_000_000)), sdk.NewCoin("denomZ", sdk.NewInt(1_000_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewDepositWithinBatchCmd()
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

func (s *IntegrationTestSuite) TestNewWithdrawWithinBatchCmd() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

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
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", sdk.NewInt(10_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid pool id",
			[]string{
				"invalidpoolid",
				sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", sdk.NewInt(10_000))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		// TODO: devise another way to handle this case.
		// Note that since liquidity module is implemented in batch-style, tx message is included in a block and return error as below error message.
		// {"height":"4","txhash":"93BF23046BC55F7D763DF6FA2FF739E2C4158663B03699E2966FB91196065857","codespace":"liquidity","code":29,"data":"","raw_log":"failed to execute message; message index: 0: bad pool coin denom","logs":[],"info":"","gas_wanted":"200000","gas_used":"52016","tx":null,"timestamp":""}
		// {
		// 	"bad pool coin",
		// 	[]string{
		// 		fmt.Sprintf("%d", uint32(1)),
		// 		sdk.NewCoins(sdk.NewCoin("poolBadPoolCoinDenom", sdk.NewInt(10_000))).String(),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		// 	},
		// 	true, nil, 0,
		// },
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewWithdrawWithinBatchCmd()
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

func (s *IntegrationTestSuite) TestNewSwapWithinBatchCmd() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

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
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", liquiditytypes.DefaultSwapTypeId),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 1.0),
				fmt.Sprintf("%.3f", 0.003),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.3f", 0.003),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"invalid swap type id",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				"invalidswaptypeid",
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.3f", 0.003),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		{
			"swap type id not available",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(2)),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.2f", 0.03),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
		},
		// TODO: devise another way to handle this case.
		// Note that since liquidity module is implemented in batch-style, tx message is included in a block and return error as below error message.
		// {"height":"5","txhash":"B3F2C9AA81357CCCEE8DBD456994FD8310F307B7E4AB5D91C84DEB234083BC21","codespace":"liquidity","code":35,"data":"","raw_log":"failed to execute message; message index: 0: bad offer coin fee","logs":[],"info":"","gas_wanted":"200000","gas_used":"73743","tx":null,"timestamp":""}
		// {
		// 	"bad offer coin fee",
		// 	[]string{
		// 		fmt.Sprintf("%d", uint32(1)),
		// 		fmt.Sprintf("%d", uint32(1)),
		// 		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000))).String(),
		// 		denomY,
		// 		fmt.Sprintf("%.2f", 0.02),
		// 		fmt.Sprintf("%.2f", 0.01),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		// 	},
		// 	true, nil, 0,
		// },
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewSwapWithinBatchCmd()
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

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
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
				var poolResp liquiditytypes.QueryLiquidityPoolResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &poolResp)
				s.Require().NoError(err)
				s.Require().Len(poolResp.GetPool().ReserveCoinDenoms, 2)
				s.Require().Equal(poolResp.GetPool().ReserveAccountAddress, "cosmos1cva80e6jcxpezd3k5dl7zwy2eg30u7ld3y0a67")
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPools() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryLiquidityPools()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var poolsResp liquiditytypes.QueryLiquidityPoolsResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &poolsResp)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPoolBatch() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryLiquidityPoolBatch()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var poolBatchResp liquiditytypes.QueryLiquidityPoolBatchResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &poolBatchResp)
				s.Require().NoError(err)
				s.Require().Equal(uint64(1), poolBatchResp.GetBatch().PoolId)
				s.Require().Equal(false, poolBatchResp.GetBatch().Executed)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsg() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// create new deposit
	_, err = liquiditytestutil.MsgDepositWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000_000)), sdk.NewCoin(denomY, sdk.NewInt(10_000_000))).String(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchDepositMsg()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var depositResp liquiditytypes.QueryPoolBatchDepositMsgResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &depositResp)
				s.Require().NoError(err)
				s.Require().Equal(val.Address.String(), depositResp.GetDeposit().Msg.DepositorAddress)
				s.Require().Equal(true, depositResp.GetDeposit().Executed)
				s.Require().Equal(true, depositResp.GetDeposit().Succeeded)
				s.Require().Equal(true, depositResp.GetDeposit().ToBeDeleted)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsgs() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// create new deposit
	_, err = liquiditytestutil.MsgDepositWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000_000)), sdk.NewCoin(denomY, sdk.NewInt(10_000_000))).String(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchDepositMsgs()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var depositsResp liquiditytypes.QueryPoolBatchDepositMsgsResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &depositsResp)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsg() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// withdraw pool coin from the pool
	poolCoinDenom := "poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D"
	_, err = liquiditytestutil.MsgWithdrawWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		sdk.NewCoins(sdk.NewCoin(poolCoinDenom, sdk.NewInt(10_000))).String(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchWithdrawMsg()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var withdrawResp liquiditytypes.QueryPoolBatchWithdrawMsgResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &withdrawResp)
				s.Require().NoError(err)
				s.Require().Equal(val.Address.String(), withdrawResp.GetWithdraw().Msg.WithdrawerAddress)
				s.Require().Equal(poolCoinDenom, withdrawResp.GetWithdraw().Msg.PoolCoin.Denom)
				s.Require().Equal(true, withdrawResp.GetWithdraw().Executed)
				s.Require().Equal(true, withdrawResp.GetWithdraw().Succeeded)
				s.Require().Equal(true, withdrawResp.GetWithdraw().ToBeDeleted)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsgs() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(100_000_000)), sdk.NewCoin(denomY, sdk.NewInt(100_000_000))).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// withdraw pool coin from the pool
	_, err = liquiditytestutil.MsgWithdrawWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", sdk.NewInt(10_000))).String(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchWithdrawMsgs()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				var withdrawsResp liquiditytypes.QueryPoolBatchWithdrawMsgsResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &withdrawsResp)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsg() {
	val := s.network.Validators[0]

	// use two different tokens that are minted to the test accounts
	// when creating a new network for integration tests.
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("node0token", s.network.Config.BondDenom)
	X := sdk.NewCoin(denomX, sdk.NewInt(100_000_000))
	Y := sdk.NewCoin(denomY, sdk.NewInt(100_000_000))

	// create a liquidity pool
	_, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(X, Y).String(),
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// swap coins from the pool
	offerCoin := sdk.NewCoin(denomX, sdk.NewInt(100_000))
	currentPrice := X.Amount.ToDec().Quo(Y.Amount.ToDec())
	_, err = liquiditytestutil.MsgSwapWithinBatchExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", uint32(1)),
		fmt.Sprintf("%d", liquiditytypes.DefaultSwapTypeId),
		offerCoin.String(),
		denomY,
		currentPrice.String(),
		fmt.Sprintf("%.3f", 0.003),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"valid case",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
		},
		{
			"with invalid pool id",
			[]string{
				"invalidpoolid",
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
		{
			"with not available pool id",
			[]string{
				fmt.Sprintf("%d", uint32(2)),
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryPoolBatchSwapMsg()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				fmt.Println("out: ", out)
				fmt.Println("err: ", err)
				var swapResp liquiditytypes.QueryPoolBatchSwapMsgResponse
				err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &swapResp)
				s.Require().NoError(err)
				s.Require().Equal(val.Address.String(), swapResp.GetSwap().Msg.SwapRequesterAddress)
				s.Require().Equal(true, swapResp.GetSwap().Executed)
				s.Require().Equal(true, swapResp.GetSwap().Succeeded) // false
				s.Require().Equal(true, swapResp.GetSwap().ToBeDeleted)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsgs() {}
