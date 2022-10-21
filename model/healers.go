package model

type HealerFrequency struct {
	Frequency int
	HealerKey string
	Healers   []string
}

type HealerDetails struct {
	FightID   int
	Report    string
	StartTime int
	EndTime   int
	Healers   []Healer
}

type Healer struct {
	Name  string
	Class string
	Spec  string
}
