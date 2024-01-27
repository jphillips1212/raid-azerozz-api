package db

import (
	"sync"

	"github.com/jphillips1212/roztools-api/model"
)

// DB interface declares how to interact with underlying database technology
type DB interface {
	GetIfReportExists(reportCode, encounterName string, fightID int) bool
	GetIfMaxReportsForEncounter(encounterID int) bool
	GetAllEncounters(encounterName string) ([]model.EncounterData, error)

	SaveHealerFrequencies(encounterName, healerKey string, healerFrequency model.RoleFrequency, wg *sync.WaitGroup)
	GetHealerFrequencies(encounterName string) ([]model.RoleFrequency, error)

	SaveEncounterData(encounterName string, encounterData model.EncounterData) error
}
