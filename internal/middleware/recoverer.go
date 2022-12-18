package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Recoverer struct {
	Logger watermill.LoggerAdapter
}

type recoveredPanicError struct {
	V          interface{}
	Stacktrace string
}

func (p recoveredPanicError) Error() string {
	return fmt.Sprintf("panic occurred: %#v, stacktrace: \n%s", p.V, p.Stacktrace)
}

func (rr Recoverer) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (messages []*message.Message, err error) {
		defer func() {
			if r := recover(); r != nil {
				panicErr := errors.WithStack(recoveredPanicError{V: r, Stacktrace: string(debug.Stack())})
				err = multierror.Append(err, panicErr)
				if rr.Logger != nil {
					rr.Logger.Error(fmt.Sprintf("panic on message %s", msg.UUID), err, watermill.LogFields{
						"payload":  string(msg.Payload),
						"metadata": msg.Metadata,
					})
				}
				msg.Ack()
			}
		}()

		return h(msg)
	}
}
