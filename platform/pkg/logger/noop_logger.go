package logger

import (
	"context"

	"go.uber.org/zap"
)

type NoopLogger struct{}

func (n *NoopLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {}

func (n *NoopLogger) With(fields ...zap.Field) Logger {
	return n
}

func (n *NoopLogger) SetLevel(level string) error {
	return nil
}
