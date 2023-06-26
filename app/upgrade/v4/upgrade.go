package v4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	tfmwKeeper "github.com/notional-labs/centauri/v3/x/transfermiddleware/keeper"
)

var listAllowedRelayAddress = []string{
	"centauri1eqv3xl0vk0md74qukfghfff4z3axsp29rr9c85",
	"centauri1av6x9sll0yx4anske424jtgxejnrgqv6j6tjjt",
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	tfmwKeeper tfmwKeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Add relayer address to store
		for _, allowedRelayAddress := range listAllowedRelayAddress {
			tfmwKeeper.SetAllowRlyAddress(ctx, allowedRelayAddress)
		}
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
