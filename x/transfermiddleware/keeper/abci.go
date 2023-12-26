package keeper

import (
	"github.com/notional-labs/composable/v6/x/transfermiddleware/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker of epochs module.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Iterate over remove list
	k.IterateRemoveListInfo(ctx, func(removeList types.RemoveParachainIBCTokenInfo) (stop bool) {
		// If pass the duration, remove parachain token info
		if ctx.BlockTime().After(removeList.RemoveTime) {
			err := k.RemoveParachainIBCInfo(ctx, removeList.NativeDenom)
			if err != nil {
				return true
			}
		}
		return false
	})
}
