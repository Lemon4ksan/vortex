# Discover library packages, excluding examples, scripts, cmd, and vendor
PKG       := $(shell go list ./... 2>/dev/null | grep -v /examples | grep -v /scripts | grep -v /cmd/ | grep -v /vendor/)
COVER_PKG := $(shell go list ./... 2>/dev/null | grep -v /examples | grep -v /scripts | grep -v /vendor/)
BIN_DIR   ?= bin
TMP_DIR   ?= .tmp
COVER_OUT ?= $(TMP_DIR)/coverage.out

# Colors for console output
CYAN  := \033[0;36m
RESET := \033[0m

.PHONY: build install test race bench cover cover-clean cover-html lint format clean help

build: ## Build vortex CLI binary (bin/vortex)
	@printf "$(CYAN)Building vortex CLI binary...$(RESET)\n"
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/vortex ./cmd/vortex

install: ## Install vortex binary into GOPATH/bin
	@printf "$(CYAN)Installing vortex binary...$(RESET)\n"
	go install ./cmd/vortex

test: ## Run quick unit tests
	@printf "$(CYAN)Running unit tests...$(RESET)\n"
	go test -timeout 60s $(PKG)

race: ## Run unit tests with race detector enabled
	@printf "$(CYAN)Running tests with race detector...$(RESET)\n"
	go test -race -timeout 120s $(PKG)

bench: ## Run silicon hardware inspection and microsecond benchmark suite
	@printf "$(CYAN)Running benchmarks...$(RESET)\n"
	go test -bench=. -benchmem -run=^$$ $(PKG)

cover: ## Calculate and print exact test coverage report
	@printf "$(CYAN)Generating exact coverage report...$(RESET)\n"
	@mkdir -p $(TMP_DIR)
	go test -coverpkg=$(COVER_PKG) -coverprofile=$(COVER_OUT) ./...
	go tool cover -func=$(COVER_OUT)

cover-clean: ## Generate clean coverage report and run deduplicated coverage analysis tool
	@printf "$(CYAN)Generating clean coverage report...$(RESET)\n"
	@mkdir -p $(TMP_DIR)
	go test -coverpkg=$(COVER_PKG) -coverprofile=$(COVER_OUT) ./...
	go tool cover -func=$(COVER_OUT)

cover-html: cover ## Generate coverage report and open interactive HTML in browser
	@printf "$(CYAN)Opening coverage report in browser...$(RESET)\n"
	go tool cover -html=$(COVER_OUT)

lint: ## Run golangci-lint check and vortex AST contract inspector
	@printf "$(CYAN)Running golangci-lint...$(RESET)\n"
	golangci-lint run ./...
	@printf "$(CYAN)Running vortex AST check...$(RESET)\n"
	go run ./cmd/vortex check ./...
	@printf "$(CYAN)Running AST borrow checker...$(RESET)\n"
	go run ./cmd/vortex borrow ./...

format: ## Format code and auto-fix linter suggestions
	@printf "$(CYAN)Formatting Go code...$(RESET)\n"
	go fmt ./...
	addlicense -c "Lemon4ksan" -l bsd -ignore "**/*.yml" .
	golangci-lint run --fix ./...

clean: ## Delete temporary files, binaries, and coverage profiles
	@printf "$(CYAN)Cleaning up temporary artifacts...$(RESET)\n"
	rm -rf $(BIN_DIR)/ $(TMP_DIR)/
	rm -f $(COVER_OUT) coverage.out profile.cov *.out *.test *.exe

help: ## Show this help message
	@printf "Usage: make [target]\n\nTargets:\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
