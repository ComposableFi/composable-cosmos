package v6

import (
	store "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/notional-labs/composable/v6/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v6"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{authz.ModuleName},
	},
}
