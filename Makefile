PKGS := $(shell go list ./... | grep -v dist/)

BIN_DIR=$(shell go env GOPATH)/bin

LINTER=$(BIN_DIR)/golangci-lint

lint:
	@test ! -f $(LINTER) && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || true
	@$(LINTER) run -v

test:
	@go test -v -cover $(PKGS)

integration:
	@go test -v -cover -tags=integration $(PKGS)