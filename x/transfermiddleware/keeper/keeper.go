package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       storetypes.StoreKey
	ics4Wrapper    porttypes.ICS4Wrapper
	bankKeeper     types.BankKeeper
	transferKeeper ibctransferkeeper.Keeper

	// the address capable of executing a AddParachainIBCTokenInfo and RemoveParachainIBCTokenInfo message. Typically, this
	// should be the x/gov module account.
	authority string
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

// TODO: testing
// AddParachainIBCTokenInfo add new parachain token information token to chain state.
func (keeper Keeper) AddParachainIBCTokenInfo(ctx sdk.Context, ibcDenom, channelId, nativeDenom string) error {
	if keeper.hasParachainIBCTokenInfo(ctx, nativeDenom) {
		return types.ErrDuplicateParachainIBCTokenInfo
	}

	info := types.ParachainIBCTokenInfo{
		IbcDenom:    ibcDenom,
		ChannelId:   channelId,
		NativeDenom: nativeDenom,
	}

	bz, err := keeper.cdc.Marshal(&info)
	if err != nil {
		return err
	}
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.GetKeyKeysParachainIBCTokenInfo(nativeDenom), bz)

	return nil
}

// TODO: testing
// RemoveParachainIBCTokenInfo remove parachain token information from chain state.
func (keeper Keeper) RemoveParachainIBCTokenInfo(ctx sdk.Context, ibcDenom, channelId, nativeDenom string) error {
	if !keeper.hasParachainIBCTokenInfo(ctx, nativeDenom) {
		return types.ErrDuplicateParachainIBCTokenInfo
	}

	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.GetKeyKeysParachainIBCTokenInfo(nativeDenom))

	return nil
}

// TODO: testing
// GetParachainIBCTokenInfo add new information about parachain token to chain state.
func (keeper Keeper) GetParachainIBCTokenInfo(ctx sdk.Context, nativeDenom string) (info types.ParachainIBCTokenInfo) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.GetKeyKeysParachainIBCTokenInfo(nativeDenom))

	keeper.cdc.Unmarshal(bz, &info)

	return info
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+exported.ModuleName+"-"+types.ModuleName)
}
