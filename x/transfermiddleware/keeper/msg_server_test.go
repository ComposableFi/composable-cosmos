package keeper_test

import (
	"testing"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/keeper"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddParachainIBCInfo(t *testing.T) {
	app, ctx := SetupTest(t)
	msgServer := keeper.NewMsgServerImpl(app.TransferMiddlewareKeeper)
	// TODO: change hard code address
	authorityAddr := "centauri10556m38z4x6pqalr9rl5ytf3cff8q46nk85k9m"
	var (
		validInfo = types.ParachainIBCTokenInfo{
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "2",
		}
		duplAssetId = types.ParachainIBCTokenInfo{
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "1",
		}
		duplIBCDenom = types.ParachainIBCTokenInfo{
			IbcDenom:    "ibc-test",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "2",
		}
		dup1NativeDenom = types.ParachainIBCTokenInfo{
			IbcDenom:    "ibc-test",
			ChannelId:   "channel-1",
			NativeDenom: "pica",
			AssetId:     "2",
		}
	)
	var notAuthorityAddr = authtypes.NewModuleAddress(types.ModuleName).String()

	testCases := map[string]struct {
		info types.ParachainIBCTokenInfo
		fromAddress string
		expectedErr error
		expectedPass bool
	}{
		"valid parachain info": {
			info: validInfo,
			fromAddress: authorityAddr,
			expectedErr: nil,
			expectedPass: true,
		},
		"not authority address": {
			info: duplAssetId,
			fromAddress: notAuthorityAddr,
			expectedErr: govtypes.ErrInvalidSigner,
			expectedPass: false,
		},
		"duplicate asset ID": {
			info: duplAssetId,
			fromAddress: authorityAddr,
			expectedErr: types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate IBC denom": {
			info: duplIBCDenom,
			fromAddress: authorityAddr,
			expectedErr: types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate native denom": {
			info: dup1NativeDenom,
			fromAddress: authorityAddr,
			expectedErr: types.ErrMultipleMapping,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := types.NewMsgAddParachainIBCTokenInfo(
				tc.fromAddress,
				tc.info.IbcDenom,
				tc.info.NativeDenom,
				tc.info.AssetId,
				tc.info.ChannelId,
			)
			_, err := msgServer.AddParachainIBCTokenInfo(
				ctx, 
				msg,
			)

			if tc.expectedPass {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, tc.expectedErr, err)
			}
		})
	}
}