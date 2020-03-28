[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](https://opensource.org/licenses/MIT) ![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup) [![Build Status](https://travis-ci.org/leocov-dev/tadpoles-backup.svg?branch=golang)](https://travis-ci.org/leocov-dev/tadpoles-backup)

# Tadpoles Image Backup

## About
This tool will allow you to save all your child's images at full resolution from _tadpoles.com_.

Current save back-ends:
* Local file system

## Install
Get latest prebuilt executable from the [releases page](https://github.com/leocov-dev/tadpoles-backup/releases).

#### Linux
```
$ sudo wget https://github.com/leocov-dev/tadpoles-backup/releases/latest/download/tadpoles-backup-linux-amd64 -O /usr/local/bin/tadpoles-backup
$ sudo chmod +x /usr/local/bin/tadpoles-backup
```

#### OS X
```
$ sudo curl -Lo /usr/local/bin/tadpoles-backup https://github.com/leocov-dev/tadpoles-backup/releases/latest/download/tadpoles-backup-darwin-amd64
$ sudo chmod +x /usr/local/bin/tadpoles-backup
```

#### Windows
```
# PowerShell:
$ Invoke-WebRequest -OutFile $env:USERPROFILE\tadpoles-backup.exe https://github.com/leocov-dev/tadpoles-backup/releases/latest/download/tadpoles-backup-windows-amd64.exe
```

## Usage

> :exclamation:**IMPORTANT**:exclamation:
>
> You **MUST** have a _tadpoles.com_ account with a tadpoles specific password.
You **CAN NOT** log in to this tool with Google Auth.
If you normally log into _tadpoles.com_ with Google/Gmail account verification read these [instructions](.github/GoogleAccountSignIn.md).

```
# Print help with command details:
$ tadpoles-backup --help

# Get account statistics (requires login)
$ tadpoles-backup stat

# Download images (requires login)
$ tadpoles-backup backup "/a/directory/on/your/machine/"
```

## Development

Install Go version specified in `.go-version` (recommended to use [goenv](https://github.com/syndbg/goenv))

```
$ make dev
$ bin/tadpoles-backup --help
```

The latest development branch is always in the format: `release-v1.0.0`.

## Notes

`tadples-backup` caches your login session cookie locally so you are not prompted to enter your password every time you use the tool.
It **DOES NOT** store or retain your actual email or password.
Instead it writes a file to your home directory with a temporary authentication cookie which lasts for 2 weeks.
This file is located in `$HOME/.tadpoles-backup/` and can be deleted whenever you choose.


## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [huckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)
