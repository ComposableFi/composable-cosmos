package types

// Minting module event types
const (
	EventTypeMint              = ModuleName
	EventTypeReward            = "reward_distributed"
	EventTypeMintSlashed       = "mint_slashed_into_comminity_pool"
	EventAddAllowedFundAddress = "add_allowed_fund"

	AttributeKeyBondedRatio      = "bonded_ratio"
	AttributeKeyInflation        = "inflation"
	AttributeKeyAnnualProvisions = "annual_provisions"
	AttributeKeyAllowedAddress   = "allowed_address"
)
