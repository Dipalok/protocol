# Makefile for building and running Go and Rust applications

# Configuration
GO := go
GO_SRC := $(wildcard src/go/*.go)
BIN_DIR := bin
BIN_GO := $(BIN_DIR)/server 
CARGO := cargo
RUST_TARGET := target/release/


# Default target (build both)
all: build-go build-rust

# Create bin directory if it doesn't exist
$(shell mkdir -p $(BIN_DIR))

# Go targets
build-go: $(BIN_GO)

$(BIN_GO): $(GO_SRC)
	$(GO) build -o $@ $(GO_SRC)

run-go: build-go
	@echo "Running Go server..."
	@./$(BIN_GO) -cert server.crt -key server.key

# Rust targets
build-rust:
	$(CARGO) build --release

run-client: build-rust
	@echo "Running Rust client..."
	@./$(RUST_TARGET)protocol  # Replace 'client' with your binary name from Cargo.toml

# Cleanup
clean:
	rm -rf $(BIN_DIR)
	$(CARGO) clean

.PHONY: all build-go build-rust run-go run-client clean