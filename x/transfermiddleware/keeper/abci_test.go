package keeper_test

import (
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
)

func (suite *TransferMiddlewareKeeperTestSuite) TestBeginBlocker() {
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

	suite.app.TransferMiddlewareKeeper.BeginBlocker(suite.ctx)

	countRemove := 0
	suite.app.TransferMiddlewareKeeper.IterateRemoveListInfo(suite.ctx, func(removeList types.RemoveParachainIBCTokenInfo) (stop bool) {
		if suite.ctx.BlockTime().After(removeList.RemoveTime) {
			countRemove++
			count--
		}
		return false
	})

	suite.Require().Equal(countRemove+count, 5)
}
