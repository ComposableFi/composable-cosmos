package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

// InitGenesis new mint genesis
func (keeper Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	keeper.SetLastTotalPower(ctx, data.LastTotalPower)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (keeper Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	minter := keeper.GetLastTotalPower(ctx)
	return types.NewGenesisState(minter)
}
