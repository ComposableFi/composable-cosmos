package v7

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/notional-labs/composable/v6/app/upgrades"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v7-prefix-change"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{authz.ModuleName},
		Deleted: []string{alliancetypes.ModuleName},
	},
}
