[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](https://opensource.org/licenses/MIT) ![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup/main) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/leocov-dev/tadpoles-backup/ci)

# Tadpoles Image Backup

## About
This tool will allow you to save all your child's images at full resolution from _tadpoles.com_. Comments and timestamp info will be applied as EXIF image metadata.

Current save back-ends:
* Local file system

## Install
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/leocov-dev/tadpoles-backup)](https://tadpoles-backup/releases/latest)

Get [latest](https://tadpoles-backup/releases/latest) prebuilt executable from the [releases](https://tadpoles-backup/releases) page.

Or use one of these commands:
#### Linux
```
$ sudo wget https://tadpoles-backup/releases/latest/download/tadpoles-backup-linux-amd64 -O /usr/local/bin/tadpoles-backup
$ sudo chmod +x /usr/local/bin/tadpoles-backup
```

#### OS X
```
$ sudo curl -Lo /usr/local/bin/tadpoles-backup https://tadpoles-backup/releases/latest/download/tadpoles-backup-darwin-amd64
$ sudo chmod +x /usr/local/bin/tadpoles-backup
```

#### Windows
```
# PowerShell:
$ Invoke-WebRequest -OutFile $env:USERPROFILE\tadpoles-backup.exe https://tadpoles-backup/releases/latest/download/tadpoles-backup-windows-amd64.exe
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

> Note for macOS
> Gatekeeper will prevent you from running unidentified apps.
> You can allow the app from system preferences or by right-click opening

### Docker

Pre-built images are available from Docker Hub at [leocov/tadpoles-backup](https://hub.docker.com/r/leocov/tadpoles-backup).

```shell
$ docker pull leocov/tadpoles-backup:latest

# list account info
$ docker run --rm -eTADPOLES_USER=<email> -eTADPOLES_PASS=<password> leocov/tadpoles-backup stat

# download new images
$ docker run --rm -eTADPOLES_USER=<email> -eTADPOLES_PASS=<password> -v$HOME/Pictures/tadpoles:/images leocov/tadpoles-backup backup /images

# enable api response caching by mapping app data directory
$ docker run --rm -eTADPOLES_USER=<email> -eTADPOLES_PASS=<password> -v$HOME/.tadpoles-backup:/app/.tadpoles-backup leocov/tadpoles-backup stat
```

You may also build the docker image locally.
```shell
# will be automatically tagged as `tadpoles-backup`
$ make docker-image
```

Docker Compose and Kubernetes [examples](examples) are available.

## Development

See the contributing guide [here](CONTRIBUTING.md).

Install the Go version defined in [go.mod](go.mod) or use [goenv](https://github.com/syndbg/goenv) to manage Go (as set by [.go-version](.go-version)).

```
# build for your platform only and run.
$ make && bin/tadpoles-backup --help
```

## Notes

`tadples-backup` caches your login session cookie locally so you are not prompted to enter your password every time you use the tool.
It **DOES NOT** store or retain your actual email or password.
Instead it writes a file to your home directory with a temporary authentication cookie which lasts for 2 weeks.
This file is located in `$HOME/.tadpoles-backup/` and can be deleted whenever you choose.


## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [ChuckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)
