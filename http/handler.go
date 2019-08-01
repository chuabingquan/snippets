package http

import (
	"net/http"
	"strings"
)

// Handler implements the http.Handler interface and acts as the main handler for the server,
// redirecting requests to sub-handlers
type Handler struct {
	UserHandler    *UserHandler
	SnippetHandler *SnippetHandler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "/")
	if len(urlSegments) < 4 {
		http.NotFound(w, r)
		return
	}

	resourceName := urlSegments[3]

	switch resourceName {
	case "users":
		h.UserHandler.ServeHTTP(w, r)
		break
	case "snippets":
		h.SnippetHandler.ServeHTTP(w, r)
		break
	default:
		http.NotFound(w, r)
	}
}
