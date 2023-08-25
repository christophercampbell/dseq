#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/dseq/${BINARY:-dseq}-linux
ID=${ID:-0}
LOG=${LOG:-dseq.log}
HOMEDIR="/dseq/node${ID}"
COMMAND="run --home ${HOMEDIR}"

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found"
	exit 1
fi

BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64 (build with 'make build-linux')"
	exit 1
fi

$BINARY start --home $HOMEDIR | tee "${HOMEDIR}/${LOG}"

chmod 777 -R /dseq

