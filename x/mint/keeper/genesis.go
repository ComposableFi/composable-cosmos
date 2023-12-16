package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/composable/v6/x/mint/types"
)

// InitGenesis new mint genesis
func (k Keeper) InitGenesis(ctx sdk.Context, ak types.AccountKeeper, data *types.GenesisState) {
	k.SetMinter(ctx, data.Minter)

	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	newCoins := sdk.NewCoins(data.IncentivesSupply)
	if err := k.MintCoins(ctx, newCoins); err != nil {
		panic(err)
	}

	ak.GetModuleAccount(ctx, types.ModuleName)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context, authKeeper types.AccountKeeper) *types.GenesisState {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	remIncentives := k.bankKeeper.GetBalance(ctx, authKeeper.GetModuleAddress(types.ModuleName), params.MintDenom)
	return types.NewGenesisState(minter, params, remIncentives)
}
