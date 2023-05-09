package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

var _ sdk.Msg = &MsgAddParachainIBCTokenInfo{}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgAddParachainIBCTokenInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgAddParachainIBCTokenInfo message.
func (msg *MsgAddParachainIBCTokenInfo) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgAddParachainIBCTokenInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	// validate channelId
	if err := host.ChannelIdentifierValidator(msg.ChannelId); err != nil {
		return err
	}

	return nil
}

var _ sdk.Msg = &MsgRemoveParachainIBCTokenInfo{}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgRemoveParachainIBCTokenInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgRemoveParachainIBCTokenInfo message.
func (msg *MsgRemoveParachainIBCTokenInfo) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgRemoveParachainIBCTokenInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	return nil
}

func NewMsgAddParachainIBCTokenInfo(
	authority string,
	ibcDenom string,
	channelId string,
	nativeDenom string,
) *MsgAddParachainIBCTokenInfo {
	return &MsgAddParachainIBCTokenInfo{
		Authority:   authority,
		IbcDenom:    ibcDenom,
		ChannelId:   channelId,
		NativeDenom: nativeDenom,
	}
}
