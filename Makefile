SHELL=/bin/sh
include Makefile.*

.PHONY: fmt
fmt:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Formatting code..."
	@echo "----------------------------------------------------------------"
	$(GO) fmt ./...
	$(GOMOD) tidy

.PHONY: lint
lint:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Linting code..."
	@echo "----------------------------------------------------------------"
	$(GOINSTALL) "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	golangci-lint run

.PHONY: test
test:
	@echo "----------------------------------------------------------------"
	@echo " ✅  Testing code..."
	@echo "----------------------------------------------------------------"
	$(GO) test ./... -coverprofile=coverage.out

.PHONY: coverage
coverage:
	@echo "----------------------------------------------------------------"
	@echo " 📊  Checking coverage..."
	@echo "----------------------------------------------------------------"
	$(GOTOOL) cover -html=coverage.out -o coverage.html
	$(GOTOOL) cover -func=coverage.out

.PHONY: deps
deps:
	@echo "----------------------------------------------------------------"
	@echo " ⬇️  Downloading dependencies..."
	@echo "----------------------------------------------------------------"
	$(GOGET) ./...

.PHONY: build
build: deps fmt
	@echo "----------------------------------------------------------------"
	@echo " 📦 Building binary..."
	@echo "----------------------------------------------------------------"
	$(GOBUILD) -o foaas-limiting main.go

.PHONY: run
run:
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Running..."
	@echo "----------------------------------------------------------------"
	./foaas-limiting serve

.PHONY: all
all: test build
