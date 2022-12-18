package sub

import (
	"context"
	"fmt"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/subscribers"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type Subscriber struct {
	ctx        context.Context
	logger     *zap.Logger
	subscriber message.Subscriber
}

func NewSubscriber(c context.Context, sub message.Subscriber, logger *zap.Logger) *Subscriber {
	return &Subscriber{
		ctx:        c,
		logger:     logger,
		subscriber: sub,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return s.subscriber.Subscribe(ctx, topic)
}

func (s *Subscriber) Close() error {
	return s.subscriber.Close()
}

func (s *Subscriber) Handler(msg *message.Message) error {
	if msg == nil {
		return fmt.Errorf("request: %w", subscribers.ErrNilDeliveredMessage)
	}

	s.logger.Sugar().Infof("Received new message %v", msg)

	return nil
}
