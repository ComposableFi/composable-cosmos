package interchaintest

import (
	"os"
	"strings"

	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

var (
	BanksyMainRepo   = "ghcr.io/notional-labs/banksy"
	BanksyICTestRepo = "ghcr.io/notional-labs/banksy-ictest"

	repo, version = GetDockerImageInfo()

	BanksyImage = ibc.DockerImage{
		Repository: repo,
		Version:    version,
		UidGid:     "1025:1025",
	}

	banksyConfig = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "banksy",
		ChainID:             "banksy-2",
		Images:              []ibc.DockerImage{BanksyImage},
		Bin:                 "banksyd",
		Bech32Prefix:        "banksy",
		Denom:               "stake",
		CoinType:            "118",
		GasPrices:           "0.0stake",
		GasAdjustment:       1.1,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ModifyGenesis:       nil,
		ConfigFileOverrides: nil,
	}
)

// GetDockerImageInfo returns the appropriate repo and branch version string for integration with the CI pipeline.
// The remote runner sets the BRANCH_CI env var. If present, interchaintest will use the docker image pushed up to the repo.
// If testing locally, user should run `make docker-build-debug` and interchaintest will use the local image.
func GetDockerImageInfo() (repo, version string) {
	branchVersion, found := os.LookupEnv("BRANCH_CI")
	repo = BanksyICTestRepo
	if !found {
		// make local-image
		repo = "banksy"
		branchVersion = "debug"
	}

	// github converts / to - for pushed docker images
	branchVersion = strings.ReplaceAll(branchVersion, "/", "-")
	return repo, branchVersion
}
