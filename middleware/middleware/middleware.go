package middleware

import (
	"net/http"
)

// HandleFunc is a function that handles HTTP requests.
// This is a simple shorthand to define easier to read functions.
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// Middleware is a special type that handles HandleFuncs.
type Middleware func(HandleFunc) HandleFunc

// Handle handles the middlewares.
// It executes the middlewares in the order presented and finishes by calling the final handler.
func Handle(final HandleFunc, middlewares ...Middleware) HandleFunc {
	if final == nil {
		panic("no final handler")
		// Or return a default handler.
	}
	// Execute the middleware in the same order and return the final func.
	// This is a confusing and tricky construct :)
	// We need to use the reverse order since we are chaining inwards.
	for i := len(middlewares) - 1; i >= 0; i-- {
		final = middlewares[i](final) // mw1(mw2(mw3(final)))
	}
	return final
}
