package custommiddleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Length     int64
	Err        error
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.Length += int64(n)
	return n, err
}

func (rw *ResponseWriter) WriteError(err error) {
	rw.Err = err
}
