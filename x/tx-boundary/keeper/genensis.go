package keeper

import (
	"github.com/notional-labs/composable/v6/x/tx-boundary/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	err := k.SetDelegateBoundary(ctx, genState.DelegateBoundary)
	if err != nil {
		panic(err) //todo: handle error
	}
	err = k.SetRedelegateBoundary(ctx, genState.RedelegateBoundary)
	if err != nil {
		panic(err) //todo: handle error
	}
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()

	genesis.DelegateBoundary = k.GetDelegateBoundary(ctx)
	genesis.RedelegateBoundary = k.GetRedelegateBoundary(ctx)

	return genesis
}
