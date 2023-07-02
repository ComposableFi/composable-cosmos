package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/centauri/v3/x/ratelimit/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// TODO: implement init genesis

}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	// TODO: implement export genesis

	return genesis
}
