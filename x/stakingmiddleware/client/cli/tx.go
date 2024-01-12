package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	txCmd.AddCommand(
		GetCmdAddRevenueFundsToStaking(),
	)

	return txCmd
}

func GetCmdAddRevenueFundsToStaking() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-revenue [amount]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgAddRevenueFundsToStakingParams{
				FromAddress: clientCtx.GetFromAddress().String(),
				Amount:      coins,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
