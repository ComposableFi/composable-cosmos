package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/keyrotation/types"
)

func (k Keeper) SetKeyRotationHistory(ctx sdk.Context, consRotationHistory types.ConsPubKeyRotationHistory) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&consRotationHistory)
	store.Set(types.GetKeyRotationHistory(consRotationHistory.OperatorAddress), bz)
}

func (k Keeper) IterateRotationHistory(ctx sdk.Context, cb func(types.ConsPubKeyRotationHistory) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyRotation)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rotationHistory types.ConsPubKeyRotationHistory
		k.cdc.MustUnmarshal(iterator.Value(), &rotationHistory)
		if cb(rotationHistory) {
			break
		}
	}
}
