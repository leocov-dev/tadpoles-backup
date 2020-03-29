#!/usr/bin/env bash

set -e

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."
echo

mapfile -t args < <(find . -type f -name "*.go")
gofmt_files=$(gofmt -e -l "${args[@]}")
if [[ -n ${gofmt_files} ]]; then
    echo "gofmt [FAIL]"
    echo "${gofmt_files}"
    echo "Run: \`make fmt\` to reformat code."
    exit 1
else
    echo "gofmt [PASS]"
fi

exit 0
