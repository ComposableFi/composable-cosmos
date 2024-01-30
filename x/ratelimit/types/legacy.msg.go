package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/notional-labs/composable/v6/x/transfermiddleware/types"
)

const (
	TypeMsgAddRateLimitLegacy    = "add_rate_limit"
	TypeMsgUpdateRateLimitLegacy = "update_rate_limit"
)

var _ sdk.Msg = &MsgAddRateLimitLegacy{}

func NewMsgAddRateLimitLegacy(
	authority string,
	denom string,
	channelID string,
	maxPercentSend math.Int,
	maxPercentRecv math.Int,
	durationHours uint64,
) *MsgAddRateLimitLegacy {
	return &MsgAddRateLimitLegacy{
		Authority:      authority,
		Denom:          denom,
		ChannelId:      channelID,
		MaxPercentSend: maxPercentSend,
		MaxPercentRecv: maxPercentRecv,
		DurationHours:  durationHours,
	}
}

// Route Implements Msg.
func (msg MsgAddRateLimitLegacy) Route() string { return types.RouterKey }

// Type Implements Msg.
func (msg MsgAddRateLimitLegacy) Type() string { return TypeMsgAddRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgAddRateLimitLegacy) GetSignBytes() []byte {
	return sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgAddRateLimitLegacy) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgAddRateLimitLegacy) ValidateBasic() error {
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

var _ sdk.Msg = &MsgUpdateRateLimitLegacy{}

func NewMsgUpdateRateLimitLegacy(
	authority string,
	denom string,
	channelID string,
	maxPercentSend math.Int,
	maxPercentRecv math.Int,
	durationHours uint64,
) *MsgUpdateRateLimitLegacy {
	return &MsgUpdateRateLimitLegacy{
		Authority:      authority,
		Denom:          denom,
		ChannelId:      channelID,
		MaxPercentSend: maxPercentSend,
		MaxPercentRecv: maxPercentRecv,
		DurationHours:  durationHours,
	}
}

// Route Implements Msg.
func (msg MsgUpdateRateLimitLegacy) Route() string { return types.RouterKey }

// Type Implements Msg.
func (msg MsgUpdateRateLimitLegacy) Type() string { return TypeMsgUpdateRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateRateLimitLegacy) GetSignBytes() []byte {
	return sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgUpdateRateLimitLegacy) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateRateLimitLegacy) ValidateBasic() error {
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
