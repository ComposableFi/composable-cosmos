package v6

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/notional-labs/composable/v6/app/keepers"
	"github.com/notional-labs/composable/v6/app/upgrades"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	cdc codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		allowed := []string{
			keepers.AccountKeeper.GetModuleAddress(wasmtypes.ModuleName).String(),
			keepers.AccountKeeper.GetModuleAddress(govtypes.ModuleName).String(),
			"centauri1u2sr0p2j75fuezu92nfxg5wm46gu22ywfgul6k", // "dzmitry lahoda CVM/MANTIS dev
		}
		wasmdParams := keepers.WasmKeeper.GetParams(ctx)
		wasmdParams.CodeUploadAccess.Permission = wasmtypes.AccessTypeAnyOfAddresses
		wasmdParams.CodeUploadAccess.Addresses = append(wasmdParams.CodeUploadAccess.Addresses, allowed...)

		err := keepers.WasmKeeper.SetParams(ctx, wasmdParams)
		if err != nil {
			return nil, err
		}

		migrator := pfmkeeper.NewMigrator(keepers.RouterKeeper, keepers.GetSubspace(pfmtypes.ModuleName))
		err = migrator.Migrate1to2(ctx)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
