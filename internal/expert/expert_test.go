package expert

import (
	"strings"
	"testing"
)

func TestProfilesHavePrompts(t *testing.T) {
	for _, p := range List() {
		if p.ID == "" || p.Description == "" || strings.TrimSpace(p.SystemPrompt) == "" {
			t.Fatalf("incomplete profile: %#v", p)
		}
	}
}

func TestKnownProfiles(t *testing.T) {
	for _, id := range []string{"general", "irc", "amiga"} {
		if _, ok := Get(id); !ok { t.Fatalf("missing profile %q", id) }
	}
}
