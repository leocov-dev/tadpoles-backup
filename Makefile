GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

ci:
	@sh -c "$(CURDIR)/scripts/ci.sh"

# bin generates release zip packages in ./dist
release: tidy fmt
	@sh -c "$(CURDIR)/scripts/release.sh"

clear:
	@rm -rf "$(CURDIR)/bin"
	@rm -rf "$(CURDIR)/dist"

# dev creates binaries for testing locally.
# These are put into $GOPATH/bin
dev: tidy fmt
	@go build -race -o "$(CURDIR)/bin/tadpoles-backup" .

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

tidy:
	@go mod tidy

.NOTPARALLEL:

.PHONY: bin default dev fmtcheck fmt tidy
