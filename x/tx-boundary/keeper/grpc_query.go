package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/centauri/v4/x/tx-boundary/types"
)

var _ types.QueryServer = Keeper{}

// DelegateBoundary returns params of the mint module.
func (k Keeper) DelegateBoundary(c context.Context, req *types.QueryDelegateBoundaryRequest) (*types.QueryDelegateBoundaryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	boundary := k.GetDelegateBoundary(ctx)

	return &types.QueryDelegateBoundaryResponse{Boundary: boundary}, nil
}

// DelegateBoundary returns params of the mint module.
func (k Keeper) RedelegateBoundary(c context.Context, _ *types.QueryRedelegateBoundaryRequest) (*types.QueryRedelegateBoundaryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	boundary := k.GetRedelegateBoundary(ctx)

	return &types.QueryRedelegateBoundaryResponse{Boundary: boundary}, nil
}
