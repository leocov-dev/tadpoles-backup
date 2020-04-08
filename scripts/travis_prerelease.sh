#!/usr/bin/env bash

set -e

# Get the parent directory of where this script is and change into it.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"
cd "$DIR" || exit

source scripts/_utils.sh

center "PRE-RELEASE SETUP"
echo

branch_re="^(dev-)(v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?)$"

if [[ -n "${TRAVIS_BRANCH}" &&
      $TRAVIS_BRANCH =~ $branch_re
]]; then
    export TRAVIS_TAG=${TRAVIS_TAG:-${BASH_REMATCH[2]}}
    echo "OVERRIDE: TRAVIS_TAG=${TRAVIS_TAG}"

    git config --global user.name "Travis CI"
    git config --global user.email "builds@travis-ci.com"
    git tag "$TRAVIS_TAG"

    echo
fi
