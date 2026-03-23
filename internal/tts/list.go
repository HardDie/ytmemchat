package tts

import (
	"fmt"
	"runtime"
)

// VoiceInfo holds detailed metadata about a system voice.
type VoiceInfo struct {
	Name     string // The identifier used in configuration (e.g., "Alex" or "Microsoft David")
	Language string // The locale code (e.g., "en-US")
	Gender   string // Voice gender ("Male", "Female")
	Details  string // Raw system description
}

// GetAvailableVoices queries the operating system and returns a list
// of all installed TTS voices.
func GetAvailableVoices() ([]VoiceInfo, error) {
	voices, err := getAvailableVoices()
	if err != nil {
		return nil, fmt.Errorf("failed to get voice info on %s: %w", runtime.GOOS, err)
	}
	return voices, nil
}
