package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/notional-labs/centauri/v4/app"
	cmd "github.com/notional-labs/centauri/v4/cmd/centaurid/cmd"
	cmdcfg "github.com/notional-labs/centauri/v4/cmd/centaurid/config"
)

func main() {
	cmdcfg.SetupConfig()
	cmdcfg.RegisterDenoms()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
