package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the tx commands for staking middleware module.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Exp transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()

	return txCmd
}
