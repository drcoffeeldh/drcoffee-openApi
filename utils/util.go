package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const TokenHeader = "Authorization"

// 构建 Canonical Query String
func BuildCanonicalQueryString(queryParam map[string][]string) string {
	// 提取并排序键
	keys := make([]string, 0, len(queryParam))
	for key := range queryParam {
		keys = append(keys, key)
	}
	sort.Strings(keys) // 按键排序

	// 构建排序后的键值对
	var queryParts []string
	for _, key := range keys {
		values := queryParam[key]
		sort.Strings(values) // 按值排序（如果有多个值）
		for _, value := range values {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}
	}

	// 将键值对拼接为字符串
	return strings.Join(queryParts, "&")
}

// BuildCanonicalHeaders 构建 Canonical Headers（动态支持 token）
func BuildCanonicalHeaders(headers map[string]string, token string) string {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var headerParts []string

	for _, key := range keys {
		headerParts = append(headerParts, fmt.Sprintf("%s:%s", key, headers[key]))
	}
	if token != "" {
		headerParts = append(headerParts, fmt.Sprintf("%s:%s", TokenHeader, headers[token]))
	}
	return strings.Join(headerParts, "\n")
}

func MD5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func GenerateHMACSHA1(secret, data string) []byte {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return h.Sum(nil)
}
