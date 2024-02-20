package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	stakingmiddleware "github.com/notional-labs/composable/v6/x/stakingmiddleware/keeper"
)

type Keeper struct {
	ibctransferkeeper.Keeper
	cdc               codec.BinaryCodec
	Stakingmiddleware *stakingmiddleware.Keeper
	authority         string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	ics4Wrapper porttypes.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	authKeeper types.AccountKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	scopedKeeper exported.ScopedKeeper,
	authority string,
	stakingmiddleware *stakingmiddleware.Keeper,
	//return type from this function is different from the staking keeper.
	//todo double check if this is correct
) Keeper {
	keeper := Keeper{
		Keeper:            ibctransferkeeper.NewKeeper(cdc, key, paramSpace, ics4Wrapper, channelKeeper, portKeeper, authKeeper, bk, scopedKeeper),
		authority:         authority,
		Stakingmiddleware: stakingmiddleware,
		cdc:               cdc,
	}
	return keeper
}
