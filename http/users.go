package http

import (
	"encoding/json"
	"net/http"

	"github.com/chuabingquan/snippets"
	"github.com/gorilla/mux"
)

// UserHandler is a sub-router that handles requests related to operations on Users
type UserHandler struct {
	*mux.Router
	UserService   snippets.UserService
	Authenticator Authenticator
}

// NewUserHandler constructs a new UserHandler given a UserService implementation
func NewUserHandler(us snippets.UserService, auth Authenticator) *UserHandler {
	h := &UserHandler{
		Router:        mux.NewRouter(),
		UserService:   us,
		Authenticator: auth,
	}

	verifyUser := verifyRoute(auth)

	h.Handle("/api/v0/users", Adapt(http.HandlerFunc(h.handleGetUsers), verifyUser)).Methods("GET")
	h.Handle("/api/v0/users/{userID}", Adapt(http.HandlerFunc(h.handleGetUserByID), verifyUser)).Methods("GET")
	h.Handle("/api/v0/users", Adapt(http.HandlerFunc(h.handleCreateUser))).Methods("POST")
	h.Handle("/api/v0/users/{userID}", Adapt(http.HandlerFunc(h.handlePatchUser), verifyUser)).Methods("PATCH")
	h.Handle("/api/v0/users/{userID}", Adapt(http.HandlerFunc(h.handleDeleteUser), verifyUser)).Methods("DELETE")

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
	userID := mux.Vars(r)["userID"]
	userInfo, err := uh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested user"})
		return
	}
	if userInfo.UserID != userID {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, requested user is not found"})
		return
	}

	user, err := uh.UserService.User(userID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested user"})
		return
	}
	if user == (snippets.User{}) {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, requested user is not found"})
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
	userID := mux.Vars(r)["userID"]
	userInfo, err := uh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when updating user"})
		return
	}
	if userInfo.UserID != userID {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, user to update is not found"})
		return
	}

	userToUpdate, err := uh.UserService.User(userID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when updating user"})
		return
	}
	if userToUpdate == (snippets.User{}) {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, user to update is not found"})
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&userToUpdate)
	if err != nil {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"JSON could not be decoded, invalid request format supplied"})
		return
	}

	// validation
	if userToUpdate.ID != userID {
		// Prevent attackers from updating another user's information by using
		// a different userID supplied in JSON from the one specified in the url params
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"JSON could not be decoded, invalid request format supplied"})
		return
	}

	err = uh.UserService.UpdateUser(userToUpdate)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when updating user"})
		return
	}

	createResponse(w, http.StatusOK, defaultResponse{"User is successfully updated"})
}

// handleDeleteUser
func (uh UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]
	userInfo, err := uh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting user"})
		return
	}
	if userInfo.UserID != userID {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, user to delete is not found"})
		return
	}

	userToDelete, err := uh.UserService.User(userID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting user"})
		return
	}
	if userToDelete == (snippets.User{}) {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, user to delete is not found"})
		return
	}

	err = uh.UserService.DeleteUser(userID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting user"})
		return
	}

	createResponse(w, http.StatusOK, defaultResponse{"User is successfully deleted"})
}
