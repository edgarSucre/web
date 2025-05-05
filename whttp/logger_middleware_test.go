package whttp_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/edgarsucre/web/whttp"
	"github.com/edgarsucre/web/whttp/header"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMIddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		ct := r.Header.Get(header.ContentType)

		_ = ct

		status := r.URL.Query().Get("status")
		intStatus, _ := strconv.Atoi(status)

		w.WriteHeader(intStatus)
		w.Write([]byte(msg))
	})

	var writer bytes.Buffer

	baseLogger := slog.New(slog.NewTextHandler(&writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	reqWithHeaders := httptest.NewRequest(http.MethodGet, "/test?status=205&msg=withHeaders", nil)
	header.SetContentType(reqWithHeaders.Header, header.ApplicationJSON, header.UTF8)
	header.SetRequestID(reqWithHeaders.Header, "request_id")

	type opts struct {
		loggerOpts []whttp.LoggerOpt
		req        *http.Request
		skipper    whttp.Skipper
	}

	tests := []struct {
		name  string
		opts  opts
		check func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			"skipped",
			opts{
				nil,
				httptest.NewRequest(http.MethodGet, "/test?status=200", nil),
				func(r *http.Request) bool { return true },
			},
			func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Body.String(), writer.String())

				assert.Equal(t, 200, rec.Code)
			},
		},
		{
			"simpleInfo",
			opts{
				nil,
				httptest.NewRequest(http.MethodGet, "/test?status=201&msg=testing", nil),
				func(r *http.Request) bool { return false },
			},
			func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, 201)
				assert.Equal(t, rec.Body.String(), "testing")

				content := writer.String()
				assert.Contains(t, content, "msg=request")
				assert.Contains(t, content, "status=201")
			},
		},
		{
			"withHeaders",
			opts{
				[]whttp.LoggerOpt{whttp.WithHeaders([]string{header.ContentType, header.RequestID})},
				reqWithHeaders,
				func(r *http.Request) bool { return false },
			},
			func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, 205)
				assert.Equal(t, rec.Body.String(), "withHeaders")

				content := writer.String()
				assert.Contains(t, content, "msg=request")
				assert.Contains(t, content, "status=205")
				assert.Contains(t, content, "headers.X-Request-ID=request_id")
				assert.Contains(t, content, "headers.Content-Type=\"application/json; charset=utf-8\"")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer.Reset()

			middleware := whttp.LoggerMiddleware(
				baseLogger,
				handler,
				tt.opts.skipper,
				tt.opts.loggerOpts...,
			)

			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, tt.opts.req)

			tt.check(t, rec)
		})
	}
}
