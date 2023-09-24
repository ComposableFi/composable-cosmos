package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group ratelimit queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryAllRateLimits(),
		GetCmdQueryRateLimit(),
		// TODO: add more commands
	)
	return cmd
}

// GetCmdQueryAllRateLimits return all available rate limits.
func GetCmdQueryAllRateLimits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-rate-limits",
		Short: "Query all rate limits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllRateLimitsRequest{}
			res, err := queryClient.AllRateLimits(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRateLimit return a rate limit by denom and channel id.
func GetCmdQueryRateLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate-limit [denom] [channel-id]",
		Short: "Query a rate limit by denom and channel id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]
			if err := sdk.ValidateDenom(denom); err != nil {
				return err
			}
			channelID := args[1]

			req := &types.QueryRateLimitRequest{
				Denom:     denom,
				ChannelID: channelID,
			}
			res, err := queryClient.RateLimit(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
