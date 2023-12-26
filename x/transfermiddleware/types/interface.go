package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type TransferMiddlewareKeeper interface {
	// GetParachainIBCTokenInfoByAssetID returns the parachain token information by assetID
	GetParachainIBCTokenInfoByAssetID(ctx sdk.Context, assetID string) ParachainIBCTokenInfo
	// GetParachainIBCTokenInfoByNativeDenom returns the parachain token information by native denom
	GetParachainIBCTokenInfoByNativeDenom(ctx sdk.Context, nativeDenom string) ParachainIBCTokenInfo
	// GetNativeDenomByIBCDenomSecondaryIndex returns the native denom by IBC denom secondary index
	GetNativeDenomByIBCDenomSecondaryIndex(ctx sdk.Context, ibcDenom string) string
	// GetKeyParachainIBCTokenRemoveListByNativeDenom returns the key of the parachain token remove list by native denom
	GetKeyParachainIBCTokenRemoveListByNativeDenom(nativeDenom string) []byte
	// GetKeyByRlyAddress returns the key of the relay address
	GetKeyByRlyAddress(rlyAddress string) []byte
	// HasParachainIBCTokenInfoByAssetID returns true if the parachain token information exists by assetID
	HasParachainIBCTokenInfoByAssetID(ctx sdk.Context, assetID string) bool
	// HasParachainIBCTokenInfoByNativeDenom returns true if the parachain token information exists by native denom
	HasParachainIBCTokenInfoByNativeDenom(ctx sdk.Context, nativeDenom string) bool
	// HasRlyAddress returns true if the relay address exists
	HasRlyAddress(ctx sdk.Context, rlyAddress string) bool
	// IterateAllowRlyAddress iterates all relay addresses
	IterateAllowRlyAddress(ctx sdk.Context, cb func(rlyAddress string) (stop bool))
}
