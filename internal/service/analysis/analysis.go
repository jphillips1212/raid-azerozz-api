package analysis

import (
	db "github.com/jphillips1212/roztools-api/internal/db"
	logs "github.com/jphillips1212/roztools-api/internal/logs"
)

type Analysis struct {
	Database db.DB
	Logs     logs.WarcraftLogs
}

var EncounterIDs = map[int]string{
	2587: "Eranog",
	2639: "Terros",
	2590: "The Primal Council",
	2592: "Sennarth The Cold Breath",
	2635: "Dathea",
	2605: "Kurog Grimtotem",
	2614: "Broodkeeper Diurna",
	2607: "Raszageth the Storm-Eater",
}
