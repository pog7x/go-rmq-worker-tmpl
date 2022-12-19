package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/httpclients"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/metrics"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type HTTPClient struct {
	logger  *zap.Logger
	client  *resty.Client
	baseURL string
}

func NewHTTPClient(logger *zap.Logger, baseURL string, timeout time.Duration) Client {
	return &HTTPClient{
		logger:  logger,
		baseURL: baseURL,
		client: resty.New().
			SetHeader("Content-Type", "application/json").
			SetTimeout(timeout).
			SetRetryCount(3).
			SetRetryWaitTime(200 * time.Millisecond).
			AddRetryCondition(func(r *resty.Response, _ error) bool {
				return r.StatusCode() >= http.StatusInternalServerError
			}),
	}
}

func (c *HTTPClient) PostRequest(ctx context.Context, payload interface{}) error {
	const (
		clientName   = "HTTPClient"
		method       = resty.MethodPost
		endpointName = "PostRequest"
		basePath     = "/path/to/post/request"
	)

	url := c.baseURL + basePath
	c.logger.Info(
		"Sending HTTP request",
		zap.String("service", clientName),
		zap.String("endpoint", endpointName),
		zap.String("url", url),
		zap.String("method", method),
		zap.Any("body", payload),
	)

	resp, err := c.client.R().SetContext(ctx).SetBody(payload).Execute(method, url)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s request error", clientName), zap.Error(err))
		return fmt.Errorf("%w: %v", httpclients.ErrExecRequest, err)
	}

	metrics.ResponseTimeExternal.With(map[string]string{
		"service": clientName,
		"code":    strconv.Itoa(resp.StatusCode()),
		"handler": endpointName,
		"method":  method,
	}).Observe(float64(resp.Time().Milliseconds()))

	if resp.IsError() {
		return httpclients.HTTPError{
			Code:    resp.StatusCode(),
			Method:  resp.Request.Method,
			Message: resp.String(),
		}
	}

	return nil
}
