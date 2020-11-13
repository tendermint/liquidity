package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// TODO: add commands after rebase sdk stable 0.40.0

// GetTxCmd returns a root CLI command handler for all x/liquidity transaction commands.
func GetTxCmd() *cobra.Command {
	liquidityTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Liquidity transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	liquidityTxCmd.AddCommand()

	return liquidityTxCmd
}
