PRJROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PRJNAME:=$(notdir $(PRJROOT))
APPNAME?=${PRJNAME}

BINDIR?=${PRJROOT}/bin
# BIN?=${BINDIR}/${APPNAME}

# VER?=$(shell cd ${PRJROOT}; git describe --all)
# GOVER:=$(shell go version)
# GOPATH:=$(shell go env GOPATH)

.PHONY: all
all: clean bindir build

.PHONY: bindir
bindir:
	mkdir -p ${BINDIR}

.PHONY: build
build: bindir
	GOBIN=${BINDIR} go install ./cmd/...

.PHONY: clean
clean:
	rm -rf ${BINDIR}

.PHONY: golangci
golangci:
	golangci-lint run --go=1.19

.PHONY: test
test: golangci
	go test -v -coverprofile=coverage.out -vet '' ./...

.PHONY: run-dev
run-dev:
	go run ./cmd/...

.PHONY: run-server
run-server: all
	${BINDIR}/${APPNAME}