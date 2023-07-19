package ibctesting

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	tmtypes "github.com/cometbft/cometbft/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	wasmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
)

// ConstructUpdateWasmClientHeader will construct a valid 08-wasm Header with a zero height
// to update the light client on the source chain.
func (chain *TestChain) ConstructUpdateWasmClientHeader(counterparty *TestChain, clientID string) (*wasmtypes.Header, error) {
	return chain.ConstructUpdateWasmClientHeaderWithTrustedHeight(counterparty, clientID, clienttypes.ZeroHeight())
}

// ConstructUpdateWasmClientHeaderWithTrustedHeight will construct a valid 08-wasm Header
// to update the light client on the source chain.
func (chain *TestChain) ConstructUpdateWasmClientHeaderWithTrustedHeight(counterparty *TestChain, clientID string, trustedHeight clienttypes.Height) (*wasmtypes.Header, error) {
	tmHeader, err := chain.ConstructUpdateTMClientHeaderWithTrustedHeight(counterparty, clientID, trustedHeight)
	if err != nil {
		return nil, err
	}

	tmWasmHeaderData, err := chain.Codec.MarshalInterface(tmHeader)
	if err != nil {
		return nil, err
	}

	height, ok := tmHeader.GetHeight().(clienttypes.Height)
	if !ok {
		return nil, fmt.Errorf("error casting exported height to clienttypes height")
	}
	wasmHeader := wasmtypes.Header{
		Data:   tmWasmHeaderData,
		Height: height,
	}

	return &wasmHeader, nil
}

func (chain *TestChain) CreateWasmClientHeader(chainID string, blockHeight int64, trustedHeight clienttypes.Height, timestamp time.Time, tmValSet, _, tmTrustedVals *tmtypes.ValidatorSet, signers []tmtypes.PrivValidator) *wasmtypes.Header {
	tmHeader := chain.CreateTMClientHeader(chainID, blockHeight, trustedHeight, timestamp, tmValSet, tmTrustedVals, signers)
	tmWasmHeaderData, err := chain.Codec.MarshalInterface(tmHeader)
	require.NoError(chain.t, err)
	height, ok := tmHeader.GetHeight().(clienttypes.Height)
	require.True(chain.t, ok)
	return &wasmtypes.Header{
		Data:   tmWasmHeaderData,
		Height: height,
	}
}
