package types_test

import (
	"testing"

	"github.com/notional-labs/centauri/v4/x/transfermiddleware/types"
	"github.com/stretchr/testify/require"
)

func TestValidateBasic(t *testing.T) {
	var (
		validInfo = types.ParachainIBCTokenInfo{
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "1",
		}
		invalidInfo = types.ParachainIBCTokenInfo{
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "asset-1",
		}
	)
	testCases := map[string]struct {
		info types.ParachainIBCTokenInfo

		expectedErr bool
	}{
		"valid parachain info": {
			info: validInfo,
		},
		"invalid parachain info": {
			info:        invalidInfo,
			expectedErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.info.ValidateBasic()

			// Assertions.
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
