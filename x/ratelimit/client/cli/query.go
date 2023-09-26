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
		GetRateLimitsByChainID(),
		GetRateLimitsByChannelID(),
		GetAllWhitelistedAddresses(),
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

			denom := args[0]
			channelID := args[1]

			if err := sdk.ValidateDenom(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

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

// GetRateLimitsByChainID return all rate limits by chain id.
func GetRateLimitsByChainID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-rate-limits [chain-id]",
		Short: "Query all rate limits by chain id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryRateLimitsByChainIDRequest{
				ChainId: args[0],
			}
			res, err := queryClient.RateLimitsByChainID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetRateLimitsByChannelID return all rate limits by channel id.
func GetRateLimitsByChannelID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-rate-limits [channel-id]",
		Short: "Query a rate limit by denom and channel id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryRateLimitsByChannelIDRequest{
				ChannelID: args[0],
			}
			res, err := queryClient.RateLimitsByChannelID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetAllWhitelistedAddresses return all whitelisted addresses.
func GetAllWhitelistedAddresses() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-whitelisted-addresses",
		Short: "Query all whitelisted addresses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllWhitelistedAddressesRequest{}
			res, err := queryClient.AllWhitelistedAddresses(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
