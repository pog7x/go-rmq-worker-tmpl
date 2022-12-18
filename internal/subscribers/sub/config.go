package sub

import (
	"context"
	"fmt"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/config"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/marshalers"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/pkg/errors"
	streadWayAmqp "github.com/streadway/amqp"
	"go.uber.org/zap"
)

func GetSubscriberByConfig(
	ctx context.Context,
	cfg *config.Config,
	wmLogger watermill.LoggerAdapter,
	logger *zap.Logger,
) (*Subscriber, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.RMQUser, cfg.RMQPassword, cfg.RMQHost, cfg.RMQPort, cfg.RMQVHost)
	connConf := amqp.ConnectionConfig{AmqpURI: uri}
	subscriberConn, err := amqp.NewConnection(connConf, wmLogger)

	if err != nil {
		return nil, errors.Wrapf(err, "subscriber connection error")
	}

	subscriberConf := amqp.Config{
		Connection: connConf,
		Marshaler:  marshalers.CustomMarshaler{},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				return cfg.RMQExchange
			},
			Type:    "topic",
			Durable: true,
		},
		Queue: amqp.QueueConfig{
			GenerateName: func(topic string) string {
				return cfg.RMQQueue
			},
			Durable: true,
			Arguments: streadWayAmqp.Table{
				"x-dead-letter-exchange":    cfg.RMQExchange,
				"x-dead-letter-routing-key": cfg.RMQRoutingKey,
				"x-message-ttl":             5,
			},
		},
		QueueBind: amqp.QueueBindConfig{
			GenerateRoutingKey: func(topic string) string {
				return cfg.RMQRoutingKey
			},
		},
		Consume: amqp.ConsumeConfig{
			Qos: amqp.QosConfig{
				PrefetchCount: cfg.RMQQosPrefetchCount,
			},
		},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	}
	subscriber, err := amqp.NewSubscriberWithConnection(subscriberConf, wmLogger, subscriberConn)

	if err != nil {
		return nil, errors.Wrapf(err, "subscriber config error")
	}

	return NewSubscriber(ctx, subscriber, logger), nil
}
