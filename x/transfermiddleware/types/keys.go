package types

const (
	// Module name store the name of the module
	ModuleName = "transmiddleware"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the feeabs module
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// Contract: Coin denoms cannot contain this character
	KeySeparator = "|"
)

var (
	KeyParachainIBCTokenInfoByNativeDenom = []byte{0x01}
	KeyParachainIBCTokenInfoByAssetID     = []byte{0x02}
	KeyIBCDenomAndNativeIndex             = []byte{0x03}
)

func GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom string) []byte {
	return append(KeyParachainIBCTokenInfoByNativeDenom, []byte(nativeDenom)...)
}

func GetKeyParachainIBCTokenInfoByAssetID(assetId string) []byte {
	return append(KeyParachainIBCTokenInfoByAssetID, []byte(assetId)...)
}

func GetKeyNativeDenomAndIbcSecondaryIndex(ibcDenom string) []byte {
	return append(KeyIBCDenomAndNativeIndex, []byte(ibcDenom)...)
}
