package types

import (
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgFundModuleAccount{}

// Route Implements Msg.
func (m MsgFundModuleAccount) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgFundModuleAccount) Type() string { return sdk.MsgTypeURL(&m) }

// GetSigners returns the expected signers for a MsgMintAndAllocateExp .
func (m MsgFundModuleAccount) GetSigners() []sdk.AccAddress {
	daoAccount, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{daoAccount}
}

// GetSignBytes Implements Msg.
func (m MsgFundModuleAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgFundModuleAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return errorsmod.Wrap(err, "from address must be valid address")
	}
	return nil
}

func NewMsgFundModuleAccount(fromAddr sdk.AccAddress, amount sdk.Coins) *MsgFundModuleAccount {
	return &MsgFundModuleAccount{
		FromAddress: fromAddr.String(),
		Amount:      amount,
	}
}

var _ sdk.Msg = &MsgAddAccountToFundModuleSet{}

// Route Implements Msg.
func (m MsgAddAccountToFundModuleSet) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgAddAccountToFundModuleSet) Type() string { return sdk.MsgTypeURL(&m) }

// GetSigners returns the expected signers for a MsgMintAndAllocateExp .
func (m MsgAddAccountToFundModuleSet) GetSigners() []sdk.AccAddress {
	daoAccount, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{daoAccount}
}

// GetSignBytes Implements Msg.
func (m MsgAddAccountToFundModuleSet) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgAddAccountToFundModuleSet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return errorsmod.Wrap(err, "authority must be valid address")
	}

	_, err = sdk.AccAddressFromBech32(m.AllowedAddress)
	if err != nil {
		return errorsmod.Wrap(err, "allowed address must be valid address")
	}

	return nil
}

func NewMsgAddAccountToFundModuleSet(authority, allowedAddress string) *MsgAddAccountToFundModuleSet {
	return &MsgAddAccountToFundModuleSet{
		Authority:      authority,
		AllowedAddress: allowedAddress,
	}
}
