package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

// TODO: add init genesis logic
// InitGenesis initializes the transfermiddleware module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	for _, tokenInfo := range genState.TokenInfos {
		k.AddParachainIBCInfo(ctx, tokenInfo.IbcDenom, tokenInfo.ChannelId, tokenInfo.NativeDenom)
	}
}

// ExportGenesis returns the x/transfermiddleware module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
