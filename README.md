# go-rmq-worker-tmpl

[![Build Status](https://github.com/pog7x/go-rmq-worker-tmpl/actions/workflows/go.yml/badge.svg)](https://github.com/pog7x/go-rmq-worker-tmpl/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/pog7x/go-rmq-worker-tmpl/blob/master/LICENSE)

## Golang RMQ asynchronous worker template

### Technologies used:
- Logger: [zap](https://github.com/uber-go/zap)
- RMQ: [streadway/amqp](https://github.com/streadway/amqp), [ThreeDotsLabs/watermill](https://github.com/ThreeDotsLabs/watermill)
- Metrics: [prometheus/client_golang](https://github.com/prometheus/client_golang), [TheZeroSlave/zapsentry](https://github.com/TheZeroSlave/zapsentry)
- Linter: [golangci/golangci-lint](https://github.com/golangci/golangci-lint)
- Tests: [stretchr/testify](https://github.com/stretchr/testify)

### Run dev with docker-compose
```bash
docker-compose -f docker-compose.dev.yml up -d 
```

### Run tests
```bash
./scripts/test.sh
```