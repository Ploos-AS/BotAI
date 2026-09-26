package expert

import "strings"

func score(text string, terms []string) int {
	text = strings.ToLower(text)
	n := 0
	for _, term := range terms {
		if strings.Contains(text, strings.ToLower(term)) { n++ }
	}
	return n
}

// Route deterministically chooses a specialist only for one unique best
// positive match. Ties and no-match cases intentionally fall back to general.
func Route(message string) Profile {
	p, ids := snapshot()
	bestID, bestScore, ties := "general", 0, 0
	for _, id := range ids {
		if id == "general" { continue }
		s := score(message, p[id].RouteTerms)
		if s > bestScore {
			bestID, bestScore, ties = id, s, 1
		} else if s > 0 && s == bestScore {
			ties++
		}
	}
	if bestScore == 0 || ties != 1 { return p["general"] }
	return p[bestID]
}
