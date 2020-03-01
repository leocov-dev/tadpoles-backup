#!/usr/bin/env bash

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."
gofmt_files=$(gofmt -l "$(find . -type f -name "*.go")")
if [[ -n ${gofmt_files} ]]; then
    echo "gofmt [FAIL]"
    echo "${gofmt_files}"
    echo "Run: \`make fmt\` to reformat code."
    exit 1
else
  echo "gofmt [PASS]"
fi

exit 0
