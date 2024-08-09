package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (c *httpClient) PostRequest(ctx context.Context, endpoint string, reqBody any, authorization string) ([]byte, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(jsonBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("%s", string(bodyBytes))
	}
	return bodyBytes, nil
}

func (c *httpClient) GetRequest(ctx context.Context, endpoint string, params map[string]interface{}, authorization string) ([]byte, error) {
	client := &http.Client{}
	query := url.Values{}
	for k, v := range params {
		query[k] = []string{fmt.Sprintf("%v", v)}
	}
	req, err := http.NewRequest(http.MethodGet, endpoint+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authorization)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("%s", string(bodyBytes))
	}
	return bodyBytes, nil
}
