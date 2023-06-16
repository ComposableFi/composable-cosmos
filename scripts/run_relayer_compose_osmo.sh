#!/usr/bin/env bash

hermes --config scripts/relayer_hermes/config_compose_osmosis.toml keys delete --chain testing --all
hermes --config scripts/relayer_hermes/config_compose_osmosis.toml keys add --chain testing --mnemonic-file scripts/relayer_hermes/gnad.json

hermes --config scripts/relayer_hermes/config_compose_osmosis.toml keys delete --chain centaurid-t1 --all
hermes --config scripts/relayer_hermes/config_compose_osmosis.toml keys add --chain centaurid-t1 --mnemonic-file scripts/relayer_hermes/alice.json

hermes --config scripts/relayer_hermes/config_compose_osmosis.toml start
