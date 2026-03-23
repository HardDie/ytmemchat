// Package alerts handles the detection and execution of media-based commands
// found within chat messages. It parses a YAML configuration of commands
// and matches them against incoming text using a specific token prefix.
package alerts

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HardDie/ytmemchat/internal/server"
	"github.com/HardDie/ytmemchat/pkg/utils"
)

// Config defines the dependencies and settings required to initialize the Alerts system.
type Config struct {
	// Token is the character prefix that triggers a command (e.g., "@" or "!").
	Token string
	// MediaPath is the local filesystem path where alert media files are stored.
	MediaPath string
	// CommandsFilePath is the path to the YAML file defining the available commands.
	CommandsFilePath string
	// Broadcast is a channel used to send payloads to the WebSocket server for client-side rendering.
	Broadcast chan server.WebsocketPayload
}

// Alerts manages the lifecycle of chat commands and coordinates media broadcasting.
type Alerts struct {
	cfg       Config
	commands  map[string]command
	broadcast chan server.WebsocketPayload
}

// New creates a new Alerts instance. It parses the commands YAML file
// and validates that there are no duplicate command names (case-insensitive).
func New(cfg Config) (*Alerts, error) {
	cmd, err := parseCommands(cfg.CommandsFilePath)
	if err != nil {
		return nil, fmt.Errorf("parseCommands(): %w", err)
	}

	a := Alerts{
		cfg:       cfg,
		commands:  make(map[string]command),
		broadcast: cfg.Broadcast,
	}
	for _, it := range cmd.Commands {
		_, ok := a.commands[strings.ToLower(it.Name)]
		if ok {
			return nil, fmt.Errorf("duplicate command: %s", it.Name)
		}
		a.commands[strings.ToLower(it.Name)] = it
	}

	return &a, nil
}

// Alert scans a message for a command token. If a valid command is found,
// it pushes a WebsocketPayload to the broadcast channel and returns true.
func (a *Alerts) Alert(msg string) bool {
	token := findToken(a.cfg.Token, msg)
	if token == "" {
		return false
	}
	cmd, ok := a.commands[strings.ToLower(token)]
	if !ok {
		return false
	}
	a.broadcast <- server.WebsocketPayload{
		Type:     server.PayloadTypeAlert,
		Filename: cmd.File,
		Volume:   utils.FromPtr(cmd.Volume, 1),
		Scale:    utils.FromPtr(cmd.Scale, 1),
	}
	return true
}

// GetMediaHandler returns an http.Handler configured to serve static files
// from the media directory under the "/media/" URL prefix.
func (a *Alerts) GetMediaHandler() http.Handler {
	mediaDir := http.Dir(a.cfg.MediaPath)
	mediaHandler := http.StripPrefix("/media/", http.FileServer(mediaDir))
	return mediaHandler
}

func findToken(token, str string) string {
	// Trim everything before token.
	_, str, _ = strings.Cut(str, string(token))
	if str == "" {
		return ""
	}
	// If string have spaces, return only first word.
	return strings.Split(str, " ")[0]
}
