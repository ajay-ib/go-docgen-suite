# Makefile for go-docgen-suite

# Format the code
.PHONY: fmt
fmt:
	@go fmt ./...

# Run tests
.PHONY: test
test:
	@go test ./...