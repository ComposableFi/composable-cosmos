package wasmbinding

import (
	"github.com/CosmWasm/wasmd/x/wasm"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin()

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	return []wasm.Option{
		queryPluginOpt,
	}
}

func RegisterStargateQueries(queryRouter baseapp.GRPCQueryRouter, cdc codec.Codec) []wasmkeeper.Option {
	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: StargateQuerier(queryRouter, cdc),
	})

	return []wasm.Option{
		queryPluginOpt,
	}
}
