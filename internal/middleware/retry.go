package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/retryable"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cenkalti/backoff/v3"
)

// Retry provides a middleware that retries the handler if retryable errors are returned.
// The retry behaviour is configurable, with exponential backoff and maximum elapsed time.
type Retry struct {
	// MaxRetries is maximum number of times a retry will be attempted.
	MaxRetries int

	// InitialInterval is the first interval between retries. Subsequent intervals will be scaled by Multiplier.
	InitialInterval time.Duration
	// MaxInterval sets the limit for the exponential backoff of retries. The interval will not be increased beyond MaxInterval.
	MaxInterval time.Duration
	// Multiplier is the factor by which the waiting interval will be multiplied between retries.
	Multiplier float64
	// MaxElapsedTime sets the time limit of how long retries will be attempted. Disabled if 0.
	MaxElapsedTime time.Duration
	// RandomizationFactor randomizes the spread of the backoff times within the interval of:
	// [currentInterval * (1 - randomization_factor), currentInterval * (1 + randomization_factor)].
	RandomizationFactor float64

	// OnRetryHook is an optional function that will be executed on each retry attempt.
	// The number of the current retry is passed as retryNum,
	OnRetryHook func(retryNum int, delay time.Duration)

	Logger watermill.LoggerAdapter
}

const (
	MaxRetriesCountKey  = "max_retries_count"
	CurrRetriesCountKey = "curr_retries_count"
)

func (r Retry) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		var retryNum int

		// access to current try inside handler issue
		msg.Metadata.Set(MaxRetriesCountKey, strconv.Itoa(r.MaxRetries))
		msg.Metadata.Set(CurrRetriesCountKey, strconv.Itoa(retryNum))

		producedMessages, err := h(msg)
		if err == nil {
			return producedMessages, nil
		}

		if !retryable.IsRetryableError(err) {
			msg.Ack()
			return nil, err
		}

		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = r.InitialInterval
		expBackoff.MaxInterval = r.MaxInterval
		expBackoff.Multiplier = r.Multiplier
		expBackoff.MaxElapsedTime = r.MaxElapsedTime
		expBackoff.RandomizationFactor = r.RandomizationFactor

		ctx := msg.Context()
		if r.MaxElapsedTime > 0 {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, r.MaxElapsedTime)
			defer cancel()
		}

		retryNum = 1
		expBackoff.Reset()
		for {
			waitTime := expBackoff.NextBackOff()
			select {
			case <-ctx.Done():
				return producedMessages, err
			case <-time.After(waitTime):
				// go on
			}

			msg.Metadata.Set(CurrRetriesCountKey, strconv.Itoa(retryNum))

			producedMessages, err = h(msg)
			if err == nil {
				return producedMessages, nil
			}

			if r.Logger != nil {
				r.Logger.Info(fmt.Sprintf("Error occurred: %v, retrying", err), watermill.LogFields{
					"retry_no":     retryNum,
					"max_retries":  r.MaxRetries,
					"wait_time":    waitTime,
					"elapsed_time": expBackoff.GetElapsedTime(),
				})
			}
			if r.OnRetryHook != nil {
				r.OnRetryHook(retryNum, waitTime)
			}

			retryNum++
			if retryNum > r.MaxRetries {
				break
			}
		}

		if r.Logger != nil {
			r.Logger.Error("message retries exceeded", nil, watermill.LogFields{
				"max_retries": r.MaxRetries,
				"message_id":  msg.UUID,
				"headers":     msg.Metadata,
				"payload":     msg.Payload,
			})
		}

		msg.Ack()
		return nil, err
	}
}
