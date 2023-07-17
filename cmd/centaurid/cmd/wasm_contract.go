package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	wasmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	wasm08types "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
)

type Checksum = wasmtypes.Checksum

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddWasmContractCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-wasm-contract [wasm file]",
		Short: "Add a wasm contract to genesis.json",
		Long:  `Add a wasm contract to genesis.json. Wasm contract don't need be added via gov module.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.Codec
			cdc := depCdc

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			wasm, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			if !ioutils.IsWasm(wasm) {
				return fmt.Errorf("invalid input file. Use wasm binary or gzip")
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			codeHash, err := CreateChecksum(wasm)
			if err != nil {
				return err
			}
			contract := wasm08types.GenesisContract{
				CodeHash:     codeHash,
				ContractCode: wasm,
			}

			wasmGenState := GetGenesisStateFromAppState(cdc, appState)

			wasmGenState.Contracts = append(wasmGenState.Contracts, contract)

			wasmGenStateBz, err := cdc.MarshalJSON(wasmGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal 08-wasm genesis state: %w", err)
			}

			appState["08-wasm"] = wasmGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetGenesisStateFromAppState return GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *wasm08types.GenesisState {
	var genesisState wasm08types.GenesisState

	if appState["08-wasm"] != nil {
		cdc.MustUnmarshalJSON(appState["08-wasm"], &genesisState)
	}

	return &genesisState
}

func CreateChecksum(wasm []byte) (Checksum, error) {
	if len(wasm) == 0 {
		return Checksum{}, fmt.Errorf("wasm bytes nil or empty")
	}
	if len(wasm) < 4 {
		return Checksum{}, fmt.Errorf("wasm bytes shorter than 4 bytes")
	}
	// magic number for Wasm is "\0asm"
	// See https://webassembly.github.io/spec/core/binary/modules.html#binary-module
	if !bytes.Equal(wasm[:4], []byte("\x00\x61\x73\x6D")) {
		return Checksum{}, fmt.Errorf("wasm bytes do not not start with Wasm magic number")
	}
	hash := sha256.Sum256(wasm)
	return Checksum(hash[:]), nil
}
