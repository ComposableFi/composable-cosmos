package keeper

import (
	"context"
	"fmt"
	stdmath "math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	Keeper
	msgServer stakingtypes.MsgServer
}

var _ stakingtypes.MsgServer = msgServer{}

func NewMsgServerImpl(stakingKeeper stakingkeeper.Keeper, customstakingkeeper Keeper) stakingtypes.MsgServer {
	return &msgServer{Keeper: customstakingkeeper, msgServer: stakingkeeper.NewMsgServerImpl(&stakingKeeper)}
}

func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	return k.msgServer.CreateValidator(goCtx, msg)
}

func (k msgServer) EditValidator(goCtx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	return k.msgServer.EditValidator(goCtx, msg)
}

func (k msgServer) Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	return k.msgServer.Delegate(goCtx, msg)
}

func (k msgServer) BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.Stakingmiddleware.GetParams(ctx)
	height := ctx.BlockHeight()

	if params.BlocksPerEpoch > uint64(stdmath.MaxInt64) {
		return nil, fmt.Errorf("blocks per epoch exceeds maximum int64 value")
	}
	if params.AllowUnbondAfterEpochProgressBlockNumber > uint64(stdmath.MaxInt64) {
		return nil, fmt.Errorf("allow unbond after epoch progress block number exceeds maximum int64 value")
	}

	epoch_progress_block_number := (height % int64(params.BlocksPerEpoch))
	if epoch_progress_block_number > int64(params.AllowUnbondAfterEpochProgressBlockNumber) || epoch_progress_block_number == 0 {
		return nil, fmt.Errorf("unbonding is not allowed at this block number in the epoch")
	}
	return k.msgServer.BeginRedelegate(goCtx, msg)
}

func (k msgServer) Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.Stakingmiddleware.GetParams(ctx)
	height := ctx.BlockHeight()

	if params.BlocksPerEpoch > uint64(stdmath.MaxInt64) {
		return nil, fmt.Errorf("blocks per epoch exceeds maximum int64 value")
	}
	if params.AllowUnbondAfterEpochProgressBlockNumber > uint64(stdmath.MaxInt64) {
		return nil, fmt.Errorf("allow unbond after epoch progress block number exceeds maximum int64 value")
	}

	epoch_progress_block_number := (height % int64(params.BlocksPerEpoch))
	if epoch_progress_block_number > int64(params.AllowUnbondAfterEpochProgressBlockNumber) || epoch_progress_block_number == 0 {
		return k.msgServer.Undelegate(goCtx, msg)
	}
	return &stakingtypes.MsgUndelegateResponse{}, nil
}

func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *stakingtypes.MsgCancelUnbondingDelegation) (*stakingtypes.MsgCancelUnbondingDelegationResponse, error) {
	return k.msgServer.CancelUnbondingDelegation(goCtx, msg)
}

func (ms msgServer) UpdateParams(goCtx context.Context, msg *stakingtypes.MsgUpdateParams) (*stakingtypes.MsgUpdateParamsResponse, error) {
	return ms.msgServer.UpdateParams(goCtx, msg)
}
