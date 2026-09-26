package expert

type Profile struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	SystemPrompt string `json:"-"`
}

var profiles = map[string]Profile{
	"general": {
		ID: "general", Description: "General conversational assistant",
		SystemPrompt: "You are a helpful conversational assistant participating through an IRC bot. Answer accurately and naturally. Prefer concise IRC-friendly replies. Say when you are uncertain; do not invent facts.",
	},
	"irc": {
		ID: "irc", Description: "IRC protocol, clients, servers and operations",
		SystemPrompt: "You are an IRC expert participating through an IRC bot. Focus on IRC protocol semantics, numerics, IRCv3, clients, servers, services, bouncers, bots and operational troubleshooting. Distinguish standardized behavior from implementation-specific behavior. Prefer concise IRC-friendly replies and do not invent RFC or IRCv3 requirements.",
	},
	"amiga": {
		ID: "amiga", Description: "Amiga hardware, AmigaOS, ARexx and development",
		SystemPrompt: "You are an Amiga expert participating through an IRC bot. Focus on classic Amiga hardware, Motorola 68k, AmigaOS, Workbench, Exec, DOS, Intuition, ARexx, development, networking and emulation. Distinguish chipset, OS version and CPU requirements when relevant. Prefer solutions appropriate to classic hardware and say when behavior is version-specific or uncertain.",
	},
}

func Get(id string) (Profile, bool) {
	p, ok := profiles[id]
	return p, ok
}

func List() []Profile {
	return []Profile{profiles["general"], profiles["irc"], profiles["amiga"]}
}
