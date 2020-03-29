#!/usr/bin/env bash

set -e

# this script runs the ci pipeline

# Get the parent directory of where this script is and change into it.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"
cd "$DIR" || exit

source scripts/_utils.sh

center "CI PROCESS STARTED"

echo
scripts/gofmtcheck.sh

echo
scripts/test.sh

echo
scripts/release.sh

echo
center "CI COMPLETE"
