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

// Create new liquidity pool with the specified pool type and deposit coins.
func NewCreatePoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-type-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Create liquidity pool and deposit coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create liquidity pool and deposit coins.

Example:
$ %s tx liquidity create-pool 1 1000000000uatom,50000000000uusd --from mykey

This example creates a liquidity pool of pool-type-id 1 and deposits 100000000stake and 100000000token.
New liquidity pools can be created only for coin combinations that do not exist in the network.
The only supported pool-type-id is 1. pool-type-id 1 requires two different coins.

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
				return fmt.Errorf("pool-type-id %s not a valid uint, input a valid pool-type-id", args[0])
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

// Deposit coins to the specified liquidity pool.
func NewDepositWithinBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [pool-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit coins to the specified liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit coins to the specified liquidity pool.

This swap request may not be processed immediately since it will be accumulated in the batch of the liquidity pool.
This will be processed with other requests at once in every end of batch.

Example:
$ %s tx liquidity deposit 1 100000000uatom,5000000000uusd --from mykey

In this example, user requests to deposit 100000000uatom and 5000000000uusd to the specified liquidity pool.
User must deposit the same coin denoms as the reserve coins.
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

// Withdraw pool coin from the specified liquidity pool.
func NewWithdrawWithinBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw pool coin from the specified liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw pool coin from the specified liquidity pool.

This swap request may not be processed immediately since it will be accumulated in the batch of the liquidity pool.
This will be processed with other requests at once in every end of batch.

Example:
$ %s tx liquidity withdraw 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from mykey

In this example, user requests to withdraw 10000 pool coin from the specified liquidity pool.
User must request the appropriate pool coin from the specified pool.
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

// Swap offer coin with demand coin from the specified liquidity pool with the given order price.
func NewSwapWithinBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [swap-type-id] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate]",
		Args:  cobra.ExactArgs(6),
		Short: "Swap offer coin with demand coin from the specified liquidity pool with the given order price",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Swap offer coin with demand coin from the specified liquidity pool with the given order price.

This swap request may not be processed immediately since it will be accumulated in the batch of the liquidity pool.
This will be processed with other requests at once in every end of batch.
Note that the order of swap requests is ignored since the universal swap price is calculated within every batch to prevent front running.

The requested swap is executed with a swap price calculated from given swap price function of the pool, the current other swap requests and the current liquidity pool coin reserve status.
Swap orders are executed only when execution swap price is equal or better than submitted order price of the swap order.

Example:
$ %s liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from mykey

In this example, we assume there exists a liquidity pool with 1000000000uatom and 50000000000uusd.
User requests to swap 50000000uusd for at least 950000uatom with the order price of 0.019 and swap fee rate of 0.003.
User must have sufficient balance half of the swap-fee-rate of the offer coin to reserve offer coin fee.

The order price is the exchange ratio of X/Y where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically.
Increasing order price means to decrease the possibility for your request to be processed and end up buying uatom at cheaper price than the pool price.

For explicit calculations, you must enter the swap-fee-rate value of the current parameter state.
In this version, swap-type-id 1 is only available. The detailed swap algorithm can be found at https://github.com/tendermint/liquidity`,
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
