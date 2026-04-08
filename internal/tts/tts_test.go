package tts

import (
	"strings"
	"testing"
)

func TestRemoveEmojis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple emoji",
			input:    "Hello World! 😊",
			expected: "Hello World!",
		},
		{
			name:     "Multiple emojis",
			input:    "🔥 Live Chat 🔥 is crazy 🚀",
			expected: "Live Chat  is crazy",
		},
		{
			name:     "Emoji with skin tone",
			input:    "High five! ✋🏾",
			expected: "High five!",
		},
		{
			name:     "Flag and complex symbols",
			input:    "Welcome from 🇺🇸! ✌️", // The peace sign here has the hidden FE0F
			expected: "Welcome from !",
		},
		{
			name:     "Text only",
			input:    "Just a normal message.",
			expected: "Just a normal message.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeEmojis(tt.input)
			// Clean up extra spaces that might be left behind after removal
			result = strings.Join(strings.Fields(result), " ")

			// Simple check
			expectedClean := strings.Join(strings.Fields(tt.expected), " ")

			if result != expectedClean {
				t.Errorf("expected %q, got %q", expectedClean, result)
			}
		})
	}
}
