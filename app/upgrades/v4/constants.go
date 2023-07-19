package v4

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/notional-labs/centauri/v4/app/upgrades"
	ibchookstypes "github.com/notional-labs/centauri/v4/x/ibc-hooks/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the Centauri upgrade.
	UpgradeName = "v4"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{wasmtypes.StoreKey, ibchookstypes.StoreKey},
		Deleted: []string{},
	},
}
