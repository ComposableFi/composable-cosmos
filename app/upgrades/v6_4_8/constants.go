package v6_4_8

import (
	"github.com/notional-labs/composable/v6/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v6_4_7"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
}
