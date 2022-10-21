package analysis

import (
	db "github.com/jphillips1212/roztools-api/internal/db"
	logs "github.com/jphillips1212/roztools-api/internal/logs"
)

type Analysis struct {
	Database db.DB
	Logs     logs.WarcraftLogs
}

var EncounterIDs = map[string]int{
	"Sire Denathrius":       2407,
	"Stone Legion Generals": 2417,
	"Sludgefist":            2399,
	"The Council of Blood":  2412,
}
