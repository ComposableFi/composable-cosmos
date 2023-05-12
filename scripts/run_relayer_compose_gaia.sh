#!/usr/bin/env bash

hermes --config scripts/relayer_hermes/config_compose_gaia.toml keys delete --chain gaiad-t1 --all
hermes --config scripts/relayer_hermes/config_compose_gaia.toml keys add --chain gaiad-t1 --mnemonic-file scripts/relayer_hermes/bob.json

hermes --config scripts/relayer_hermes/config_compose_gaia.toml keys delete --chain banksyd-t1 --all
hermes --config scripts/relayer_hermes/config_compose_gaia.toml keys add --chain banksyd-t1 --mnemonic-file scripts/relayer_hermes/alice.json

hermes --config scripts/relayer_hermes/config_compose_gaia.toml start
