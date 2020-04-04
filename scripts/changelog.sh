#!/usr/bin/env bash

set -e

# this script generates a CHANGELOG.md file

# Get the parent directory of where this script is and change into it.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"
cd "$DIR" || exit

re="^(.+)[\/](.+)$"

if [[ -n "${TRAVIS_TAG}" &&
      -n "${CHANGELOG_GITHUB_TOKEN}" &&
      $TRAVIS_REPO_SLUG =~ $re
]]; then

        user=${BASH_REMATCH[1]}
        repo=${BASH_REMATCH[2]}
        github_changelog_generator -u "$user" -p "$repo" --exclude-tags-regex "^archive\/.+$" --no-verbose --output "$DIR/CHANGELOG.md"

        git config --global user.email "builds@travis-ci.com"
        git config --global user.name "Travis CI"
        git add CHANGELOG.md
        git commit -m "[skip travis] update CHANGELOG.md"
        git push master

fi
