package response_writer

import "net/http"

type rateLimiterResponseWriter interface {
	WriteResponse(w *http.ResponseWriter) error
	WriteError(w *http.ResponseWriter, err error) error
}
