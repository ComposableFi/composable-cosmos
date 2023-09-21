package v6

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
	"github.com/notional-labs/centauri/v6/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the Centauri upgrade.
	UpgradeName = "v6"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{ccvconsumertypes.ModuleName},
	},
}
