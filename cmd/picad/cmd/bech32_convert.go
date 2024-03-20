package cmd

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/spf13/cobra"
)

var flagBech32Prefix = "prefix"

// AddBech32ConvertCommand returns bech32-convert cobra Command.
func AddBech32ConvertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bech32-convert [address]",
		Short: "Convert any bech32 string to the pica prefix",
		Long: `Convert any bech32 string to the pica prefix

Example:
	picad debug bech32-convert akash1a6zlyvpnksx8wr6wz8wemur2xe8zyh0ytz6d88

	picad debug bech32-convert stride1673f0t8p893rqyqe420mgwwz92ac4qv6synvx2 --prefix osmo
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32prefix, err := cmd.Flags().GetString(flagBech32Prefix)
			if err != nil {
				return err
			}

			_, bz, err := bech32.DecodeAndConvert(args[0])
			if err != nil {
				return err
			}

			bech32Addr, err := bech32.ConvertAndEncode(bech32prefix, bz)
			if err != nil {
				panic(err)
			}

			cmd.Println(bech32Addr)

			return nil
		},
	}

	cmd.Flags().StringP(flagBech32Prefix, "p", "picasso", "Bech32 Prefix to encode to")

	return cmd
}

// addDebugCommands injects custom debug commands into another command as children.
func addDebugCommands(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(AddBech32ConvertCommand())
	return cmd
}
