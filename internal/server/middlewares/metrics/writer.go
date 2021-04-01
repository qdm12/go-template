package metrics

import "net/http"

// statefulWriter wraps the HTTP writer in order to report
// the HTTP status code and the number of bytes written.
type statefulWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statefulWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statefulWriter) Write(b []byte) (n int, err error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err = w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}
