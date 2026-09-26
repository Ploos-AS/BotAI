package expert

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	oldProfiles, oldOrder := profiles, order
	defer func(){ profiles, order = oldProfiles, oldOrder }()

	path := filepath.Join(t.TempDir(), "experts.json")
	data := `{"experts":[
		{"id":"general","description":"General","system_prompt":"General prompt"},
		{"id":"retro","description":"Retro computers","system_prompt":"Retro prompt","route_terms":["retro","8-bit"]}
	]}`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil { t.Fatal(err) }
	if err := LoadFile(path); err != nil { t.Fatal(err) }
	if got := Route("I like 8-bit retro computers").ID; got != "retro" { t.Fatalf("route=%q", got) }
	if len(List()) != 2 { t.Fatalf("experts=%d", len(List())) }
}

func TestLoadFileRejectsMissingGeneral(t *testing.T) {
	path := filepath.Join(t.TempDir(), "experts.json")
	if err := os.WriteFile(path, []byte(`{"experts":[{"id":"irc","description":"IRC","system_prompt":"IRC"}]}`), 0600); err != nil { t.Fatal(err) }
	if err := LoadFile(path); err == nil { t.Fatal("expected missing general error") }
}

func TestLoadFileRejectsDuplicate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "experts.json")
	data := `{"experts":[{"id":"general","description":"a","system_prompt":"a"},{"id":"GENERAL","description":"b","system_prompt":"b"}]}`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil { t.Fatal(err) }
	if err := LoadFile(path); err == nil { t.Fatal("expected duplicate error") }
}
