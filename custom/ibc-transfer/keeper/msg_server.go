package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibctransfermiddlewaretypes "github.com/notional-labs/composable/v6/x/ibctransfermiddleware/types"
)

type msgServer struct {
	Keeper
	bank      types.BankKeeper
	msgServer types.MsgServer
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.Keeper.IbcTransfermiddleware.GetParams(ctx)
	if params.ChannelFees != nil && len(params.ChannelFees) > 0 {
		channelFee := findChannelParams(params.ChannelFees, msg.SourceChannel)
		if channelFee != nil {
			coin := findCoinByDenom(channelFee.AllowedTokens, msg.Token.Denom)
			if coin != nil {
				return &types.MsgTransferResponse{}, nil
			}
			minFee := coin.MinFee.Amount
			charge := minFee
			if charge.GT(msg.Token.Amount) {
				charge = msg.Token.Amount
			}

			newAmount := msg.Token.Amount.Sub(charge)

			if newAmount.IsPositive() {
				percentageCharge := newAmount.QuoRaw(coin.Percentage)
				newAmount = newAmount.Sub(percentageCharge)
				charge = charge.Add(percentageCharge)
			}

			msgSender, err := sdk.AccAddressFromBech32(msg.Sender)
			if err != nil {
				return nil, err
			}

			feeAddress, err := sdk.AccAddressFromBech32(channelFee.FeeAddress)
			if err != nil {
				return nil, err
			}

			k.bank.SendCoins(ctx, msgSender, feeAddress, sdk.NewCoins(sdk.NewCoin(msg.Token.Denom, charge)))

			if newAmount.IsZero() {
				return &types.MsgTransferResponse{}, nil
			}
			msg.Token.Amount = newAmount
		}

	}

	return k.msgServer.Transfer(goCtx, msg)
}

func findChannelParams(channelFees []*ibctransfermiddlewaretypes.ChannelFee, targetChannelID string) *ibctransfermiddlewaretypes.ChannelFee {
	for _, fee := range channelFees {
		if fee.Channel == targetChannelID {
			return fee
		}
	}
	return nil // If the channel is not found
}

func findCoinByDenom(allowedTokens []*ibctransfermiddlewaretypes.CoinItem, denom string) *ibctransfermiddlewaretypes.CoinItem {
	for _, coin := range allowedTokens {
		if coin.MinFee.Denom == denom {
			return coin
		}
	}
	return nil // If the denom is not found
}
