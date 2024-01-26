package legacy

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

const (
	TypeMsgAddRateLimit    = "add_rate_limit"
	TypeMsgUpdateRateLimit = "update_rate_limit"
	TypeMsgRemoveRateLimit = "remove_rate_limit"
	TypeMsgResetRateLimit  = "reset_rate_limit"
)

var _ sdk.Msg = &MsgAddRateLimit{}

func NewMsgAddRateLimit(
	authority string,
	denom string,
	channelID string,
	maxPercentSend math.Int,
	maxPercentRecv math.Int,
	durationHours uint64,
) *MsgAddRateLimit {
	return &MsgAddRateLimit{
		Authority:      authority,
		Denom:          denom,
		ChannelId:      channelID,
		MaxPercentSend: maxPercentSend,
		MaxPercentRecv: maxPercentRecv,
		DurationHours:  durationHours,
	}
}

// Route Implements Msg.
func (msg MsgAddRateLimit) Route() string { return "" }

// Type Implements Msg.
func (msg MsgAddRateLimit) Type() string { return TypeMsgAddRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgAddRateLimit) GetSignBytes() []byte {
	return nil
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgAddRateLimit) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgAddRateLimit) ValidateBasic() error {
	// validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	// validate channelIDs
	if err := host.ChannelIdentifierValidator(msg.ChannelId); err != nil {
		return err
	}

	if msg.MaxPercentSend.GT(math.NewInt(100)) || msg.MaxPercentSend.LT(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max-percent-send percent must be between 0 and 100 (inclusively), Provided: %v", msg.MaxPercentSend)
	}

	if msg.MaxPercentRecv.GT(math.NewInt(100)) || msg.MaxPercentRecv.LT(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max-percent-recv percent must be between 0 and 100 (inclusively), Provided: %v", msg.MaxPercentRecv)
	}

	if msg.MaxPercentRecv.IsZero() && msg.MaxPercentSend.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "either the max send or max receive threshold must be greater than 0")
	}

	if msg.MinRateLimitAmount.LTE(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "mint rate limit amount must be greater than 0")
	}

	if msg.DurationHours == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duration can not be zero")
	}

	return nil
}
