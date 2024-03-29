GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

test:
	@bash -c "scripts/test.sh"

# bin generates release zip packages in ./dist
release: tidy
	@bash -c "scripts/release.sh"

clean:
	@rm -rf "bin"
	@rm -rf "dist"

dev: tidy fmt
	@go build -race -o "bin/tadpoles-backup" .

docker-image:
	@docker build -t tadpoles-backup .

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@bash -c "scripts/gofmtcheck.sh"

tidy:
	@go mod tidy

.NOTPARALLEL:

.PHONY: bin default dev fmtcheck fmt tidy
