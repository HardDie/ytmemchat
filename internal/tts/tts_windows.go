//go:build windows

package tts

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// speak remains the same as it only needs the voice Name.
func speak(text string, voiceName string) error {
	// ... (Speak logic is unchanged)
	text = strings.ReplaceAll(text, "'", "''")

	powershellScript := fmt.Sprintf(`
		Add-Type -AssemblyName System.Speech;
		$synth = New-Object System.Speech.Synthesis.SpeechSynthesizer;
		$synth.SelectVoice('%s'); 
		$synth.Speak('%s');       
	`, voiceName, text)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", powershellScript)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PowerShell execution failed (Voice: %s): %v. Output: %s", voiceName, err, string(output))
	}
	return nil
}

// getAvailableVoices now parses JSON output from PowerShell for detailed VoiceInfo.
func getAvailableVoices() ([]VoiceInfo, error) {
	// PowerShell command to list VoiceInfo properties and format them as JSON.
	// Note: We use the actual property names from the .NET object.
	powershellCommand := `
		Add-Type -AssemblyName System.Speech;
		(Get-SpeechSynthesizer).GetInstalledVoices() | 
		Select-Object -ExpandProperty VoiceInfo | 
		Select-Object Name, Culture, Gender, Description | 
		ConvertTo-Json -Compress
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", powershellCommand)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell for voice list: %w", err)
	}

	// Define a temporary struct matching the JSON keys
	type PsVoiceInfo struct {
		Name        string
		Culture     string // This will contain the language code (e.g., "en-US")
		Gender      string // e.g., "Male" or "Female"
		Description string
	}

	var psVoices []PsVoiceInfo
	// The output is a JSON array of objects
	if err := json.Unmarshal(output, &psVoices); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output from PowerShell: %w. Raw output: %s", err, string(output))
	}

	// Convert the temporary struct to our public VoiceInfo struct
	var voices []VoiceInfo
	for _, pv := range psVoices {
		voices = append(voices, VoiceInfo{
			Name:     pv.Name,
			Language: pv.Culture,
			Gender:   pv.Gender,
			Details:  pv.Description,
		})
	}

	return voices, nil
}
