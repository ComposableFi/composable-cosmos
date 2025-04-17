package keeper

import (
	"fmt"
	stdmath "math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/mint/types"
)

// BlockProvisionWithCheck returns the provisions for a block after checking for valid BlocksPerYear
func (k Keeper) BlockProvisionWithCheck(ctx sdk.Context, minter types.Minter, params types.Params) (sdk.Coin, error) {
	// Ensure BlocksPerYear can be safely converted to int64
	if params.BlocksPerYear > uint64(stdmath.MaxInt64) {
		return sdk.Coin{}, fmt.Errorf("blocks per year exceeds maximum int64 value")
	}

	provision, err := minter.BlockProvision(params)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("failed to calculate block provision: %w", err)
	}

	return provision, nil
}
