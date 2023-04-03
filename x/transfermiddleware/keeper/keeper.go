package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper returns a new instance of the x/ibchooks keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	codec codec.BinaryCodec,
) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      codec,
	}
}
