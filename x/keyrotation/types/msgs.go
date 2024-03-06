package types

import (
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRotateConsPubKey = "rotate_cons_pubkey"
)

var (
	_ sdk.Msg = &MsgRotateConsPubKey{}
)

func NewMsgRotateConsPubKey(
	valAddr sdk.ValAddress,
	pubKey cryptotypes.PubKey,
) (*MsgRotateConsPubKey, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgRotateConsPubKey{
		ValidatorAddress: valAddr.String(),
		Pubkey:           pkAny,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgRotateConsPubKey) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgRotateConsPubKey) Type() string { return TypeMsgRotateConsPubKey }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgRotateConsPubKey) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)

	valAccAddr := sdk.AccAddress(valAddr)

	return []sdk.AccAddress{valAccAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgRotateConsPubKey) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgRotateConsPubKey) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.Pubkey == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty pubkey")
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgRotateConsPubKey) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}
