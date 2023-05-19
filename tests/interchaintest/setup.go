package interchaintest

import (
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

var (
	BanksyMainRepo = "ghcr.io/notional-labs/banksy"

	BanksyImage = ibc.DockerImage{
		Repository: "ghcr.io/notional-labs/banksy",
		Version:    "2.0.1",
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
