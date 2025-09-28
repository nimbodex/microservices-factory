package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
)

type logger struct {
	zapLogger *zap.Logger
}

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	SetLevel(level string) error
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{
		zapLogger: l.zapLogger.With(fields...),
	}
}

func (l *logger) SetLevel(level string) error {
	parsedLevel := parseLevel(level)
	dynamicLevel.SetLevel(parsedLevel)
	return nil
}

func GetLogger() Logger {
	if globalLogger == nil {
		panic("logger not initialized")
	}
	return globalLogger
}

func Init(levelStr string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))
		encoderCfg := buildProductionEncoderConfig()
		var encoder zapcore.Encoder
		if asJSON {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}
		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			dynamicLevel,
		)
		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})
	return nil
}

func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return config
}
