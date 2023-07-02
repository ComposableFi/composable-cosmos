package keeper

import (
	"context"

	"github.com/notional-labs/centauri/v3/x/ratelimit/types"
)

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

type msgServer struct {
	Keeper
}

func (k Keeper) AddRateLimit(goCtx context.Context, msg *types.MsgAddRateLimit) (*types.MsgAddRateLimitResponse, error) {
	return &types.MsgAddRateLimitResponse{}, nil
}

func (k Keeper) UpdateRateLimit(goCtx context.Context, msg *types.MsgUpdateRateLimit) (*types.MsgUpdateRateLimitResponse, error) {
	return &types.MsgUpdateRateLimitResponse{}, nil
}

func (k Keeper) RemoveRateLimit(goCtx context.Context, msg *types.MsgRemoveRateLimit) (*types.MsgRemoveRateLimitResponse, error) {
	return &types.MsgRemoveRateLimitResponse{}, nil
}

func (k Keeper) ResetRateLimit(goCtx context.Context, msg *types.MsgResetRateLimit) (*types.MsgResetRateLimitResponse, error) {

}
