package http

import "net/http"

// Adapter adds functionality before/after the http.Handler it takes in by wrapping it in another http.Handler that will be returned
type Adapter func(http.Handler) http.Handler

// Adapt returns a http.Handler that when invoked, executes the output of the adapters in the sequence it was given before running the main http.Handler
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}

	return h
}
