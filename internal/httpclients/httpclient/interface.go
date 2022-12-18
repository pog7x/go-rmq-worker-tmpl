package httpclient

import "context"

type Client interface {
	PostRequest(ctx context.Context) error
}
