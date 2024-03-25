package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/notional-labs/composable/v6/x/ibctransfermiddleware/types"
)

var _ types.MsgServer = msgServer{}

// msgServer is a wrapper of Keeper.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the x/ibctransfermiddleware MsgServer interface.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

// UpdateParams updates the params.
func (ms msgServer) UpdateCustomIbcParams(goCtx context.Context, req *types.MsgUpdateCustomIbcParams) (*types.MsgUpdateParamsCustomIbcResponse, error) {
	if ms.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsCustomIbcResponse{}, nil
}

// AddIBCFeeConfig(MsgAddIBCFeeConfig) returns (MsgAddIBCFeeConfigResponse);
func (ms msgServer) AddIBCFeeConfig(goCtx context.Context, req *types.MsgAddIBCFeeConfig) (*types.MsgAddIBCFeeConfigResponse, error) {
	if ms.authority != req.Authority {
		// return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; Nikita expected %s, got %s", ms.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(req.FeeAddress)
	if err != nil {
		return nil, err
	}

	params := ms.Keeper.GetParams(ctx)
	channelFee := findChannelParams(params.ChannelFees, req.ChannelID)
	if channelFee != nil {
		channelFee.FeeAddress = req.FeeAddress
		channelFee.MinTimeoutTimestamp = req.MinTimeoutTimestamp
	} else {
		channelFee := &types.ChannelFee{
			Channel:             req.ChannelID,
			FeeAddress:          req.FeeAddress,
			MinTimeoutTimestamp: req.MinTimeoutTimestamp,
			AllowedTokens:       []*types.CoinItem{},
		}
		params.ChannelFees = append(params.ChannelFees, channelFee)
	}
	errSetParams := ms.Keeper.SetParams(ctx, params)
	if errSetParams != nil {
		return nil, errSetParams
	}
	return &types.MsgAddIBCFeeConfigResponse{}, nil
}

func (ms msgServer) RemoveIBCFeeConfig(goCtx context.Context, req *types.MsgRemoveIBCFeeConfig) (*types.MsgRemoveIBCFeeConfigResponse, error) {
	if ms.authority != req.Authority {
		// return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.Keeper.GetParams(ctx)
	//remove channel id in list of channel fees
	for i, fee := range params.ChannelFees {
		if fee.Channel == req.ChannelID {
			params.ChannelFees = append(params.ChannelFees[:i], params.ChannelFees[i+1:]...)
			break
		}
	}
	errSetParams := ms.Keeper.SetParams(ctx, params)
	if errSetParams != nil {
		return nil, errSetParams
	}

	return &types.MsgRemoveIBCFeeConfigResponse{}, nil
}

func findChannelParams(channelFees []*types.ChannelFee, targetChannelID string) *types.ChannelFee {
	for _, fee := range channelFees {
		if fee.Channel == targetChannelID {
			return fee
		}
	}
	return nil // If the channel is not found
}
