package keeper

import (
	"context"

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

	// Send Fund to account module
	err = ms.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, req.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgFundModuleAccountResponse{}, nil
}
