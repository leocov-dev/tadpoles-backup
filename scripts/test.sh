#!/usr/bin/env bash

set -e

echo "==> Running tests..."
echo

_root=$(go list .)

for d in $(go list ./... | grep -v vendor); do
    echo "Testing: $d"

    _covertarget="coverage/$d"
    mkdir -p "$_covertarget"

    go test -coverprofile="$_covertarget/coverage.out" "$d"
done
