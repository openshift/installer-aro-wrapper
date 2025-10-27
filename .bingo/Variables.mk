# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.9. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for bingo variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BINGO)
#	@echo "Running bingo"
#	@$(BINGO) <flags/args..>
#
BINGO := $(GOBIN)/bingo-v0.9.0
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.9.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.9.0 "github.com/bwplotka/bingo"

ENUMER := $(GOBIN)/enumer-v1.6.1
$(ENUMER): $(BINGO_DIR)/enumer.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/enumer-v1.6.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=enumer.mod -o=$(GOBIN)/enumer-v1.6.1 "github.com/dmarkham/enumer"

FIPS_DETECT := $(GOBIN)/fips-detect-v0.0.0-20230309083406-7157dae5bafd
$(FIPS_DETECT): $(BINGO_DIR)/fips-detect.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/fips-detect-v0.0.0-20230309083406-7157dae5bafd"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=fips-detect.mod -o=$(GOBIN)/fips-detect-v0.0.0-20230309083406-7157dae5bafd "github.com/acardace/fips-detect"

GOJQ := $(GOBIN)/gojq-v0.12.17
$(GOJQ): $(BINGO_DIR)/gojq.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/gojq-v0.12.17"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=gojq.mod -o=$(GOBIN)/gojq-v0.12.17 "github.com/itchyny/gojq/cmd/gojq"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v2.2.1
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v2.2.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v2.2.1 "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"

GOTESTSUM := $(GOBIN)/gotestsum-v1.12.0
$(GOTESTSUM): $(BINGO_DIR)/gotestsum.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/gotestsum-v1.12.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=gotestsum.mod -o=$(GOBIN)/gotestsum-v1.12.0 "gotest.tools/gotestsum"

MOCKGEN := $(GOBIN)/mockgen-v1.7.0-rc.1
$(MOCKGEN): $(BINGO_DIR)/mockgen.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/mockgen-v1.7.0-rc.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=mockgen.mod -o=$(GOBIN)/mockgen-v1.7.0-rc.1 "github.com/golang/mock/mockgen"

