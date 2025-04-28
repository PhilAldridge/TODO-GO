package logging

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	traceIDKey contextKey = "traceID"
	loggerKey  contextKey = "logger"
)

func generateTraceID() string {
	return string(uuid.NewString())
}

func WithTraceIDAndLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        traceID := generateTraceID()

        logger := slog.With("traceID", traceID)
        ctx := context.WithValue(r.Context(), traceIDKey, traceID)
        ctx = context.WithValue(ctx, loggerKey, logger)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger := LoggerFromContext(r.Context())
        traceID := TraceIDFromContext(r.Context())

        start := time.Now()
        logger.Info("request started", "method", r.Method, "path", r.URL.Path)

        // Call the next handler
        next.ServeHTTP(w, r)

        duration := time.Since(start)
        logger.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", duration, "traceID", traceID)
    })
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}

func TraceIDFromContext(ctx context.Context) string {
    if traceID, ok := ctx.Value(traceIDKey).(string); ok {
        return traceID
    }
    return "unknown"
}