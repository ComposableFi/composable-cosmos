package interchaintest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	helpers "github.com/notional-labs/composable-testnet/tests/interchaintest/helpers"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestIBCHooks ensures the ibc-hooks middleware from osmosis works.
func TestIBCHooks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with centauri and centauri2
	numVals := 1
	numFullNodes := 0

	cfg2 := centauriConfig.Clone()
	cfg2.Name = "centauri-counterparty"
	cfg2.ChainID = "centauri-counterparty-2"

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "centauri",
			ChainConfig:   centauriConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "centauri",
			ChainConfig:   cfg2,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	const (
		path = "ibc-path"
	)

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	client, network := interchaintest.DockerSetup(t)

	centauri, centauri2 := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	relayerType, relayerName := ibc.CosmosRly, "relay"

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		relayerType,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(centauri).
		AddChain(centauri2).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  centauri,
			Chain2:  centauri2,
			Relayer: r,
			Path:    path,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, centauri, centauri2)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, centauri, centauri2)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	centauriUser, centauri2User := users[0], users[1]

	centauriUserAddr := centauriUser.FormattedAddress()

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, centauri.Config().ChainID, centauri2.Config().ChainID)
	require.NoError(t, err)

	_, contractAddr := helpers.SetupContract(t, ctx, centauri2, centauri2User.KeyName(), "bytecode/counter.wasm", `{"count":0}`)

	// do an ibc transfer through the memo to the other chain.
	transfer := ibc.WalletAmount{
		Address: contractAddr,
		Denom:   centauri.Config().Denom,
		Amount:  int64(1),
	}

	memo := ibc.TransferOptions{
		Memo: fmt.Sprintf(`{"wasm":{"contract":"%s","msg":%s}}`, contractAddr, `{"increment":{}}`),
	}

	// Initial transfer. Account is created by the wasm execute is not so we must do this twice to properly set up
	transferTx, err := centauri.SendIBCTransfer(ctx, channel.ChannelID, centauriUser.KeyName(), transfer, memo)
	require.NoError(t, err)
	centauriHeight, err := centauri.Height(ctx)
	require.NoError(t, err)

	// TODO: Remove when the relayer is fixed
	r.Flush(ctx, eRep, path, channel.ChannelID)
	_, err = testutil.PollForAck(ctx, centauri, centauriHeight-5, centauriHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Second time, this will make the counter == 1 since the account is now created.
	transferTx, err = centauri.SendIBCTransfer(ctx, channel.ChannelID, centauriUser.KeyName(), transfer, memo)
	require.NoError(t, err)
	centauriHeight, err = centauri.Height(ctx)
	require.NoError(t, err)

	// TODO: Remove when the relayer is fixed
	r.Flush(ctx, eRep, path, channel.ChannelID)
	_, err = testutil.PollForAck(ctx, centauri, centauriHeight-5, centauriHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Get the address on the other chain's side
	addr := helpers.GetIBCHooksUserAddress(t, ctx, centauri, channel.ChannelID, centauriUserAddr)
	require.NotEmpty(t, addr)

	// Get funds on the receiving chain
	funds := helpers.GetIBCHookTotalFunds(t, ctx, centauri2, contractAddr, addr)
	require.Equal(t, int(1), len(funds.Data.TotalFunds))

	var ibcDenom string
	for _, coin := range funds.Data.TotalFunds {
		if strings.HasPrefix(coin.Denom, "ibc/") {
			ibcDenom = coin.Denom
			break
		}
	}
	require.NotEmpty(t, ibcDenom)

	// ensure the count also increased to 1 as expected.
	count := helpers.GetIBCHookCount(t, ctx, centauri2, contractAddr, addr)
	require.Equal(t, int64(1), count.Data.Count)
}
