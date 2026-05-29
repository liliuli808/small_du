package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	client    *http.Client
	userAgent string
	cookie    string
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(userAgent, cookie string, timeout int) *HTTPClient {
	if timeout <= 0 {
		timeout = 15
	}
	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		userAgent: userAgent,
		cookie:    cookie,
	}
}

// Get 发送GET请求
func (c *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	return c.client.Do(req)
}

// GetJSON 发送GET请求并解析JSON
func (c *HTTPClient) GetJSON(ctx context.Context, url string, v interface{}) error {
	resp, err := c.Get(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// GetString 发送GET请求并返回字符串
func (c *HTTPClient) GetString(ctx context.Context, url string) (string, error) {
	resp, err := c.Get(ctx, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// setHeaders 设置请求头
func (c *HTTPClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json, text/html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}
}
