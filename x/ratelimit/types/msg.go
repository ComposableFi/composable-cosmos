package types

import (
	errorsmod "cosmossdk.io/errors"
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
	maxPercentSend sdk.Int,
	maxPercentRecv sdk.Int,
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
func (msg MsgAddRateLimit) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAddRateLimit) Type() string { return TypeMsgAddRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgAddRateLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
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

	// validate channelIds
	if err := host.ChannelIdentifierValidator(msg.ChannelId); err != nil {
		return err
	}

	if msg.MaxPercentSend.GT(sdk.NewInt(100)) || msg.MaxPercentSend.LT(sdk.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max-percent-send percent must be between 0 and 100 (inclusively), Provided: %v", msg.MaxPercentSend)
	}

	if msg.MaxPercentRecv.GT(sdk.NewInt(100)) || msg.MaxPercentRecv.LT(sdk.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max-percent-recv percent must be between 0 and 100 (inclusively), Provided: %v", msg.MaxPercentRecv)
	}

	if msg.MaxPercentRecv.IsZero() && msg.MaxPercentSend.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "either the max send or max receive threshold must be greater than 0")
	}

	if msg.DurationHours == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duration can not be zero")
	}

	return nil
}

var _ sdk.Msg = &MsgUpdateRateLimit{}

func NewMsgUpdateRateLimit(
	authority string,
	denom string,
	channelID string,
	maxPercentSend sdk.Int,
	maxPercentRecv sdk.Int,
	durationHours uint64,
) *MsgUpdateRateLimit {
	return &MsgUpdateRateLimit{
		Authority:      authority,
		Denom:          denom,
		ChannelId:      channelID,
		MaxPercentSend: maxPercentSend,
		MaxPercentRecv: maxPercentRecv,
		DurationHours:  durationHours,
	}
}

// Route Implements Msg.
func (msg MsgUpdateRateLimit) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgUpdateRateLimit) Type() string { return TypeMsgUpdateRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateRateLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgUpdateRateLimit) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateRateLimit) ValidateBasic() error {
	// validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	// validate channelIds
	if err := host.ChannelIdentifierValidator(msg.ChannelId); err != nil {
		return err
	}

	if msg.MaxPercentSend.GT(sdk.NewInt(100)) || msg.MaxPercentSend.LT(sdk.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max-percent-send percent must be between 0 and 100 (inclusively), Provided: %v", msg.MaxPercentSend)
	}

	if msg.MaxPercentRecv.GT(sdk.NewInt(100)) || msg.MaxPercentRecv.LT(sdk.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max-percent-recv percent must be between 0 and 100 (inclusively), Provided: %v", msg.MaxPercentRecv)
	}

	if msg.MaxPercentRecv.IsZero() && msg.MaxPercentSend.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "either the max send or max receive threshold must be greater than 0")
	}

	if msg.DurationHours == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duration can not be zero")
	}

	return nil
}

var _ sdk.Msg = &MsgRemoveRateLimit{}

func NewMsgRemoveRateLimit(
	authority string,
	denom string,
	channelID string,
) *MsgRemoveRateLimit {
	return &MsgRemoveRateLimit{
		Authority: authority,
		Denom:     denom,
		ChannelId: channelID,
	}
}

// Route Implements Msg.
func (msg MsgRemoveRateLimit) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRemoveRateLimit) Type() string { return TypeMsgRemoveRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgRemoveRateLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgRemoveRateLimit) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgRemoveRateLimit) ValidateBasic() error {
	// validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	// validate channelIds
	if err := host.ChannelIdentifierValidator(msg.ChannelId); err != nil {
		return err
	}

	return nil
}

var _ sdk.Msg = &MsgResetRateLimit{}

func NewMsgResetRateLimit(
	authority string,
	denom string,
	channelID string,
) *MsgResetRateLimit {
	return &MsgResetRateLimit{
		Authority: authority,
		Denom:     denom,
		ChannelId: channelID,
	}
}

// Route Implements Msg.
func (msg MsgResetRateLimit) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgResetRateLimit) Type() string { return TypeMsgResetRateLimit }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgResetRateLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgResetRateLimit) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgResetRateLimit) ValidateBasic() error {
	// validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	// validate channelIds
	if err := host.ChannelIdentifierValidator(msg.ChannelId); err != nil {
		return err
	}

	return nil
}
