package middleware

import (
	"errors"
	"testing"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/retryable"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	testCases := []struct {
		name            string
		retries         int
		err             error
		expectedRetries int
	}{
		{
			name:            "not retryable",
			retries:         3,
			err:             errors.New("some"),
			expectedRetries: 1,
		},
		{
			name:            "retryable",
			retries:         3,
			err:             retryable.NewRetryableError(errors.New("some")),
			expectedRetries: 4,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			retry := Retry{
				MaxRetries: tc.retries,
			}

			runCount := 0
			h := retry.Middleware(func(msg *message.Message) (messages []*message.Message, e error) {
				runCount++
				return nil, tc.err
			})

			_, err := h(message.NewMessage("foobar", nil))

			assert.Equal(t, tc.expectedRetries, runCount)
			assert.EqualError(t, err, tc.err.Error())
		})
	}
}
