[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](https://opensource.org/licenses/MIT) ![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup) [![Build Status](https://travis-ci.org/leocov-dev/tadpoles-backup.svg?branch=golang)](https://travis-ci.org/leocov-dev/tadpoles-backup)

# Tadpoles Image Backup

#### **This is still a work in progress! - Non-functional**

## About
This tool will allow you to save all your child's images at full resolution from _tadpoles.com_.  

Current save back-ends:
* Local file system
* ~~Amazon S3~~
* ~~Backblaze B2~~

## Install
Get a prebuilt executable from the releases page.  Download and extract `tadpoles-backup` to a place of your choosing.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/leocov-dev/tadpoles-backup)

## Usage

> You **MUST** have a _tadpoles.com_ account with a valid password. 
You **CAN NOT** log in to this tool with Google Auth.
If you normally log into _tadpoles.com_ with Gmail/Google account verification read these [instructions](.github/GoogleAccountSignIn.md).

```bash
# Print help with command details:
$ tadpoles-backup --help

# Get account statistics (requires login)
$ tadpoles-backup stat

# Download images (requires login)
$ tadpoles-backup backup "/a/directory/on/your/machine/"
```

## Development

Install Go version specified in `.go-version` (recommended to use [goenv](https://github.com/syndbg/goenv))

```bash
$ make dev
$ bin/tadpoles-backup --help
```

## Notes

`tadples-backup` caches your login session cookie locally so you are not prompted repeatedly to enter your password. 
It **DOES NOT** store or retain your email or password!

It writes a file to your home directory with a temporary authentication cookie which lasts a few weeks.
This file is located at `$HOME/.tadpoles-backup-cookie` and can be deleted whenever you choose.


## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [huckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)