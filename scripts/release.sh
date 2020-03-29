#!/usr/bin/env bash

set -e

# This script builds the application from source for multiple platforms.
echo "==> Build for release..."

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"
RELEASE_TAG="${RELEASE_TAG:-v0.0.0-dev}"

# Change into that directory
cd "$DIR" || exit

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-linux darwin windows}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

# Delete the old dir
echo
echo "==> Removing old directory..."
rm -rf bin/*
mkdir -p bin/
rm -rf dist/*
mkdir -p dist/

if ! command -v gox >/dev/null; then
    echo "==> Installing gox..."
    go get -u github.com/mitchellh/gox
fi

# Instruct gox to build statically linked binaries
export CGO_ENABLED=0

LD_FLAGS="-s -w -X 'github.com/leocov-dev/tadpoles-backup/config.Version=${RELEASE_TAG}'"

# Ensure all remote modules are downloaded and cached before build so that
# the concurrent builds launched by gox won't race to redundantly download them.
go mod download

# Build!
echo
echo "==> Building..."
echo "ldflags: ${LD_FLAGS}"

BIN_NAME=${PWD##*/}
BUILD_PREFIX="${BIN_NAME}-"

gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -osarch="${XC_EXCLUDE_OSARCH}" \
    -ldflags "${LD_FLAGS}" \
    -output "dist/${BUILD_PREFIX}{{.OS}}-{{.Arch}}" \
    .

# Copy our OS/Arch to the bin/ directory
# only when not running in CI
DEV_PLATFORM="./dist/${BUILD_PREFIX}$(go env GOOS)-$(go env GOARCH)"
if [[ -f "${DEV_PLATFORM}" && -z "${CI}" ]]; then
    echo
    echo "==> Moving ${DEV_PLATFORM} to bin/"
    cp "${DEV_PLATFORM}" "./bin/${BIN_NAME}"
fi

# Packaging operations
# only if not a pull request
if [[ "${TRAVIS_PULL_REQUEST}" != "false" ]]; then
    echo
    echo "==> Packaging..."
    echo
    for file in ./dist/*; do
        echo "--> ${file}"

        echo "    calculate hash..."
        sha256sum "./${file}" >"./${file}.sha256"
    done

    # Done!
    echo
    echo "==> Results:"
    echo
    ls -hl dist/*
fi
