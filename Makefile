BINARY_NODE=selsichain
BINARY_CLI=selsichain-cli

build:
	@echo "Building SelsiChain Node..."
	@mkdir -p bin
	go build -o bin/$(BINARY_NODE) cmd/selsichain/main.go
	@echo "Building SelsiChain CLI..."
	go build -o bin/$(BINARY_CLI) cmd/cli/main.go

run: build
	@echo "Starting SelsiChain node..."
	./bin/$(BINARY_NODE)

cli: build
	@echo "Starting SelsiChain CLI..."
	./bin/$(BINARY_CLI)

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

test:
	@echo "Running tests..."
	go test ./...

.PHONY: build run cli clean test
