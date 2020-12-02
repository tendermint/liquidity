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
		GetCmdQueryLiquidityPool(),
		GetCmdQueryParams(),
	)

	return liquidityQueryCmd
}

func GetCmdQueryLiquidityPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidity-pool [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a liquidity pool
Example:
$ %s query liquidity liquidity-pool 1
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
				return err
			}

			_, err = queryClient.LiquidityPool(
				context.Background(),
				&types.QueryLiquidityPoolRequest{PoolId:poolId},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch poolId %d: %s", poolId, err)
			}

			//params := types.NewQueryLiquidityPoolParams(poolId)
			params := &types.QueryLiquidityPoolRequest{PoolId:poolId}
			res, err = queryClient.LiquidityPool(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.LiquidityPool)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
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
