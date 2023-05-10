package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	testCases := []struct {
		msg          string
		malleate     func()
		recvIsSource bool // the receiving chain is the source of the coin originally
		expPass      bool
	}{
		{
			"success receive on source chain",
			func() {},
			true,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			path := NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			if tc.recvIsSource {
				transferAmount := sdk.NewInt(1)
				coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
				timeoutHeight := clienttypes.NewHeight(1, 110)
				msg := ibctransfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, coinToSendToB, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
				_, err := suite.chainA.SendMsgs(msg)
				suite.Require().NoError(err)
				suite.Require().NoError(path.EndpointB.UpdateClient())

				// then
				suite.Require().Equal(1, len(suite.chainA.PendingSendPackets))
				suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

				// and when relay to chain B and handle Ack on chain A
				err = suite.coordinator.RelayAndAckPendingPackets(path)
				suite.Require().NoError(err)

				// then
				suite.Require().Equal(0, len(suite.chainA.PendingSendPackets))
				suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))
			}
		})
	}
}
