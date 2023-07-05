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
	KeyRlyAddress                         = []byte{0x04}
	KeyTotalTransfered                    = []byte{0x05}
)

func GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom string) []byte {
	return append(KeyParachainIBCTokenInfoByNativeDenom, []byte(nativeDenom)...)
}

func GetKeyParachainIBCTokenInfoByAssetID(assetID string) []byte {
	return append(KeyParachainIBCTokenInfoByAssetID, []byte(assetID)...)
}

func GetKeyNativeDenomAndIbcSecondaryIndex(ibcDenom string) []byte {
	return append(KeyIBCDenomAndNativeIndex, []byte(ibcDenom)...)
}

func GetKeyByRlyAddress(rlyAddress string) []byte {
	return append(KeyRlyAddress, []byte(rlyAddress)...)
}
