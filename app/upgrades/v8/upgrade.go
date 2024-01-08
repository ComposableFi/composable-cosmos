package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/notional-labs/composable/v6/app/keepers"
	"github.com/notional-labs/composable/v6/app/upgrades"
	customstmiddleware "github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	_ codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Add params for custom middleware
		custommiddlewareparams := customstmiddleware.Params{BlocksPerEpoch: 100, AllowUnbondAfterEpochProgressBlockNumber: 0}
		keepers.StakingMiddlewareKeeper.SetParams(ctx, custommiddlewareparams)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
