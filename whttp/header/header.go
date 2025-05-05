package header

import "net/http"

type (
	ContentTypeHeader string
)

const (
	ContentType                                 = "Content-Type"
	ApplicationFormUrlEncoded ContentTypeHeader = "application/x-www-form-urlencoded"
	ApplicationJSON           ContentTypeHeader = "application/json"
	Css                       ContentTypeHeader = "text/css"
	Html                      ContentTypeHeader = "text/html"
	Javascript                ContentTypeHeader = "text/javascript"
	MultiPartForm             ContentTypeHeader = "multipart/form-data"
	UTF8                      ContentTypeHeader = "charset=utf-8"
)

func SetContentType(he http.Header, types ...ContentTypeHeader) {
	for _, h := range types {
		he.Add(ContentType, string(h))
	}

}

const RequestID = "X-Request-ID"

func SetRequestID(he http.Header, id string) {
	// avoid header capitalization
	he[RequestID] = []string{id}
}
