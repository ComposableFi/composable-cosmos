package v6_4

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/notional-labs/composable/v6/app/upgrades"
	customstmiddleware "github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v6_4"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{customstmiddleware.StoreKey},
		Deleted: []string{},
	},
}
