package httpclients

import (
	"errors"
	"fmt"
)

var (
	ErrExecRequest    = errors.New("exec request error")
	ErrDecodeResponse = errors.New("decode response error")
)

type HTTPError struct {
	Code    int
	Method  string
	Message string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s %d: %s", e.Method, e.Code, e.Message)
}
