package keeper

import (
	"context"

	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

type msgServer struct {
	Keeper
	msgServer types.MsgServer
}

var _ types.MsgServer = msgServer{}

// // TODO - Add the stakingkeeper.Keeper as a parameter to the NewMsgServerImpl function
// func NewMsgServerImpl(stakingKeeper stakingkeeper.Keeper, customstakingkeeper Keeper) types.MsgServer {
// 	return &msgServer{Keeper: customstakingkeeper, msgServer: stakingkeeper.NewMsgServerImpl(&stakingKeeper)}
// }

func (k msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	return k.msgServer.Transfer(goCtx, msg)
}
