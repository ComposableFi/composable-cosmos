package types

const (
	// Module name store the name of the module
	ModuleName = "transfermiddleware"

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
)

func GetKeyKeysParachainIBCTokenInfo(nativeDenom string) []byte {
	return append(KeysParachainIBCTokenInfo, []byte(nativeDenom)...)
}
