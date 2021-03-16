package cli

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	liquidityQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the liquidity module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	liquidityQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryLiquidityPool(),
		GetCmdQueryLiquidityPools(),
		GetCmdQueryLiquidityPoolBatch(),
		GetCmdQueryLiquidityPoolsBatches(),
		GetCmdQueryPoolBatchDepositMsgs(),
		GetCmdQueryPoolBatchDepositMsg(),
		GetCmdQueryPoolBatchWithdrawMsgs(),
		GetCmdQueryPoolBatchWithdrawMsg(),
		GetCmdQueryPoolBatchSwapMsgs(),
		GetCmdQueryPoolBatchSwapMsg(),
	)

	return liquidityQueryCmd
}

//GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current liquidity parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as liquidity parameters.

Example:
$ %s query liquidity params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryLiquidityPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a liquidity pool
Example:
$ %s query liquidity pool 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			// Query the pool
			res, err := queryClient.LiquidityPool(
				context.Background(),
				&types.QueryLiquidityPoolRequest{PoolId: poolId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch poolId %d: %s", poolId, err)
			}

			params := &types.QueryLiquidityPoolRequest{PoolId: poolId}
			res, err = queryClient.LiquidityPool(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryLiquidityPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Args:  cobra.NoArgs,
		Short: "Query for all liquidity pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all liquidity pools on a network.
Example:
$ %s query liquidity pools
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.LiquidityPools(context.Background(), &types.QueryLiquidityPoolsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryLiquidityPoolBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a liquidity pool batch of the pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a liquidity pool batch
Example:
$ %s query liquidity batch 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			// Query the pool
			result, err := queryClient.LiquidityPoolBatch(
				context.Background(),
				&types.QueryLiquidityPoolBatchRequest{PoolId: poolId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch poolId %d: %s", poolId, err)
			}

			params := &types.QueryLiquidityPoolBatchRequest{PoolId: poolId}
			result, err = queryClient.LiquidityPoolBatch(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryLiquidityPoolsBatches() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batches",
		Args:  cobra.NoArgs,
		Short: "Query for all liquidity pools batches",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all liquidity pools batches on a network.
Example:
$ %s query liquidity batches
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.LiquidityPoolsBatches(context.Background(), &types.QueryLiquidityPoolsBatchesRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchDepositMsgs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all deposit messages on the batch of the liquidity pool specified pool-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all deposit messages on the batch of the liquidity pool specified pool-id

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ %s query liquidity deposits 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			result, err := queryClient.PoolBatchDepositMsgs(context.Background(), &types.QueryPoolBatchDepositMsgsRequest{
				PoolId: poolId, Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchDepositMsg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [pool-id] [msg-index]",
		Args:  cobra.ExactArgs(2),
		Short: "Query for the deposit message on the batch of the liquidity pool specified pool-id and msg-index",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for the deposit message on the batch of the liquidity pool specified pool-id and msg-index

if the batch message are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ %s query liquidity deposit 1 20
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			msgIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("msg-index %s not a valid uint, please input a valid msg-index", args[1])
			}

			result, err := queryClient.PoolBatchDepositMsg(context.Background(), &types.QueryPoolBatchDepositMsgRequest{
				PoolId: poolId, MsgIndex: msgIndex})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchWithdrawMsgs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraws [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all withdraw messages on the batch of the liquidity pool specified pool-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all withdraws messages on the batch of the liquidity pool specified pool-id

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ %s query liquidity withdraws 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			result, err := queryClient.PoolBatchWithdrawMsgs(context.Background(), &types.QueryPoolBatchWithdrawMsgsRequest{
				PoolId: poolId, Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchWithdrawMsg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [pool-id] [msg-index]",
		Args:  cobra.ExactArgs(2),
		Short: "Query for the withdraw message on the batch of the liquidity pool specified pool-id and msg-index",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for the withdraw message on the batch of the liquidity pool specified pool-id and msg-index

if the batch message are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ %s query liquidity withdraw 1 20
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			msgIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("msg-index %s not a valid uint, please input a valid msg-index", args[1])
			}

			result, err := queryClient.PoolBatchWithdrawMsg(context.Background(), &types.QueryPoolBatchWithdrawMsgRequest{
				PoolId: poolId, MsgIndex: msgIndex})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchSwapMsgs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swaps [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all swap messages on the batch of the liquidity pool specified pool-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all swap messages on the batch of the liquidity pool specified pool-id

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ %s query liquidity swaps 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			result, err := queryClient.PoolBatchSwapMsgs(context.Background(), &types.QueryPoolBatchSwapMsgsRequest{
				PoolId: poolId, Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchSwapMsg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [msg-index]",
		Args:  cobra.ExactArgs(2),
		Short: "Query for the swap message on the batch of the liquidity pool specified pool-id and msg-index",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for the swap message on the batch of the liquidity pool specified pool-id and msg-index

if the batch message are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ %s query liquidity swap 1 20
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			msgIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("msg-index %s not a valid uint, please input a valid msg-index", args[1])
			}

			result, err := queryClient.PoolBatchSwapMsg(context.Background(), &types.QueryPoolBatchSwapMsgRequest{
				PoolId: poolId, MsgIndex: msgIndex})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
