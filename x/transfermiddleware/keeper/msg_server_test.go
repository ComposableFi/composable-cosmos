package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/centauri/v4/x/transfermiddleware/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// TODO: change hard code address
const (
	authorityAddr = "centauri10556m38z4x6pqalr9rl5ytf3cff8q46nk85k9m"
)

func (suite *TransferMiddlewareKeeperTestSuite) TestMsgAddParachainIBCInfo() {
	suite.SetupTest()
	t := suite.T()

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
			_, err := suite.msgServer.AddParachainIBCTokenInfo(
				suite.ctx,
				msg,
			)

			if tc.expectedPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorIs(tc.expectedErr, err)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestRemoveParachainIBCTokenInfo() {
	suite.SetupTest()
	t := suite.T()

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

			_, err := suite.msgServer.RemoveParachainIBCTokenInfo(
				suite.ctx,
				msg,
			)

			if tc.expectedPass {
				suite.Require().NoError(err)
				found := false

				suite.app.TransferMiddlewareKeeper.IterateRemoveListInfo(suite.ctx, func(removeInfo types.RemoveParachainIBCTokenInfo) (stop bool) {
					if removeInfo.NativeDenom == tc.denom {
						found = true
					}
					return false
				})

				suite.Require().True(found)
			} else {
				suite.Require().ErrorIs(tc.expectedErr, err)
			}
		})
	}
}

func (suite *TransferMiddlewareKeeperTestSuite) TestMsgAddRlyAddress() {
	suite.SetupTest()
	t := suite.T()

	notAuthorityAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	allowedAddress := "allowed"
	notAllowedAddress := "not_allowed"
	msg := types.NewMsgAddRlyAddress(
		authorityAddr,
		authorityAddr,
	)

	_, err := suite.msgServer.AddRlyAddress(
		suite.ctx,
		msg,
	)
	suite.Require().NoError(err)

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

			_, err := suite.msgServer.AddRlyAddress(
				suite.ctx,
				msg,
			)
			if tc.expectedPass {
				suite.Require().NoError(err)
				res := suite.app.TransferMiddlewareKeeper.HasAllowRlyAddress(suite.ctx, tc.address)
				suite.Require().True(res)
			} else {
				suite.Require().Error(tc.expectedErr, err)
			}
		})
	}
}
