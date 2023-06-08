package centauri

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	bech32authmigration "github.com/notional-labs/centauri/v2/bech32-migration/auth"
	bech32govmigration "github.com/notional-labs/centauri/v2/bech32-migration/gov"
	bech32slashingmigration "github.com/notional-labs/centauri/v2/bech32-migration/slashing"
	bech32stakingmigration "github.com/notional-labs/centauri/v2/bech32-migration/staking"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keys map[string]*storetypes.KVStoreKey, codec codec.Codec,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// set old prefix
		ctx.Logger().Info("Second step: Migrate addresses stored in bech32 form to use new prefix")
		bech32stakingmigration.MigrateAddressBech32(ctx, keys[stakingtypes.StoreKey], codec)
		bech32slashingmigration.MigrateAddressBech32(ctx, keys[slashingtypes.StoreKey], codec)
		bech32govmigration.MigrateAddressBech32(ctx, keys[govtypes.StoreKey], codec)
		bech32authmigration.MigrateAddressBech32(ctx, keys[authtypes.StoreKey], codec)

		// set new prefix
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
