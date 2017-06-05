package lib

import (
	"log"
	"net/http"
	"time"
)

// ----- LOGGER

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// WrapHandlerWithLogging for loggin purpose
func WrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// log.Printf("--> %s %s", r.Method, r.URL.Path)
		start := time.Now()

		lrw := newLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode
		// log.Printf("<-- %d %s [ %s ]", statusCode, http.StatusText(statusCode), time.Since(start))
		// or
		if r.URL.Path != "/favicon.ico" {
			log.Printf("%d %s %s [ %s ]", statusCode, r.Method, r.URL.Path, time.Since(start))
		}
	})
}
