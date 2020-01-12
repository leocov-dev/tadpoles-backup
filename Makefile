GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

# bin generates release zip packages in ./dist
bin: fmtcheck
	@sh -c "$(CURDIR)/scripts/build.sh"

# dev creates binaries for testing locally.
# These are put into ./bin/ as well as $GOPATH/bin
dev: fmtcheck
	go install .

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.NOTPARALLEL:

.PHONY: bin default dev fmtcheck fmt
