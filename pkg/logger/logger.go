package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RequestID   = "requestID"
	ServiceName = "service"
	Mode        = "mode"
)

type (
	LoggerKey    struct{}
	RequestIDKey struct{}
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Stop() error
}

type logger struct {
	serviceName string
	logger      *zap.Logger
}

func (l logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName))

	if ctx.Value(RequestIDKey{}) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestIDKey{}).(string)))
	}
	l.logger.Debug(msg, fields...)
}

func (l logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName))

	if ctx.Value(RequestIDKey{}) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestIDKey{}).(string)))
	}
	l.logger.Info(msg, fields...)
}

func (l logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName))

	if ctx.Value(RequestIDKey{}) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestIDKey{}).(string)))
	}
	l.logger.Warn(msg, fields...)
}

func (l logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName))

	if ctx.Value(RequestIDKey{}) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestIDKey{}).(string)))
	}
	l.logger.Error(msg, fields...)
}

func (l logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName))

	if ctx.Value(RequestIDKey{}) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestIDKey{}).(string)))
	}
	l.logger.Fatal(msg, fields...)
}

func New(serviceName string) (Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger configuration: %w", err)
	}

	return &logger{
		serviceName: serviceName,
		logger:      zapLogger,
	}, nil
}

func (l logger) Stop() error {
	_ = l.logger.Sync()

	return nil
}

func GetLoggerFromCtx(ctx context.Context) Logger {
	return ctx.Value(LoggerKey{}).(Logger)
}

func SetToCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, LoggerKey{}, l)
}
