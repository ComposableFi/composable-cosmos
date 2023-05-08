package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

func (keeper Keeper) hasParachainIBCTokenInfo(ctx sdk.Context, nativeDenom string) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(types.GetKeyKeysParachainIBCTokenInfo(nativeDenom))
}

func (keeper Keeper) handleOverrideSendPacketTransferLogic(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string, sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	var fungibleTokenPacketData transfertypes.FungibleTokenPacketData

	err = keeper.cdc.Unmarshal(data, &fungibleTokenPacketData)
	if err != nil {
		return 0, err
	}
	sender, err := sdk.AccAddressFromBech32(fungibleTokenPacketData.Sender)
	if err != nil {
		return 0, err
	}

	// check if denom in fungibleTokenPacketData is native denom in parachain info and
	parachainInfo := keeper.GetParachainIBCTokenInfo(ctx, fungibleTokenPacketData.Denom)
	if parachainInfo.ChannelId != sourceChannel {
		keeper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// burn native token in escrow address
	transferAmount, ok := sdk.NewIntFromString(fungibleTokenPacketData.Amount)

	// TODO: remove this panic and replace by err hanlde
	if !ok {
		panic("cannot parse string amount")
	}
	nativeTransferToken := sdk.NewCoin(fungibleTokenPacketData.Denom, transferAmount)
	ibcTransferToken := sdk.NewCoin(parachainInfo.IbcDenom, transferAmount)

	// burn native token
	keeper.bankKeeper.BurnCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(nativeTransferToken))

	// release lock IBC token and send it to sender
	keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, sender, sdk.NewCoins(ibcTransferToken))

	// new msg transfer from transfer to parachain
	transferMsg := transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    sourceChannel,
		Token:            ibcTransferToken,
		Sender:           fungibleTokenPacketData.Sender,
		Receiver:         fungibleTokenPacketData.Receiver,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Memo:             fungibleTokenPacketData.Memo,
	}
	res, err := keeper.executeTransferMsg(ctx, &transferMsg)
	if err != nil {
		return 0, err
	}

	return res.Sequence, nil
}

func (k Keeper) executeTransferMsg(ctx sdk.Context, transferMsg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error) {
	if err := transferMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("bad msg %v", err.Error())
	}
	return k.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)

}

// SendPacket wraps IBC ChannelKeeper's SendPacket function
func (k Keeper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string, sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	return k.handleOverrideSendPacketTransferLogic(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// WriteAcknowledgement wraps IBC ICS4Wrapper WriteAcknowledgement function.
// ICS29 WriteAcknowledgement is used for asynchronous acknowledgements.
func (k *Keeper) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, acknowledgement ibcexported.Acknowledgement) error {
	return k.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, acknowledgement)
}

// WriteAcknowledgement wraps IBC ICS4Wrapper GetAppVersion function.
func (k *Keeper) GetAppVersion(
	ctx sdk.Context,
	portID,
	channelID string,
) (string, bool) {
	return k.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}
