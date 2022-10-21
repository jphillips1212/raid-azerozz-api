package db

import (
	"sync"

	"github.com/jphillips1212/roztools-api/model"
)

// DB interface declares how to interact with underlying database technology
type DB interface {
	SaveHealerFrequencies(encounterName, healerKey string, healerFrequency model.HealerFrequency, wg *sync.WaitGroup)
	GetHealerFrequencies(encounterName string) ([]model.HealerFrequency, error)

	SaveHealerComposition(encounterName string, healerDetails model.HealerDetails) error
	GetIfReportExists(reportCode, encounterName string, fightID int) bool
	GetAllHealerCompositions(encounterName string) ([]model.HealerDetails, error)
}
