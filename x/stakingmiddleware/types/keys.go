package types

// MinterKey is the key to use for the keeper store.
var (
	ParamsKey = []byte{0x01} // key for global staking middleware params in the keeper store
)

const (
	// module name
	ModuleName = "stakingmiddleware"

	// StoreKey is the default store key for mint

	StoreKey = "customstmiddleware" // not using the module name because of collisions with key "ibc"
)
