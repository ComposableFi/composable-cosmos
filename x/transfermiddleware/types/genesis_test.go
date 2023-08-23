package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	info1 = ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test-1",
		ChannelID:   "channel-1",
		NativeDenom: "native-1",
		AssetId:     "1",
	}
	info2 = ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test-2",
		ChannelID:   "channel-2",
		NativeDenom: "native-2",
		AssetId:     "2",
	}
	invalidAssetID = ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test-3",
		ChannelID:   "channel-3",
		NativeDenom: "native-3",
		AssetId:     "asset-3",
	}
	dup1 = ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test-4",
		ChannelID:   "channel-4",
		NativeDenom: "native-1",
		AssetId:     "1",
	}
)

func TestGenesisState_Validate(t *testing.T) {
	var (
		singleInfo = []ParachainIBCTokenInfo{
			info1,
		}
		multiInfo = []ParachainIBCTokenInfo{
			info1,
			info2,
		}
		singleInvalid = []ParachainIBCTokenInfo{
			invalidAssetID,
		}
		mixedInvalid = []ParachainIBCTokenInfo{
			info1,
			invalidAssetID,
		}
		duplicateInfos = []ParachainIBCTokenInfo{
			info1,
			dup1,
		}
	)

	testCases := map[string]struct {
		infos []ParachainIBCTokenInfo

		expectedErr bool
	}{
		"valid single parachain info": {
			infos: singleInfo,
		},
		"valid mutiple parachain infos": {
			infos: multiInfo,
		},
		"invalid single parachain info": {
			infos:       singleInvalid,
			expectedErr: true,
		},
		"invalid mutiple parachain infos": {
			infos:       mixedInvalid,
			expectedErr: true,
		},
		"duplicate parachain info": {
			infos:       duplicateInfos,
			expectedErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Setup.

			// System under test.
			err := validateTokenInfos(tc.infos)

			// Assertions.
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
