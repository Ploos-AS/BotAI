package expert

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ConfigProfile struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	SystemPrompt string   `json:"system_prompt"`
	RouteTerms   []string `json:"route_terms,omitempty"`
}

func LoadFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil { return fmt.Errorf("read expert config: %w", err) }
	var cfg struct {
		Experts []ConfigProfile `json:"experts"`
	}
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil { return fmt.Errorf("decode expert config: %w", err) }
	if len(cfg.Experts) == 0 { return errors.New("expert config contains no experts") }

	next := make(map[string]Profile, len(cfg.Experts))
	nextOrder := make([]string, 0, len(cfg.Experts))
	for _, e := range cfg.Experts {
		id := strings.ToLower(strings.TrimSpace(e.ID))
		if id == "" || id == "auto" { return fmt.Errorf("invalid expert id %q", e.ID) }
		if _, exists := next[id]; exists { return fmt.Errorf("duplicate expert id %q", id) }
		desc := strings.TrimSpace(e.Description)
		prompt := strings.TrimSpace(e.SystemPrompt)
		if desc == "" || prompt == "" { return fmt.Errorf("expert %q requires description and system_prompt", id) }
		terms := make([]string, 0, len(e.RouteTerms))
		for _, term := range e.RouteTerms {
			term = strings.TrimSpace(term)
			if term != "" { terms = append(terms, term) }
		}
		next[id] = Profile{ID:id, Description:desc, SystemPrompt:prompt, RouteTerms:terms}
		nextOrder = append(nextOrder, id)
	}
	if _, ok := next["general"]; !ok { return errors.New("expert config must define general") }

	profiles = next
	order = nextOrder
	return nil
}
