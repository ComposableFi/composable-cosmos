package keeper_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/keeper"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

// TODO: change hard code address
const (
	authorityAddr = "centauri10556m38z4x6pqalr9rl5ytf3cff8q46nk85k9m"
)

func setupMsgServer(k keeper.Keeper) types.MsgServer {
	return keeper.NewMsgServerImpl(k)
}

func TestMsgAddParachainIBCInfo(t *testing.T) {
	app, ctx := SetupTest(t)
	msgServer := setupMsgServer(app.TransferMiddlewareKeeper)
	notAuthorityAddr := authtypes.NewModuleAddress(types.ModuleName).String()

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

	testCases := map[string]struct {
		info         types.ParachainIBCTokenInfo
		fromAddress  string
		expectedErr  error
		expectedPass bool
	}{
		"valid parachain info": {
			info:         validInfo,
			fromAddress:  authorityAddr,
			expectedErr:  nil,
			expectedPass: true,
		},
		"not authority address": {
			info:         duplAssetId,
			fromAddress:  notAuthorityAddr,
			expectedErr:  govtypes.ErrInvalidSigner,
			expectedPass: false,
		},
		"duplicate asset ID": {
			info:         duplAssetId,
			fromAddress:  authorityAddr,
			expectedErr:  types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate IBC denom": {
			info:         duplIBCDenom,
			fromAddress:  authorityAddr,
			expectedErr:  types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate native denom": {
			info:         dup1NativeDenom,
			fromAddress:  authorityAddr,
			expectedErr:  types.ErrMultipleMapping,
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

func TestRemoveParachainIBCTokenInfo(t *testing.T) {
	app, ctx := SetupTest(t)
	msgServer := setupMsgServer(app.TransferMiddlewareKeeper)
	notAuthorityAddr := authtypes.NewModuleAddress(types.ModuleName).String()

	infos := [2]types.ParachainIBCTokenInfo{
		{
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "2",
		},
		{
			IbcDenom:    "ibc-test-2",
			ChannelId:   "channel-1",
			NativeDenom: "native-2",
			AssetId:     "3",
		},
	}
	for _, info := range infos {
		app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}
	testCases := map[string]struct {
		denom        string
		fromAddress  string
		expectedErr  error
		expectedPass bool
	}{
		"not authority address": {
			denom:        "native-1",
			fromAddress:  notAuthorityAddr,
			expectedErr:  govtypes.ErrInvalidSigner,
			expectedPass: false,
		},
		"valid denom": {
			denom:        "native-1",
			fromAddress:  authorityAddr,
			expectedErr:  nil,
			expectedPass: true,
		},
		"not existed denom": {
			denom:        "native-xxxxx",
			fromAddress:  authorityAddr,
			expectedErr:  sdkerrors.ErrKeyNotFound,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := types.NewMsgRemoveParachainIBCTokenInfo(
				tc.fromAddress,
				tc.denom,
			)

			_, err := msgServer.RemoveParachainIBCTokenInfo(
				ctx,
				msg,
			)

			if tc.expectedPass {
				require.NoError(t, err)
				found := false

				app.TransferMiddlewareKeeper.IterateRemoveListInfo(ctx, func(removeInfo types.RemoveParachainIBCTokenInfo) (stop bool) {
					if removeInfo.NativeDenom == tc.denom {
						found = true
					}
					return false
				})

				require.True(t, found)
			} else {
				require.ErrorIs(t, tc.expectedErr, err)
			}
		})
	}
}

func TestMsgAddRlyAddress(t *testing.T) {
	app, ctx := SetupTest(t)
	msgServer := setupMsgServer(app.TransferMiddlewareKeeper)
	notAuthorityAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	allowedAddress := "allowed"
	notAllowedAddress := "not_allowed"
	msg := types.NewMsgAddRlyAddress(
		authorityAddr,
		authorityAddr,
	)

	_, err := msgServer.AddRlyAddress(
		ctx,
		msg,
	)
	require.NoError(t, err)

	testCases := map[string]struct {
		address      string
		fromAddress  string
		expectedErr  error
		expectedPass bool
	}{
		"allowed address": {
			address:      allowedAddress,
			fromAddress:  authorityAddr,
			expectedErr:  nil,
			expectedPass: true,
		},
		"dupl allowed address": {
			address:      authorityAddr,
			fromAddress:  authorityAddr,
			expectedErr:  types.DuplRlyAddress,
			expectedPass: false,
		},
		"not allowed address": {
			address:      notAllowedAddress,
			fromAddress:  notAuthorityAddr,
			expectedErr:  govtypes.ErrInvalidSigner,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := types.NewMsgAddRlyAddress(
				tc.fromAddress,
				tc.address,
			)

			_, err := msgServer.AddRlyAddress(
				ctx,
				msg,
			)
			if tc.expectedPass {
				require.NoError(t, err)
				res := app.TransferMiddlewareKeeper.HasAllowRlyAddress(ctx, tc.address)
				require.True(t, res)
			} else {
				require.Error(t, tc.expectedErr, err)
			}
		})
	}
}
