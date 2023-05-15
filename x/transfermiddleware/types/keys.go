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
	KeysParachainIBCTokenInfo = []byte{0x01}
	KeyIBCDenomAndNativeIndex = []byte{0x02}
	KeyInflightPacket         = []byte{0x03}
)

func GetKeyParachainIBCTokenInfo(nativeDenom string) []byte {
	return append(KeysParachainIBCTokenInfo, []byte(nativeDenom)...)
}

func GetKeyIBCDenomAndNativeIndex(IBCdenom string) []byte {
	return append(KeysParachainIBCTokenInfo, []byte(IBCdenom)...)
}

func GetKeyInFlightPacketKey(packetBz []byte) []byte {
	return append(KeysParachainIBCTokenInfo, packetBz...)
}
