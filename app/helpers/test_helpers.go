package helpers

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	abcitypes1 "github.com/cometbft/cometbft/proto/tendermint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banksy "github.com/notional-labs/banksy/v2/app"
	"github.com/stretchr/testify/require"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "fee-app"
)

// DefaultConsensusParams defines the default Tendermint consensus params used
// in feeapp testing.
var DefaultConsensusParams = &abcitypes1.ConsensusParams{
	Block: &abcitypes1.BlockParams{
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

func (EmptyAppOptions) Get(o string) interface{} { return nil }

func NewContextForApp(app banksy.BanksyApp) sdk.Context {
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  1,
	})
	return ctx
}

func Setup(t *testing.T, isCheckTx bool, invCheckPeriod uint) *banksy.BanksyApp {
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

func setup(withGenesis bool, invCheckPeriod uint) (*banksy.BanksyApp, banksy.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := banksy.MakeEncodingConfig()
	app := banksy.NewBanksyApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		banksy.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		EmptyAppOptions{},
	)
	if withGenesis {
		return app, banksy.NewDefaultGenesisState()
	}

	return app, banksy.GenesisState{}
}
