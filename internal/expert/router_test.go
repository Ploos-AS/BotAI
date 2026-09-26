package expert

import "testing"

func TestRoute(t *testing.T) {
	tests := []struct{ text, want string }{
		{"Why do I get IRC numeric 433 for my nick?", "irc"},
		{"How do I create an ARexx port on AmigaOS?", "amiga"},
		{"How do I configure systemd on Debian Linux?", "linux"},
		{"How should I debug IPv6 routing and DNS?", "networking"},
		{"What hardening and least privilege should I use?", "security"},
		{"How does the VIC-II work on a Commodore 64?", "c64"},
		{"How does GEMDOS work on an Atari ST?", "atari"},
		{"What should I cook tonight?", "general"},
		{"Can an Amiga IRC client use SASL?", "irc"},
	}
	for _, tc := range tests {
		if got := Route(tc.text).ID; got != tc.want {
			t.Fatalf("Route(%q)=%q want %q", tc.text, got, tc.want)
		}
	}
}

func TestListContainsAllExperts(t *testing.T) {
	want := []string{"general","irc","amiga","linux","networking","security","c64","atari"}
	got := List()
	if len(got) != len(want) { t.Fatalf("len=%d want=%d", len(got), len(want)) }
	for i, id := range want {
		if got[i].ID != id { t.Fatalf("expert[%d]=%q want=%q", i, got[i].ID, id) }
		if id != "general" && len(got[i].RouteTerms) == 0 { t.Fatalf("%s has no route terms", id) }
	}
}
