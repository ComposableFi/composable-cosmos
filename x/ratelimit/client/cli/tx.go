package cli

import (
	"fmt"

	"github.com/notional-labs/composable/v6/x/ratelimit/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the tx commands for router
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		Short:                      fmt.Sprintf("Tx commands for the %s module", types.ModuleName),
	}

	txCmd.AddCommand()

	return txCmd
}
