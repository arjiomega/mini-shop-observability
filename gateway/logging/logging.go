package logging

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	logger *slog.Logger
}

func New(serviceName, instanceID string) *Logger {
	return &Logger{
		logger: slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		).With(
			slog.String("service", serviceName),
			slog.String("service_instance", instanceID),
		),
	}
}

func (l *Logger) Info(ctx context.Context, message string, args ...any) {
	l.logger.InfoContext(ctx, message, withTrace(ctx, args...)...)
}

func (l *Logger) Warn(ctx context.Context, message string, args ...any) {
	l.logger.WarnContext(ctx, message, withTrace(ctx, args...)...)
}

func (l *Logger) Error(ctx context.Context, message string, args ...any) {
	l.logger.ErrorContext(ctx, message, withTrace(ctx, args...)...)
}

func (l *Logger) Fatal(ctx context.Context, message string, args ...any) {
	l.logger.ErrorContext(
		ctx,
		message,
		withTrace(ctx, args...)...,
	)

	os.Exit(1)
}

func withTrace(ctx context.Context, args ...any) []any {
	sc := trace.SpanContextFromContext(ctx)

	if sc.IsValid() {
		args = append(args,
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}

	return args
}
