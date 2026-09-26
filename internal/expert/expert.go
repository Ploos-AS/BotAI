package expert

type Profile struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

var profiles = map[string]Profile{
	"general": {ID: "general", Description: "General conversational assistant"},
	"irc":     {ID: "irc", Description: "IRC protocol, clients, servers and operations"},
	"amiga":   {ID: "amiga", Description: "Amiga hardware, AmigaOS, ARexx and development"},
}

func Get(id string) (Profile, bool) {
	p, ok := profiles[id]
	return p, ok
}

func List() []Profile {
	return []Profile{profiles["general"], profiles["irc"], profiles["amiga"]}
}
