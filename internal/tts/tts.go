// Package tts provides Text-to-Speech capabilities by leveraging native system tools.
// It supports Windows (via PowerShell/System.Speech), macOS (via 'say'),
// and Linux (via 'espeak').
package tts

import (
	"fmt"
	"runtime"

	"github.com/HardDie/ytmemchat/internal/server"
	"github.com/HardDie/ytmemchat/pkg/utils"
)

// TTS handles the synthesis and broadcasting of chat messages as audio.
type TTS struct {
	cfg       Config
	broadcast chan server.WebsocketPayload
}

// New creates a new TTS instance with the provided configuration.
func New(cfg Config) *TTS {
	return &TTS{
		cfg:       cfg,
		broadcast: cfg.Broadcast,
	}
}

// --- Public API Functions ---

// Speak immediately plays the provided text through the system's
// default audio output (local playback).
func Speak(text string, voiceName string) error {
	return speak(text, voiceName)
}

// SynthesizeToBuffer synthesizes the given text using the specified voice
// and returns the raw audio data as a byte slice and the file format (e.g., "mp3", "wav").
func SynthesizeToBuffer(text string, voiceName string) ([]byte, string, error) {
	// The actual implementation is provided by tts_windows.go or tts_unix.go
	audioData, format, err := synthesize(text, voiceName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to synthesize audio on %s: %w", runtime.GOOS, err)
	}
	return audioData, format, nil
}

// SynthesizeAudio converts text to audio bytes and broadcasts it
// via the WebSocket channel for remote client playback.
func (t *TTS) SynthesizeAudio(text string) error {
	// The actual implementation is provided by tts_windows.go or tts_unix.go
	audioData, _, err := synthesize(text, t.cfg.VoiceName)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}
	t.broadcast <- server.WebsocketPayload{
		Type:    server.PayloadTypeTTS,
		Payload: audioData,
		Volume:  utils.FromPtr(t.cfg.Volume, 1),
	}
	return nil
}
