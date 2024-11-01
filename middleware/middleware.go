package middleware

import (
	"log"
	"net/http"
	"time"
)

// Get statusCode
type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWritter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		// Aqui esta la magia, hemos sobreescrito WriteHeader, lo que hemos enviado es algo
		// que se parece a un writter, pero es mas. En cierto momento se usa WriteHeader
		// Y eso hace que status Code se rellene.
		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
