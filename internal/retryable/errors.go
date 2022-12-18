package retryable

import (
	"errors"
	"fmt"
)

type retryableError struct {
	err error
}

func NewRetryableError(err error) error {
	return &retryableError{err: err}
}

func IsRetryableError(err error) bool {
	var re *retryableError
	return errors.As(err, &re)
}

func (e *retryableError) Error() string {
	if e.err == nil {
		return "retryable: nil"
	}
	return fmt.Sprintf("retryable: %s", e.err.Error())
}

func (e *retryableError) Unwrap() error {
	return e.err
}
