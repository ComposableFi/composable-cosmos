package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrValidationMsg  = errorsmod.Register(ModuleName, 1, "invalid msg")
	ErrInvalidCoin    = errorsmod.Register(ModuleName, 2, "invalid coin")
	ErrInvalidAddress = errorsmod.Register(ModuleName, 3, "invalid address")
)
