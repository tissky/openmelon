.PHONY: build test clean help

build:
	go build -o openmelon ./cmd/openmelon/

test:
	go test ./...

clean:
	rm -f openmelon

help:
	@echo "Available targets:"
	@echo "  build   - Build the OpenMelon CLI"
	@echo "  test    - Run all tests"
	@echo "  clean   - Remove build artifacts"
