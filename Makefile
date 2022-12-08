gofumpt := mvdan.cc/gofumpt@v0.4.0
gosimports := github.com/rinchsan/gosimports/cmd/gosimports@v0.3.4
golangci_lint := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
goreleaser := github.com/goreleaser/goreleaser@v1.13.1

.PHONY: build
build:
	@mkdir -p build
	@go test -c . -o build/http-wasm-tck

.PHONY: test
test:
	@go test -v ./e2e

.PHONY: lint
lint:
	@CGO_ENABLED=0 go run $(golangci_lint) run --timeout 5m

.PHONY: format
format:
	@go run $(gofumpt) -l -w .
	@go run $(gosimports) -local github.com/http-wasm/ -w $(shell find . -name '*.go' -type f)

.PHONY: check
check: lint format
	@go mod tidy
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi

.PHONY: clean
clean: ## Ensure a clean build
	@go clean -testcache
	@rm -rf build
	@rm -rf dist

.PHONY: snapshot
snapshot:
	@go run $(goreleaser) build --snapshot --rm-dist