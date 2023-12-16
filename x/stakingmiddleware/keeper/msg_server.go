package keeper

import (
	"context"

	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

var _ types.MsgServer = msgServer{}

// msgServer is a wrapper of Keeper.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the x/mint MsgServer interface.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

// UpdateParams updates the params.
func (ms msgServer) SetPower(goCtx context.Context, req *types.MsgSetPower) (*types.MsgUpdateParamsResponse, error) {

	return &types.MsgSetPowerResponse{}, nil
}
