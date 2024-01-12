package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

var _ types.MsgServer = msgServer{}

// msgServer is a wrapper of Keeper.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the x/stakingmiddleware MsgServer interface.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

// UpdateParams updates the params.
func (ms msgServer) UpdateEpochParams(goCtx context.Context, req *types.MsgUpdateEpochParams) (*types.MsgUpdateParamsEpochResponse, error) {
	if ms.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsEpochResponse{}, nil
}

// UpdateParams updates the params.
func (ms msgServer) AddRevenueFundsToStaking(goCtx context.Context, req *types.MsgAddRevenueFundsToStakingParams) (*types.MsgAddRevenueFundsToStakingResponse, error) {
	// Unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check sender address
	sender, err := sdk.AccAddressFromBech32(req.FromAddress)
	if err != nil {
		return nil, err
	}

	rewardDenom := ms.GetRewardDenom(ctx)

	// Check that reward is 1 coin rewardDenom
	if len(req.Amount.Denoms()) != 1 || req.Amount[0].Denom != rewardDenom.Denom {
		return nil, errorsmod.Wrapf(types.ErrInvalidCoin, "Invalid coin")
	}

	// Send Fund to account module
	moduleAccountAccAddress := ms.GetModuleAccountAccAddress(ctx)
	err = ms.bankKeeper.SendCoins(ctx, sender, moduleAccountAccAddress, req.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddRevenueFundsToStakingResponse{}, nil
}
