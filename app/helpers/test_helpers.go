package helpers

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/ibc-go/v7/testing/mock"
	centauri "github.com/notional-labs/centauri/v4/app"
	"github.com/stretchr/testify/require"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "fee-app"
)

// DefaultConsensusParams defines the default Tendermint consensus params used
// in feeapp testing.
var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type EmptyAppOptions struct{}

func (EmptyAppOptions) Get(_ string) interface{} { return nil }

func NewContextForApp(app centauri.CentauriApp) sdk.Context {
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  1,
	})
	return ctx
}

func Setup(t *testing.T, isCheckTx bool, invCheckPeriod uint) *centauri.CentauriApp {
	t.Helper()
	app, genesisState := setup(!isCheckTx, invCheckPeriod)
	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func setup(withGenesis bool, invCheckPeriod uint, opts ...wasm.Option) (*centauri.CentauriApp, centauri.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := centauri.MakeEncodingConfig()
	app := centauri.NewCentauriApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		wasmtypes.EnableAllProposals,
		map[int64]bool{},
		centauri.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		EmptyAppOptions{},
		opts,
	)
	if withGenesis {
		return app, centauri.NewDefaultGenesisState()
	}

	return app, centauri.GenesisState{}
}

// SetupWithGenesisValSet initializes a new CentauriApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the CentauriApp from first genesis
// account. A Nop logger is set in CentauriApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, chainID string, balances ...banktypes.Balance) *centauri.CentauriApp {
	t.Helper()
	app, genesisState := setup(true, 5)
	genesisState, err := simtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)
	baseapp.SetChainID(chainID)(app.GetBaseApp())

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         chainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	// app.BaseApp.Set
	app.Commit()
	app.BeginBlock(
		abci.RequestBeginBlock{
			Header: tmproto.Header{
				ChainID:            chainID,
				Height:             app.LastBlockHeight() + 1,
				AppHash:            app.LastCommitID().Hash,
				ValidatorsHash:     valSet.Hash(),
				NextValidatorsHash: valSet.Hash(),
			},
		},
	)

	return app
}

func SetupCentauriAppWithValSet(t *testing.T) *centauri.CentauriApp {
	t.Helper()
	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	amount, ok := sdk.NewIntFromString("10000000000000000000")
	require.True(t, ok)

	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount)),
	}

	centauriApp := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, "notional", balance)
	return centauriApp
}
