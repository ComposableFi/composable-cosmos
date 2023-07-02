package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/notional-labs/centauri/v3/x/ratelimit/types"
)

type (
	Keeper struct {
		storeKey   storetypes.StoreKey
		cdc        codec.BinaryCodec
		paramstore paramtypes.Subspace

		bankKeeper    types.BankKeeper
		channelKeeper types.ChannelKeeper
		ics4Wrapper   types.ICS4Wrapper
	}
)
