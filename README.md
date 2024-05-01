# Media Backup for Childcare Services

[![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)](https://opensource.org/licenses/MIT)
![Go version for branch](https://img.shields.io/github/go-mod/go-version/leocov-dev/tadpoles-backup/main)
![CI Status](https://img.shields.io/github/actions/workflow/status/leocov-dev/tadpoles-backup/ci.yml)

## About
This tool will allow you to save all your child's images and videos at full resolution from various service providers. Comments and timestamp info will be applied as EXIF image metadata where possible.

Providers:
* Tadpoles
* Bright Horizons

---
## Install
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/leocov-dev/tadpoles-backup)](https://github.com/leocov-dev/tadpoles-backup/releases/latest)

Get the [latest](https://github.com/leocov-dev/tadpoles-backup/releases/latest) prebuilt
executable from the [releases](https://github.com/leocov-dev/tadpoles-backup/releases) page.
Download the zip file for your system/architecture.

> **macOS** Gatekeeper will prevent you from running unsigned apps.
> You can allow the app from system preferences or by right-clicking
> the file and choosing open from the menu.

---
## Usage

```
# Print help with command details:
$ tadpoles-backup --help

# Get account statistics
$ tadpoles-backup --provider <service-provider> stat

# Download media (only new files not present in the target dir are downloaded)
$ tadpoles-backup --provider <service-provider> backup <a-local-directory>

# Clear Saved Login
$ tadpoles-backup --provider <service-provider> clear login
```

### Provider Notes

#### Tadpoles

You **MUST** have a _www.tadpoles.com_ account with a tadpoles specific password.
You **CAN NOT** log in to this tool with Google Auth.
If you normally log into _tadpoles.com_ with Google/Gmail account verification you will need to
request a password reset with the command:
```shell
# this simply requests a reset email be sent to you
# it does not change or access your password
$ tadpoles-backup --provider tadpoles reset-password <email>
```

The tool stores your _www.tadpoles.com_ authentication cookie for future use so that you don't need to enter your password every time.
This cookie lasts for about 2 weeks. Your email and password are never stored.

#### Bright Horizons

Due to how the system provides download data the `backup` command can't use cached data for speed-up.
Every run of the `backup` command will fetch all reports and may take some time.

The tool stores your _mybrightday.brighthorizons.com_ api-key for future use so that you don't need to enter your password every time.
This api-key may only expire if you change your password. Your email and password are never stored.

---
## Container Image
Pre-built images are available from Docker Hub

[![Docker Image Version (latest by date)](https://img.shields.io/docker/v/leocov/tadpoles-backup?label=latest&sort=date)](https://hub.docker.com/r/leocov/tadpoles-backup)

```shell
$ docker pull leocov/tadpoles-backup:latest

# list account info
$ docker run --rm -eUSERNAME=<email> -ePASSWORD=<password> leocov/tadpoles-backup stat

# download new images
$ docker run --rm -eUSERNAME=<email> -ePASSWORD=<password> -v$HOME/Pictures/tadpoles:/images leocov/tadpoles-backup backup /images

# enable api response caching by mapping app data directory
$ docker run --rm -eUSERNAME=<email> -ePASSWORD=<password> -v$HOME/.tadpoles-backup:/app/.tadpoles-backup leocov/tadpoles-backup stat
```

You may also build the docker image locally.
```shell
# will be automatically tagged as `tadpoles-backup`
$ make docker-image
```

### Docker Compose / Kubernetes

Please note that this utility is intended to run as a scheduled job.

[Examples](examples) are available.

#### Kubernetes

This [example](examples/kubernetes) configures a `CronJob` that will run on a schedule. It's best to configure
this so that only 1 job instance will run at a time. The example uses `kustomize` for
configuration to provide authentication environment vars as a secret.

#### Docker Compose

This [example](examples/docker-compose.yml) configures a basic service with env
vars defining the login values. Its important to remember that this service will
exit after each run.

---
## Development

See the contributing guide [here](CONTRIBUTING.md).

### Basic Setup

Install the Go version defined in [go.mod](go.mod) or use [goenv](https://github.com/syndbg/goenv) to manage Go (as set by [.go-version](.go-version)).

### Dev build
```shell
# build for your platform only and run.
$ make && bin/tadpoles-backup --help
```

### Testing

Run all unit tests with helper utility. This will build a coverage report as
`coverage.html`
```shell
make test
```


---
## Inspired By
* [twneale/tadpoles](https://github.com/twneale/tadpoles)
* [ChuckMac/tadpoles-scraper](https://github.com/ChuckMac/tadpoles-scraper)

## Thanks to
* @arthurnn - for assistance with Docker image
* @AndyRPH - for assistance with Bright Horizons support
* @s0rcy - for assistance with Tadpoles password reset
