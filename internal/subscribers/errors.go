package subscribers

import (
	"fmt"
	"net/http"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/httpclients"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/retryable"

	"github.com/pkg/errors"
)

var (
	ErrNilDeliveredMessage = errors.New("empty delivered message")
)

func WrapHTTPClientError(msg string, err error) error {
	var errHTTP httpclients.HTTPError
	switch {
	case errors.As(err, &errHTTP) &&
		(errHTTP.Code == http.StatusRequestTimeout ||
			errHTTP.Code == http.StatusTooManyRequests ||
			errHTTP.Code >= http.StatusInternalServerError):
		err = retryable.NewRetryableError(errHTTP)
	case errors.Is(err, httpclients.ErrExecRequest):
		err = retryable.NewRetryableError(err)
	}
	return fmt.Errorf("%s: %w", msg, err)
}
