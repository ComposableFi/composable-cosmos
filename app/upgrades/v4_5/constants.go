package v45

import "github.com/notional-labs/composable/v5/app/upgrades"

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName   = "v4_5"
	UpgradeHeight = 967554
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
