# go-rmq-worker-tmpl

[![Build Status](https://github.com/pog7x/go-rmq-worker-tmpl/actions/workflows/go.yml/badge.svg)](https://github.com/pog7x/go-rmq-worker-tmpl/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/pog7x/go-rmq-worker-tmpl/blob/master/LICENSE)

Production-ready template for building asynchronous RabbitMQ workers in Go.

## Features

- **Message routing** via [ThreeDotsLabs/watermill](https://github.com/ThreeDotsLabs/watermill) + [watermill-amqp](https://github.com/ThreeDotsLabs/watermill-amqp)
- **Middleware pipeline**: structured logging, configurable retry with backoff, panic recovery, Prometheus metrics
- **Diagnostic HTTP server**: `/ping` healthcheck, `/metrics` (Prometheus), `/debug/pprof/*` profiling endpoints
- **Structured logging** via [uber-go/zap](https://github.com/uber-go/zap) with optional [Sentry](https://github.com/TheZeroSlave/zapsentry) integration
- **Graceful shutdown** on OS signal
- **Environment-based configuration** via [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)

## Project Structure

```
.
├── cmd/                    # Entrypoint
├── internal/
│   ├── app/
│   │   ├── config/         # Configuration (envconfig)
│   │   ├── logger/         # Zap logger setup
│   │   └── server/         # Diagnostic HTTP server
│   ├── middleware/         # Logging, retry, recoverer, metrics
│   ├── publishers/pub/     # AMQP publisher
│   ├── subscribers/sub/    # AMQP subscriber + message handler
│   ├── marshalers/         # Message marshaling
│   ├── metrics/            # Prometheus metrics
│   ├── httpclients/        # HTTP client wrapper
│   └── retryable/          # Retryable error types
├── scripts/
└── docker-compose.dev.yml
```

## Configuration

All configuration is done via environment variables:

| Variable                  | Required | Default        | Description                        |
|---------------------------|----------|----------------|------------------------------------|
| `RMQ_USER`                | yes      | —              | RabbitMQ username                  |
| `RMQ_PASSWORD`            | yes      | —              | RabbitMQ password                  |
| `RMQ_HOST`                | yes      | —              | RabbitMQ host                      |
| `RMQ_PORT`                | yes      | —              | RabbitMQ port                      |
| `RMQ_VHOST`               | yes      | —              | RabbitMQ virtual host              |
| `RMQ_EXCHANGE`            | yes      | —              | Exchange name                      |
| `RMQ_QUEUE`               | yes      | —              | Queue name                         |
| `RMQ_ROUTING_KEY`         | yes      | —              | Routing key                        |
| `RMQ_QOS_PREFETCH_COUNT`  | no       | `1`            | QoS prefetch count                 |
| `MAX_RETRIES`             | no       | `5`            | Max message processing retries     |
| `RETRIES_INTERVAL`        | no       | `3s`           | Initial retry backoff interval     |
| `SERVER_LISTEN_ADDR`      | no       | `0.0.0.0:8080` | Diagnostic HTTP server address     |
| `SERVER_READ_TIMEOUT`     | no       | `3s`           | HTTP server read timeout           |
| `SERVER_WRITE_TIMEOUT`    | no       | `3s`           | HTTP server write timeout          |
| `LOG_LEVEL`               | no       | `debug`        | Log level                          |
| `SENTRY_DSN`              | no       | —              | Sentry DSN for error reporting     |

## Getting Started

### Run with Docker Compose (development)

```bash
docker-compose -f docker-compose.dev.yml up -d
```

### Build and run locally

```bash
make build
make run-server
```

### Run in dev mode (without building binary)

```bash
make run-dev
```

## Development

### Run linter + tests

```bash
make test
```

Or just tests (without lint):

```bash
./scripts/test.sh
```

### Run linter only

```bash
make golangci
```

## Diagnostic Endpoints

| Endpoint              | Description                  |
|-----------------------|------------------------------|
| `GET /ping`           | Health check                 |
| `GET /metrics`        | Prometheus metrics           |
| `GET /debug/pprof/*`  | pprof profiling endpoints    |

## Extending the Worker

The message handler is located in `internal/subscribers/sub/subscriber.go`. Implement your business logic inside the `Handler` method:

```go
func (s *Subscriber) Handler(msg *message.Message) error {
    // your logic here
    return nil
}
```

To make an error retryable (triggers the retry middleware), wrap it with the retryable error type from `internal/retryable/errors.go`.

## Related

- [Branch with cobra + viper CLI](https://github.com/pog7x/go-rmq-worker-tmpl/tree/master-cobra) — variant using [spf13/cobra](https://github.com/spf13/cobra) + [spf13/viper](https://github.com/spf13/viper)

## License

[MIT](LICENSE)
