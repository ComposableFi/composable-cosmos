package types

// ParamsKey is the key to use for the keeper store.
var (
	ParamsKey = []byte{0x01} // key for global staking middleware params in the keeper store
)

const (
	// module name
	ModuleName = "ibctransferparamsmodule"

	// StoreKey is the default store key for ibctransfermiddleware module that store params when apply validator set changes and when allow to unbond/redelegate

	StoreKey = "customibcparams" // not using the module name because of collisions with key "staking"

	RouterKey = ModuleName
)
