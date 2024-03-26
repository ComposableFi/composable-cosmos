package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgAddIBCFeeConfig{}

const (
	TypeMsgAddIBCFeeConfig       = "add_config"
	TypeMsgRemoveIBCFeeConfig    = "remove_config"
	TypeMsgAddAllowedIbcToken    = "add_allowed_ibc_token"
	TypeMsgRemoveAllowedIbcToken = "remove_allowed_ibc_token"
)

func NewMsgAddIBCFeeConfig(
	authority string,
	channelID string,
	feeAddress string,
	minTimeoutTimestamp int64,
) *MsgAddIBCFeeConfig {
	return &MsgAddIBCFeeConfig{
		Authority:           authority,
		ChannelID:           channelID,
		FeeAddress:          feeAddress,
		MinTimeoutTimestamp: minTimeoutTimestamp,
	}
}

// Route Implements Msg.
func (msg MsgAddIBCFeeConfig) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAddIBCFeeConfig) Type() string { return TypeMsgAddIBCFeeConfig }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgAddIBCFeeConfig) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgAddIBCFeeConfig) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgAddIBCFeeConfig) ValidateBasic() error {
	// validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	return nil
}

var _ sdk.Msg = &MsgRemoveIBCFeeConfig{}

func NewMsgRemoveIBCFeeConfig(
	authority string,
	channelID string,
) *MsgRemoveIBCFeeConfig {
	return &MsgRemoveIBCFeeConfig{
		Authority: authority,
		ChannelID: channelID,
	}
}

// Route Implements Msg.
func (msg MsgRemoveIBCFeeConfig) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRemoveIBCFeeConfig) Type() string { return TypeMsgRemoveIBCFeeConfig }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgRemoveIBCFeeConfig) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgRemoveParachainIBCTokenInfo message.
func (msg *MsgRemoveIBCFeeConfig) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgRemoveIBCFeeConfig) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	return nil
}

var _ sdk.Msg = &MsgAddAllowedIbcToken{}

func NewMsgAddAllowedIbcToken(
	authority string,
	channelID string,
	denom string,
	amount int64,
	percentage int64,
) *MsgAddAllowedIbcToken {
	return &MsgAddAllowedIbcToken{
		Authority:  authority,
		ChannelID:  channelID,
		Denom:      denom,
		Amount:     amount,
		Percentage: percentage,
	}
}

// Route Implements Msg.
func (msg MsgAddAllowedIbcToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAddAllowedIbcToken) Type() string { return TypeMsgAddAllowedIbcToken }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgAddAllowedIbcToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgAddAllowedIbcToken) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgAddAllowedIbcToken) ValidateBasic() error {
	// validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	return nil
}

var _ sdk.Msg = &MsgRemoveAllowedIbcToken{}

func NewMsgRemoveAllowedIbcToken(
	authority string,
	channelID string,
	denom string,
) *MsgRemoveAllowedIbcToken {
	return &MsgRemoveAllowedIbcToken{
		Authority: authority,
		ChannelID: channelID,
		Denom:     denom,
	}
}

// Route Implements Msg.
func (msg MsgRemoveAllowedIbcToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRemoveAllowedIbcToken) Type() string { return TypeMsgRemoveAllowedIbcToken }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgRemoveAllowedIbcToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgRemoveParachainIBCTokenInfo message.
func (msg *MsgRemoveAllowedIbcToken) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgRemoveAllowedIbcToken) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	return nil
}
