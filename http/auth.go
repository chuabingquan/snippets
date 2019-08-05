package http

import (
	"encoding/json"
	"net/http"

	"github.com/chuabingquan/snippets"
	"github.com/gorilla/mux"
)

// AuthHandler is a sub-router that handles requests related to resource access-control
type AuthHandler struct {
	*mux.Router
	AuthService   snippets.AuthenticationService
	Authenticator Authenticator
	UserService   snippets.UserService
}

// NewAuthHandler serves as a constructor for an AuthHandler
func NewAuthHandler(as snippets.AuthenticationService, us snippets.UserService, auth Authenticator) *AuthHandler {
	h := &AuthHandler{
		Router:        mux.NewRouter(),
		AuthService:   as,
		Authenticator: auth,
		UserService:   us,
	}

	h.Handle("/api/v0/auth/login", Adapt(http.HandlerFunc(h.handleLogin))).Methods("POST")

	return h
}

// handleLogin
func (ah AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials map[string]string

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"Invalid request body"})
		return
	}

	isAuthenticated, err := ah.AuthService.Authenticate(credentials["username"], credentials["password"])
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when authenticating login"})
		return
	}
	if !isAuthenticated {
		createResponse(w, http.StatusUnauthorized, defaultResponse{"Invalid credentials supplied"})
		return
	}

	user, err := ah.UserService.UserByUsername(credentials["username"])
	if err != nil || user == (snippets.User{}) {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when authenticating login"})
		return
	}

	tokenString, err := ah.Authenticator.GenerateToken(snippets.AuthorizationInfo{
		UserID: user.ID,
	})
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when authenticating login"})
		return
	}

	createResponse(w, http.StatusOK, struct {
		Token string `json:"accessToken"`
	}{tokenString})
}
