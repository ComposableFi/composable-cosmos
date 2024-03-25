package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/ratelimit module sentinel errors
var (
	ErrChannelFeeNotFound = errorsmod.Register(ModuleName, 1, "channel fee not found for channel")
)
