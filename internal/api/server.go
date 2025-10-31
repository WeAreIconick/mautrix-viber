// Package api server provides REST API server for external bridge management.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/mautrix-viber/internal/database"
)

// Server provides REST API endpoints for bridge management.
// Renamed from APIServer to avoid stuttering (api.APIServer).
type Server struct {
	db *database.DB
}

// NewAPIServer creates a new API server.
// Deprecated: use NewServer instead to avoid stuttering (api.APIServer -> api.Server).
func NewAPIServer(db *database.DB) *Server {
	return &Server{db: db}
}

// NewServer creates a new API server.
func NewServer(db *database.DB) *Server {
	return &Server{db: db}
}

// RegisterRoutes registers API routes.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/users", s.handleUsers)
	mux.HandleFunc("/api/v1/rooms", s.handleRooms)
	mux.HandleFunc("/api/v1/link", s.handleLink)
	mux.HandleFunc("/api/v1/unlink", s.handleUnlink)
	mux.HandleFunc("/api/v1/status", s.handleStatus)
}

// handleUsers handles user management API.
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if s.db == nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"users": []interface{}{},
			"error": "database not configured",
		})
		return
	}

	// Query database for linked users
	users, err := s.db.ListLinkedUsers(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list users: %v", err), http.StatusInternalServerError)
		return
	}

	userList := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		userList = append(userList, map[string]interface{}{
			"viber_id":       u.ViberID,
			"viber_name":     u.ViberName,
			"matrix_user_id": u.MatrixUserID,
			"linked_at":      u.UpdatedAt,
		})
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"users": userList,
	})
}

// handleRooms handles room management API.
func (s *Server) handleRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if s.db == nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"rooms": []interface{}{},
			"error": "database not configured",
		})
		return
	}

	// Query database for room mappings
	rooms, err := s.db.ListRoomMappings(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list rooms: %v", err), http.StatusInternalServerError)
		return
	}

	roomList := make([]map[string]interface{}, 0, len(rooms))
	for _, r := range rooms {
		roomList = append(roomList, map[string]interface{}{
			"matrix_room_id": r.MatrixRoomID,
			"viber_chat_id":  r.ViberChatID,
			"created_at":     r.CreatedAt,
		})
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"rooms": roomList,
	})
}

// handleLink handles linking users via API.
func (s *Server) handleLink(w http.ResponseWriter, r *http.Request) {
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

	if s.db == nil {
		http.Error(w, "database not configured", http.StatusInternalServerError)
		return
	}

	// Validate input
	if req.MatrixUserID == "" || req.ViberUserID == "" {
		http.Error(w, "matrix_user_id and viber_user_id are required", http.StatusBadRequest)
		return
	}

	// Link the user
	if err := s.db.LinkViberUser(r.Context(), req.ViberUserID, req.MatrixUserID); err != nil {
		http.Error(w, fmt.Sprintf("failed to link user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "linked",
		"matrix_user_id": req.MatrixUserID,
		"viber_user_id":  req.ViberUserID,
	})
}

// handleUnlink handles unlinking users via API.
func (s *Server) handleUnlink(w http.ResponseWriter, r *http.Request) {
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

	if s.db == nil {
		http.Error(w, "database not configured", http.StatusInternalServerError)
		return
	}

	if req.MatrixUserID == "" {
		http.Error(w, "matrix_user_id is required", http.StatusBadRequest)
		return
	}

	// Unlink by setting matrix_user_id to NULL
	if err := s.db.UnlinkMatrixUser(r.Context(), req.MatrixUserID); err != nil {
		http.Error(w, fmt.Sprintf("failed to unlink user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "unlinked",
		"matrix_user_id": req.MatrixUserID,
	})
}

// handleStatus handles status API.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	InfoHandler(w, r)
}
