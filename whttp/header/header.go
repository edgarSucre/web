package header

import "net/http"

type (
	ContentTypeHeader string
)

const (
	contentType                                 = "Content-Type"
	ApplicationFormUrlEncoded ContentTypeHeader = "application/x-www-form-urlencoded"
	ApplicationJSON           ContentTypeHeader = "application/json"
	Css                       ContentTypeHeader = "text/css"
	Javascript                ContentTypeHeader = "text/javascript"
	MultiPartForm             ContentTypeHeader = "multipart/form-data"
	UTF8                      ContentTypeHeader = "charset=utf-8"
)

func SetContentType(he http.Header, types ...ContentTypeHeader) {
	for _, h := range types {
		he.Add(contentType, string(h))
	}

}
