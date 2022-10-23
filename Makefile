GOFMT_FILES?=$$(find . -type f -name '*.go')

default: dev

test:
	@sh -c "$(CURDIR)/scripts/test.sh"

node-install:
	@echo "\n==> NPM Install"
	@cd utils && npm install --no-progress > ../node.log

node-build:
	@echo "\n==> Build html templates"
	@cd utils && npm run build >> ../node.log
	@go generate ./main.go

# bin generates release zip packages in ./dist
release: tidy
	@sh -c "$(CURDIR)/scripts/release.sh"

clean:
	@rm -rf "$(CURDIR)/bin"
	@rm -rf "$(CURDIR)/dist"

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
