package whttp

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/edgarsucre/web"
	"github.com/edgarsucre/web/whttp/header"
)

type LoggerOpt func(r *http.Request) slog.Attr

func WithHeaders(headers []string) LoggerOpt {
	return func(r *http.Request) slog.Attr {
		attrs := []any{}
		for _, headerKey := range headers {
			headerContent := strings.Join(r.Header[headerKey], "; ")
			attrs = append(attrs, slog.String(headerKey, headerContent))
		}

		return slog.Group("headers", attrs...)
	}
}

func LoggerMiddleware(
	logger *slog.Logger,
	next http.Handler,
	skipper Skipper,
	opts ...LoggerOpt,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipper != nil && skipper(r) {
			next.ServeHTTP(w, r)
			return
		}

		attrs := make([]any, len(opts))
		for i, fn := range opts {
			attrs[i] = fn(r)
		}

		request := slog.Group("request",
			"method", r.Method,
			"path", r.URL.Path,
			"requestID", r.Header.Get(header.RequestID),
		)

		attrs = append(attrs, request)

		ctx := r.Context()
		ctx = context.WithValue(ctx, web.LoggerKey, logger.With(request))

		lw := &LoggerWriter{ResponseWriter: w}
		next.ServeHTTP(lw, r.WithContext(ctx))

		attrs = append(attrs, slog.Int("status", lw.status))

		logger.Info("http request", attrs...)
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
