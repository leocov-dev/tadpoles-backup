# Tadpoles Image Backup

#### **This is still a work in progress! - Non-functional**

Inspired by but reworked from scratch to make use of the REST API behind the `www.tadpoles.com`. 

I started writing this in Python but have since switched to Go.

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