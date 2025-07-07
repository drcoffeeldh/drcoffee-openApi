package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type OpenAPIClient struct {
	ClientID   string
	SecretKey  string
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

func NewClient(clientID, secretKey, baseURL string) *OpenAPIClient {
	return &OpenAPIClient{
		ClientID:  clientID,
		SecretKey: secretKey,
		BaseURL:   baseURL,
		Timeout:   10 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *OpenAPIClient) Get(endpoint string, queryParams map[string][]string, customHeaders map[string]string) (map[string]interface{}, error) {
	return c.request("GET", endpoint, queryParams, nil, customHeaders)
}

func (c *OpenAPIClient) Post(endpoint string, body interface{}, queryParams map[string][]string, customHeaders map[string]string) (map[string]interface{}, error) {
	return c.request("POST", endpoint, queryParams, body, customHeaders)
}

func (c *OpenAPIClient) request(method, endpoint string, queryParams map[string][]string, body interface{}, customHeaders map[string]string) (map[string]interface{}, error) {
	// 1. 处理请求体
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
	} else {
		bodyBytes = []byte("") // 修改为真正的空字符串
	}

	// 2. 生成时间戳
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	// 3. 提取 Authorization Token (如果存在)
	authToken := ""
	if val, ok := customHeaders["Authorization"]; ok {
		authToken = val
	}

	// 4. 生成签名
	// 将提取出的 authToken 传递给 GenerateSignature 的 token 参数
	signature, err := GenerateSignature(c.SecretKey, method, endpoint, c.ClientID, timestamp, string(bodyBytes), queryParams, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %v", err)
	}

	// 5. 构造 URL
	reqURL, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// 6. 添加查询参数
	if queryParams != nil {
		query := reqURL.Query()
		for key, values := range queryParams {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		reqURL.RawQuery = query.Encode()
	}

	// 7. 创建请求
	req, err := http.NewRequest(method, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 8. 添加默认请求头
	req.Header.Set("X-DRC-Client-ID", c.ClientID)
	req.Header.Set("X-DRC-Sign", signature)
	req.Header.Set("X-DRC-Timestamp", timestamp)
	req.Header.Set("Content-Type", "application/json") // 确保 Content-Type 始终设置

	// 9. 添加自定义请求头 (此时 Authorization 应该已经被提取并处理，不会再次添加)
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	// 10. 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 11. 处理响应
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed: %s - %s", resp.Status, string(responseBody))
	}

	// 12. 解析响应体
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return responseData, nil
}
