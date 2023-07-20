package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/centauri/v3/app"
	helpers "github.com/notional-labs/centauri/v3/app/helpers"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
	"github.com/stretchr/testify/require"
)

func SetupTest(t *testing.T) (*app.CentauriApp, sdk.Context) {
	app := helpers.SetupCentauriAppWithValSet(t)
	ctx := helpers.NewContextForApp(*app)
	tokenInfos := make([]types.ParachainIBCTokenInfo, 1)
	tokenInfos[0] = types.ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test",
		ChannelId:   "channel-0",
		NativeDenom: "pica",
		AssetId:     "1",
	}

	app.TransferMiddlewareKeeper.InitGenesis(ctx, types.GenesisState{
		TokenInfos: tokenInfos,
	})

	return app, ctx
}

func TestAddParachainIBCInfo(t *testing.T) {
	app, ctx := SetupTest(t)
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
		info types.ParachainIBCTokenInfo
		expectedErr error
		expectedPass bool
	}{
		"valid parachain info": {
			info: validInfo,
			expectedErr: nil,
			expectedPass: true,
		},
		"duplicate asset ID": {
			info: duplAssetId,
			expectedErr: types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate IBC denom": {
			info: duplIBCDenom,
			expectedErr: types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate native denom": {
			info: dup1NativeDenom,
			expectedErr: types.ErrMultipleMapping,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := app.TransferMiddlewareKeeper.AddParachainIBCInfo(
				ctx, 
				tc.info.IbcDenom,
				tc.info.ChannelId,
				tc.info.NativeDenom,
				tc.info.AssetId,
			)

			if tc.expectedPass {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, tc.expectedErr, err)
			}
		})
	}
}

func TestAddParachainIBCInfoToRemoveList(t *testing.T) {
	app, ctx := SetupTest(t)

	validInfo := types.ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test-1",
		ChannelId:   "channel-1",
		NativeDenom: "native-1",
		AssetId:     "2",
	}

	app.TransferMiddlewareKeeper.AddParachainIBCInfo(
		ctx, 
		validInfo.IbcDenom,
		validInfo.ChannelId,
		validInfo.NativeDenom,
		validInfo.AssetId,
	)

	testCases := map[string]struct {
		denom string
		expectedErr error
		expectedPass bool
	}{
		"valid denom": {
			denom: "pica",
			expectedErr: nil,
			expectedPass: true,
		},
		"not existed denom": {
			denom: "native-xxxxx",
			expectedErr: sdkerrors.ErrKeyNotFound,
			expectedPass: false,
		},
	}
	params := app.TransferMiddlewareKeeper.GetParams(ctx)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			time, err := app.TransferMiddlewareKeeper.AddParachainIBCInfoToRemoveList(ctx, tc.denom)
			removeTime := ctx.BlockTime().Add(params.Duration)

			if tc.expectedPass {
				require.Equal(t, removeTime, time)
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}

func TestIterateRemoveListInfo(t *testing.T) {
	app, ctx := SetupTest(t)

	infos := [5]types.ParachainIBCTokenInfo{
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
		{
			IbcDenom:    "ibc-test-3",
			ChannelId:   "channel-1",
			NativeDenom: "native-3",
			AssetId:     "4",
		},
		{
			IbcDenom:    "ibc-test-4",
			ChannelId:   "channel-2",
			NativeDenom: "native-4",
			AssetId:     "5",
		},
		{
			IbcDenom:    "ibc-test-5",
			ChannelId:   "channel-2",
			NativeDenom: "native-5",
			AssetId:     "6",
		},
	}

	count := 0
	for _, info := range infos {
		app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			ctx, 
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}

	app.TransferMiddlewareKeeper.IterateRemoveListInfo(ctx, func(_ types.RemoveParachainIBCTokenInfo) (stop bool) {
		count++
		return false
	})

	require.Equal(t, count, 0)
	
	for _, info := range infos {
		app.TransferMiddlewareKeeper.AddParachainIBCInfoToRemoveList(
			ctx, 
			info.NativeDenom,
		)
	}

	app.TransferMiddlewareKeeper.IterateRemoveListInfo(ctx, func(_ types.RemoveParachainIBCTokenInfo) (stop bool) {
		count++
		return false
	})

	require.Equal(t, count, 5)
}

func TestRemoveParachainIBCInfo(t *testing.T) {
	app, ctx := SetupTest(t)

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
		denom string
		expectedErr error
		expectedPass bool
	}{
		"valid denom": {
			denom: "native-1",
			expectedErr: nil,
			expectedPass: true,
		},
		"not existed denom": {
			denom: "native-xxxxx",
			expectedErr: types.NotRegisteredNativeDenom,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := app.TransferMiddlewareKeeper.RemoveParachainIBCInfo(ctx, tc.denom)

			if tc.expectedPass {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}

func TestAllowRlyAddress(t *testing.T) {
	app, ctx := SetupTest(t)
	allowedAddress := "allowed"
	notAllowedAddress := "not_allowed"
	app.TransferMiddlewareKeeper.SetAllowRlyAddress(ctx, allowedAddress)
	testCases := map[string]struct {
		address string
		expectedRes bool
	}{
		"allowed address": {
			address: allowedAddress,
			expectedRes: true,
		},
		"not allowed address": {
			address: notAllowedAddress,
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := app.TransferMiddlewareKeeper.HasAllowRlyAddress(ctx, tc.address)

			require.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestHasParachainIBCInfoByNativeDenom(t *testing.T) {
	app, ctx := SetupTest(t)

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
		denom string
		expectedRes bool
	}{
		"has info by default native denom": {
			denom: "pica",
			expectedRes: true,
		},
		"has info by native denom 1": {
			denom: "native-1",
			expectedRes: true,
		},
		"has info by native denom 2": {
			denom: "native-1",
			expectedRes: true,
		},
		"not have info by native denom": {
			denom: "native-xxxxx",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := app.TransferMiddlewareKeeper.HasParachainIBCTokenInfoByNativeDenom(ctx, tc.denom)

			require.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestHasParachainIBCInfoByAssetID(t *testing.T) {
	app, ctx := SetupTest(t)

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
		assetID string
		expectedRes bool
	}{
		"has info by default asset ID": {
			assetID: "1",
			expectedRes: true,
		},
		"has info by asset ID 2": {
			assetID: "2",
			expectedRes: true,
		},
		"has info by asset ID 3": {
			assetID: "3",
			expectedRes: true,
		},
		"not have info by asset ID": {
			assetID: "4",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := app.TransferMiddlewareKeeper.HasParachainIBCTokenInfoByAssetID(ctx, tc.assetID)

			require.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestGetParachainIBCTokenInfo(t *testing.T) {
	app, ctx := SetupTest(t)

	infos := map[string]types.ParachainIBCTokenInfo{
		"2": {
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-2",
			AssetId:     "2",
		},
		"3": {
			IbcDenom:    "ibc-test-2",
			ChannelId:   "channel-1",
			NativeDenom: "native-3",
			AssetId:     "3",
		},
		"4": {
			IbcDenom:    "ibc-test-3",
			ChannelId:   "channel-1",
			NativeDenom: "native-4",
			AssetId:     "5",
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
		denom string
		assetID string
		expectedRes bool
	}{
		"valid info of token 2": {
			denom: "native-2",
			assetID: "2",
			expectedRes: true,
		},
		"valid info of token 3": {
			denom: "native-3",
			assetID: "3",
			expectedRes: true,
		},
		"not have info": {
			denom: "native-xxxxx",
			assetID: "4",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tokenInfoRaw := infos[tc.assetID]
			tokenInfoByDenom := app.TransferMiddlewareKeeper.GetParachainIBCTokenInfoByNativeDenom(ctx, tc.denom)
			tokenInfoByAssetID := app.TransferMiddlewareKeeper.GetParachainIBCTokenInfoByAssetID(ctx, tc.assetID)

			if (tc.expectedRes) {
				require.Equal(t, tokenInfoRaw, tokenInfoByDenom)
				require.Equal(t, tokenInfoRaw, tokenInfoByAssetID)
			} else {
				require.NotEqual(t, tokenInfoRaw, tokenInfoByDenom)
				require.NotEqual(t, tokenInfoRaw, tokenInfoByAssetID)
			}
		})
	}
}

func TestGetNativeDenomByIBCDenomSecondaryIndex(t *testing.T) {
	app, ctx := SetupTest(t)

	infos := map[string]types.ParachainIBCTokenInfo{
		"2": {
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-2",
			AssetId:     "2",
		},
		"3": {
			IbcDenom:    "ibc-test-2",
			ChannelId:   "channel-1",
			NativeDenom: "native-3",
			AssetId:     "3",
		},
		"4": {
			IbcDenom:    "ibc-test-3",
			ChannelId:   "channel-1",
			NativeDenom: "native-4",
			AssetId:     "5",
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
		ibcDenom string
		assetID string
		expectedRes bool
	}{
		"valid info of token 2": {
			ibcDenom: "ibc-test-1",
			assetID: "2",
			expectedRes: true,
		},
		"valid info of token 3": {
			ibcDenom: "ibc-test-2",
			assetID: "3",
			expectedRes: true,
		},
		"not have info": {
			ibcDenom: "ibc-test-xxxxx",
			assetID: "4",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			nativeDenomRaw := infos[tc.assetID].NativeDenom
			nativeDenomByIbcDenom := app.TransferMiddlewareKeeper.GetNativeDenomByIBCDenomSecondaryIndex(ctx, tc.ibcDenom)

			if (tc.expectedRes) {
				require.Equal(t, nativeDenomRaw, nativeDenomByIbcDenom)
			} else {
				require.NotEqual(t, nativeDenomRaw, nativeDenomByIbcDenom)
			}
		})
	}
}

func TestLogger(t *testing.T) {
	app, ctx := SetupTest(t)
	require.Equal(t, ctx.Logger().With("module", "x/ibc-transfermiddleware"), app.TransferMiddlewareKeeper.Logger(ctx))
}

// TODO: TestGetTotalEscrowedToken with IBC relay