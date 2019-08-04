package http

import (
	"net/http"

	"github.com/chuabingquan/snippets"
)

// Authenticator defines a set of operations that deals with the generation,
// validation, and extraction of information from authentication tokens
type Authenticator interface {
	GetAuthorizationInfo(r *http.Request) (snippets.AuthorizationInfo, error)
	GenerateToken(info snippets.AuthorizationInfo) (string, error)
	Authenticate(r *http.Request) (bool, error)
}

// verifyRoute is a middleware that permits/reject entry to a route depending
// on whether a user is successfully authenticated or not
func verifyRoute(a Authenticator) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok, err := a.Authenticate(r)
			if err != nil {
				createResponse(w, http.StatusBadRequest, defaultResponse{
					"Invalid token format supplied"})
				return
			}
			if !ok {
				createResponse(w, http.StatusUnauthorized, defaultResponse{
					"Invalid token supplied"})
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
