package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Confirm the channel value is not zero
	channelValue := k.GetChannelValue(ctx, msg.Denom)
	if channelValue.IsZero() {
		return nil, errors.Wrap(types.ErrZeroChannelValue, "zero channel value")
	}

	// Confirm the rate limit does not already exist
	_, found := k.GetRateLimit(ctx, msg.Denom, msg.ChannelId)
	if found {
		return nil, errors.Wrap(types.ErrRateLimitAlreadyExists, "rate limit already exists")
	}

	// Create and store the rate limit object
	path := types.Path{
		Denom:     msg.Denom,
		ChannelId: msg.ChannelId,
	}
	quota := types.Quota{
		MaxPercentSend: msg.MaxPercentSend,
		MaxPercentRecv: msg.MaxPercentRecv,
		DurationHours:  msg.DurationHours,
	}
	flow := types.Flow{
		Inflow:       math.ZeroInt(),
		Outflow:      math.ZeroInt(),
		ChannelValue: channelValue,
	}

	k.SetRateLimit(ctx, types.RateLimit{
		Path:  &path,
		Quota: &quota,
		Flow:  &flow,
	})

	return &types.MsgAddRateLimitResponse{}, nil
}

func (k Keeper) UpdateRateLimit(goCtx context.Context, msg *types.MsgUpdateRateLimit) (*types.MsgUpdateRateLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Confirm the rate limit exists
	_, found := k.GetRateLimit(ctx, msg.Denom, msg.ChannelId)
	if !found {
		return nil, errors.Wrap(types.ErrRateLimitNotFound, "rate limit not found")
	}

	// Update the rate limit object with the new quota information
	// The flow should also get reset to 0
	path := types.Path{
		Denom:     msg.Denom,
		ChannelId: msg.ChannelId,
	}
	quota := types.Quota{
		MaxPercentSend: msg.MaxPercentSend,
		MaxPercentRecv: msg.MaxPercentRecv,
		DurationHours:  msg.DurationHours,
	}
	flow := types.Flow{
		Inflow:       math.ZeroInt(),
		Outflow:      math.ZeroInt(),
		ChannelValue: k.GetChannelValue(ctx, msg.Denom),
	}

	k.SetRateLimit(ctx, types.RateLimit{
		Path:  &path,
		Quota: &quota,
		Flow:  &flow,
	})

	return &types.MsgUpdateRateLimitResponse{}, nil
}

func (k Keeper) RemoveRateLimit(goCtx context.Context, msg *types.MsgRemoveRateLimit) (*types.MsgRemoveRateLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	_, found := k.GetRateLimit(ctx, msg.Denom, msg.ChannelId)
	if !found {
		return nil, errors.Wrap(types.ErrRateLimitNotFound, "rate limit not found")
	}

	k.removeRateLimit(ctx, msg.Denom, msg.ChannelId)
	return &types.MsgRemoveRateLimitResponse{}, nil
}

func (k Keeper) ResetRateLimit(goCtx context.Context, msg *types.MsgResetRateLimit) (*types.MsgResetRateLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	err := k.resetRateLimit(ctx, msg.Denom, msg.ChannelId)
	if err != nil {
		return nil, err
	}
	return &types.MsgResetRateLimitResponse{}, nil
}
