SHELL := /bin/bash
.PHONY: run routes build test test-quiet

run:
	go run ./cmd/api

# Usage: make routes            (all routes)
#        make routes ARGS="-method POST"
#        make routes ARGS="-path auth"
routes:
	@go run ./cmd/routes $(ARGS)

build:
	go build ./...

# Verbose: every test and subtest with PASS/FAIL plus a total.
# Usage: make test ARGS="-run TestAuth"
test:
	@set -o pipefail; go test -v -count=1 $(ARGS) ./... 2>&1 | awk -f scripts/test-report.awk

# Package-level summary only.
test-quiet:
	@go test ./...
