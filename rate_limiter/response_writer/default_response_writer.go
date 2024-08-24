package response_writer

import "net/http"

type RateLimiterDefaultResponseWriter struct {
	message    string
	statusCode int
}

func NewRateLimiterDefaultResponseWriter() *RateLimiterDefaultResponseWriter {
	responseWriter := &RateLimiterDefaultResponseWriter{}
	responseWriter.statusCode = 429
	responseWriter.message = "you have reached the maximum number of requests or actions allowed within a certain time frame"
	return responseWriter
}

func (rw *RateLimiterDefaultResponseWriter) WriteResponse(w *http.ResponseWriter) error {
	(*w).WriteHeader(rw.statusCode)
	(*w).Write([]byte(rw.message))
	return nil
}

func (rw *RateLimiterDefaultResponseWriter) WriteError(w *http.ResponseWriter, err error) error {
	(*w).WriteHeader(500)
	(*w).Write([]byte("internal server error"))
	return nil
}
