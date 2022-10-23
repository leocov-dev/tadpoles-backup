#!/usr/bin/env bash

set -e

# This script builds the application from source for multiple platforms.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"
cd "$DIR" || exit

# Starting...
echo
if [[ -z $RELEASE_TAG ]]; then
    echo "==> Build for UNTAGGED release"
else
    echo "==> Build for release ${RELEASE_TAG}"
fi

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"arm arm64 amd64"}
XC_OS=${XC_OS:-"linux darwin windows"}
XC_EXCLUDE_OSARCH="!darwin/arm !windows/arm"

if [[ -d ./bin || -d ./dist ]]; then
    echo
    echo "==> Clearing Old Builds..."
    rm -rf bin/*
    rm -rf dist/*
fi

mkdir -p dist/

if ! command -v gox >/dev/null; then
    echo
    echo "==> Installing gox..."
    go get -u github.com/mitchellh/gox
fi

# Instruct gox to build statically linked binaries
export CGO_ENABLED=0

LD_FLAGS="-s -w -X 'tadpoles-backup/config.Version=${RELEASE_TAG}'"

# Ensure all remote modules are downloaded and cached before build so that
# the concurrent builds launched by gox won't race to redundantly download them.
go mod download

# Build!
echo
echo "==> Building..."
echo "ldflags: ${LD_FLAGS}"

BIN_NAME="tadpoles-backup"
BUILD_PREFIX="${BIN_NAME}"

gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -osarch="${XC_EXCLUDE_OSARCH}" \
    -ldflags "${LD_FLAGS}" \
    -output "dist/${BUILD_PREFIX}-{{.OS}}-{{.Arch}}" \
    .

# Copy our OS/Arch to the bin/ directory
# only when not running in CI
DEV_PLATFORM="./dist/${BUILD_PREFIX}$(go env GOOS)-$(go env GOARCH)"
if [[ -f "${DEV_PLATFORM}" && -z "${CI}" ]]; then
    mkdir -p bin/
    echo
    echo "==> Copy ${DEV_PLATFORM} to bin/"
    cp "${DEV_PLATFORM}" "./bin/${BIN_NAME}"
fi

# Packaging operations
if [[ -n "${CI}" ]]; then
  echo
  echo "==> Packaging..."
  echo
  for file in ./dist/*; do
    echo "--> ${file}"

    echo "    calculate hash..."
    sha256sum "./${file}" >"./${file}.sha256"

    # IF there is a file extension, find and store it
    if [[ ${file##*/} =~ \. ]]; then
      ext="${file##*.}"
      if [[ -n ${ext} ]]; then
        ext=".${ext}"
      fi
    fi

    build="./dist/${BUILD_PREFIX}${ext}"
    sha="./dist/${BUILD_PREFIX}${ext}.sha256"

    echo "build: ${build}"
    echo "sha:   ${sha}"

    mv -f "${file}" "${build}"
    mv -f "${file}.sha256" "${sha}"

    if [[ -n $ext ]]; then
      zip -j "${file%.*}" "${build}" "${sha}"
    else
      zip -j "${file}" "${build}" "${sha}"
    fi

    rm -f "${build}" "${sha}"
  done

  # Done!
  echo
  echo "==> Results:"
  echo
  ls -hl dist/*
fi
