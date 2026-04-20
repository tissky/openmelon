.PHONY: build test lint clean help

build:
	go build -o skillplus-engine ./cmd/skillplus-engine/

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f skillplus-engine

help:
	@echo "Available targets:"
	@echo "  build   - Build the engine binary"
	@echo "  test    - Run all tests"
	@echo "  lint    - Run linter"
	@echo "  clean   - Remove build artifacts"
