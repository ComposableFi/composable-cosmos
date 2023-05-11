package keeper_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

var keyCounter uint64

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress(_ testing.TB) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress(t testing.TB) string {
	return RandomAccountAddress(t).String()
}

func submitAndExecuteProposal(t *testing.T, ctx sdk.Context, content v1.MsgExecLegacyContent) {

}

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
				// setup transfermiddleware state
				authorityAddress := suite.chainB.GetTestSupport().TransferMiddleware().GetAuthority(suite.chainB.GetContext())
				msgAddParamInfo := types.NewMsgAddParachainIBCTokenInfo()

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
				myActorAddress := RandomBech32AccountAddress(suite.T())
				// fmt.Printf("%s\n", myActorAddress)
				// suite.Require().True(false)
				// suite.chainA.GetTestSupport().TransferMiddleware().AddParachainIBCInfo(suite.chainA.GetContext())
			}
		})
	}
}
