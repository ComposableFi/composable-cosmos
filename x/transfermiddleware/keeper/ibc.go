package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	coretypes "github.com/cosmos/ibc-go/v7/modules/core/types"
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

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	if err := data.ValidateBasic(); err != nil {
		return errorsmod.Wrapf(err, "error validating ICS-20 transfer packet data")
	}

	if !k.transferKeeper.GetReceiveEnabled(ctx) {
		return transfertypes.ErrReceiveDisabled
	}

	// decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to decode receiver address: %s", data.Receiver)
	}

	// parse the transfer amount
	transferAmount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		return errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount: %s", data.Amount)
	}

	labels := []metrics.Label{
		telemetry.NewLabel(coretypes.LabelSourcePort, packet.GetSourcePort()),
		telemetry.NewLabel(coretypes.LabelSourceChannel, packet.GetSourceChannel()),
	}

	paraTokenInfo := k.GetParachainIBCTokenInfo(ctx, data.Denom)

	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		return errorsmod.Wrap(errors.New("Not source chain"), "sender chain is not the source")
	}

	// sender chain is the source, mint vouchers

	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := paraTokenInfo.NativeDenom
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix

	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

	voucherDenom := denomTrace.IBCDenom()
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypeDenomTrace,
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, voucherDenom),
		),
	)
	voucher := sdk.NewCoin(voucherDenom, transferAmount)

	// mint new tokens if the source of the transfer is the same chain
	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(voucher),
	); err != nil {
		return errorsmod.Wrap(err, "failed to mint tokens")
	}

	// send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, receiver, sdk.NewCoins(voucher),
	); err != nil {
		return errorsmod.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
	}

	defer func() {
		if transferAmount.IsInt64() {
			telemetry.SetGaugeWithLabels(
				[]string{"ibc", types.ModuleName, "packet", "receive"},
				float32(transferAmount.Int64()),
				[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, data.Denom)},
			)
		}

		telemetry.IncrCounterWithLabels(
			[]string{"ibc", types.ModuleName, "receive"},
			1,
			append(
				labels, telemetry.NewLabel(coretypes.LabelSource, "false"),
			),
		)
	}()

	return nil
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

	return k.ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
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
