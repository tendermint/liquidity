package cli

// DONTCOVER

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
		GetCmdQueryLiquidityPoolsBatch(),
		GetCmdQueryPoolBatchDeposit(),
		GetCmdQueryPoolBatchWithdraw(),
		GetCmdQueryPoolBatchSwap(),
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.Params)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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

			return clientCtx.PrintOutput(res)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			return clientCtx.PrintOutput(result)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			// Query the pool
			res, err := queryClient.LiquidityPoolBatch(
				context.Background(),
				&types.QueryLiquidityPoolBatchRequest{PoolId: poolId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch poolId %d: %s", poolId, err)
			}

			params := &types.QueryLiquidityPoolBatchRequest{PoolId: poolId}
			res, err = queryClient.LiquidityPoolBatch(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryLiquidityPoolsBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batches",
		Args:  cobra.NoArgs,
		Short: "Query for all liquidity pools batch",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all liquidity pools batch on a network.
Example:
$ %s query liquidity batches
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.LiquidityPoolsBatch(context.Background(), &types.QueryLiquidityPoolsBatchRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Args:  cobra.NoArgs,
		Short: "Query for all deposit messages on the batch of the liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all deposit messages on the batch of the liquidity pool

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.

Example:
$ %s query liquidity deposits
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.PoolBatchDepositMsgs(context.Background(), &types.QueryPoolBatchDepositMsgsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraws",
		Args:  cobra.NoArgs,
		Short: "Query for all withdraw messages on the batch of the liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all withdraws messages on the batch of the liquidity pool

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.

Example:
$ %s query liquidity withdraws
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.PoolBatchWithdrawMsgs(context.Background(), &types.QueryPoolBatchWithdrawMsgsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryPoolBatchSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swaps",
		Args:  cobra.NoArgs,
		Short: "Query for all swap messages on the batch of the liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all swap messages on the batch of the liquidity pool

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.

Example:
$ %s query liquidity swaps
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.PoolBatchSwapMsgs(context.Background(), &types.QueryPoolBatchSwapMsgsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(result)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
