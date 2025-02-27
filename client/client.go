package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"
)

type OpenAPIClient struct {
	ClientID   string
	SecretKey  string
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

func NewClient(clientID, secretKey, baseUrl string) *OpenAPIClient {
	return &OpenAPIClient{
		ClientID:  clientID,
		SecretKey: secretKey,
		BaseURL:   baseUrl,
		Timeout:   10 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *OpenAPIClient) Get(endpoint string, queryParams map[string][]string) ([]byte, error) {
	return c.request("GET", endpoint, queryParams, nil)
}

func (c *OpenAPIClient) Post(endpoint string, body []byte) ([]byte, error) {
	return c.request("POST", endpoint, nil, body)
}

func (c *OpenAPIClient) request(method, endpoint string, queryParams map[string][]string, body []byte) ([]byte, error) {
	// 1. 生成时间戳
	timestamp := time.Now().Format(time.RFC3339)

	// 2. 生成签名
	signature, err := GenerateSignature(c.SecretKey, method, endpoint, c.ClientID, timestamp, string(body), queryParams, "")
	if err != nil {
		return nil, err
	}

	// 3. 构造请求
	req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 4. 添加请求头
	req.Header.Set("X-DRC-Client-ID", c.ClientID)
	req.Header.Set("X-DRC-Sign", signature)
	req.Header.Set("X-DRC-Timestamp", timestamp)
	req.Header.Set("Content-Type", "application/json")

	// 5. 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 6. 处理响应
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API request failed: " + resp.Status)
	}

	return io.ReadAll(resp.Body)
}
