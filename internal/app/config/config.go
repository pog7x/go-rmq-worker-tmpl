package config

import "time"

type Config struct {
	ServerListenAddr   string        `mapstructure:"SERVER_LISTEN_ADDR"`
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	LogLevel           string        `mapstructure:"LOG_LEVEL"`
	SentryDSN          string        `mapstructure:"SENTRY_DSN"`

	MaxRetries      int           `mapstructure:"MAX_RETRIES"`
	RetriesInterval time.Duration `mapstructure:"RETRIES_INTERVAL"`

	RMQUser             string `mapstructure:"RMQ_USER"`
	RMQPassword         string `mapstructure:"RMQ_PASSWORD"`
	RMQHost             string `mapstructure:"RMQ_HOST"`
	RMQPort             int    `mapstructure:"RMQ_PORT"`
	RMQQosPrefetchCount int    `mapstructure:"RMQ_QOS_PREFETCH_COUNT"`

	RMQVHost string `mapstructure:"RMQ_VHOST"`

	RMQExchange   string `mapstructure:"RMQ_EXCHANGE"`
	RMQQueue      string `mapstructure:"RMQ_QUEUE"`
	RMQRoutingKey string `mapstructure:"RMQ_ROUTING_KEY"`
}

var Configuration = new(Config)
