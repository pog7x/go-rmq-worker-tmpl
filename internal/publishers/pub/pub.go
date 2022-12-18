package pub

import (
	"fmt"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/config"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	streadWayAmqp "github.com/streadway/amqp"
)

func GetPublisherByConfig(cfg *config.Config, wmLogger watermill.LoggerAdapter) (message.Publisher, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.RMQUser, cfg.RMQPassword, cfg.RMQHost, cfg.RMQPort, cfg.RMQVHost)
	connConf := amqp.ConnectionConfig{AmqpURI: uri}
	publisherConn, err := amqp.NewConnection(connConf, wmLogger)

	if err != nil {
		return nil, errors.Wrapf(err, "publisher connection error")
	}

	publisherConf := amqp.Config{
		Connection: connConf,
		Marshaler:  amqp.DefaultMarshaler{},
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
		Publish: amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string {
				return cfg.RMQRoutingKey
			},
		},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	}

	publisher, err := amqp.NewPublisherWithConnection(publisherConf, wmLogger, publisherConn)
	if err != nil {
		return nil, errors.Wrapf(err, "auto call publisher config error")
	}

	err = declareQueueForPublishConfig(wmLogger, publisherConf, publisherConn)
	if err != nil {
		return nil, fmt.Errorf("declare queue for publish: %v", err)
	}

	return publisher, nil
}

func declareQueueForPublishConfig(
	wmLogger watermill.LoggerAdapter,
	publisherConf amqp.Config,
	conn *amqp.ConnectionWrapper,
) error {
	ch, err := conn.Connection().Channel()
	if err != nil {
		return fmt.Errorf("unable to connect to channel %v", ch)
	}

	defer func() {
		err = ch.Close()

		if err != nil {
			wmLogger.Error("unable to close channel", err, map[string]interface{}{})
		}
	}()

	topic := ""

	logFields := watermill.LogFields{"topic": topic}

	queueName := publisherConf.Queue.GenerateName(topic)
	logFields["amqp_queue_name"] = queueName

	exchangeName := publisherConf.Exchange.GenerateName(topic)
	logFields["amqp_exchange_name"] = exchangeName

	err = publisherConf.TopologyBuilder.BuildTopology(ch, queueName, exchangeName, publisherConf, wmLogger)
	if err != nil {
		return fmt.Errorf("unable to build topology %v", ch)
	}

	return nil
}
