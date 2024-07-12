FROM golangci/golangci-lint:latest AS linter
FROM golang:1.22.2

COPY --from=linter /usr/bin/golangci-lint /bin/golangci-lint
RUN mkdir /go-rmq-worker-tmpl
COPY . /go-rmq-worker-tmpl
WORKDIR /go-rmq-worker-tmpl

RUN make
