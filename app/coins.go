package app

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
)

// BaseDenomUnit defines the base denomination unit for Banksy.
// 1 pica = 1x10^{BaseDenomUnit} ppica
var BaseDenomUnit = 12

// PowerReduction defines the default power reduction value for staking
var PowerReduction = sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil))
