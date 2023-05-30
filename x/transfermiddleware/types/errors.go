package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrDuplicateParachainIBCTokenInfo = sdkerrors.Register(ModuleName, 1, "duplicate ParachainIBC Token Info")
	InvalidIBCDenom                   = sdkerrors.Register(ModuleName, 2, "invalid ibc denom")
	NotFungibleTokenPacketData        = sdkerrors.Register(ModuleName, 3, "not fungible token packet data")
)
