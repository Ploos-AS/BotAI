package expert

import "strings"

var ircTerms = []string{
	"irc", "ircv3", "sasl", "nickserv", "chanserv", "bouncer", "bnc",
	"numerics", "numeric 433", "ping pong", "privmsg", "notice", "whois",
	"channel mode", "ircd",
}

var amigaTerms = []string{
	"amiga", "amigaos", "workbench", "arexx", "kickstart", "68000", "68020",
	"68030", "68040", "68060", "m68k", "chip ram", "fast ram", "intuition",
	"exec.library", "dos.library", "bsdsocket.library", "a500", "a600", "a1200",
	"a2000", "a3000", "a4000", "ocs", "ecs", "aga",
}

func score(text string, terms []string) int {
	text = strings.ToLower(text)
	n := 0
	for _, term := range terms {
		if strings.Contains(text, term) { n++ }
	}
	return n
}

// Route deterministically chooses a specialist only for a positive domain match.
// Ties and no-match cases intentionally fall back to general.
func Route(message string) Profile {
	irc := score(message, ircTerms)
	amiga := score(message, amigaTerms)
	if irc > amiga && irc > 0 { return profiles["irc"] }
	if amiga > irc && amiga > 0 { return profiles["amiga"] }
	return profiles["general"]
}
