package v6

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	consumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
	"github.com/notional-labs/centauri/v6/app/keepers"
	"github.com/notional-labs/centauri/v6/app/upgrades"

	"github.com/spf13/cast"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	appOpts servertypes.AppOptions,
	cdc codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting upgrade v6...")

		consumerKeeper := keepers.ConsumerKeeper

		nodeHome := cast.ToString(appOpts.Get(flags.FlagHome))
		consumerUpgradeGenFile := nodeHome + "/config/ccv.json"
		appState, _, err := genutiltypes.GenesisStateFromGenFile(consumerUpgradeGenFile)
		if err != nil {
			panic("Unable to read consumer genesis")
		}

		var consumerGenesis = consumertypes.GenesisState{}
		cdc.MustUnmarshalJSON(appState[consumertypes.ModuleName], &consumerGenesis)

		consumerGenesis.PreCCV = true
		consumerGenesis.Params.SoftOptOutThreshold = "0.05"
		consumerKeeper.InitGenesis(ctx, &consumerGenesis)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
