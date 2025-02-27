package client

import (
	"encoding/base64"
	"fmt"
	"openapi-sdk-go/utils"
)

func GenerateSignature(secret, method, uri, clientID, timestamp, body string, queryParams map[string][]string, token string) (string, error) {
	bodyMD5 := utils.MD5Hash(body)
	canonicalHeaders := utils.BuildCanonicalHeaders(map[string]string{
		"X-DRC-Client-ID": clientID,
		"X-DRC-Timestamp": timestamp,
	}, token)
	canonicalQueryString := utils.BuildCanonicalQueryString(queryParams)
	canonicalString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, uri, canonicalQueryString, canonicalHeaders, bodyMD5)
	signature := utils.GenerateHMACSHA1(secret, canonicalString)
	return base64.StdEncoding.EncodeToString(signature), nil
}
