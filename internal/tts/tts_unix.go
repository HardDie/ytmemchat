//go:build !windows

package tts

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// speak remains the same as it only needs the voice Name.
func speak(text string, voiceName string) error {
	// ... (Speak logic is unchanged)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("say", "-v", voiceName, text)
	case "linux":
		cmd = exec.Command("espeak", "-v", voiceName, text)
	default:
		return fmt.Errorf("unsupported operating system for native TTS: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("TTS command failed on %s: %w", runtime.GOOS, err)
	}
	return nil
}

// getAvailableVoices lists voices on macOS/Linux and parses details.
func getAvailableVoices() ([]VoiceInfo, error) {
	switch runtime.GOOS {
	case "darwin": // macOS
		// Command: 'say -v ?' lists available voices with details.
		cmd := exec.Command("say", "-v", "?")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute 'say -v ?' on macOS: %w", err)
		}

		// The output has lines like: "Alex      en_US # Male"
		// Regex captures Name, Language code (en_US), and Detail/Gender
		re := regexp.MustCompile(`^(\w+)\s+([a-z]{2}_[A-Z]{2})[^\n]*#\s*([^\n]+)`)
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")

		var voices []VoiceInfo
		for _, line := range lines {
			if matches := re.FindStringSubmatch(line); len(matches) == 4 {
				// matches[1] = Name, matches[2] = Language, matches[3] = Details/Gender
				details := strings.TrimSpace(matches[3])
				gender := ""
				if strings.Contains(details, "Male") {
					gender = "Male"
				} else if strings.Contains(details, "Female") {
					gender = "Female"
				}

				voices = append(voices, VoiceInfo{
					Name:     matches[1],
					Language: matches[2],
					Gender:   gender,
					Details:  details,
				})
			}
		}
		return voices, nil

	case "linux": // Linux (requires 'espeak' to be installed)
		// Command: 'espeak --voices' lists all voices.
		// Output columns: | ID | Name | Language | Gender | Age |
		// We'll use a simpler version to get language and name.
		cmd := exec.Command("espeak", "--voices")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute 'espeak --voices' on Linux. Is 'espeak' installed? Error: %w", err)
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		var voices []VoiceInfo

		// Skip the header line(s).
		for i, line := range lines {
			if i < 1 || strings.Contains(line, "Language") {
				continue
			}

			// Fields are usually separated by excessive whitespace.
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				// Fields: | ID | Name | Language | Gender | Age |
				voiceName := fields[1]
				language := fields[2]
				gender := fields[3]

				voices = append(voices, VoiceInfo{
					Name:     voiceName,
					Language: language,
					Gender:   gender,
					Details:  fmt.Sprintf("Age: %s", fields[4]),
				})
			}
		}
		return voices, nil

	default:
		return nil, fmt.Errorf("voice listing not supported on %s", runtime.GOOS)
	}
}
