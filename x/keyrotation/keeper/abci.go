package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {

}
