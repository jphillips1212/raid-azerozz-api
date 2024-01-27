package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	analysis "github.com/jphillips1212/roztools-api/internal/service/analysis"
)

type GenerateEncounterRequest struct {
	EncounterID int  `json:"encounter_id"`
	Persist     bool `json:"persist"`
}

func (handler Handler) GenerateEncounterData(w http.ResponseWriter, r *http.Request) {
	var encounterReq *GenerateEncounterRequest
	if err := json.NewDecoder(r.Body).Decode(&encounterReq); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding body for generating encounter data: [%v]", err), 400)
		return
	}

	asis := analysis.Analysis{
		Database: handler.Database,
		Logs:     handler.Logs,
	}

	service := Service{
		asis,
	}

	if _, ok := analysis.EncounterIDs[encounterReq.EncounterID]; !ok {
		http.Error(w, fmt.Sprintf("An encounter ID doesn't exist for the ID provided: [%v]", encounterReq.EncounterID), 404)
		return
	}

	// TODO: Check Reports Database for 2500 entries already saved and skip if true
	service.ScrapeAndSaveBoss(encounterReq.EncounterID, encounterReq.Persist)
	fmt.Printf("\nSuccessfully scraped WarcraftLogs and saved to database for : [%s]\n", analysis.EncounterIDs[encounterReq.EncounterID])

	// TODO: Check Analysis Database for 2500 total recorded kills already analysed and skip if true
	service.AnalyseAndSaveBoss(encounterReq.EncounterID)
	fmt.Printf("\nSuccessfully analysed xxx amount of logs for encounter : [%s]\n", analysis.EncounterIDs[encounterReq.EncounterID])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
