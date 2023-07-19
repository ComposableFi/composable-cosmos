package types

// MinterKey is the key to use for the keeper store.
var (
	MinterKey         = []byte{0x00}
	AllowedAddressKey = []byte{0x01}
)

const (
	// module name
	ModuleName = "mint"

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
)

func GetAllowedAddressStoreKey(address string) []byte {
	return append(AllowedAddressKey, []byte(address)...)
}
