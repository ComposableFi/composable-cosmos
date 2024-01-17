package types

// ParamsKey is the key to use for the keeper store.
var (
	ParamsKey      = []byte{0x01} // key for global staking middleware params in the keeper store
	RewardDenomKey = []byte{0x02}
)

const (
	// module name
	ModuleName       = "stakingmiddleware"
	RewardModuleName = "fee_collector"

	// StoreKey is the default store key for stakingmiddleware module that store params when apply validator set changes and when allow to unbond/redelegate

	StoreKey = "customstakingparams" // not using the module name because of collisions with key "staking"
)
