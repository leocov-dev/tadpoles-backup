GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

# bin generates release zip packages in ./dist
release: tidy fmt
	@sh -c "$(CURDIR)/scripts/build.sh"

# dev creates binaries for testing locally.
# These are put into ./bin/ as well as $GOPATH/bin
dev: tidy fmt
	@go install ./...

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

tidy:
	@go mod tidy

.NOTPARALLEL:

.PHONY: bin default dev fmtcheck fmt tidy
