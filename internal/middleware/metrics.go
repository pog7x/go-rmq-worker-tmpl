package middleware

import (
	"time"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/metrics"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/retryable"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct{}

// own custom metrics middleware,
// issue: limited watermill metrics middleware without the possibility to customize
func (m Metrics) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (messages []*message.Message, err error) {
		labels := prometheus.Labels{
			"handler": message.HandlerNameFromCtx(msg.Context()),
			"err":     "false",
			"retries": "false",
		}

		defer func(now time.Time) {
			if err != nil {
				labels["err"] = "true"
				if retryable.IsRetryableError(err) {
					labels["retries"] = "true"
				}
			}
			metrics.HandlerDuration.With(labels).Observe(time.Since(now).Seconds())
		}(time.Now())

		return h(msg)
	}
}
