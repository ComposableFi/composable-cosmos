package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/keyrotation/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) RotateConsPubKey(goCtx context.Context, msg *types.MsgRotateConsPubKey) (*types.MsgRotateConsPubKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// check to see if the validator not exist
	if _, found := k.sk.GetValidator(ctx, valAddr); !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator not exists")
	}

	return &types.MsgRotateConsPubKeyResponse{}, nil
}
