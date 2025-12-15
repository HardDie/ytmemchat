package alerts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindToken(t *testing.T) {
	tests := map[string]struct {
		Token string
		Str   string
		Want  string
	}{
		"simple": {
			Token: "@",
			Str:   "@some",
			Want:  "some",
		},
		"only symbol": {
			Token: "@",
			Str:   "@",
			Want:  "",
		},
		"last symbol": {
			Token: "@",
			Str:   "some @",
			Want:  "",
		},
		"in the beginning": {
			Token: "@",
			Str:   "@some more text",
			Want:  "some",
		},
		"in the middle": {
			Token: "@",
			Str:   "check @some more",
			Want:  "some",
		},
		"in the end": {
			Token: "@",
			Str:   "check @some",
			Want:  "some",
		},
		"not found": {
			Token: "@",
			Str:   "check some",
			Want:  "",
		},
		"token in word": {
			Token: "@",
			Str:   "check @so@me",
			Want:  "so@me",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := findToken(tc.Token, tc.Str)
			assert.Equal(t, tc.Want, got)
		})
	}
}
