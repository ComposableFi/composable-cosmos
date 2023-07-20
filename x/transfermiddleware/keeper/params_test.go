package keeper_test

import (
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
)

func (suite *TransferMiddlewareKeeperTestSuite) TestGetParams() {
	suite.SetupTest()
	params := types.DefaultParams()

	suite.app.TransferMiddlewareKeeper.SetParams(suite.ctx, params)

	suite.Require().Equal(params, suite.app.TransferMiddlewareKeeper.GetParams(suite.ctx))
}
