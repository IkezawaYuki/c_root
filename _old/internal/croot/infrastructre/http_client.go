package infrastructre

import (
	"context"
)

type HttpClient interface {
	GetRequest(ctx context.Context, endpoint string, params map[string]interface{}, authorization string) ([]byte, error)
	PostRequest(ctx context.Context, endpoint string, reqBody any, authorization string) ([]byte, error)
}

type httpClient struct {
}

func NewHttpClient() HttpClient {
	return &httpClient{}
}
