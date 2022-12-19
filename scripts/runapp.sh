#!/usr/bin/env bash

set -xe

echo 'LAUNCHING APP' && \
exec /go-rmq-worker-tmpl/bin/main runworker
