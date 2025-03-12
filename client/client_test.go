package client

import (
	"testing"
	"time"
)

// 测试 GET 请求
func TestGetRequest(t *testing.T) {
	// 创建客户端
	client := NewClient("client-20250220135412-d448dea560a6f93ef564f58a2898c23998978204eb6d3dcc0bf7e1195c2fd4",
		"SEACE3Z0JB5ZSRUY83AJF1EZI70FBEWSUSMAVOMNXFE", "https://dev.t8s.dr-coffee.cn")

	// 发送 GET 请求
	queryParams := map[string][]string{
		"page":  {"1"},
		"limit": {"10"},
	}
	response, err := client.Get("/v1/api-gateway/internal/token/endpoint", queryParams, nil)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	// 验证响应
	if response["status"] != "success" || response["message"] != "GET request received" {
		t.Errorf("Expected response to be 'success', got '%v'", response)
	}
}

// 测试 POST 请求
func TestPostRequest(t *testing.T) {

	// 创建客户端
	client := NewClient("client-20250220135412-d448dea560a6f93ef564f58a2898c23998978204eb6d3dcc0bf7e1195c2fd4",
		"SEACE3Z0JB5ZSRUY83AJF1EZI70FBEWSUSMAVOMNXFE", "https://dev.t8s.dr-coffee.cn")

	// 发送 POST 请求
	body := map[string]interface{}{
		"name":  "John Doe",
		"email": "john.doe@example.com",
	}
	response, err := client.Post("/v1/api-gateway/internal/token/endpoint", body, nil, nil)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	// 验证响应
	if response["status"] != "success" || response["message"] != "POST request received" {
		t.Errorf("Expected response to be 'success', got '%v'", response)
	}
}

// 测试请求失败场景
func TestRequestFailure(t *testing.T) {
	// 创建客户端
	client := NewClient("client-123", "your-secret-key", "http://localhost:80")

	// 发送 GET 请求
	_, err := client.Get("/v1/resource", nil, nil)
	if err == nil {
		t.Fatal("Expected request to fail, but it succeeded")
	}

	// 验证错误信息
	expectedError := "API request failed: 400 Bad Request - {\"status\":\"error\",\"message\":\"Invalid request\"}\n"
	if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', got '%v'", expectedError, err)
	}
}

// 测试签名生成
func TestGenerateSignature(t *testing.T) {
	secret := "your-secret-key"
	method := "POST"
	uri := "/v1/users"
	clientID := "client-123"
	timestamp := time.Now().Format(time.RFC3339)
	body := `{"name":"John Doe","email":"john.doe@example.com"}`
	queryParams := map[string][]string{
		"page":  {"1"},
		"limit": {"10"},
	}

	// 生成签名
	signature, err := GenerateSignature(secret, method, uri, clientID, timestamp, body, queryParams, "")
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	// 验证签名格式
	if signature == "" {
		t.Error("Expected signature to be non-empty, got empty string")
	}
}
