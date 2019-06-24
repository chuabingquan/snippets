package http

import (
	"encoding/json"
	"net/http"
)

// defaultResponse represents the default structure of a response body for cases where the
// response is only informational
type defaultResponse struct {
	Message string `json:"message"`
}

// createResponse constructs and returns a HTTP response that contains JSON
func createResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
	return
}
