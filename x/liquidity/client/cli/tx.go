package cli

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

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
		NewDepositToLiquidityPoolCmd(),
		NewWithdrawFromLiquidityPoolCmd(),
		NewSwapCmd(),
	)

	return liquidityTxCmd
}

func NewCreateLiquidityPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-type-index] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Create Liquidity pool with the specified pool-type, deposit coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create Liquidity pool with the specified pool-type-index, deposit coins for reserve

Example:
$ %s tx liquidity create-pool 1 100000000acoin,100000000bcoin --from mykey

Currently, only the default pool-type-index 1 is available on this version
the number of deposit coins must be two in the pool-type-index 1

{"pool_type_index":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
			depositCoins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			err = depositCoins.Validate()
			if err != nil {
				return err
			}

			if poolTypeIndex != 1 {
				return types.ErrPoolTypeNotExists
			}

			if depositCoins.Len() != 2 {
				return fmt.Errorf("the number of deposit coins must be two in the pool-type-index 1")
			}

			msg := types.NewMsgCreateLiquidityPool(poolCreator, uint32(poolTypeIndex), depositCoins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins
func NewDepositToLiquidityPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [pool-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins for reserve

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ %s tx liquidity deposit 1 100000000acoin,100000000bcoin --from mykey

You should deposit the same coin as the reserve coin.
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			depositor := clientCtx.GetFromAddress()

			// Get pool type index
			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			// Get deposit coins
			depositCoins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			err = depositCoins.Validate()
			if err != nil {
				return err
			}

			if depositCoins.Len() != 2 {
				return fmt.Errorf("the number of deposit coins must be two in the pool-type-index 1")
			}

			msg := types.NewMsgDepositToLiquidityPool(depositor, poolId, depositCoins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool
func NewWithdrawFromLiquidityPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ %s tx liquidity withdraw 1 1000cosmos1d9w9j3rq5aunkrkdm86paduz4attl78thlj07f --from mykey

You should request the matched pool-coin as the pool.
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			withdrawer := clientCtx.GetFromAddress()

			// Get pool type index
			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			// Get pool coin of the target pool
			poolCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			err = poolCoin.Validate()
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFromLiquidityPool(withdrawer, poolId, poolCoin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Swap offer to the Liquidity pool with the specified the pool info with offer-coin, order-price
func NewSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [swap-type] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate]",
		Args:  cobra.ExactArgs(6),
		Short: "Swap offer to the Liquidity pool with the specified the pool info with offer-coin, order-price, swap-fee-rate",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Swap offer to the Liquidity pool with the specified pool-id, swap-type demand-coin-denom 
with the coin and the price you're offering

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ %s tx liquidity swap 2 1 100000000acoin bcoin 1.15 0.003 --from mykey

You should request the same each field as the pool.

Must have sufficient balance half the of the swapFee Rate of the offer coin to reserve offer coin fee.

For explicit calculations, you must enter the params.swap_fee_rate value of the current parameter state.

Currently, only the default pool-type, swap-type 1 is available on this version
The detailed swap algorithm can be found here.
https://github.com/tendermint/liquidity
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			swapRequester := clientCtx.GetFromAddress()

			// Get pool id
			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, please input a valid pool-id", args[0])
			}

			// Get swap type
			swapType, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return fmt.Errorf("swap-type %s not a valid uint, please input a valid swap-type", args[2])
			}

			if swapType != 1 {
				return types.ErrSwapTypeNotExists
			}

			// Get offer coin
			offerCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			err = offerCoin.Validate()
			if err != nil {
				return err
			}

			err = sdk.ValidateDenom(args[3])
			if err != nil {
				return err
			}

			if err != nil {
				return fmt.Errorf("pool-type-index %s not a valid uint, please input a valid pool-type-index", args[1])
			}

			orderPrice, err := sdk.NewDecFromStr(args[4])
			if err != nil {
				return err
			}

			swapFeeRate, err := sdk.NewDecFromStr(args[5])
			if err != nil {
				return err
			}

			msg := types.NewMsgSwap(swapRequester, poolId, uint32(swapType), offerCoin, args[3], orderPrice, swapFeeRate)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
