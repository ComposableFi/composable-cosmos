package interchaintest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestValidator is a basic test to accrue enough token to join active validator set, gets slashed for missing or tombstoned for double signing
func TestValidator(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx := context.Background()

	// Create chain factory with Centauri
	numVals := 5
	numFullNodes := 3

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "centauri",
			ChainConfig:   centauriConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	centauri := chains[0].(*cosmos.CosmosChain)

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().AddChain(centauri)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	err = ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,

		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	err = testutil.WaitForBlocks(ctx, 1, centauri)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 1, centauri)
	require.NoError(t, err)

	err = centauri.Validators[1].StopContainer(ctx)
	require.NoError(t, err)

	// _, _, err = centauri.Validators[1].ExecBin(ctx, "status")
	// require.Error(t, err)
	err = testutil.WaitForBlocks(ctx, 101, centauri)
	require.NoError(t, err)

	validators, err := centauri.QueryValidators(ctx)
	require.NoError(t, err)

	// slashingParams, err := centauri.QuerySlashingParams(ctx)
	require.NoError(t, err)

	fmt.Println("validators", string(validators[1].ConsensusPubkey.Value))

	var defaultTime = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	// sdk.ConsAddress(validators[1].ConsensusPubkey)
	infos, err := centauri.QuerySigningInfos(ctx)
	for _, info := range infos {
		if info.JailedUntil != defaultTime {
			fmt.Println("Jailed Validator", info.Address)
		}
	}
	require.NoError(t, err)
}
