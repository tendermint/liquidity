package cli

// DONTCOVER

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// GetTxCmd returns a root CLI command handler for all x/liquidity transaction commands.
func GetTxCmd() *cobra.Command {
	liquidityTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Liquidity transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	liquidityTxCmd.AddCommand(
		NewCreateLiquidityPoolCmd(),
	)

	return liquidityTxCmd
}

func NewCreateLiquidityPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-liquidity-pool [pool-type-index] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Create Liquidity pool with the specified pool-type, deposit coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create Liquidity pool with the specified pool-type-index, deposit coins for reserve

Example:
$ %s tx liquidity create-liquidity-pool 1 100000000acoin,100000000bcoin --from mykey

Currently, only the default pool-type-index 1 is available
the number of deposit coins must be two in the pool-type-index 1

{"pool_type_index":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			poolCreator := clientCtx.GetFromAddress()

			// Get pool type index
			poolTypeIndex, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("pool-type-index %s not a valid uint, please input a valid pool-type-index", args[0])
			}

			// Get deposit coins
			depositCoins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			if depositCoins.Len() != 2 {
				return fmt.Errorf("the number of deposit coins must be two in the pool-type-index 1")
			}

			reserveCoinDenoms := []string{depositCoins[0].Denom, depositCoins[1].Denom}
			msg := types.NewMsgCreateLiquidityPool(poolCreator, uint32(poolTypeIndex), reserveCoinDenoms, depositCoins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}


			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}