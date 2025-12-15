package alerts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/HardDie/ytmemchat/pkg/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Token            string
	MediaPath        string
	CommandsFilePath string
}

type Alerts struct {
	cfg      Config
	commands map[string]command
}

func New(cfg Config) (*Alerts, error) {
	cmd, err := parseCommands(cfg.CommandsFilePath)
	if err != nil {
		return nil, fmt.Errorf("parseCommands(): %w", err)
	}

	a := Alerts{
		cfg:      cfg,
		commands: make(map[string]command),
	}
	for _, it := range cmd.Commands {
		_, ok := a.commands[strings.ToLower(it.Name)]
		if ok {
			return nil, fmt.Errorf("duplicate command: %s", it.Name)
		}
		a.commands[strings.ToLower(it.Name)] = it
	}

	{
		data, _ := json.MarshalIndent(a.commands, "", "  ")
		log.Println(string(data))
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
	broadcast <- WebhookPayload{
		Filename: cmd.File,
		Volume:   utils.FromPtr(cmd.Volume, 1),
		Scale:    utils.FromPtr(cmd.Scale, 1),
	}
	return true
}

func someMain() {
	// Start the broadcaster in a goroutine
	go broadcaster()

	// 1. Media File Server: Serves files from the local 'media/' directory under the URL path '/media/'.
	// e.g., files in ./media are accessible via http://localhost:8080/media/...
	mediaDir := http.Dir("./media")
	mediaHandler := http.StripPrefix("/media/", http.FileServer(mediaDir))
	http.Handle("/media/", mediaHandler)

	// 2. HTTP/HTML Route: Handles the root path (/)
	http.HandleFunc("/", htmlHandler)

	//// 3. Webhook Route: Handles incoming webhook POST requests
	//http.HandleFunc("/webhook", webhookHandler)

	// 4. WebSocket Route: Handles real-time client connections
	http.HandleFunc("/ws", wsHandler)

	port := ":8080"
	log.Printf("Go service listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

type command struct {
	Name   string
	File   string
	Volume *float64
	Scale  *float64
}
type commands struct {
	Commands []command
}

func parseCommands(path string) (*commands, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %w", err)
	}
	defer file.Close()

	var cmd commands
	err = yaml.NewDecoder(file).Decode(&cmd)
	if err != nil {
		return nil, fmt.Errorf("yaml.Decode(): %w", err)
	}

	return &cmd, nil
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
