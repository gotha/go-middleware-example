package main

import (
	"context"
	"fmt"
	"net/http"
)

type key int

const (
	loggingKey   key = iota
	someOtherKey key = iota
)

func firstMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), someOtherKey, "myValue")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func secondMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appendToContext(r.Context(), loggingKey, "logKey1", "logValue1")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func thirdMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := appendToContext(r.Context(), loggingKey, "logKey2", "logValue2")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggingMiddleware() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		fmt.Printf("%+v\n", "-----------")
		fmt.Printf("%+v\n", ctx.Value(loggingKey))
		fmt.Printf("%+v\n", "-----------")
		fmt.Printf("%+v\n", ctx.Value(someOtherKey))
		fmt.Printf("%+v\n", "-----------")

	})
}

func appendToContext(ctx context.Context, contextKey key, key string, value interface{}) context.Context {
	contextValue := ctx.Value(contextKey)
	data := make(map[string]interface{}, 0)
	if contextValue != nil {
		data = contextValue.(map[string]interface{})
	}
	data[key] = value

	newContext := context.WithValue(ctx, contextKey, data)

	return newContext

}
