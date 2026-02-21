package middlewares

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		log.Printf("[REQUEST] method=%s path=%s remote=%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		next(recorder, r)

		duration := time.Since(start)

		log.Printf("[RESPONSE] method=%s path=%s status=%d duration=%s",
			r.Method,
			r.RequestURI,
			recorder.statusCode,
			duration,
		)
	}
}
