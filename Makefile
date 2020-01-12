GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')

default: dev

# bin generates the releaseable binaries for Terraform
bin: fmtcheck
	@TF_RELEASE=1 sh -c "$(CURDIR)/scripts/build.sh"

# dev creates binaries for testing Terraform locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: fmtcheck
	go install .

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.NOTPARALLEL:

.PHONY: bin default dev e2etest fmtcheck generate
