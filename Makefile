GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

test:
	@rm -f coverage/coverage.out coverage/coverage.html
	@go test -v -coverpkg=./... -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html

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

.PHONY: bin default dev fmtcheck fmt tidy test
