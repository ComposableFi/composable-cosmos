package interchaintest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	//simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
)

// Spin up a banksyd chain, push a contract, and get that contract code from chain
func TestPushWasmClientCode(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	client, network := interchaintest.DockerSetup(t)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	// Override config files to support an ~2.5MB contract
	configFileOverrides := make(map[string]any)

	appTomlOverrides := make(testutil.Toml)
	configTomlOverrides := make(testutil.Toml)

	apiOverrides := make(testutil.Toml)
	apiOverrides["rpc-max-body-bytes"] = 1350000000
	appTomlOverrides["api"] = apiOverrides

	rpcOverrides := make(testutil.Toml)
	rpcOverrides["max_body_bytes"] = 1350000000
	rpcOverrides["max_header_bytes"] = 1400000000
	configTomlOverrides["rpc"] = rpcOverrides

	//mempoolOverrides := make(testutil.Toml)
	//mempoolOverrides["max_tx_bytes"] = 6000000
	//configTomlOverrides["mempool"] = mempoolOverrides

	configFileOverrides["config/app.toml"] = appTomlOverrides
	configFileOverrides["config/config.toml"] = configTomlOverrides

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{ChainConfig: ibc.ChainConfig{
			Type:    "cosmos",
			Name:    "banksy",
			ChainID: "banksyd",
			Images: []ibc.DockerImage{
				{
					Repository: "ghcr.io/notional-labs/banksy",
					Version:    "2.0.1",
					UidGid:     "1025:1025",
				},
			},
			Bin:            "banksyd",
			Bech32Prefix:   "banksy",
			Denom:          "stake",
			GasPrices:      "0.00stake",
			GasAdjustment:  1.3,
			TrustingPeriod: "504h",
			//EncodingConfig: WasmClientEncoding(),
			NoHostMount:         true,
			ConfigFileOverrides: configFileOverrides,
		},
		},
	})

	t.Logf("Calling cf.Chains")
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	banksyd := chains[0]

	t.Logf("NewInterchain")
	ic := interchaintest.NewInterchain().
		AddChain(banksyd)

	t.Logf("Interchain build options")
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  true, // Skip path creation, so we can have granular control over the process
	}))

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create and Fund User Wallets
	fundAmount := int64(100_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(fundAmount), banksyd)
	banksyd1User := users[0]

	err = testutil.WaitForBlocks(ctx, 2, banksyd)
	require.NoError(t, err)

	banksyd1UserBalInitial, err := banksyd.GetBalance(ctx, banksyd1User.FormattedAddress(), banksyd.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount, banksyd1UserBalInitial)

	err = testutil.WaitForBlocks(ctx, 2, banksyd)
	require.NoError(t, err)

	banksydChain := banksyd.(*cosmos.CosmosChain)

	codeHash, err := banksydChain.StoreClientContract(ctx, banksyd1User.KeyName(), "ics10_grandpa_cw.wasm")
	t.Logf("Contract codeHash: %s", codeHash)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 5, banksyd)
	require.NoError(t, err)

	var getCodeQueryMsgRsp GetCodeQueryMsgResponse
	err = banksydChain.QueryClientContractCode(ctx, codeHash, &getCodeQueryMsgRsp)
	codeHashByte32 := sha256.Sum256(getCodeQueryMsgRsp.Code)
	codeHash2 := hex.EncodeToString(codeHashByte32[:])
	t.Logf("Contract codeHash from code: %s", codeHash2)
	require.NoError(t, err)
	require.NotEmpty(t, getCodeQueryMsgRsp.Code)
	require.Equal(t, codeHash, codeHash2)
}

type GetCodeQueryMsgResponse struct {
	Code []byte `json:"code"`
}
