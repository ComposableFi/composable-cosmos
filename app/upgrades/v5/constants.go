package v5

import "github.com/notional-labs/centauri/v4/app/upgrades"

const (
	// UpgradeName defines the on-chain upgrade name for the Composable v5 upgrade.
	UpgradeName = "v5"

	// UpgradeHeight defines the block height at which the Composable v6 upgrade is
	// triggered.
	UpgradeHeight = 1_116_808
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
