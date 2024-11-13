package presenter

import "net/http"

type WrapResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Length     int64
	Err        error
}

func NewWrapResponseWriter(w http.ResponseWriter) *WrapResponseWriter {
	return &WrapResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
}

func GetWrapResponseWriter(w http.ResponseWriter) *WrapResponseWriter {
	if rw, ok := w.(*WrapResponseWriter); ok {
		return rw
	}
	return NewWrapResponseWriter(w)
}

func (rw *WrapResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *WrapResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.Length += int64(n)
	return n, err
}

func (rw *WrapResponseWriter) WriteError(err error) {
	rw.Err = err
}
