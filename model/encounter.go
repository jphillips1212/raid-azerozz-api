package model

type EncounterDataChannel struct {
	Error error
	EncounterData
}

type EncounterData struct {
	StartTime    int
	EndTime      int
	Code         string
	FightID      int
	GuildID      int
	GuildName    string
	GuildFaction int
	Healers      []Player
	Tanks        []Player
}
