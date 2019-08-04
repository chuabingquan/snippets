package http

import (
	"net/http"
	"strings"

	"github.com/chuabingquan/snippets"
)

// Authenticator defines a set of operations that deals with the generation,
// validation, and extraction of information from authentication tokens
type Authenticator interface {
	GetAuthorizationInfo(tokenString string) (snippets.AuthorizationInfo, error)
	GenerateToken(info snippets.AuthorizationInfo) (string, error)
	Authenticate(tokenString string) (bool, error)
}

// verifyRoute is a middleware that permits/reject entry to a route depending
// on whether a user is successfully authenticated or not
func verifyRoute(a Authenticator) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			tokenParts := strings.Split(tokenString, " ")

			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				createResponse(w, http.StatusBadRequest, defaultResponse{
					"Invalid token format supplied"})
				return
			}

			ok, err := a.Authenticate(tokenParts[1])
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
