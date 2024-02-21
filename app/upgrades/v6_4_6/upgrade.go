package v6_4_6

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/notional-labs/composable/v6/app/keepers"
	"github.com/notional-labs/composable/v6/app/upgrades"
	bech32authmigration "github.com/notional-labs/composable/v6/bech32-migration/auth"
	bech32govmigration "github.com/notional-labs/composable/v6/bech32-migration/gov"
	bech32slashingmigration "github.com/notional-labs/composable/v6/bech32-migration/slashing"
	bech32stakingmigration "github.com/notional-labs/composable/v6/bech32-migration/staking"
	"github.com/notional-labs/composable/v6/bech32-migration/utils"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	codec codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		keys := keepers.GetKVStoreKey()
		// Migration prefix
		ctx.Logger().Info("First step: Migrate addresses stored in bech32 form to use new prefix")
		sdk.GetConfig().SetBech32PrefixForAccount(utils.NewBech32PrefixAccAddr, utils.NewBech32PrefixAccPub)
		sdk.GetConfig().SetBech32PrefixForValidator(utils.NewBech32PrefixValAddr, utils.NewBech32PrefixValPub)
		sdk.GetConfig().SetBech32PrefixForConsensusNode(utils.NewBech32PrefixConsAddr, utils.NewBech32PrefixConsPub)

		bech32stakingmigration.MigrateAddressBech32(ctx, keys[stakingtypes.StoreKey], codec)
		bech32slashingmigration.MigrateAddressBech32(ctx, keys[slashingtypes.StoreKey], codec)
		bech32govmigration.MigrateAddressBech32(ctx, keys[govtypes.StoreKey], codec)
		bech32authmigration.MigrateAddressBech32(ctx, keys[authtypes.StoreKey], codec)
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
