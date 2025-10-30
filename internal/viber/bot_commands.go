// Package viber bot_commands implements Viber bot command parsing and Matrix command bridge.
package viber

import (
	"context"
	"fmt"
	"strings"
)

// BotCommand represents a parsed bot command.
type BotCommand struct {
	Command string
	Args    []string
	Sender  string
	ChatID  string
}

// CommandHandler handles bot commands.
type CommandHandler func(ctx context.Context, cmd BotCommand) (string, error)

// BotCommandManager manages Viber bot commands.
type BotCommandManager struct {
	handlers map[string]CommandHandler
	prefix   string // Command prefix (e.g., "/" or "!")
}

// NewBotCommandManager creates a new bot command manager.
func NewBotCommandManager(prefix string) *BotCommandManager {
	return &BotCommandManager{
		handlers: make(map[string]CommandHandler),
		prefix:   prefix,
	}
}

// RegisterCommand registers a command handler.
func (bcm *BotCommandManager) RegisterCommand(command string, handler CommandHandler) {
	bcm.handlers[strings.ToLower(command)] = handler
}

// ParseCommand parses a command from text.
func (bcm *BotCommandManager) ParseCommand(text string) *BotCommand {
	if !strings.HasPrefix(text, bcm.prefix) {
		return nil
	}
	
	// Remove prefix and split
	text = strings.TrimPrefix(text, bcm.prefix)
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return nil
	}
	
	return &BotCommand{
		Command: strings.ToLower(parts[0]),
		Args:    parts[1:],
	}
}

// HandleCommand handles a parsed command.
func (bcm *BotCommandManager) HandleCommand(ctx context.Context, cmd BotCommand) (string, error) {
	handler, ok := bcm.handlers[cmd.Command]
	if !ok {
		return "", fmt.Errorf("unknown command: %s", cmd.Command)
	}
	
	return handler(ctx, cmd)
}

// HandleMessage checks if a message is a command and handles it.
func (bcm *BotCommandManager) HandleMessage(ctx context.Context, text, senderID, chatID string) (bool, string, error) {
	cmd := bcm.ParseCommand(text)
	if cmd == nil {
		return false, "", nil
	}
	
	cmd.Sender = senderID
	cmd.ChatID = chatID
	
	response, err := bcm.HandleCommand(ctx, *cmd)
	if err != nil {
		return true, "", err
	}
	
	return true, response, nil
}

// RegisterDefaultCommands registers default bot commands.
func (bcm *BotCommandManager) RegisterDefaultCommands() {
	// Help command
	bcm.RegisterCommand("help", func(ctx context.Context, cmd BotCommand) (string, error) {
		var commands []string
		for cmdName := range bcm.handlers {
			commands = append(commands, bcm.prefix+cmdName)
		}
		return fmt.Sprintf("Available commands: %s", strings.Join(commands, ", ")), nil
	})
	
	// Status command
	bcm.RegisterCommand("status", func(ctx context.Context, cmd BotCommand) (string, error) {
		return "Bridge is running", nil
	})
	
	// Ping command
	bcm.RegisterCommand("ping", func(ctx context.Context, cmd BotCommand) (string, error) {
		return "pong", nil
	})
}

