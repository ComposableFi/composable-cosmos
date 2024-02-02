#!/usr/bin/env sh

BINARY=/centaurid/${BINARY:-cosmovisor}
ID=${ID:-0}
LOG=${LOG:-centaurid.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'centaurid'"
	exit 1
fi

BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"

if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

export CENTAURID_HOME="/centaurid/node${ID}/centaurid"

if [ -d "$(dirname "${CENTAURID_HOME}"/"${LOG}")" ]; then
    "${BINARY}" run "$@" --home "${CENTAURID_HOME}" | tee "${CENTAURID_HOME}/${LOG}"
else
    "${BINARY}" run "$@" --home "${CENTAURID_HOME}"
fi