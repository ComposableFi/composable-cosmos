package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/notional-labs/composable/v6/x/ibctransfermiddleware/types"
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
		AddIBCFeeConfig(),
		RemoveIBCFeeConfig(),
	)

	return txCmd
}

func AddIBCFeeConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-config [addr ]",
		Short:   "add ibc fee config",
		Args:    cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
		Example: fmt.Sprintf("%s tx ibctransfermiddleware add-config [channel] [feeAddress] [minTimeoutTimestamp]", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			channel := args[0]
			feeAddress := args[1]
			minTimeoutTimestamp := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddress := clientCtx.GetFromAddress().String()

			//string to int64 conversion
			minTimeoutTimestampInt, err := strconv.ParseInt(minTimeoutTimestamp, 10, 64)

			if err != nil {
				return err
			}

			msg := types.NewMsgAddIBCFeeConfig(
				fromAddress,
				channel,
				feeAddress,
				minTimeoutTimestampInt,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func RemoveIBCFeeConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove-config [addr ]",
		Short:   "add ibc fee config",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Example: fmt.Sprintf("%s tx ibctransfermiddleware remove-config [channel]", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			channel := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddress := clientCtx.GetFromAddress().String()

			msg := types.NewMsgRemoveIBCFeeConfig(
				fromAddress,
				channel,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
