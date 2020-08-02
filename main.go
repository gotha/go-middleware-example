package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("%+v\n", "hi there")

	r := mux.NewRouter()

	mainHandler := Chain(
		firstMiddleware,
		secondMiddleware,
		thirdMiddleware,
		newMainHandler,
		loggingMiddleware,
	)(nopMiddleware())

	// chain middlewarez the native way
	//var mainHandler http.Handler
	//mainHandler = loggingMiddleware(dummyHandler())
	//mainHandler = newMainHandler(mainHandler)
	//mainHandler = thirdMiddleware(mainHandler)
	//mainHandler = secondMiddleware(mainHandler)
	//mainHandler = firstMiddleware(mainHandler)
	r.Handle("/", mainHandler)

	http.Handle("/", r)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 9999),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	srv.SetKeepAlivesEnabled(false)

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
