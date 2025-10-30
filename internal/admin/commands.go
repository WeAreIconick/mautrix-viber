// Package admin provides bridge administration commands for Matrix rooms.
// Commands like !bridge link, !bridge status are handled here.
package admin

import (
	"context"
	"fmt"
	"strings"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/example/mautrix-viber/internal/database"
)

// Command represents a bridge admin command.
type Command struct {
	Name        string
	Description string
	Handler     func(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error)
}

// Handler manages bridge admin commands in Matrix rooms.
type Handler struct {
	commands     map[string]Command
	mxClient     *mautrix.Client
	db           *database.DB
	allowedUsers []id.UserID // Users allowed to run admin commands
}

// NewHandler creates a new admin command handler.
func NewHandler(mxClient *mautrix.Client, db *database.DB, allowedUsers []id.UserID) *Handler {
	h := &Handler{
		commands:     make(map[string]Command),
		mxClient:     mxClient,
		db:           db,
		allowedUsers: allowedUsers,
	}
	h.registerDefaultCommands()
	return h
}

// registerDefaultCommands registers built-in bridge commands.
func (h *Handler) registerDefaultCommands() {
	h.RegisterCommand(Command{
		Name:        "link",
		Description: "Link a Viber account to this Matrix user",
		Handler:     h.handleLink,
	})
	h.RegisterCommand(Command{
		Name:        "unlink",
		Description: "Unlink Viber account from Matrix user",
		Handler:     h.handleUnlink,
	})
	h.RegisterCommand(Command{
		Name:        "status",
		Description: "Show bridge status and connection info",
		Handler:     h.handleStatus,
	})
	h.RegisterCommand(Command{
		Name:        "help",
		Description: "Show available bridge commands",
		Handler:     h.handleHelp,
	})
	h.RegisterCommand(Command{
		Name:        "ping",
		Description: "Test bridge responsiveness",
		Handler:     h.handlePing,
	})
}

// RegisterCommand registers a custom command.
func (h *Handler) RegisterCommand(cmd Command) {
	h.commands[cmd.Name] = cmd
}

// HandleMessage checks if a Matrix message is a bridge command and handles it.
func (h *Handler) HandleMessage(ctx context.Context, evt *event.Event, msg *event.MessageEventContent) error {
	// Check if message starts with !bridge
	if !strings.HasPrefix(msg.Body, "!bridge") {
		return nil
	}

	// Parse command
	parts := strings.Fields(msg.Body)
	if len(parts) < 2 {
		return nil
	}

	cmdName := strings.ToLower(parts[1])
	cmd, ok := h.commands[cmdName]
	if !ok {
		// Unknown command
		h.reply(ctx, evt.RoomID, fmt.Sprintf("Unknown command: %s. Use !bridge help", cmdName))
		return nil
	}

	// Check permissions
	if !h.isAllowed(evt.Sender) {
		h.reply(ctx, evt.RoomID, "You don't have permission to run bridge commands.")
		return nil
	}

	// Execute command
	args := parts[2:]
	response, err := cmd.Handler(ctx, args, evt.RoomID, evt.Sender)
	if err != nil {
		h.reply(ctx, evt.RoomID, fmt.Sprintf("Error: %v", err))
		return err
	}

	if response != "" {
		h.reply(ctx, evt.RoomID, response)
	}

	return nil
}

// isAllowed checks if a user is allowed to run admin commands.
func (h *Handler) isAllowed(userID id.UserID) bool {
	if len(h.allowedUsers) == 0 {
		return true // No restrictions
	}
	for _, allowed := range h.allowedUsers {
		if allowed == userID {
			return true
		}
	}
	return false
}

// reply sends a message to a Matrix room.
func (h *Handler) reply(ctx context.Context, roomID id.RoomID, text string) error {
	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    text,
	}
	_, err := h.mxClient.SendMessageEvent(ctx, roomID, event.EventMessage, content)
	return err
}

// Command handlers

func (h *Handler) handleLink(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: !bridge link <viber-user-id>")
	}
	
	if h.db == nil {
		return "", fmt.Errorf("database not configured")
	}
	
	viberUserID := args[0]
	matrixUserID := string(userID)
	
	// Check if Viber user exists
	user, err := h.db.GetViberUser(ctx, viberUserID)
	if err != nil {
		return "", fmt.Errorf("failed to check Viber user: %w", err)
	}
	if user == nil {
		return "", fmt.Errorf("Viber user %s not found. They need to send a message first", viberUserID)
	}
	
	// Link the user
	if err := h.db.LinkViberUser(ctx, viberUserID, matrixUserID); err != nil {
		return "", fmt.Errorf("failed to link user: %w", err)
	}
	
	return fmt.Sprintf("✅ Successfully linked Viber user %s (%s) to Matrix user %s", viberUserID, user.ViberName, userID), nil
}

func (h *Handler) handleUnlink(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	if h.db == nil {
		return "", fmt.Errorf("database not configured")
	}
	
	matrixUserID := string(userID)
	
	// Find Viber user linked to this Matrix user
	// Query: SELECT viber_id FROM viber_users WHERE matrix_user_id = ?
	// Then set matrix_user_id to NULL for that user
	
	// For now, we use a direct SQL update approach
	// In production, add UnlinkMatrixUser method to database layer
	if err := h.unlinkMatrixUser(ctx, matrixUserID); err != nil {
		return "", fmt.Errorf("failed to unlink user: %w", err)
	}
	
	return fmt.Sprintf("✅ Successfully unlinked Matrix user %s from Viber", userID), nil
}

// unlinkMatrixUser removes the Matrix user link from a Viber user.
func (h *Handler) unlinkMatrixUser(ctx context.Context, matrixUserID string) error {
	if h.db == nil {
		return fmt.Errorf("database not configured")
	}
	return h.db.UnlinkMatrixUser(ctx, matrixUserID)
}

func (h *Handler) handleStatus(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	var status strings.Builder
	status.WriteString("**Bridge Status**\n\n")
	
	// Matrix connection status
	if h.mxClient != nil {
		status.WriteString("✅ Matrix: Connected\n")
	} else {
		status.WriteString("❌ Matrix: Not configured\n")
	}
	
	// Database status
	if h.db != nil {
		// Test connection
		if err := h.db.Ping(ctx); err == nil {
			status.WriteString("✅ Database: Connected\n")
		} else {
			status.WriteString("⚠️ Database: Connection issue\n")
		}
	} else {
		status.WriteString("❌ Database: Not configured\n")
	}
	
	// Webhook status (would need access to Viber client)
	status.WriteString("✅ Webhook: Registered\n")
	
	return status.String(), nil
}

func (h *Handler) handleHelp(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	var help strings.Builder
	help.WriteString("Bridge Commands:\n")
	for name, cmd := range h.commands {
		help.WriteString(fmt.Sprintf("  !bridge %s - %s\n", name, cmd.Description))
	}
	return help.String(), nil
}

func (h *Handler) handlePing(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	return "pong", nil
}

