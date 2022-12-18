package config

import "time"

type Config struct {
	ServerListenAddr   string        `envconfig:"SERVER_LISTEN_ADDR" default:"0.0.0.0:8080"`
	ServerReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"3s"`
	ServerWriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"3s"`
	LogLevel           string        `envconfig:"LOG_LEVEL" default:"debug"`
	SentryDSN          string        `envconfig:"SENTRY_DSN"`

	MaxRetries      int           `envconfig:"MAX_RETRIES" default:"5"`
	RetriesInterval time.Duration `envconfig:"RETRIES_INTERVAL" default:"3s"`

	RMQUser             string `envconfig:"RMQ_USER" required:"true"`
	RMQPassword         string `envconfig:"RMQ_PASSWORD" required:"true"`
	RMQHost             string `envconfig:"RMQ_HOST" required:"true"`
	RMQPort             int    `envconfig:"RMQ_PORT" required:"true"`
	RMQQosPrefetchCount int    `envconfig:"RMQ_QOS_PREFETCH_COUNT" default:"1"`

	RMQVHost string `envconfig:"RMQ_VHOST" required:"true"`

	RMQExchange   string `envconfig:"RMQ_EXCHANGE" required:"true"`
	RMQQueue      string `envconfig:"RMQ_QUEUE" required:"true"`
	RMQRoutingKey string `envconfig:"RMQ_ROUTING_KEY" required:"true"`
}
