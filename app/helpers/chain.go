package testutil

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/testing/mock"
	legacyibctesting "github.com/cosmos/interchain-security/v3/legacy_ibc_testing/testing"
	"github.com/notional-labs/centauri/v6/app/ibctesting/simapp"
	"github.com/stretchr/testify/require"
)

const DefaultDenom = "ppica"

// NewTestChainWithValSet copypasted and modified to use neutron denom from here https://github.com/cosmos/ibc-go/blob/af9b461c63274b9ce5917beb89a2c92e865419df/testing/chain.go#L94
func NewTestChainWithValSet(t *testing.T, coord *legacyibctesting.Coordinator, appIniter legacyibctesting.AppIniter, chainID string, valSet *tmtypes.ValidatorSet, signers map[string]tmtypes.PrivValidator) *legacyibctesting.TestChain {
	genAccs := []authtypes.GenesisAccount{}
	genBals := []banktypes.Balance{}
	senderAccs := []legacyibctesting.SenderAccount{}

	// generate genesis accounts
	for i := 0; i < legacyibctesting.MaxAccounts; i++ {
		senderPrivKey := secp256k1.GenPrivKey()
		acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), uint64(i), 0)
		amount, ok := sdk.NewIntFromString("10000000000000000000")
		require.True(t, ok)

		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(DefaultDenom, amount)),
		}

		genAccs = append(genAccs, acc)
		genBals = append(genBals, balance)

		senderAcc := legacyibctesting.SenderAccount{
			SenderAccount: acc,
			SenderPrivKey: senderPrivKey,
		}

		senderAccs = append(senderAccs, senderAcc)
	}

	app := legacyibctesting.SetupWithGenesisValSet(t, appIniter, valSet, genAccs, chainID, sdk.DefaultPowerReduction, genBals...)

	// create current header and call begin block
	header := tmproto.Header{
		ChainID: chainID,
		Height:  1,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := app.GetTxConfig()

	// create an account to send transactions from
	chain := &legacyibctesting.TestChain{
		T:              t,
		Coordinator:    coord,
		ChainID:        chainID,
		App:            app,
		CurrentHeader:  header,
		QueryServer:    app.GetIBCKeeper(),
		TxConfig:       txConfig,
		Codec:          app.AppCodec(),
		Vals:           valSet,
		NextVals:       valSet,
		Signers:        signers,
		SentPackets:    make(map[string]channeltypes.Packet),
		SenderPrivKey:  senderAccs[0].SenderPrivKey,
		SenderAccount:  senderAccs[0].SenderAccount,
		SenderAccounts: senderAccs,
	}

	coord.CommitBlock(chain)

	return chain
}

// NewTestChain initializes a new test chain with a default of 4 validators
// Use this function if the tests do not need custom control over the validator set
// Copypasted from https://github.com/cosmos/ibc-go/blob/af9b461c63274b9ce5917beb89a2c92e865419df/testing/chain.go#L159
func NewTestChain(t *testing.T, coord *legacyibctesting.Coordinator, appIniter legacyibctesting.AppIniter, chainID string) *legacyibctesting.TestChain {
	// generate validators private/public key
	var (
		validatorsPerChain = 4
		validators         []*tmtypes.Validator
		signersByAddress   = make(map[string]tmtypes.PrivValidator, validatorsPerChain)
	)

	for i := 0; i < validatorsPerChain; i++ {
		privVal := mock.NewPV()
		pubKey, err := privVal.GetPubKey()
		require.NoError(t, err)
		validators = append(validators, tmtypes.NewValidator(pubKey, 1))
		signersByAddress[pubKey.Address().String()] = privVal
	}

	// construct validator set;
	// Note that the validators are sorted by voting power
	// or, if equal, by address lexical order
	valSet := tmtypes.NewValidatorSet(validators)

	return NewTestChainWithValSet(t, coord, appIniter, chainID, valSet, signersByAddress)
}

func GetSimApp(chain *legacyibctesting.TestChain) *simapp.SimApp {
	app, ok := chain.App.(*simapp.SimApp)
	require.True(chain.T, ok)

	return app
}
