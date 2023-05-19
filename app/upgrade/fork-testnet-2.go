package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	FaucetAmount  = sdk.NewInt(1000000000000000)
	TokenDenom    = "upica"
	FaucetAddress = "banksy1wxl09qwpwe94c5702ttw5hj8h4lr4k5cleakh6"
	// TODO: need to calculate halt height from halt time.
	HaltHeight = int64(15000)
)
