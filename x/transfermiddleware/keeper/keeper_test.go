package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/centauri/v3/app"
	helpers "github.com/notional-labs/centauri/v3/app/helpers"

	"github.com/notional-labs/centauri/v3/x/transfermiddleware/keeper"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
	"github.com/stretchr/testify/suite"
)

type TransferMiddlewareKeeperTestSuite struct {
	suite.Suite
	app       *app.CentauriApp
	ctx       sdk.Context
	msgServer types.MsgServer
}

func (suite *TransferMiddlewareKeeperTestSuite) SetupTest() {
	suite.app = helpers.SetupCentauriAppWithValSet(suite.T())
	suite.ctx = helpers.NewContextForApp(*suite.app)
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.TransferMiddlewareKeeper)

	tokenInfos := make([]types.ParachainIBCTokenInfo, 1)
	tokenInfos[0] = types.ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test",
		ChannelId:   "channel-0",
		NativeDenom: "pica",
		AssetId:     "1",
	}

	suite.app.TransferMiddlewareKeeper.InitGenesis(suite.ctx, types.GenesisState{
		TokenInfos: tokenInfos,
	})
}

func (suite *TransferMiddlewareKeeperTestSuite) TestAddParachainIBCInfo() {
	suite.SetupTest()
	t := suite.T()
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
		expectedErr  error
		expectedPass bool
	}{
		"valid parachain info": {
			info:         validInfo,
			expectedErr:  nil,
			expectedPass: true,
		},
		"duplicate asset ID": {
			info:         duplAssetId,
			expectedErr:  types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate IBC denom": {
			info:         duplIBCDenom,
			expectedErr:  types.ErrMultipleMapping,
			expectedPass: false,
		},
		"duplicate native denom": {
			info:         dup1NativeDenom,
			expectedErr:  types.ErrMultipleMapping,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
				suite.ctx,
				tc.info.IbcDenom,
				tc.info.ChannelId,
				tc.info.NativeDenom,
				tc.info.AssetId,
			)

			if tc.expectedPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorIs(tc.expectedErr, err)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestAddParachainIBCInfoToRemoveList() {
	suite.SetupTest()
	t := suite.T()
	validInfo := types.ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test-1",
		ChannelId:   "channel-1",
		NativeDenom: "native-1",
		AssetId:     "2",
	}

	suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
		suite.ctx,
		validInfo.IbcDenom,
		validInfo.ChannelId,
		validInfo.NativeDenom,
		validInfo.AssetId,
	)

	testCases := map[string]struct {
		denom        string
		expectedErr  error
		expectedPass bool
	}{
		"valid denom": {
			denom:        "pica",
			expectedErr:  nil,
			expectedPass: true,
		},
		"not existed denom": {
			denom:        "native-xxxxx",
			expectedErr:  sdkerrors.ErrKeyNotFound,
			expectedPass: false,
		},
	}
	params := suite.app.TransferMiddlewareKeeper.GetParams(suite.ctx)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			time, err := suite.app.TransferMiddlewareKeeper.AddParachainIBCInfoToRemoveList(suite.ctx, tc.denom)
			removeTime := suite.ctx.BlockTime().Add(params.Duration)

			if tc.expectedPass {
				suite.Require().Equal(removeTime, time)
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorIs(err, tc.expectedErr)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestIterateRemoveListInfo() {
	suite.SetupTest()

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
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			suite.ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}

	suite.app.TransferMiddlewareKeeper.IterateRemoveListInfo(suite.ctx, func(_ types.RemoveParachainIBCTokenInfo) (stop bool) {
		count++
		return false
	})

	suite.Require().Equal(count, 0)

	for _, info := range infos {
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfoToRemoveList(
			suite.ctx,
			info.NativeDenom,
		)
	}

	suite.app.TransferMiddlewareKeeper.IterateRemoveListInfo(suite.ctx, func(_ types.RemoveParachainIBCTokenInfo) (stop bool) {
		count++
		return false
	})

	suite.Require().Equal(count, 5)
}

func (suite *TransferMiddlewareKeeperTestSuite) TestRemoveParachainIBCInfo() {
	suite.SetupTest()
	t := suite.T()

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
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			suite.ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}
	testCases := map[string]struct {
		denom        string
		expectedErr  error
		expectedPass bool
	}{
		"valid denom": {
			denom:        "native-1",
			expectedErr:  nil,
			expectedPass: true,
		},
		"not existed denom": {
			denom:        "native-xxxxx",
			expectedErr:  types.NotRegisteredNativeDenom,
			expectedPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := suite.app.TransferMiddlewareKeeper.RemoveParachainIBCInfo(suite.ctx, tc.denom)

			if tc.expectedPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorIs(err, tc.expectedErr)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestAllowRlyAddress() {
	suite.SetupTest()
	t := suite.T()

	allowedAddress := "allowed"
	notAllowedAddress := "not_allowed"
	suite.app.TransferMiddlewareKeeper.SetAllowRlyAddress(suite.ctx, allowedAddress)
	testCases := map[string]struct {
		address     string
		expectedRes bool
	}{
		"allowed address": {
			address:     allowedAddress,
			expectedRes: true,
		},
		"not allowed address": {
			address:     notAllowedAddress,
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := suite.app.TransferMiddlewareKeeper.HasAllowRlyAddress(suite.ctx, tc.address)

			suite.Require().Equal(tc.expectedRes, res)
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestHasParachainIBCInfoByNativeDenom() {
	suite.SetupTest()
	t := suite.T()

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
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			suite.ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}
	testCases := map[string]struct {
		denom       string
		expectedRes bool
	}{
		"has info by default native denom": {
			denom:       "pica",
			expectedRes: true,
		},
		"has info by native denom 1": {
			denom:       "native-1",
			expectedRes: true,
		},
		"has info by native denom 2": {
			denom:       "native-1",
			expectedRes: true,
		},
		"not have info by native denom": {
			denom:       "native-xxxxx",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := suite.app.TransferMiddlewareKeeper.HasParachainIBCTokenInfoByNativeDenom(suite.ctx, tc.denom)

			suite.Require().Equal(tc.expectedRes, res)
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestHasParachainIBCInfoByAssetID() {
	suite.SetupTest()
	t := suite.T()

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
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			suite.ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}
	testCases := map[string]struct {
		assetID     string
		expectedRes bool
	}{
		"has info by default asset ID": {
			assetID:     "1",
			expectedRes: true,
		},
		"has info by asset ID 2": {
			assetID:     "2",
			expectedRes: true,
		},
		"has info by asset ID 3": {
			assetID:     "3",
			expectedRes: true,
		},
		"not have info by asset ID": {
			assetID:     "4",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := suite.app.TransferMiddlewareKeeper.HasParachainIBCTokenInfoByAssetID(suite.ctx, tc.assetID)

			suite.Require().Equal(tc.expectedRes, res)
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestGetParachainIBCTokenInfo() {
	suite.SetupTest()
	t := suite.T()

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
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			suite.ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}
	testCases := map[string]struct {
		denom       string
		assetID     string
		expectedRes bool
	}{
		"valid info of token 2": {
			denom:       "native-2",
			assetID:     "2",
			expectedRes: true,
		},
		"valid info of token 3": {
			denom:       "native-3",
			assetID:     "3",
			expectedRes: true,
		},
		"not have info": {
			denom:       "native-xxxxx",
			assetID:     "4",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tokenInfoRaw := infos[tc.assetID]
			tokenInfoByDenom := suite.app.TransferMiddlewareKeeper.GetParachainIBCTokenInfoByNativeDenom(suite.ctx, tc.denom)
			tokenInfoByAssetID := suite.app.TransferMiddlewareKeeper.GetParachainIBCTokenInfoByAssetID(suite.ctx, tc.assetID)

			if tc.expectedRes {
				suite.Require().Equal(tokenInfoRaw, tokenInfoByDenom)
				suite.Require().Equal(tokenInfoRaw, tokenInfoByAssetID)
			} else {
				suite.Require().NotEqual(tokenInfoRaw, tokenInfoByDenom)
				suite.Require().NotEqual(tokenInfoRaw, tokenInfoByAssetID)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestGetNativeDenomByIBCDenomSecondaryIndex() {
	suite.SetupTest()
	t := suite.T()

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
		suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			suite.ctx,
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}
	testCases := map[string]struct {
		ibcDenom    string
		assetID     string
		expectedRes bool
	}{
		"valid info of token 2": {
			ibcDenom:    "ibc-test-1",
			assetID:     "2",
			expectedRes: true,
		},
		"valid info of token 3": {
			ibcDenom:    "ibc-test-2",
			assetID:     "3",
			expectedRes: true,
		},
		"not have info": {
			ibcDenom:    "ibc-test-xxxxx",
			assetID:     "4",
			expectedRes: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			nativeDenomRaw := infos[tc.assetID].NativeDenom
			nativeDenomByIbcDenom := suite.app.TransferMiddlewareKeeper.GetNativeDenomByIBCDenomSecondaryIndex(suite.ctx, tc.ibcDenom)

			if tc.expectedRes {
				suite.Require().Equal(nativeDenomRaw, nativeDenomByIbcDenom)
			} else {
				suite.Require().NotEqual(nativeDenomRaw, nativeDenomByIbcDenom)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestLogger() {
	suite.SetupTest()
	suite.Require().Equal(suite.ctx.Logger().With("module", "x/ibc-transfermiddleware"), suite.app.TransferMiddlewareKeeper.Logger(suite.ctx))
}