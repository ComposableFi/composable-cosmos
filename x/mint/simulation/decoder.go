package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/notional-labs/composable/v6/x/mint/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding mint type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		if len(kvA.Key) == 0 {
			panic("invalid key length")
		}

		switch {
		case bytes.Equal(kvA.Key, types.MinterKey):
			var minterA, minterB types.Minter
			if err := cdc.Unmarshal(kvA.Value, &minterA); err != nil {
				panic(fmt.Sprintf("failed to unmarshal minter A: %v", err))
			}
			if err := cdc.Unmarshal(kvB.Value, &minterB); err != nil {
				panic(fmt.Sprintf("failed to unmarshal minter B: %v", err))
			}
			return fmt.Sprintf("%v\n%v", minterA, minterB)

		case bytes.Equal(kvA.Key, types.ParamsKey):
			var paramsA, paramsB types.Params
			if err := cdc.Unmarshal(kvA.Value, &paramsA); err != nil {
				panic(fmt.Sprintf("failed to unmarshal params A: %v", err))
			}
			if err := cdc.Unmarshal(kvB.Value, &paramsB); err != nil {
				panic(fmt.Sprintf("failed to unmarshal params B: %v", err))
			}
			return fmt.Sprintf("%v\n%v", paramsA, paramsB)

		default:
			panic(fmt.Sprintf("invalid mint key %X", kvA.Key))
		}
	}
}
