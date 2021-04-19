package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquiditycli "github.com/tendermint/liquidity/x/liquidity/client/cli"
)

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

// MsgCreatePoolExec creates a transaction for creating liquidity pool.
func MsgCreatePoolExec(clientCtx client.Context, from, poolId, depositCoins string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolId,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewCreatePoolCmd(), args)
}

// MsgDepositWithinBatchExec creates a transaction to deposit new amounts to the pool.
func MsgDepositWithinBatchExec(clientCtx client.Context, from, poolId, depositCoins string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolId,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewDepositWithinBatchCmd(), args)
}

// MsgWithdrawWithinBatchExec creates a transaction to withraw pool coin amount from the pool.
func MsgWithdrawWithinBatchExec(clientCtx client.Context, from, poolId, poolCoin string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolId,
		poolCoin,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewWithdrawWithinBatchCmd(), args)
}

// MsgSwapWithinBatchExec creates a transaction to swap coins in the pool.
func MsgSwapWithinBatchExec(clientCtx client.Context, from, poolId, swapTypeId,
	offerCoin, demandCoinDenom, orderPrice, swapFeeRate string, extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolId,
		swapTypeId,
		offerCoin,
		demandCoinDenom,
		orderPrice,
		swapFeeRate,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewSwapWithinBatchCmd(), args)
}
