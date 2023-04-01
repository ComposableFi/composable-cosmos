package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/notional-labs/composable-testnet/v2/app"
	cmd "github.com/notional-labs/composable-testnet/v2/cmd/polytoped/cmd"
	cmdcfg "github.com/notional-labs/composable-testnet/v2/cmd/polytoped/config"
)

func main() {
	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
