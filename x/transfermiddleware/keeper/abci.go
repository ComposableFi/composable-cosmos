package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
)

// BeginBlocker of epochs module.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Iterate over remove list
	k.IterateRemoveListInfo(ctx, func(removeList types.RemoveParachainIBCTokenInfo) (stop bool) {
		// If pass the duration, remove parachain token info
		if removeList.RemoveTime.After(ctx.BlockTime()) {
			k.RemoveParachainIBCInfo(ctx, removeList.NativeDenom)
		}

		return false
	})

}
