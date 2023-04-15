package httpclient

import "context"

type Client interface {
	PostRequest(ctx context.Context, payload interface{}) error
}
