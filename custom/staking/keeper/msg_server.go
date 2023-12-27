package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	Keeper
	msgServer types.MsgServer
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(stakingKeeper stakingkeeper.Keeper, customstakingkeeper Keeper) types.MsgServer {
	return &msgServer{Keeper: customstakingkeeper, msgServer: stakingkeeper.NewMsgServerImpl(&stakingKeeper)}
}

func (k msgServer) CreateValidator(goCtx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	return k.msgServer.CreateValidator(goCtx, msg)
}

func (k msgServer) EditValidator(goCtx context.Context, msg *types.MsgEditValidator) (*types.MsgEditValidatorResponse, error) {
	return k.msgServer.EditValidator(goCtx, msg)
}

func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	return k.msgServer.Delegate(goCtx, msg)
}

func (k msgServer) BeginRedelegate(goCtx context.Context, msg *types.MsgBeginRedelegate) (*types.MsgBeginRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.Keeper.Stakingmiddleware.GetParams(ctx)
	height := ctx.BlockHeight()
	epoch_progress_block_number := (height % int64(params.BlocksPerEpoch))
	if epoch_progress_block_number > int64(params.AllowUnbondAfterEpochProgressBlockNumber) || epoch_progress_block_number == 0 {
		return k.msgServer.BeginRedelegate(goCtx, msg)
	}
	return &types.MsgBeginRedelegateResponse{}, nil
}

func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.Keeper.Stakingmiddleware.GetParams(ctx)
	height := ctx.BlockHeight()
	epoch_progress_block_number := (height % int64(params.BlocksPerEpoch))
	if epoch_progress_block_number > int64(params.AllowUnbondAfterEpochProgressBlockNumber) || epoch_progress_block_number == 0 {
		return k.msgServer.Undelegate(goCtx, msg)
	}
	return &types.MsgUndelegateResponse{}, nil
}

func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *types.MsgCancelUnbondingDelegation) (*types.MsgCancelUnbondingDelegationResponse, error) {
	return k.msgServer.CancelUnbondingDelegation(goCtx, msg)
}

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return ms.msgServer.UpdateParams(goCtx, msg)
}
