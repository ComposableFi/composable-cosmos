package interchaintest

import (
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
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
		SkipGenTx:           false,
		PreGenesis:          nil,
		ModifyGenesis:       nil,
		ConfigFileOverrides: nil,
		EncodingConfig:      banksyEncoding(),
	}

	pathBanksyPicasso   = "banksy-picasso"
	genesisWalletAmount = int64(10_000_000)
)

// banksyEncoding registers the banksy specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func banksyEncoding() *simappparams.EncodingConfig {
	cfg := cosmos.DefaultEncoding()
	// register custom types

	return &cfg
}
