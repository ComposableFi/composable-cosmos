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
		//find the channel fee with a matching channel
		if channelFee != nil {
			//find the coin with a matching denom
			coin := findCoinByDenom(channelFee.AllowedTokens, msg.Token.Denom)
			if coin != nil {
				//token not allowed by this channel. should ignore the transfer
				return &types.MsgTransferResponse{}, nil
			}
			minFee := coin.MinFee.Amount
			charge := minFee
			if charge.GT(msg.Token.Amount) {
				charge = msg.Token.Amount
			}

			newAmount := msg.Token.Amount.Sub(charge)

			if newAmount.IsPositive() {
				//if Percentage = 100 it means we charge 1 % of the amount
				percentageCharge := newAmount.QuoRaw(coin.Percentage)
				newAmount = newAmount.Sub(percentageCharge)
				charge = charge.Add(percentageCharge)
			}

			//address from string
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
				//if the new amount is zero, then the transfer should be ignored
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
