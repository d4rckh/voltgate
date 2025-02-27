package proxy

import "net/http"

type ResponseWriterTrap struct {
	http.ResponseWriter
	StatusCode  int
	ContentSize int
	Body        []byte
}

func (w *ResponseWriterTrap) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriterTrap) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.ContentSize += size
	w.Body = append(w.Body, b...)
	return size, err
}
