package tts

import (
	"fmt"
	"runtime"
)

// VoiceInfo holds detailed information about a text-to-speech voice.
type VoiceInfo struct {
	Name     string // The command-line identifier (e.g., "Microsoft David Desktop")
	Language string // The language or locale (e.g., "en-US", "fr-FR")
	Gender   string // The gender of the voice (e.g., "Male", "Female")
	Details  string // Any other raw details provided by the system tool
}

// Speak pronounces the given text using the specified voice Name.
func Speak(text string, voiceName string) error {
	return speak(text, voiceName)
}

// GetAvailableVoices returns a list of detailed VoiceInfo structs available on the current system.
func GetAvailableVoices() ([]VoiceInfo, error) {
	voices, err := getAvailableVoices()
	if err != nil {
		return nil, fmt.Errorf("failed to get voice info on %s: %w", runtime.GOOS, err)
	}
	return voices, nil
}
