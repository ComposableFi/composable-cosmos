package types

import (
	"fmt"

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
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}

	if err := m.Amount.Validate(); err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if !m.Amount.IsValid() {
		return fmt.Errorf("invalid coin amount: %s", m.Amount)
	}

	if m.Amount.IsZero() {
		return fmt.Errorf("amount cannot be zero")
	}

	if m.Amount.IsAnyNegative() {
		return fmt.Errorf("amount cannot be negative: %s", m.Amount)
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
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}

	if _, err := sdk.AccAddressFromBech32(m.AllowedAddress); err != nil {
		return fmt.Errorf("invalid allowed address: %w", err)
	}

	return nil
}

func NewMsgAddAccountToFundModuleSet(authority, allowedAddress string) *MsgAddAccountToFundModuleSet {
	return &MsgAddAccountToFundModuleSet{
		Authority:      authority,
		AllowedAddress: allowedAddress,
	}
}

var _ sdk.Msg = &MsgUpdateParams{}

// Route Implements Msg.
func (m MsgUpdateParams) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgUpdateParams) Type() string { return sdk.MsgTypeURL(&m) }

// GetSigners returns the expected signers for a MsgMintAndAllocateExp .
func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	daoAccount, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{daoAccount}
}

// GetSignBytes Implements Msg.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&m))
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}

	if err := m.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}
