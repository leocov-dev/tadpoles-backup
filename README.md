[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) ![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup/golang) [![Build Status](https://travis-ci.org/leocov-dev/tadpoles-backup.svg?branch=golang)](https://travis-ci.org/leocov-dev/tadpoles-backup)

# Tadpoles Image Backup

#### **This is still a work in progress! - Non-functional**

## About
This tool will allow you to save all your child's images at full resolution from `www.tadpoles.com`.  It can be be configured with multiple save back-ends.

Current save back-ends:
* Local file system (non-functional)

## Usage

Get the latest release and run the executable file `tadpoles-backup`
```bash
# Print help with command details:
$ tadpoles-backup --help
```

## Development

Install Go version specified in `.go-version` (recommended to use [goenv](https://github.com/syndbg/goenv))

```bash
# Build in development mode:
$ make
```

## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [huckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)