package model

type RoleFrequency struct {
	Frequency      int
	CompositionKey string
	Composition    []string
}

type GuildFrequency struct {
	Frequency int
	GuildID   int
	GuildName string
}

type Player struct {
	Name  string
	Class string
	Spec  string
}
