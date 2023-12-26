package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/notional-labs/composable/v6/app"
	cmd "github.com/notional-labs/composable/v6/cmd/centaurid/cmd"
	cmdcfg "github.com/notional-labs/composable/v6/cmd/centaurid/config"
)

func main() {
	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "CENTAURID", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
