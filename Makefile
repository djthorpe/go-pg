# Path parameters
GO ?= $(shell which go)
BUILDDIR ?= build
CMDDIR=$(wildcard cmd/*)

# Build flags
BUILD_MODULE = $(shell cat go.mod | head -1 | cut -d ' ' -f 2)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitSource=${BUILD_MODULE}
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitTag=$(shell git describe --tags --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitBranch=$(shell git name-rev HEAD --name-only --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitHash=$(shell git rev-parse HEAD)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_FLAGS = -ldflags "-s -w ${BUILD_LD_FLAGS}"

# All targets
all: tidy $(CMDDIR)

# Rules for building
.PHONY: $(CMDDIR)
$(CMDDIR): mkdir
	@echo 'building $@'
	@$(GO) build $(BUILD_FLAGS) -o ${BUILDDIR}/$(shell basename $@) ./$@

# Rules for testing
.PHONY: test
test: tidy
	@echo 'running tests...'
	@$(GO) test .
	@$(GO) test ./pkg/...

# Other rules
.PHONY: mkdir
mkdir:
	@install -d $(BUILDDIR)

.PHONY: tidy
tidy: 
	@echo 'go tidy'
	@$(GO) mod tidy

.PHONY: clean
clean: tidy
	@echo 'clean'
	@rm -fr $(BUILDDIR)
	@$(GO) clean