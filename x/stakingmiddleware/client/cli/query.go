package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

// GetQueryCmd returns the cli query commands for the staking middleware module.
func GetQueryCmd() *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the minting module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	mintingQueryCmd.AddCommand(
		GetCmdQueryParams(),
	)

	return mintingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPowerRequest{}
			res, err := queryClient.Power(cmd.Context(), params)
			if err != nil {
				return err
			}
			_ = res

			// return nil
			return clientCtx.PrintString(fmt.Sprintf("%s\n", "hello world"))
			// return clientCtx.PrintProto(&res.Power)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
