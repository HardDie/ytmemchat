package tts

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/HardDie/ytmemchat/internal/server"
	"github.com/HardDie/ytmemchat/pkg/utils"
)

type TTS struct {
	cfg       Config
	broadcast chan server.WebsocketPayload
}

func New(cfg Config) *TTS {
	return &TTS{
		cfg:       cfg,
		broadcast: cfg.Broadcast,
	}
}

// LastAudioBuffer holds the most recently synthesized audio data and its format.
type LastAudioBuffer struct {
	Data        []byte
	ContentType string // e.g., "audio/wav", "audio/aiff"
	Format      string // e.g., "wav", "aiff"
}

// --- Internal Buffer State ---
var lastAudioBuffer = LastAudioBuffer{}
var bufferMutex sync.RWMutex

// --- Internal Function (Used by platform files) ---

// setSynthesizedAudio is a private function used by platform-specific files
// to safely update the global audio buffer after synthesis.
func setSynthesizedAudio(data []byte, format string) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	lastAudioBuffer.Data = data
	lastAudioBuffer.Format = format
	lastAudioBuffer.ContentType = getContentType(format)
}

// getContentType maps the short format string to the standard MIME type.
func getContentType(format string) string {
	if format == "aiff" {
		return "audio/aiff"
	}
	if format == "wav" {
		return "audio/wav"
	}
	return fmt.Sprintf("audio/%s", format)
}

// --- Public API Functions ---

// Speak pronounces the given text using the specified voice Name.
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

// SynthesizeAudio synthesizes the given text using the specified voice.
// The result is stored internally and accessible via GetLastSynthesizedAudio.
func (t *TTS) SynthesizeAudio(text string) error {
	// The actual implementation is provided by tts_windows.go or tts_unix.go
	if _, _, err := synthesize(text, t.cfg.VoiceName); err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}
	t.broadcast <- server.WebsocketPayload{
		Type:   server.PayloadTypeTTS,
		Volume: utils.FromPtr(t.cfg.Volume, 1),
	}
	return nil
}

// Handler to serve the audio buffer directly from memory.
func (t *TTS) GetPlaybackHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve audio directly from the package
	audioBuffer := GetLastSynthesizedAudio()

	if len(audioBuffer.Data) == 0 {
		http.Error(w, "No audio data available.", http.StatusNotFound)
		return
	}

	// Set required headers for streaming
	w.Header().Set("Content-Type", audioBuffer.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(audioBuffer.Data)))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // Prevent caching

	if _, err := w.Write(audioBuffer.Data); err != nil {
		log.Printf("Error writing audio data: %v", err)
	}
}

// GetLastSynthesizedAudio returns a copy of the most recently synthesized
// audio buffer and its metadata.
func GetLastSynthesizedAudio() LastAudioBuffer {
	bufferMutex.RLock()
	defer bufferMutex.RUnlock()

	// Return a copy of the buffer data to prevent external modification
	// Note: Creating a full copy of the []byte buffer is safer but can be slow
	// for very large files. For TTS, it is generally safe.
	dataCopy := make([]byte, len(lastAudioBuffer.Data))
	copy(dataCopy, lastAudioBuffer.Data)

	return LastAudioBuffer{
		Data:        dataCopy,
		ContentType: lastAudioBuffer.ContentType,
		Format:      lastAudioBuffer.Format,
	}
}
