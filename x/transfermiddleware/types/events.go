package types

// staking module event types
const (
	EventAddParachainIBCTokenInfo    = "add-parachain-token-info"    // #nosec G101
	EventRemoveParachainIBCTokenInfo = "remove-parachain-token-info" // #nosec G101

	AttributeKeyNativeDenom = "native-denom"
	AttributeKeyIbcDenom    = "ibc-denom"
	AttributeKeyAssetID     = "asset-id"
)
