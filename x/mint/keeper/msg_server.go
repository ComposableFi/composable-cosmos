package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/centauri/v3/x/mint/types"
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

func (ms msgServer) FundModuleAccount(goCtx context.Context, req *types.MsgFundModuleAccount) (*types.MsgFundModuleAccountResponse, error) {
	// Unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Check sender address
	sender, err := sdk.AccAddressFromBech32(req.FromAddress)
	if err != nil {
		return nil, err
	}

	params := ms.GetParams(ctx)

	if len(req.Amount.Denoms()) > 1 || req.Amount[0].Denom != params.MintDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalidCoin, "Invalid fund")
	}

	// Send Fund to account module
	moduleAccountAccAddress := ms.GetModuleAccountAccAddress(ctx)
	err = ms.bankKeeper.SendCoins(ctx, sender, moduleAccountAccAddress, req.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgFundModuleAccountResponse{}, nil
}
