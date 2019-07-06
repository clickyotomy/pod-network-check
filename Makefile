.SILENT:
.DEFAULT_GOAL := dev

# Notes:
# * Targets "dev-init" and "dev-dep" should be run independently during setup
#   or package upgrades. Delete the dependency files and start fresh before
#   running them.
#   Reference: https://github.com/kubernetes/client-go/blob/master/INSTALL.md.

# Common for all builds.
GO111MODULE := on
CGO_ENABLED := 0

CMD := pod-network-check
DEV_BIN_DIR := ./bin
SHIP_BIN_DIR := /tmp

# For building a purely static `go' binary.
# Extra flags for the linker (static, no PIC).
LDFLAGS_EXT := -static -fno-PIC

# `LDFLAGS' with symbols.
# We can use `internal' for `-linkmode' because we're not using any external
# C libraries.
LDFLAGS_SYM := -ldflags '-linkmode internal -extldflags "${LDFLAGS_EXT}"'

# `LDFLAGS' without symbols (stripped).
# We can use `internal' for `-linkmode' because we're not using any external
# linkers (like gcc, musl-gcc, etc.).
LDFLAGS_STR := -ldflags '-w -s -linkmode internal -extldflags "${LDFLAGS_EXT}"'

# For external tags (`go' and `docker').
TAG ?= dev

# Build tags to use `netgo' and `osusergo', because we're not using `cgo'
# or any "libc-backed" libraries for compiling.
# * netgo         Forces pure `go' resolver.
# * osusergo      Forces pure `go` inplementation.
# * static_build  For static builds, maybe redundant.
# * $TAG          Custom build tag.
TAGS := -tags "netgo osusergo static_build ${TAG}"
INST_SFX := -installsuffix "kube"

# All `docker' build flags.
DOCKER_FLAGS := --force-rm --rm --build-arg "TAG=${TAG}"

# All `go` build flags.
GO_FLAGS_DEV := ${TAGS} ${INST_SFX} ${LDFLAGS_SYM}
GO_FLAGS_SHIP := ${TAGS} ${INST_SFX} ${LDFLAGS_STR}

# For setting up from scratch.
dev-init:
	@mkdir -p ${DEV_BIN_DIR}
	@go mod init

dev-dep:
	@go get -u "k8s.io/client-go@master"

# Pre-build targets.
mod:
	@go mod download

chk:
	@go mod verify

tidy:
	@go mod tidy

# Temporarily switch off modules because it's not a package dependency.
lint:
	@GO111MODULE="off" go get -u "golang.org/x/lint/golint"
	@golint ./...

fmt:
	@go fmt ./...

vet:
	@go vet ./...

fix:
	@go fix ./...

# Build the binary (development).
build:
	@go build ${GO_FLAGS_DEV} -o ${DEV_BIN_DIR}/${CMD} ./...

# Develpment builds.
dev: mod tidy chk lint fmt vet fix build

# For local development (linux binary).
dev-linux: export GOOS=linux
dev-linux: export GOARCH=amd64
dev-linux: dev

# For local docker builds.
dev-docker:
	@docker build ${DOCKER_FLAGS} --tag "${CMD}:${TAG}" .

# Clean-up build artifacts (common).
clean:
	@go clean -i -r -testcache -modcache

# Clean-up local binary.
dev-clean: clean
	@rm ${DEV_BIN_DIR}/*

# Production builds.
ship: mod chk lint vet build
	@go build ${GO_FLAGS_SHIP} -o ${SHIP_BIN_DIR}/${CMD} ./...

# Clean-up production binary.
ship-clean: clean
	@rm ${SHIP_BIN_DIR}/*

.PHONY: dev-init dev-dep mod tidy chk lint fmt vet fix build dev \
        dev-linux dev-docker clean dev-clean ship ship-clean
