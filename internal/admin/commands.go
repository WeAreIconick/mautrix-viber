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
)

// Command represents a bridge admin command.
type Command struct {
	Name        string
	Description string
	Handler     func(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error)
}

// Handler manages bridge admin commands in Matrix rooms.
type Handler struct {
	commands    map[string]Command
	mxClient    *mautrix.Client
	allowedUsers []id.UserID // Users allowed to run admin commands
}

// NewHandler creates a new admin command handler.
func NewHandler(mxClient *mautrix.Client, allowedUsers []id.UserID) *Handler {
	h := &Handler{
		commands:     make(map[string]Command),
		mxClient:     mxClient,
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
	// TODO: Implement actual linking logic with database
	return fmt.Sprintf("Linking Viber user %s to Matrix user %s... (not implemented yet)", args[0], userID), nil
}

func (h *Handler) handleUnlink(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	// TODO: Implement unlinking logic
	return fmt.Sprintf("Unlinking Viber account from Matrix user %s... (not implemented yet)", userID), nil
}

func (h *Handler) handleStatus(ctx context.Context, args []string, roomID id.RoomID, userID id.UserID) (string, error) {
	// TODO: Get actual bridge status
	status := `Bridge Status:
- Matrix: Connected
- Viber: Connected
- Webhook: Registered
- Messages bridged: 0
- Users linked: 0`
	return status, nil
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

