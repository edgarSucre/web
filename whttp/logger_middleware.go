package whttp

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/edgarsucre/web"
)

type LoggerOpt func(logger *slog.Logger, r *http.Request) *slog.Logger

func WithHeaders(headers []string) LoggerOpt {
	return func(logger *slog.Logger, r *http.Request) *slog.Logger {
		attrs := []any{}
		for _, headerKey := range headers {
			headerContent := strings.Join(r.Header[headerKey], "; ")
			attrs = append(attrs, slog.String(headerKey, headerContent))
		}

		g := slog.Group("headers", attrs...)
		logger = logger.With(g)

		return logger
	}
}

func LoggerMiddleware(
	logger *slog.Logger,
	next http.Handler,
	skipper Skipper,
	opts ...LoggerOpt,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipper(r) {
			next.ServeHTTP(w, r)
			return
		}

		logger = logger.With(
			"method", r.Method,
			"path", r.URL.Path,
		)

		for _, opt := range opts {
			logger = opt(logger, r)
		}

		lw := &LoggerWriter{ResponseWriter: w}

		ctx := r.Context()
		ctx = context.WithValue(ctx, web.LoggerKey, logger)

		next.ServeHTTP(lw, r.WithContext(ctx))

		logger.Info("request", "status", lw.status)
	})
}

type LoggerWriter struct {
	http.ResponseWriter
	buf    bytes.Buffer
	status int
}

func (lw *LoggerWriter) Write(data []byte) (int, error) {
	writer := io.MultiWriter(lw.ResponseWriter, &lw.buf)

	return writer.Write(data)
}

func (lw *LoggerWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}
