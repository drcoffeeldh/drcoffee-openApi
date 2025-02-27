package model

import "time"

const (
	timestampTolerance = 1 * time.Minute // 建议设置为5分钟
	TokenHeader        = "Authorization"
	ClientIdHeader     = "X-DRC-Client-ID"
	SignHeader         = "X-DRC-Sign"
	TimestampHeader    = "X-DRC-Timestamp"
)
