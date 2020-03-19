[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](https://opensource.org/licenses/MIT) ![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup) [![Build Status](https://travis-ci.org/leocov-dev/tadpoles-backup.svg?branch=golang)](https://travis-ci.org/leocov-dev/tadpoles-backup)

# Tadpoles Image Backup

## About
This tool will allow you to save all your child's images at full resolution from _tadpoles.com_.

Current save back-ends:
* Local file system

## Install
Get a prebuilt executable from the [releases page](https://github.com/leocov-dev/tadpoles-backup/releases).

#### Linux
```
$ sudo wget https://github.com/leocov-dev/tadpoles-backup/releases/download/v1.0.0/tadpoles-backup-v1.0.0_linux_amd64 -O /usr/local/bin/tadpoles-backup
$ sudo chmod +x /usr/local/bin/tadpoles-backup
```

#### OS X
```
$ sudo curl -Lo /usr/local/bin/tadpoles-backup https://github.com/leocov-dev/tadpoles-backup/releases/download/v1.0.0/tadpoles-backup-v1.0.0_darwin_amd64
$ sudo chmod +x /usr/local/bin/tadpoles-backup
```

#### Windows
```
# PowerShell:
$ Invoke-WebRequest -OutFile $env:USERPROFILE\tadpoles-backup https://github.com/leocov-dev/tadpoles-backup/releases/download/v1.0.0/tadpoles-backup-v1.0.0_windows_amd64
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

## Notes

`tadples-backup` caches your login session cookie locally so you are not prompted to enter your password every time you use the tool.
It **DOES NOT** store or retain your actual email or password.
Instead it writes a file to your home directory with a temporary authentication cookie which lasts for 2 weeks.
This file is located at `$HOME/.tadpoles-backup-cookie` and can be deleted whenever you choose.


## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [huckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)
