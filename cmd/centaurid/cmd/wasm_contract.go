package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	wasm08types "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
)

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
			f, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				log.Fatal(err)
			}
			contract := wasm08types.GenesisContract{
				CodeHash:     h.Sum(nil),
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
