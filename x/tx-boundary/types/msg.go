package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateDelegateBoundary{}

// Route Implements Msg.
func (m MsgUpdateDelegateBoundary) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgUpdateDelegateBoundary) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgUpdateDelegateBoundary) GetSigners() []sdk.AccAddress {
	authorityAddr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authorityAddr}
}

// GetSignBytes Implements Msg.
func (m MsgUpdateDelegateBoundary) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgUpdateDelegateBoundary) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	return nil
}

func NewMsgUpdateDelegateBoundary(boundary Boundary, authority string) *MsgUpdateDelegateBoundary {
	return &MsgUpdateDelegateBoundary{
		Authority: authority,
		Boundary:  boundary,
	}
}

var _ sdk.Msg = &MsgUpdateRedelegateBoundary{}

// Route Implements Msg.
func (m MsgUpdateRedelegateBoundary) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgUpdateRedelegateBoundary) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgUpdateRedelegateBoundary) GetSigners() []sdk.AccAddress {
	authorityAddr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authorityAddr}
}

// GetSignBytes Implements Msg.
func (m MsgUpdateRedelegateBoundary) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgUpdateRedelegateBoundary) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return nil
}

func NewMsgUpdateRedelegateBoundary(boundary Boundary, authority string) *MsgUpdateRedelegateBoundary {
	return &MsgUpdateRedelegateBoundary{
		Authority: authority,
		Boundary:  boundary,
	}
}
