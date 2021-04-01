package cli_test

import (
	"context"
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

			s.network.WaitForNextBlock()

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

			nextBlockErr := s.network.WaitForNextBlock()
			s.Require().NoError(nextBlockErr)

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
		// {
		// 	"valid transaction",
		// 	[]string{
		// 		fmt.Sprintf("%d", uint32(1)),
		// 		sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", sdk.NewInt(10_000))).String(),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		// 	},
		// 	false, &sdk.TxResponse{}, 0,
		// },
		// {
		// 	"invalid pool id",
		// 	[]string{
		// 		"invalidpoolid",
		// 		sdk.NewCoins(sdk.NewCoin("poolC33A77E752C183913636A37FE1388ACA22FE7BED792BEB2E72EF2DA857703D8D", sdk.NewInt(10_000))).String(),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		// 	},
		// 	true, nil, 0,
		// },
		// TODO: needs debugging. This occurs when it gets panic?
		// panic: test timed out after 30s and goroutine runs consencutively
		{
			"bad pool coin",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin("poolBadPoolCoinDenom", sdk.NewInt(10_000))).String(),
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
			cmd := cli.NewWithdrawWithinBatchCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			nextBlockErr := s.network.WaitForNextBlock()
			s.Require().NoError(nextBlockErr)

			if tc.expectErr {
				fmt.Println("out: ", out)
				fmt.Println("err: ", err)
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
				fmt.Sprintf("%d", uint32(1)),
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
			"bad offer coin fee",
			[]string{
				fmt.Sprintf("%d", uint32(1)),
				fmt.Sprintf("%d", uint32(1)),
				sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000))).String(),
				denomY,
				fmt.Sprintf("%.2f", 0.02),
				fmt.Sprintf("%.2f", 0.01),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, nil, 0,
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
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewSwapWithinBatchCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			_, nextBlockErr := s.network.WaitForHeight(1)
			s.Require().Error(nextBlockErr)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				fmt.Println("out: ", out)

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
	denomX, denomY := liquiditytypes.AlphabeticalDenomPair("stoken", s.network.Config.BondDenom)

	// create a liquidity pool
	abcd, err := liquiditytestutil.MsgCreatePoolExec(
		val.ClientCtx,
		val.Address.String(),
		fmt.Sprintf("%d", liquiditytypes.DefaultPoolTypeId),
		sdk.NewCoins(sdk.NewCoin(denomX, sdk.NewInt(10_000)), sdk.NewCoin(denomY, sdk.NewInt(10_000))).String(),
	)
	s.Require().NoError(err)

	fmt.Println("abcd: ", abcd)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// Query the pool
	queryClient := liquiditytypes.NewQueryClient(val.ClientCtx)

	res, err := queryClient.LiquidityPool(
		context.Background(),
		&liquiditytypes.QueryLiquidityPoolRequest{
			PoolId: uint64(1),
		},
	)

	fmt.Println("res: ", res)
	fmt.Println("err: ", err)

	// params := &liquiditytypes.QueryLiquidityPoolRequest{PoolId: uint32(1)}
	// res, err = queryClient.LiquidityPool(context.Background(), params)
	// if err != nil {
	// 	return err
	// }

	// testCases := []struct {
	// 	name      string
	// 	args      []string
	// 	expectErr bool
	// }{
	// 	{
	// 		"valid case",
	// 		[]string{
	// 			fmt.Sprintf("%d", uint32(1)),
	// 			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	// 		},
	// 		false,
	// 	},
	// 	{
	// 		"with invalid pool id",
	// 		[]string{
	// 			"invalidpoolid",
	// 			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	// 		},
	// 		true,
	// 	},
	// 	{
	// 		"with not available pool id",
	// 		[]string{
	// 			fmt.Sprintf("%d", uint32(2)),
	// 			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	// 		},
	// 		true,
	// 	},
	// }

	// for _, tc := range testCases {
	// 	tc := tc

	// 	s.Run(tc.name, func() {
	// 		cmd := cli.GetCmdQueryLiquidityPool()
	// 		clientCtx := val.ClientCtx

	// 		out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
	// 		if tc.expectErr {
	// 			fmt.Println("err: ", err)
	// 			s.Require().Error(err)
	// 		} else {
	// 			fmt.Println("tc.name: ", tc.name)
	// 			fmt.Println("out.String: ", strings.TrimSpace(out.String()))
	// 			s.Require().NoError(err)
	// 		}
	// 	})
	// }
}

func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPools()        {}
func (s *IntegrationTestSuite) TestGetCmdQueryLiquidityPoolBatch()    {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsg()   {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchDepositMsgs()  {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsg()  {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchWithdrawMsgs() {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsg()      {}
func (s *IntegrationTestSuite) TestGetCmdQueryPoolBatchSwapMsgs()     {}
