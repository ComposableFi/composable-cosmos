package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	coretypes "github.com/cosmos/ibc-go/v7/modules/core/types"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	fmt.Println("dsahdhasdhjasdjhasjhdsajhg")
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
		// sender chain is not the source, unescrow tokens

		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := data.Denom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom := unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}
		token := sdk.NewCoin(denom, transferAmount)

		if k.bankKeeper.BlockedAddr(receiver) {
			return errorsmod.Wrapf(errors.New("Unauthorized"), "%s is not allowed to receive funds", receiver)
		}

		// unescrow tokens
		escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
		if err := k.bankKeeper.SendCoins(ctx, escrowAddress, receiver, sdk.NewCoins(token)); err != nil {
			// NOTE: this error is only expected to occur given an unexpected bug or a malicious
			// counterparty module. The bug may occur in bank or any part of the code that allows
			// the escrow address to be drained. A malicious counterparty module could drain the
			// escrow address by allowing more tokens to be sent back then were escrowed.
			return errorsmod.Wrap(err, "unable to unescrow tokens, this may be caused by a malicious counterparty module or a bug: please open an issue on counterparty module")
		}

		defer func() {
			if transferAmount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"ibc", types.ModuleName, "packet", "receive"},
					float32(transferAmount.Int64()),
					[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, unprefixedDenom)},
				)
			}

			telemetry.IncrCounterWithLabels(
				[]string{"ibc", types.ModuleName, "receive"},
				1,
				append(
					labels, telemetry.NewLabel(coretypes.LabelSource, "true"),
				),
			)
		}()

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
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypeDenomTrace,
			sdk.NewAttribute(transfertypes.AttributeKeyTraceHash, traceHash.String()),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, voucherDenom),
		),
	)
	voucher := sdk.NewCoin(voucherDenom, transferAmount)

	// mint new tokens if the source of the transfer is the same chain
	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(voucher),
	); err != nil {
		return errorsmod.Wrap(err, "failed to mint IBC tokens")
	}

	// lock ibc token if srcChannel is paraChannel
	if packet.GetSourceChannel() == paraTokenInfo.ChannelId {
		// escrow ibc token
		escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, escrowAddress, sdk.NewCoins(voucher),
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
