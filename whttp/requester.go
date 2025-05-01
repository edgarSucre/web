package whttp

import (
	"context"
	"io"
	"net/http"
)

type Requester interface {
	MakeRequest(
		ctx context.Context,
		method string,
		url string,
		body io.Reader,
		header http.Header,
	) (io.Reader, int, error)
}
