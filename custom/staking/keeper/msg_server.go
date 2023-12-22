package keeper

import (
	"context"

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
	return k.msgServer.BeginRedelegate(goCtx, msg)
}

func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	return k.msgServer.Undelegate(goCtx, msg)
}

func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *types.MsgCancelUnbondingDelegation) (*types.MsgCancelUnbondingDelegationResponse, error) {
	return k.msgServer.CancelUnbondingDelegation(goCtx, msg)
}

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return ms.msgServer.UpdateParams(goCtx, msg)
}
