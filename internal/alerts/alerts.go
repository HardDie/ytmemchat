package alerts

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HardDie/ytmemchat/internal/server"
	"github.com/HardDie/ytmemchat/pkg/utils"
)

type Config struct {
	Token            string
	MediaPath        string
	CommandsFilePath string
	Broadcast        chan server.WebsocketPayload
}

type Alerts struct {
	cfg       Config
	commands  map[string]command
	broadcast chan server.WebsocketPayload
}

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
