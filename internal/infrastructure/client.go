package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func GetMysqlConnection() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DatabaseUser,
		config.Env.DatabasePass,
		config.Env.DatabaseHost,
		config.Env.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func GetRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Env.RedisAddr,
		DB:   0,
	})
}

type HttpClient struct {
	baseURL string
}

func NewHttpClient(baseURL string) HttpClient {
	return HttpClient{
		baseURL: baseURL,
	}
}

func (c *HttpClient) PostRequest(ctx context.Context, path string, reqBody any, authorization string) ([]byte, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(jsonBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, body)
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

func (c *HttpClient) GetRequest(ctx context.Context, path string, params map[string]interface{}, authorization string) ([]byte, error) {
	client := &http.Client{}
	query := url.Values{}
	for k, v := range params {
		query[k] = []string{fmt.Sprintf("%v", v)}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+query.Encode(), nil)
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

func (c *HttpClient) UploadFile(ctx context.Context, endpoint, path string, authorization string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", authorization)
	client := &http.Client{}
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
