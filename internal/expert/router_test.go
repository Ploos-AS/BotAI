package expert

import "testing"

func TestRoute(t *testing.T) {
	tests := []struct{ text, want string }{
		{"Why do I get IRC numeric 433 for my nick?", "irc"},
		{"How do I create an ARexx port on AmigaOS?", "amiga"},
		{"What should I cook tonight?", "general"},
		{"Can an Amiga IRC client use SASL?", "general"},
	}
	for _, tc := range tests {
		if got := Route(tc.text).ID; got != tc.want {
			t.Fatalf("Route(%q)=%q want %q", tc.text, got, tc.want)
		}
	}
}
