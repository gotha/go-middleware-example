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

// Middleware - define what middleware is
type Middleware func(http.Handler) http.Handler

// Chain - function for chaining middlewares in correct order
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fmt.Printf("%+v\n", "-----------")
		fmt.Printf("%+v\n", ctx.Value(loggingKey))
		fmt.Printf("%+v\n", "-----------")
		fmt.Printf("%+v\n", ctx.Value(someOtherKey))
		fmt.Printf("%+v\n", "-----------")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func nopMiddleware() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func appendToContext(ctx context.Context, contextKey key, key string, value interface{}) context.Context {
	contextValue := ctx.Value(contextKey)
	data := make(map[string]interface{})
	if contextValue != nil {
		data = contextValue.(map[string]interface{})
	}
	data[key] = value

	newContext := context.WithValue(ctx, contextKey, data)

	return newContext
}
