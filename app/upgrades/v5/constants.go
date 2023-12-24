package v5

import (
	store "cosmossdk.io/store/types"
	"github.com/notional-labs/composable/v6/app/upgrades"
	txboundary "github.com/notional-labs/composable/v6/x/tx-boundary/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v5"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{txboundary.ModuleName},
	},
}
