package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

// GetTxCmd returns the tx commands for router
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "ratelimit",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		Short:                      fmt.Sprintf("Tx commands for the %s module", types.ModuleName),
	}

	txCmd.AddCommand()

	return txCmd
}
