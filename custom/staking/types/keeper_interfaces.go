package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MintKeeper interface {
	SetLastTotalPower(ctx sdk.Context, power math.Int)
}
