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
}

// NewSnippetHandler constructs a new SnippetHandler given a SnippetService implementation
func NewSnippetHandler(ss snippets.SnippetService) *SnippetHandler {
	h := &SnippetHandler{
		Router:         mux.NewRouter(),
		SnippetService: ss,
	}

	h.Handle("/api/v0/snippets", Adapt(http.HandlerFunc(h.handleGetSnippets))).Methods("GET")
	h.Handle("/api/v0/snippets/{snippetID}", Adapt(http.HandlerFunc(h.handleGetSnippetByID))).Methods("GET")
	h.Handle("/api/v0/snippets", Adapt(http.HandlerFunc(h.handleCreateSnippet))).Methods("POST")
	h.Handle("/api/v0/snippets/{snippetID}", Adapt(http.HandlerFunc(h.handlePatchSnippet))).Methods("PATCH")
	h.Handle("/api/v0/snippets/{snippetID}", Adapt(http.HandlerFunc(h.handleDeleteSnippet))).Methods("DELETE")

	return h
}

// handleGetSnippets
func (sh SnippetHandler) handleGetSnippets(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	snippets, err := sh.SnippetService.Snippets(userID)
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
	userID := r.Header.Get("Authorization")

	snippet, err := sh.SnippetService.Snippet(userID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested snippet"})
		return
	}
	createResponse(w, http.StatusOK, snippet)
}

// handleCreateSnippet
func (sh SnippetHandler) handleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")

	var newSnippet snippets.Snippet
	err := json.NewDecoder(r.Body).Decode(&newSnippet)
	if err != nil {
		createResponse(w, http.StatusBadRequest, defaultResponse{
			"Invalid request body"})
		return
	}

	// Owner should ideally be assigned to userId from access token
	newSnippet.Owner = userID

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
	userID := r.Header.Get("Authorization")
	snippetID := mux.Vars(r)["snippetID"]

	snippetToUpdate, err := sh.SnippetService.Snippet(userID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when getting requested snippet"})
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
	userID := r.Header.Get("Authorization")
	snippetID := mux.Vars(r)["snippetID"]

	err := sh.SnippetService.DeleteSnippet(userID, snippetID)
	if err != nil {
		createResponse(w, http.StatusInternalServerError, defaultResponse{
			"An unexpected error occurred when deleting snippet"})
		return
	}

	createResponse(w, http.StatusOK, defaultResponse{"Snippet is successfully deleted"})
}
