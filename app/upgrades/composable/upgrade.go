package composable

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/notional-labs/composable/v5/app/keepers"
	"github.com/notional-labs/composable/v5/app/upgrades"

	bech32authmigration "github.com/notional-labs/composable/v5/bech32-migration/auth"
	bech32govmigration "github.com/notional-labs/composable/v5/bech32-migration/gov"
	bech32slashingmigration "github.com/notional-labs/composable/v5/bech32-migration/slashing"
	bech32stakingmigration "github.com/notional-labs/composable/v5/bech32-migration/staking"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/notional-labs/composable/v5/bech32-migration/utils"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	cdc codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Migration prefix
		ctx.Logger().Info("First step: Migrate addresses stored in bech32 form to use new prefix")
		keys := keepers.GetKVStoreKey()
		bech32stakingmigration.MigrateAddressBech32(ctx, keys[stakingtypes.StoreKey], cdc)
		bech32slashingmigration.MigrateAddressBech32(ctx, keys[slashingtypes.StoreKey], cdc)
		bech32govmigration.MigrateAddressBech32(ctx, keys[govtypes.StoreKey], cdc)
		bech32authmigration.MigrateAddressBech32(ctx, keys[authtypes.StoreKey], cdc)

		// allowed relayer address
		tfmdk := keepers.TransferMiddlewareKeeper
		tfmdk.IterateAllowRlyAddress(ctx, func(rlyAddress string) bool {
			// Delete old address
			tfmdk.DeleteAllowRlyAddress(ctx, rlyAddress)

			// add new address
			newRlyAddress := utils.ConvertAccAddr(rlyAddress)
			tfmdk.SetAllowRlyAddress(ctx, newRlyAddress)
			return false
		})

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
