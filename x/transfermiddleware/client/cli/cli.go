package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for router
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "transfermiddleware",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		GetCmdParaTokenInfo(),
	)

	return queryCmd
}

// GetCmdParaTokenInfo returns the command handler for transfer-middleware para-token-info querying.
func GetCmdParaTokenInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "para-token-info",
		Short:   "Query the current transfer middleware para-token-info based on denom",
		Long:    "Query the current transfer middleware para-token-info based on denom",
		Args:    cobra.ExactValidArgs(1),
		Example: fmt.Sprintf("%s query transfermiddleware para-token-info atom", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ParaTokenInfo(cmd.Context(), &types.QueryParaTokenInfoRequest{
				NativeDenom: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewTxCmd returns the transaction commands for router
func NewTxCmd() *cobra.Command {
	return nil
}
