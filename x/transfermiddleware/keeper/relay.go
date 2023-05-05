package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

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

	paraTokenInfo := k.GetParachainIBCTokenInfo(ctx, data.Denom)

	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		// Do nothing
		return nil
	}

	// sender chain is the source, mint vouchers

	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + data.Denom

	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

	traceHash := denomTrace.Hash()
	if !k.transferKeeper.HasDenomTrace(ctx, traceHash) {
		k.transferKeeper.SetDenomTrace(ctx, denomTrace)
	}

	voucherDenom := denomTrace.IBCDenom()
	voucher := sdk.NewCoin(voucherDenom, transferAmount)

	// lock ibc token if srcChannel is paraChannel
	if packet.GetSourceChannel() == paraTokenInfo.ChannelId {
		// escrow ibc token
		escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())

		if err := k.bankKeeper.SendCoins(
			ctx, receiver, escrowAddress, sdk.NewCoins(voucher),
		); err != nil {
			return errorsmod.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
		}

		// mint native token
		denom := paraTokenInfo.NativeDenom
		voucher = sdk.NewCoin(denom, transferAmount)

		if err := k.bankKeeper.MintCoins(
			ctx, types.ModuleName, sdk.NewCoins(voucher),
		); err != nil {
			return errorsmod.Wrap(err, "failed to mint IBC tokens")
		}

		// send to receiver
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(voucher),
		); err != nil {
			return errorsmod.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
		}
	} else {
		// send to receiver
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(voucher),
		); err != nil {
			return errorsmod.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
		}
	}

	return nil
}
