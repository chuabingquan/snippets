package http

import (
	"encoding/json"
	"net/http"

	"github.com/chuabingquan/snippets"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UserHandler is a sub-router that handles requests related to operations on Users
type UserHandler struct {
	*mux.Router
	UserService   snippets.UserService
	HashUtilities snippets.HashUtilities
}

// NewUserHandler constructs a new UserHandler given a UserService implementation
func NewUserHandler(us snippets.UserService, hu snippets.HashUtilities) *UserHandler {
	h := &UserHandler{
		Router:        mux.NewRouter(),
		UserService:   us,
		HashUtilities: hu,
	}

	h.Handle("/api/v0/users", Adapt(http.HandlerFunc(h.handleGetUsers))).Methods("GET")
	h.Handle("/api/v0/users/{userID}", Adapt(http.HandlerFunc(h.handleGetUserByID))).Methods("GET")
	h.Handle("/api/v0/users", Adapt(http.HandlerFunc(h.handleCreateUser))).Methods("POST")
	h.Handle("/api/v0/users/{userID}", Adapt(http.HandlerFunc(h.handlePatchUser))).Methods("PATCH")
	h.Handle("/api/v0/users/{userID}", Adapt(http.HandlerFunc(h.handleDeleteUser))).Methods("DELETE")

	return h
}

// handleGetUsers
func (uh UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.UserService.Users()
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when retrieving users"})
		return
	}
	createResponse(w, http.StatusOK, users)
}

// handleGetUserByID
func (uh UserHandler) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	user, err := uh.UserService.User(userID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested user"})
		return
	}
	createResponse(w, http.StatusOK, user)
}

// handleCreateUser
func (uh UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser snippets.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"Invalid request body"})
		return
	}

	newUser.ID = uuid.New().String()
	hash, err := uh.HashUtilities.HashAndSalt(newUser.Password)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when creating user"})
		return
	}
	newUser.PasswordHash = hash

	err = uh.UserService.CreateUser(newUser)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when creating user"})
		return
	}

	createResponse(w, http.StatusCreated, defaultResponse{"User is successfully created"})
}

// handlePatchUser
func (uh UserHandler) handlePatchUser(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT HANDLEPATCHUSER
	w.WriteHeader(http.StatusOK)
}

// handleDeleteUser
func (uh UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	err := uh.UserService.DeleteUser(userID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting user"})
		return
	}

	createResponse(w, http.StatusOK, defaultResponse{"User is successfully deleted"})
}
