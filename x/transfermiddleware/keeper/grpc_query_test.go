package keeper_test

import (
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
)

func (suite *TransferMiddlewareKeeperTestSuite) TestParaTokenInfo() {
	suite.SetupTest()
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

	info, err := suite.app.TransferMiddlewareKeeper.ParaTokenInfo(suite.ctx, &types.QueryParaTokenInfoRequest{NativeDenom: "pica"})

	suite.Require().NoError(err)
	suite.Require().Equal("1", info.AssetId)
	suite.Require().Equal("pica", info.NativeDenom)
	suite.Require().Equal("ibc-test", info.IbcDenom)
	suite.Require().Equal("channel-0", info.ChannelId)
}

func (suite *TransferMiddlewareKeeperTestSuite) TestEscrowAddress() {
	suite.SetupTest()

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

	escrowResponse, err := suite.app.TransferMiddlewareKeeper.EscrowAddress(suite.ctx, &types.QueryEscrowAddressRequest{ChannelId: "channel-0"})
	suite.Require().NoError(err)
	suite.Require().Equal(escrowResponse.EscrowAddress, transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-0").String())
}
