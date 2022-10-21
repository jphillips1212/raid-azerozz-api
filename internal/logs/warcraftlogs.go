package warcraftlogs

import (
	wl "github.com/jphillips1212/roztools-api/internal/logs/warcraftlogs"
)

// WarcraftLogs interface declares how to interact with underlying fight logging technology
// This interface isn't abstracted from WarcraftLogs explicitly as the returning models are reusing the graphql queries
type WarcraftLogs interface {
	GetReportsForEncounter(fightId, page int) (*wl.QueryFightByID, error)

	GetEncounterComposition(reportCode string, fightId int, startTime int, endTime int) ([]wl.EncounterActor, error)
	GetFightTimes(reportCode string, fightId int) (wl.FightTimes, error)
}
