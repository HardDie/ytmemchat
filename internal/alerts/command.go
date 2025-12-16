package alerts

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

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
