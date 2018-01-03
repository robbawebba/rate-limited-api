package main

import (
	"net/http"

	"github.com/robbawebba/rate-limited-api/rate"
)

var limiter = rate.NewLimiter(2, 5)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
