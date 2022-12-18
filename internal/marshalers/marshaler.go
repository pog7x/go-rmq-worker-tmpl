package marshalers

import (
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	streadWay "github.com/streadway/amqp"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type CustomMarshaler struct {
	amqp.DefaultMarshaler
}

func (m CustomMarshaler) Unmarshal(amqpMsg streadWay.Delivery) (*message.Message, error) {
	msgUUIDStr, err := m.unmarshalMessageUUID(amqpMsg)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(msgUUIDStr, amqpMsg.Body)
	msg.Metadata = make(message.Metadata, len(amqpMsg.Headers)-1)

	for key, value := range amqpMsg.Headers {
		if key == m.computeMessageUUIDHeaderKey() || key == "x-death" {
			continue
		}

		var ok bool
		msg.Metadata[key], ok = value.(string)
		if !ok {
			return nil, errors.Errorf("metadata %s is not a string, but %#v", key, value)
		}
	}

	return msg, nil
}

func (m CustomMarshaler) unmarshalMessageUUID(amqpMsg streadWay.Delivery) (string, error) {
	var msgUUIDStr string

	msgUUID, hasMsgUUID := amqpMsg.Headers[m.computeMessageUUIDHeaderKey()]
	if !hasMsgUUID {
		return "", nil
	}

	msgUUIDStr, hasMsgUUID = msgUUID.(string)
	if !hasMsgUUID {
		return "", errors.Errorf("message UUID is not a string, but: %#v", msgUUID)
	}

	return msgUUIDStr, nil
}

func (m CustomMarshaler) computeMessageUUIDHeaderKey() string {
	if m.MessageUUIDHeaderKey != "" {
		return m.MessageUUIDHeaderKey
	}

	return amqp.DefaultMessageUUIDHeaderKey
}
