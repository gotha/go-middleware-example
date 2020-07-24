package main

import (
	"net/http"
)

func newMainHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello handler"))

		next.ServeHTTP(w, r)
	})
}
