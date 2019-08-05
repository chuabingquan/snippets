package http

import (
	"encoding/json"
	"net/http"

	"github.com/chuabingquan/snippets"
	"github.com/gorilla/mux"
)

// SnippetHandler is a sub-router that handles requests related to operations on Snippets
type SnippetHandler struct {
	*mux.Router
	SnippetService snippets.SnippetService
	Authenticator  Authenticator
}

// NewSnippetHandler constructs a new SnippetHandler given a SnippetService implementation
func NewSnippetHandler(ss snippets.SnippetService, auth Authenticator) *SnippetHandler {
	h := &SnippetHandler{
		Router:         mux.NewRouter(),
		SnippetService: ss,
		Authenticator:  auth,
	}

	verifyUser := verifyRoute(auth)

	h.Handle("/api/v0/snippets", Adapt(http.HandlerFunc(h.handleGetSnippets), verifyUser)).Methods("GET")
	h.Handle("/api/v0/snippets/{snippetID}", Adapt(http.HandlerFunc(h.handleGetSnippetByID), verifyUser)).Methods("GET")
	h.Handle("/api/v0/snippets", Adapt(http.HandlerFunc(h.handleCreateSnippet), verifyUser)).Methods("POST")
	h.Handle("/api/v0/snippets/{snippetID}", Adapt(http.HandlerFunc(h.handlePatchSnippet), verifyUser)).Methods("PATCH")
	h.Handle("/api/v0/snippets/{snippetID}", Adapt(http.HandlerFunc(h.handleDeleteSnippet), verifyUser)).Methods("DELETE")

	return h
}

// handleGetSnippets
func (sh SnippetHandler) handleGetSnippets(w http.ResponseWriter, r *http.Request) {
	userInfo, err := sh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when retrieving snippets"})
		return
	}

	snippets, err := sh.SnippetService.Snippets(userInfo.UserID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when retrieving snippets"})
		return
	}
	createResponse(w, http.StatusOK, snippets)
}

// handleGetSnippetByID
func (sh SnippetHandler) handleGetSnippetByID(w http.ResponseWriter, r *http.Request) {
	snippetID := mux.Vars(r)["snippetID"]
	userInfo, err := sh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested snippet"})
		return
	}

	snippet, err := sh.SnippetService.Snippet(userInfo.UserID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested snippet"})
		return
	}
	if snippet == (snippets.Snippet{}) {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, requested snippet is not found"})
		return
	}
	createResponse(w, http.StatusOK, snippet)
}

// handleCreateSnippet
func (sh SnippetHandler) handleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	userInfo, err := sh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when creating snippet"})
		return
	}

	var newSnippet snippets.Snippet
	err = json.NewDecoder(r.Body).Decode(&newSnippet)
	if err != nil {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"Invalid request body"})
		return
	}

	newSnippet.Owner = userInfo.UserID

	err = sh.SnippetService.CreateSnippet(newSnippet)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when creating snippet"})
		return
	}
	createResponse(w, http.StatusCreated, defaultResponse{"Snippet is successfully created"})
}

// handlePatchSnippet
func (sh SnippetHandler) handlePatchSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID := mux.Vars(r)["snippetID"]
	userInfo, err := sh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when updating snippet"})
		return
	}

	snippetToUpdate, err := sh.SnippetService.Snippet(userInfo.UserID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested snippet"})
		return
	}
	if snippetToUpdate == (snippets.Snippet{}) {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, snippet to update is not found"})
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&snippetToUpdate)
	if err != nil {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"JSON could not be decoded, invalid request format supplied"})
		return
	}

	if snippetToUpdate.ID != snippetID {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"JSON could not be decoded, invalid request format supplied"})
		return
	}

	err = sh.SnippetService.UpdateSnippet(snippetToUpdate)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when updating snippet"})
		return
	}

	createResponse(w, http.StatusOK, defaultResponse{"Snippet is successfully updated"})
}

// handleDeleteSnippet
func (sh SnippetHandler) handleDeleteSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID := mux.Vars(r)["snippetID"]
	userInfo, err := sh.Authenticator.GetAuthorizationInfo(r)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting snippet"})
		return
	}

	snippetToDelete, err := sh.SnippetService.Snippet(userInfo.UserID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting snippet"})
		return
	}
	if snippetToDelete == (snippets.Snippet{}) {
		createResponse(w, http.StatusNotFound, defaultResponse{
			"Error, snippet to delete is not found"})
		return
	}

	err = sh.SnippetService.DeleteSnippet(userInfo.UserID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting snippet"})
		return
	}

	createResponse(w, http.StatusOK, defaultResponse{"Snippet is successfully deleted"})
}
