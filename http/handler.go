package http

import (
	"net/http"
	"strings"
)

// Handler implements the http.Handler interface and acts as the main handler for the server,
// redirecting requests to sub-handlers
type Handler struct {
	UserHandler *UserHandler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "/")
	if len(urlSegments) < 4 {
		http.NotFound(w, r)
		return
	}

	mainResourceName := urlSegments[3]

	if mainResourceName == "users" {
		h.UserHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
