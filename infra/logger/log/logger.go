package log

import (
	"context"
	"github/guiferpa/bank/domain/log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func (l *Logger) Error(ctx context.Context, msg string) {
	cctx, _ := ctx.Value(log.LoggerContextKey).(*log.LoggerContext)

	l.logger.Error(msg,
		zap.String("request_id", cctx.RequestID),
	)
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	cctx, _ := ctx.Value(log.LoggerContextKey).(*log.LoggerContext)

	l.logger.Warn(msg,
		zap.String("request_id", cctx.RequestID),
	)
}

func (l *Logger) Info(ctx context.Context, msg string) {
	cctx, _ := ctx.Value(log.LoggerContextKey).(*log.LoggerContext)

	l.logger.Info(msg,
		zap.String("request_id", cctx.RequestID),
	)
}

func NewLogger() *Logger {
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				MessageKey: "message",

				LevelKey:    "level",
				EncodeLevel: zapcore.LowercaseLevelEncoder,

				TimeKey:    "time",
				EncodeTime: zapcore.ISO8601TimeEncoder,

				CallerKey:    "caller",
				EncodeCaller: zapcore.FullCallerEncoder,
			}),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

	return &Logger{logger}
}
