# Path parameters
GO ?= $(shell which go  2>/dev/null)
DOCKER ?= $(shell which docker 2>/dev/null)
WASMBUILD ?= $(shell which wasmbuild 2>/dev/null)
BUILDDIR ?= build
CMDDIR=$(wildcard cmd/*)
WASMDIR=$(wildcard wasm/*)

# Set OS and Architecture
ARCH ?= $(shell arch | tr A-Z a-z | sed 's/x86_64/amd64/' | sed 's/i386/amd64/' | sed 's/armv7l/arm/' | sed 's/aarch64/arm64/')
OS ?= $(shell uname | tr A-Z a-z)
VERSION ?= $(shell git describe --tags --always | sed 's/^v//')

# Build flags
BUILD_MODULE = $(shell cat go.mod | head -1 | cut -d ' ' -f 2)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitSource=${BUILD_MODULE}
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitTag=$(shell git describe --tags --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitBranch=$(shell git name-rev HEAD --name-only --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitHash=$(shell git rev-parse HEAD)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_FLAGS = -ldflags="-s -w ${BUILD_LD_FLAGS}"

# Docker
DOCKER_REPO ?= ghcr.io/mutablelogic/pgmanager
DOCKER_SOURCE ?= ${BUILD_MODULE}
DOCKER_TAG = ${DOCKER_REPO}-${OS}-${ARCH}:${VERSION}

# All targets
all: tidy $(CMDDIR)

# Rules for building
.PHONY: $(CMDDIR)
$(CMDDIR): go-dep mkdir
	@echo 'go build $@'
	@rm -rf ${BUILDDIR}/$(shell basename $@)
	@$(GO) build -tags frontend $(BUILD_FLAGS) -o ${BUILDDIR}/$(shell basename $@) ./$@

# Rules for building
.PHONY: $(WASMDIR)
$(WASMDIR): go-dep wasmbuild-dep mkdir
	@echo 'wasmbuild $@'
	@$(GO) get github.com/djthorpe/go-wasmbuild/pkg/bootstrap github.com/djthorpe/go-wasmbuild/pkg/bootstrap/extra github.com/djthorpe/go-wasmbuild/pkg/mvc
	@${WASMBUILD} build --go-flags='$(BUILD_FLAGS)' -o ${BUILDDIR}/wasm/$(shell basename $@) ./$@ && \
		mv ${BUILDDIR}/wasm/$(shell basename $@)/wasm_exec.html ${BUILDDIR}/wasm/$(shell basename $@)/index.html && \
		cp etc/embed.go ${BUILDDIR}/wasm/$(shell basename $@)/

# Build pgmanager with embedded frontend
.PHONY: pgmanager
pgmanager: wasm/pgmanager cmd/pgmanager

# Build the docker image
.PHONY: docker
docker: docker-dep ${NPM_DIR}
	@echo build docker image ${DOCKER_TAG} OS=${OS} ARCH=${ARCH} SOURCE=${DOCKER_SOURCE} VERSION=${VERSION}
	@${DOCKER} build \
		--tag ${DOCKER_TAG} \
		--build-arg ARCH=${ARCH} \
		--build-arg OS=${OS} \
		--build-arg SOURCE=${DOCKER_SOURCE} \
		--build-arg VERSION=${VERSION} \
		-f etc/Dockerfile .

# Push docker container
.PHONY: docker-push
docker-push: docker-dep 
	@echo push docker image: ${DOCKER_TAG}
	@${DOCKER} push ${DOCKER_TAG}

# Print out the version
.PHONY: docker-version
docker-version: docker-dep 
	@echo "tag=${VERSION}"

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

.PHONY: go-dep tidy
tidy: mkdir
	@echo 'go tidy'
	@install -d ${BUILDDIR}/wasm/pgmanager
	@cp -n etc/embed.go ${BUILDDIR}/wasm/pgmanager/ 2>/dev/null || true
	@echo 'module github.com/mutablelogic/go-pg/build/wasm/pgmanager' > ${BUILDDIR}/wasm/pgmanager/go.mod
	@$(GO) mod tidy

.PHONY: clean
clean: tidy
	@echo 'clean'
	@rm -fr $(BUILDDIR)
	@$(GO) clean

###############################################################################
# DEPENDENCIES

.PHONY: go-dep
go-dep:
	@test -f "${GO}" && test -x "${GO}"  || (echo "Missing go binary" && exit 1)

.PHONY: docker-dep
docker-dep:
	@test -f "${DOCKER}" && test -x "${DOCKER}"  || (echo "Missing docker binary" && exit 1)

.PHONY: wasmbuild-dep
wasmbuild-dep:
	@test -f "${WASMBUILD}" && test -x "${WASMBUILD}"  || (echo "Missing wasmbuild binary" && exit 1)