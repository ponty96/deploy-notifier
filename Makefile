APP_NAME = deploy-notifier
APP_VSN ?= $(shell git rev-parse --short HEAD)

.PHONY: help
help: #: Show this help message
	@echo "$(APP_NAME):$(APP_VSN)"
	@awk '/^[A-Za-z_ -]*:.*#:/ {printf("%c[1;32m%-15s%c[0m", 27, $$1, 27); for(i=3; i<=NF; i++) { printf("%s ", $$i); } printf("\n"); }' Makefile* | sort

CGO_ENABLED ?= 0
GO = CGO_ENABLED=$(CGO_ENABLED) go
GO_BUILD_FLAGS = -ldflags "-X main.Version=${APP_VSN}"

### Dev

.PHONY: run
run: #: Run the application
	# SAMPLE_ENV="localhost" \
	$(GO) run $(GO_BUILD_FLAGS) `ls -1 *.go | grep -v _test.go`

.PHONY: lint
code-check:
	golangci-lint run -E golint --timeout 2m --exclude SA1019

### Build

.PHONY: build
build: #: Build the app locally
build: clean
	$(GO) build $(GO_BUILD_FLAGS) -o ./$(APP_NAME)

.PHONY: clean
clean: #: Clean up build artifacts
clean:
	$(RM) ./$(APP_NAME)

### Test

.PHONY: test
test: #: Run Go unit tests
test:
	GO111MODULE=on $(GO) test -v ./...
