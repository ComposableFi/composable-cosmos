package v4

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/notional-labs/composable/v6/app/keepers"
	"github.com/notional-labs/composable/v6/app/upgrades"
	tfmdtypes "github.com/notional-labs/composable/v6/x/transfermiddleware/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	_ codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Add params for transfer middleware
		transmiddlewareParams := tfmdtypes.DefaultParams()
		keepers.TransferMiddlewareKeeper.SetParams(ctx, transmiddlewareParams)

		// Add params for wasmd
		var wasmdParams wasmtypes.Params
		wasmdParams.CodeUploadAccess = wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeNobody}
		wasmdParams.InstantiateDefaultPermission = wasmtypes.AccessTypeNobody
		keepers.WasmKeeper.SetParams(ctx, wasmdParams) //nolint:errcheck // TODO: handle error

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
