package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var KeyRemoveParaTokenDuration = []byte("RemoveParaTokenDuration")

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(removeDuration time.Duration) Params {
	return Params{
		Duration: removeDuration,
	}
}

// DefaultParams is the default parameter configuration for the transfermiddleware module.
func DefaultParams() Params {
	return NewParams(168 * time.Hour) // 168h (1 week)
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRemoveParaTokenDuration, &p.Duration, validateTimeDuration),
	}
}

func validateTimeDuration(i interface{}) error {
	_, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
