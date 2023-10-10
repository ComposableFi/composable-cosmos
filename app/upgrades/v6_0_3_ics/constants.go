package v6

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/notional-labs/centauri/v6/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the Centauri upgrade.
	UpgradeName = "v6.0.3-ics"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
