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
		NewCreatePoolCmd(),
		NewDepositWithinBatchCmd(),
		NewWithdrawWithinBatchCmd(),
		NewSwapWithinBatchCmd(),
	)

	return liquidityTxCmd
}

func NewCreatePoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-type-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Create Liquidity pool with the specified pool-type, deposit-coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create Liquidity pool with the specified pool-type-id, deposit-coins for reserve

Example:
$ %s tx liquidity create-pool 1 100000000stake,100000000token --from mykey

Currently, only the default pool-type-id 1 is available on this version
the number of deposit coins must be two in the pool-type-id 1

{"id":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}
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
			poolTypeId, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("pool-type-id %s not a valid uint, please input a valid pool-type-id", args[0])
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

			if poolTypeId != 1 {
				return types.ErrPoolTypeNotExists
			}

			if depositCoins.Len() != 2 {
				return fmt.Errorf("the number of deposit coins must be two in the pool-type-id 1")
			}

			msg := types.NewMsgCreatePool(poolCreator, uint32(poolTypeId), depositCoins)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit-coins
func NewDepositWithinBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [pool-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit-coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit-coins for reserve

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ %s tx liquidity deposit 1 100000000stake,100000000token --from mykey

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
				return fmt.Errorf("the number of deposit coins must be two in the pool-type-id 1")
			}

			msg := types.NewMsgDepositWithinBatch(depositor, poolId, depositCoins)
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
func NewWithdrawWithinBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ %s tx liquidity withdraw 1 1000pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07 --from mykey

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

			msg := types.NewMsgWithdrawWithinBatch(withdrawer, poolId, poolCoin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Swap offer submit to the batch to the Liquidity pool with the specified pool-id with offer-coin, order-price, etc
func NewSwapWithinBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [swap-type-id] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate]",
		Args:  cobra.ExactArgs(6),
		Short: "Swap offer submit to the batch to the Liquidity pool with the specified pool-id with offer-coin, order-price, etc",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Swap offer to the Liquidity pool with the specified pool-id, swap-type-id demand-coin-denom 
with the coin and the price you're offering and current swap-fee-rate

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ %s tx liquidity swap 2 1 100000000stake token 0.9 0.003 --from mykey

You should request the same each field as the pool.

Must have sufficient balance half the of the swapFee Rate of the offer coin to reserve offer coin fee.

For explicit calculations, you must enter the params.swap_fee_rate value of the current parameter state.

Currently, only the default pool-type-id, swap-type-id 1 is available on this version
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
			swapTypeId, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return fmt.Errorf("swap-type-id %s not a valid uint, please input a valid swap-type-id", args[2])
			}

			if swapTypeId != 1 {
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
				return fmt.Errorf("pool-type-id %s not a valid uint, please input a valid pool-type-id", args[1])
			}

			orderPrice, err := sdk.NewDecFromStr(args[4])
			if err != nil {
				return err
			}

			swapFeeRate, err := sdk.NewDecFromStr(args[5])
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapWithinBatch(swapRequester, poolId, uint32(swapTypeId), offerCoin, args[3], orderPrice, swapFeeRate)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
