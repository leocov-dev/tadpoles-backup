#!/usr/bin/env bash

set -e

# This script builds the application from source for multiple platforms.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"
RELEASE_TAG="${RELEASE_TAG:-v0.0.0-dev}"

# Change into that directory
cd "$DIR" || exit

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

# Delete the old dir
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
BUILD_PREFIX="${BIN_NAME}-${RELEASE_TAG}_"

gox \
  -os="${XC_OS}" \
  -arch="${XC_ARCH}" \
  -osarch="${XC_EXCLUDE_OSARCH}" \
  -ldflags "${LD_FLAGS}" \
  -output "dist/${BUILD_PREFIX}{{.OS}}_{{.Arch}}/${BIN_NAME}" \
  .

# Copy our OS/Arch to the bin/ directory
# only when not running in CI
DEV_PLATFORM="./dist/${BUILD_PREFIX}$(go env GOOS)_$(go env GOARCH)"
if [[ -d "${DEV_PLATFORM}" && -z "${CI}" ]]; then
  echo
  echo "==> Moving ${DEV_PLATFORM} to bin/"
  for F in $(find "${DEV_PLATFORM}" -mindepth 1 -maxdepth 1 -type f); do
    cp "${F}" bin/
  done
fi

# Zip and copy to the dist dir
echo
echo "==> Packaging..."
for PLATFORM in $(find ./dist -mindepth 1 -maxdepth 1 -type d); do
  OSARCH=$(basename "${PLATFORM}")
  echo "--> ${OSARCH}"

  pushd "${PLATFORM}" >/dev/null 2>&1 || exit

  zip "../${OSARCH}.zip" ./*
  zip -uj "../${OSARCH}.zip" ../../LICENSE

  echo "  calculate hash..."
  sha256sum "../${OSARCH}.zip" > "../${OSARCH}.zip.sha256"

  rm -rf "../${OSARCH}"
  popd >/dev/null 2>&1 || exit
done

# Done!
echo
echo "==> Results:"
ls -hl dist/*
