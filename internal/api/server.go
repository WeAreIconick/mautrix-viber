// Package api server provides REST API server for external bridge management.
package api

import (
	"encoding/json"
	"net/http"
)

// APIServer provides REST API endpoints for bridge management.
type APIServer struct {
	// Bridge management methods would be injected here
}

// NewAPIServer creates a new API server.
func NewAPIServer() *APIServer {
	return &APIServer{}
}

// RegisterRoutes registers API routes.
func (s *APIServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/users", s.handleUsers)
	mux.HandleFunc("/api/v1/rooms", s.handleRooms)
	mux.HandleFunc("/api/v1/link", s.handleLink)
	mux.HandleFunc("/api/v1/unlink", s.handleUnlink)
	mux.HandleFunc("/api/v1/status", s.handleStatus)
}

// handleUsers handles user management API.
func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// TODO: Return list of linked users
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": []interface{}{},
	})
}

// handleRooms handles room management API.
func (s *APIServer) handleRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// TODO: Return list of mapped rooms
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rooms": []interface{}{},
	})
}

// handleLink handles linking users via API.
func (s *APIServer) handleLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		MatrixUserID string `json:"matrix_user_id"`
		ViberUserID  string `json:"viber_user_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// TODO: Implement linking logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "linked",
		"matrix_user_id": req.MatrixUserID,
		"viber_user_id": req.ViberUserID,
	})
}

// handleUnlink handles unlinking users via API.
func (s *APIServer) handleUnlink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		MatrixUserID string `json:"matrix_user_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// TODO: Implement unlinking logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "unlinked",
	})
}

// handleStatus handles status API.
func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	InfoHandler(w, r)
}

