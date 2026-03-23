package alerts

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// command represents a single media trigger.
type command struct {
	Name   string   // The trigger word (without the token)
	File   string   // The filename in the media directory
	Volume *float64 // Optional: Audio volume override (0.0 to 1.0)
	Scale  *float64 // Optional: Visual scale override
}

// commands is a wrapper for unmarshaling the YAML configuration.
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
