package v4

import (
	store "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"

	"github.com/notional-labs/composable/v6/app/upgrades"
	ibchookstypes "github.com/notional-labs/composable/v6/x/ibc-hooks/types"
	ratelimitmoduletypes "github.com/notional-labs/composable/v6/x/ratelimit/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v4"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{wasmtypes.StoreKey, ibchookstypes.StoreKey, ratelimitmoduletypes.StoreKey, icahosttypes.StoreKey},
		Deleted: []string{},
	},
}
