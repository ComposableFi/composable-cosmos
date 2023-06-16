package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
)

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	fmt.Printf("-------------------------test1\n")
	// decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to decode receiver address: %s", data.Receiver)
	}

	fmt.Printf("-------------------------test2\n")
	// parse the transfer amount
	transferAmount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		return errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount: %s", data.Amount)
	}

	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		// Do nothing
		return nil
	}

	fmt.Printf("-------------------------test3\n")
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

	paraTokenInfo := k.GetParachainIBCTokenInfoByAssetID(ctx, data.Denom)

	fmt.Printf("-------------------------test4\n")
	fmt.Printf("k.GetNativeDenomByIBCDenomSecondaryIndex(ctx, denomTrace.IBCDenom()): %v\n", k.GetNativeDenomByIBCDenomSecondaryIndex(ctx, denomTrace.IBCDenom()))
	fmt.Printf("paraTokenInfo.NativeDenom: %v\n", paraTokenInfo.NativeDenom)
	if k.GetNativeDenomByIBCDenomSecondaryIndex(ctx, denomTrace.IBCDenom()) != paraTokenInfo.NativeDenom {
		return nil
	}

	fmt.Printf("-------------------------test5\n")
	// lock ibc token if dstChannel is paraChannel
	if packet.GetDestChannel() == paraTokenInfo.ChannelId {
		fmt.Printf("-------------------------test2\n")
		// escrow ibc token
		escrowAddress := transfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())

		if err := k.bankKeeper.SendCoins(
			ctx, receiver, escrowAddress, sdk.NewCoins(voucher),
		); err != nil {
			return errorsmod.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
		}

		// mint native token
		denom := paraTokenInfo.NativeDenom
		nativeCoin := sdk.NewCoin(denom, transferAmount)

		if err := k.bankKeeper.MintCoins(
			ctx, types.ModuleName, sdk.NewCoins(nativeCoin),
		); err != nil {
			return errorsmod.Wrap(err, "failed to mint IBC tokens")
		}
		fmt.Printf("-------------------------test2: %v\n", nativeCoin)
		// send to receiver
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(nativeCoin),
		); err != nil {
			return errorsmod.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
		}
	}

	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	return k.refundToken(ctx, packet, data)
}
