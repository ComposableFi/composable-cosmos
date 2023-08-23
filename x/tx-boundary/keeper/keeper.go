package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/centauri/v4/x/tx-boundary/types"
)

// Keeper struct
type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper returns keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, authority string) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		authority: authority,
	}
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// TODO: Duplicate fnc here
// SetDelegateBoundary sets the delegate boundary.
func (k Keeper) SetDelegateBoundary(ctx sdk.Context, boundary types.Boundary) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&boundary)
	store.Set(types.DelegateBoundaryKey, bz)
}

// TODO: Duplicate fnc here
// GetDelegateBoundary sets the delegate boundary.
func (k Keeper) GetDelegateBoundary(ctx sdk.Context) (boundary types.Boundary) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DelegateBoundaryKey)
	if bz == nil {
		panic("stored delegate boundary should not have been nil")
	}

	k.cdc.MustUnmarshal(bz, &boundary)
	return
}

// SetRedelegateBoundary sets the delegate boundary.
func (k Keeper) SetRedelegateBoundary(ctx sdk.Context, boundary types.Boundary) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&boundary)
	store.Set(types.RedelegateBoundaryKey, bz)
}

// GetRedelegateBoundary sets the delegate boundary.
func (k Keeper) GetRedelegateBoundary(ctx sdk.Context) (boundary types.Boundary) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RedelegateBoundaryKey)
	if bz == nil {
		panic("stored redelegate boundary should not have been nil")
	}

	k.cdc.MustUnmarshal(bz, &boundary)
	return
}

// SetDelegateCount set the number of delegate tx for a given address
func (k Keeper) SetLimitPerAddr(ctx sdk.Context, addr sdk.AccAddress, limit_per_addr types.LimitPerAddr) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&limit_per_addr)
	store.Set(addr, bz)
}

// GetDelegateCount get the number of delegate tx for a given address
func (k Keeper) GetLimitPerAddr(ctx sdk.Context, addr sdk.AccAddress) (limit_per_addr types.LimitPerAddr) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(addr) == false {
		return types.LimitPerAddr{
			DelegateCount:     0,
			ReledegateCount:   0,
			LatestUpdateBlock: 0,
		}
	}
	bz := store.Get(addr)
	k.cdc.MustUnmarshal(bz, &limit_per_addr)
	return
}

func (k Keeper) UpdateLimitPerAddr(ctx sdk.Context, addr sdk.AccAddress) {
	limit_per_addr := k.GetLimitPerAddr(ctx, addr)
	if limit_per_addr.LatestUpdateBlock == 0 {
		return
	}
	boundary := k.GetDelegateBoundary(ctx)
	if limit_per_addr.LatestUpdateBlock+boundary.BlocksPerGeneration >= ctx.BlockHeight() {
		// Calculate the generated tx number from the duration between latest update block and curent block height
		var generatedTx int64
		duration := limit_per_addr.LatestUpdateBlock + boundary.BlocksPerGeneration - ctx.BlockHeight()
		if duration/boundary.BlocksPerGeneration > 5 {
			generatedTx = 5
		} else {
			generatedTx = duration / boundary.BlocksPerGeneration
		}

		// Update the delegate tx limit
		if uint64(generatedTx) > limit_per_addr.DelegateCount {
			limit_per_addr.DelegateCount = 0
		} else {
			limit_per_addr.DelegateCount -= uint64(generatedTx)
		}
		// Update the redelegate tx limit
		if uint64(generatedTx) > limit_per_addr.ReledegateCount {
			limit_per_addr.ReledegateCount = 0
		} else {
			limit_per_addr.ReledegateCount -= uint64(generatedTx)
		}
		// Update LatestUpdateBlock
		limit_per_addr.LatestUpdateBlock = ctx.BlockHeight()
		return
	}
	return
}
