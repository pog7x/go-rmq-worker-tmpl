package middleware

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Logging struct {
	Logger watermill.LoggerAdapter
}

func (l Logging) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		if l.Logger != nil {
			l.Logger.Info("incoming message", map[string]interface{}{
				"message_id": msg.UUID,
				"payload":    string(msg.Payload),
				"metadata":   msg.Metadata,
			})
		}

		processedMessages, err := h(msg)
		if err != nil {
			if l.Logger != nil {
				l.Logger.Info("handler error", map[string]interface{}{
					"message_id": msg.UUID,
					"error":      err.Error(),
				})
			}
		}
		return processedMessages, err
	}
}
