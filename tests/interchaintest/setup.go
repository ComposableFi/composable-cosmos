package interchaintest

import (
	"os"
	"strings"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

var (
	CentauriMainRepo   = "ghcr.io/notional-labs/centauri"
	CentauriICTestRepo = "ghcr.io/notional-labs/centauri-ictest"

	repo, version = GetDockerImageInfo()

	IBCRelayerImage   = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion = "justin-localhost-ibc"
	CentauriImage     = ibc.DockerImage{
		Repository: repo,
		Version:    version,
		UidGid:     "1025:1025",
	}

	centauriConfig = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "centauri",
		ChainID:             "centauri-2",
		Images:              []ibc.DockerImage{CentauriImage},
		Bin:                 "centaurid",
		Bech32Prefix:        "centauri",
		Denom:               "stake",
		CoinType:            "118",
		GasPrices:           "0.0stake",
		GasAdjustment:       1.1,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ModifyGenesis:       nil,
		ConfigFileOverrides: nil,
		EncodingConfig:      centauriEncoding(),
	}
	genesisWalletAmount = int64(10_000_000)
)

// centauriEncoding registers the Centauri specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func centauriEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	wasmtypes.RegisterInterfaces(cfg.InterfaceRegistry)

	//github.com/cosmos/cosmos-sdk/types/module/testutil

	return &cfg
}

// GetDockerImageInfo returns the appropriate repo and branch version string for integration with the CI pipeline.
// The remote runner sets the BRANCH_CI env var. If present, interchaintest will use the docker image pushed up to the repo.
// If testing locally, user should run `make docker-build-debug` and interchaintest will use the local image.
func GetDockerImageInfo() (repo, version string) {
	branchVersion, found := os.LookupEnv("BRANCH_CI")
	repo = CentauriICTestRepo
	if !found {
		// make local-image
		repo = "centauri"
		branchVersion = "debug"
	}

	// github converts / to - for pushed docker images
	branchVersion = strings.ReplaceAll(branchVersion, "/", "-")
	return repo, branchVersion
}
