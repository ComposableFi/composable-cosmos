package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSetPower{}

// Route Implements Msg.
func (m MsgSetPower) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgSetPower) Type() string { return sdk.MsgTypeURL(&m) }

// GetSigners returns the expected signers for a MsgMintAndAllocateExp .
func (m MsgSetPower) GetSigners() []sdk.AccAddress {
	daoAccount, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{daoAccount}
}

// GetSignBytes Implements Msg.
func (m MsgSetPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgSetPower) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return errorsmod.Wrap(err, "from address must be valid address")
	}
	return nil
}

func NewMsgSetPower(fromAddr sdk.AccAddress) *MsgSetPower {
	return &MsgSetPower{
		FromAddress: fromAddr.String(),
	}
}
