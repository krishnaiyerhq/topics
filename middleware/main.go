package main

import (
	"krishnaiyerhq/golang/topics/middleware/middleware"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	timeout = 10 * time.Second
)

// requestCount holds the number of requests. This will be reset when the program is closed.
var requestCount atomic.Uint64

// General Middlware construct
// func (HandleFunc) HandleFunc.
// func (second HandleFunc) (first HandleFunc).
// The first HandleFunc is first called, which reads/manipulates the request and then calls the second.

// logMiddleware logs the request.
func logMiddleware(next middleware.HandleFunc) middleware.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First log the request.
		log.Printf("Request received: %s\n", r.RequestURI)
		// Then call the next or second handler.
		next(w, r)
	}
}

// countMiddleware counts the number of requests.
// This function is safe for concurrent use.
func countMiddleware(next middleware.HandleFunc) middleware.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First log the request.
		requestCount.Add(1)
		log.Printf("Request count (incl current): %d\n", requestCount.Load())
		// Then call the next or second handler.
		next(w, r)
	}
}

// traceMiddleware calculates the amount of time a request takes.
func traceMiddleware(next middleware.HandleFunc) middleware.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			// When the function returns from handling the request, print the diff in start and end time.
			log.Printf("Request time: %s, %s\n", r.RequestURI, time.Now().Sub(start))
		}()
		// Then call the next or second handler.
		next(w, r)
	}
}

// handleRoot handles "/".
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Root"))
}

// handleUsers handles "/users" prefix.
func handleUsers(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.WriteHeader(http.StatusAccepted) // HTTP 202
	w.Write([]byte("Hello Users"))
}

func main() {
	// mux is a request multiplexer.
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware.Handle(handleRoot, logMiddleware, countMiddleware, traceMiddleware))
	mux.HandleFunc("/users", middleware.Handle(handleUsers, countMiddleware, logMiddleware, traceMiddleware))

	server := http.Server{
		Handler:      mux,
		Addr:         "localhost:8080",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	log.Printf("Start server: %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
