package log

import "context"

type LoggerContextKeyType string

const LoggerContextKey LoggerContextKeyType = "__LOGGER_CTX_KEY__"

type LoggerContext struct {
	RequestID string      `json:"request_id"`
	Payload   interface{} `json:"payload"`
}

type LoggerRepository interface {
	Error(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
}
