package keeper_test

import (
	"github.com/notional-labs/centauri/v4/x/transfermiddleware/types"
)

func (suite *TransferMiddlewareKeeperTestSuite) TestTFMInitGenesis() {
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

	info := suite.app.TransferMiddlewareKeeper.GetParachainIBCTokenInfoByNativeDenom(suite.ctx, "pica")
	suite.Require().Equal(info, suite.app.TransferMiddlewareKeeper.GetParachainIBCTokenInfoByNativeDenom(suite.ctx, "pica"))
	suite.Require().Equal("1", info.AssetId)
	suite.Require().Equal("pica", info.NativeDenom)
	suite.Require().Equal("ibc-test", info.IbcDenom)
	suite.Require().Equal("channel-0", info.ChannelId)
}

func (suite *TransferMiddlewareKeeperTestSuite) TestTFMExportGenesis() {
	suite.SetupTest()

	err := suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(suite.ctx, "ibc-test2", "channel-0", "paca", "2")
	suite.Require().NoError(err)
	err = suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(suite.ctx, "ibc-test3", "channel-1", "poke", "3")
	suite.Require().NoError(err)
	genesis := suite.app.TransferMiddlewareKeeper.ExportGenesis(suite.ctx)

	suite.Require().Equal("1", genesis.TokenInfos[0].AssetId)
	suite.Require().Equal("pica", genesis.TokenInfos[0].NativeDenom)
	suite.Require().Equal("channel-0", genesis.TokenInfos[0].ChannelId)
	suite.Require().Equal("ibc-test", genesis.TokenInfos[0].IbcDenom)

	suite.Require().Equal("2", genesis.TokenInfos[1].AssetId)
	suite.Require().Equal("paca", genesis.TokenInfos[1].NativeDenom)
	suite.Require().Equal("channel-0", genesis.TokenInfos[1].ChannelId)
	suite.Require().Equal("ibc-test2", genesis.TokenInfos[1].IbcDenom)

	suite.Require().Equal("3", genesis.TokenInfos[2].AssetId)
	suite.Require().Equal("poke", genesis.TokenInfos[2].NativeDenom)
	suite.Require().Equal("channel-1", genesis.TokenInfos[2].ChannelId)
	suite.Require().Equal("ibc-test3", genesis.TokenInfos[2].IbcDenom)
}

func (suite *TransferMiddlewareKeeperTestSuite) TestIterateParaTokenInfos() {
	suite.SetupTest()

	err := suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(suite.ctx, "ibc-test2", "channel-0", "paca", "2")
	suite.Require().NoError(err)
	err = suite.app.TransferMiddlewareKeeper.AddParachainIBCInfo(suite.ctx, "ibc-test3", "channel-1", "poke", "3")
	suite.Require().NoError(err)

	infos := []types.ParachainIBCTokenInfo{}

	suite.app.TransferMiddlewareKeeper.IterateParaTokenInfos(suite.ctx, func(index int64, info types.ParachainIBCTokenInfo) (stop bool) {
		infos = append(infos, info)
		return false
	})

	suite.Require().Equal("1", infos[0].AssetId)
	suite.Require().Equal("pica", infos[0].NativeDenom)
	suite.Require().Equal("channel-0", infos[0].ChannelId)
	suite.Require().Equal("ibc-test", infos[0].IbcDenom)

	suite.Require().Equal("2", infos[1].AssetId)
	suite.Require().Equal("paca", infos[1].NativeDenom)
	suite.Require().Equal("channel-0", infos[1].ChannelId)
	suite.Require().Equal("ibc-test2", infos[1].IbcDenom)

	suite.Require().Equal("3", infos[2].AssetId)
	suite.Require().Equal("poke", infos[2].NativeDenom)
	suite.Require().Equal("channel-1", infos[2].ChannelId)
	suite.Require().Equal("ibc-test3", infos[2].IbcDenom)
}
