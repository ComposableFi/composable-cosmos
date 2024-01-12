package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidCoin = errorsmod.Register(ModuleName, 1, "invalid coin")
)
