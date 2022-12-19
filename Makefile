PRJROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PRJNAME:=$(notdir $(PRJROOT))
APPNAME?=${PRJNAME}

BINDIR?=${PRJROOT}/bin

.PHONY: all
all: clean bindir build

.PHONY: bindir
bindir:
	mkdir -p ${BINDIR}

.PHONY: build
build: bindir
	GOBIN=${BINDIR} go install ./main.go

.PHONY: clean
clean:
	rm -rf ${BINDIR}

.PHONY: golangci
golangci:
	golangci-lint run --go=1.23

.PHONY: test
test: golangci
	go test -v -coverprofile=coverage.out -vet '' ./...

.PHONY: run-dev
run-dev:
	go run ./main.go runworker -c=./internal/app/config/.config.dev.yml

.PHONY: run-server
run-server: all
	${BINDIR}/${APPNAME}