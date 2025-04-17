package keeper

import (
	"context"
	"fmt"
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	custombankkeeper "github.com/notional-labs/composable/v6/custom/bank/keeper"
	ibctransfermiddlewaretypes "github.com/notional-labs/composable/v6/x/ibctransfermiddleware/types"
)

type msgServer struct {
	Keeper
	bank      types.BankKeeper
	msgServer types.MsgServer
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(ibcKeeper Keeper, bankKeeper custombankkeeper.Keeper) types.MsgServer {
	return &msgServer{Keeper: ibcKeeper, bank: bankKeeper, msgServer: ibcKeeper.Keeper}
}

// Transfer is the server API around the Transfer method of the IBC transfer module.
// It checks if the sender is allowed to transfer the token and if the channel has fees.
// If the channel has fees, it will charge the sender and send the fees to the fee address.
// If the sender is not allowed to transfer the token because this tokens does not exists in the allowed tokens list, it just return without doing anything.
// If the sender is allowed to transfer the token, it will call the original transfer method.
// If the transfer amount is less than the minimum fee, it will charge the full transfer amount.
// If the transfer amount is greater than the minimum fee, it will charge the minimum fee and the percentage fee.
func (k msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.IbcTransfermiddleware.GetParams(ctx)
	charge_coin := sdk.NewCoin(msg.Token.Denom, sdk.ZeroInt())
	if params.ChannelFees != nil && len(params.ChannelFees) > 0 {
		channelFee := findChannelParams(params.ChannelFees, msg.SourceChannel)
		if channelFee != nil {
			if channelFee.MinTimeoutTimestamp > 0 {
				// check if the timeout timestamp is in the future
				if msg.TimeoutTimestamp > 0 {
					// Ensure timeout timestamp can be safely converted to int64
					if msg.TimeoutTimestamp > uint64(math.MaxInt64) {
						return nil, fmt.Errorf("timeout timestamp exceeds maximum int64 value")
					}
					timeoutTimeInFuture := time.Unix(0, int64(msg.TimeoutTimestamp))
					if timeoutTimeInFuture.Before(ctx.BlockTime()) {
						return nil, fmt.Errorf("timeout timestamp is in the past")
					}
				}

				// check if the channel is fee enabled
				senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
				if err != nil {
					return nil, err
				}
				balance := k.bank.GetBalance(ctx, senderAddr, channelFee.FeeAddress)
				if balance.Amount.LT(sdk.OneInt()) {
					return nil, fmt.Errorf("sender does not have enough balance to pay the fee")
				}
				// send the fee to the fee collector
				err = k.bank.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, sdk.NewCoins(sdk.NewCoin(channelFee.FeeAddress, sdk.OneInt())))
				if err != nil {
					return nil, err
				}
			}
			coin := findCoinByDenom(channelFee.AllowedTokens, msg.Token.Denom)
			if coin == nil {
				return nil, fmt.Errorf("token not allowed to be transferred in this channel")
			}

			minFee := coin.MinFee.Amount
			priority := GetPriority(msg.Memo)
			if priority != nil {
				p := findPriority(coin.TxPriorityFee, *priority)
				if p != nil && coin.MinFee.Denom == p.PriorityFee.Denom {
					minFee = minFee.Add(p.PriorityFee.Amount)
				}
			}

			charge := minFee
			if charge.GT(msg.Token.Amount) {
				charge = msg.Token.Amount
			}

			newAmount := msg.Token.Amount.Sub(charge)

			if newAmount.IsPositive() && coin.Percentage != 0 {
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

			charge_coin = sdk.NewCoin(msg.Token.Denom, charge)
			send_err := k.bank.SendCoins(ctx, msgSender, feeAddress, sdk.NewCoins(charge_coin))
			if send_err != nil {
				return nil, send_err
			}

			if newAmount.LTE(sdk.ZeroInt()) {
				return &types.MsgTransferResponse{}, nil
			}
			msg.Token.Amount = newAmount
		}
	}
	ret, err := k.msgServer.Transfer(goCtx, msg)
	if err == nil && ret != nil && !charge_coin.IsZero() {
		k.IbcTransfermiddleware.SetSequenceFee(ctx, ret.Sequence, charge_coin)
	}
	return ret, err
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
