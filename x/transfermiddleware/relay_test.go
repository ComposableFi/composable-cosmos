package transfermiddleware_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	customibctesting "github.com/notional-labs/banksy/v2/app/ibctesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromIBCTransferToContract(t *testing.T) {
	// scenario: given two chains,
	//           with a contract on chain B
	//           then the contract can handle the receiving side of an ics20 transfer
	//           that was started on chain A via ibc transfer module

	transferAmount := sdk.NewInt(1)
	specs := map[string]struct {
		expChainABalanceDiff sdk.Int
		expChainBBalanceDiff sdk.Int
	}{
		"ack": {
			expChainABalanceDiff: transferAmount.Neg(),
			expChainBBalanceDiff: transferAmount,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			var (
				coordinator = customibctesting.NewCoordinator(t, 2)
				chainA      = coordinator.GetChain(customibctesting.GetChainID(0))
				chainB      = coordinator.GetChain(customibctesting.GetChainID(1))
			)
			coordinator.CommitBlock(chainA, chainB)

			path := customibctesting.NewPath(chainA, chainB)
			path.EndpointA.ChannelConfig = &ibctesting.ChannelConfig{
				PortID:  "transfer",
				Version: ibctransfertypes.Version,
				Order:   channeltypes.UNORDERED,
			}
			path.EndpointB.ChannelConfig = &ibctesting.ChannelConfig{
				PortID:  "transfer",
				Version: ibctransfertypes.Version,
				Order:   channeltypes.UNORDERED,
			}

			coordinator.SetupConnections(path)
			coordinator.CreateChannels(path)
			// Setup chainB
			chainBtransMiddleware := chainB.TransferMiddleware()
			err := chainBtransMiddleware.AddParachainIBCTokenInfo(chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
			require.NoError(t, err)
			// Setup chainA
			originalChainABalance := chainA.Balance(chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
			originalChainBBalance := chainB.Balance(chainB.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
			// when transfer via sdk transfer from A (module) -> B (contract)
			coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
			timeoutHeight := clienttypes.NewHeight(1, 110)
			msg := ibctransfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, coinToSendToB, chainA.SenderAccount.GetAddress().String(), chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
			_, err = chainA.SendMsgs(msg)
			require.NoError(t, err)
			require.NoError(t, path.EndpointB.UpdateClient())

			// then
			require.Equal(t, 1, len(chainA.PendingSendPackets))
			require.Equal(t, 0, len(chainB.PendingSendPackets))

			// and when relay to chain B and handle Ack on chain A
			err = coordinator.RelayAndAckPendingPackets(path)
			require.NoError(t, err)

			// then
			require.Equal(t, 0, len(chainA.PendingSendPackets))
			require.Equal(t, 0, len(chainB.PendingSendPackets))

			// and source chain balance was decreased
			newChainABalance := chainA.Balance(chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
			assert.Equal(t, originalChainABalance.Amount.Add(spec.expChainABalanceDiff), newChainABalance.Amount)

			// and dest chain balance contains voucher
			expBalance := originalChainBBalance.Add(sdk.NewCoin(sdk.DefaultBondDenom, transferAmount))
			gotBalance := chainB.Balance(chainB.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
			assert.Equal(t, expBalance, gotBalance)
		})
	}
}
