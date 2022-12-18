# go-rmq-worker-tmpl

[![Build Status](https://github.com/pog7x/go-rmq-worker-tmpl/actions/workflows/go.yml/badge.svg)](https://github.com/pog7x/go-rmq-worker-tmpl/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/pog7x/go-rmq-worker-tmpl/blob/master/LICENSE)

## Golang RMQ asynchronous worker template

### Technologies used:
- Logger: zap logger
- RMQ: amqp, ThreeDotsLabs/watermill
- Metrics: prometheus/client_golang, logrus-sentry
- Linter: golangci/golangci-lint
- Tests: testify

### Run dev with docker-compose
```bash
docker-compose -f docker-compose.dev.yml up -d 
```

### Run tests
```bash
./scripts/test.sh
```