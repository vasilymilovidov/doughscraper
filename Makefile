.PHONY: build_go build_rust

BINARY_GO := doughscraper
BINARY_RUST := pitchdetector
INSTALL_PATH := /usr/local/bin/

build_go:
	@echo "Building Go binary"
	@cd go && go build -o ../$(BINARY_GO)

build_rust:
	@echo "Building Rust binary"
	@cd rust/pitchdetector && cargo build --release
	@cp rust/pitchdetector/target/release/$(BINARY_RUST) .

all: build_go build_rust
