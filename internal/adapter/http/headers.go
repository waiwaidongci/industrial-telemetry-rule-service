package httpadapter

import (
	"context"
	"net/http"
	"time"
)

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
func requestTimeout(next http.Handler, d time.Duration) http.Handler {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func maxBody(next http.Handler, size int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, size)
		next.ServeHTTP(w, r)
	})
}
func allowMethods(next http.Handler, methods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("Allow", methods[0])
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
