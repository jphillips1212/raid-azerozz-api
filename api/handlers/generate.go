package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	analysis "github.com/jphillips1212/roztools-api/internal/service/analysis"
)

type HealerCompRequest struct {
	EncounterName string `json:"encounter_name"`
	Persist       bool   `json:"persist"`
}

func (handler Handler) GenerateHealerComposition(w http.ResponseWriter, r *http.Request) {
	var compReq *HealerCompRequest
	if err := json.NewDecoder(r.Body).Decode(&compReq); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding body for generating healer composition: [%v]", err), 400)
	}

	asis := analysis.Analysis{
		Database: handler.Database,
		Logs:     handler.Logs,
	}

	service := Service{
		asis,
	}

	service.Service.SaveHealerLogs(compReq.EncounterName, compReq.Persist)
	service.Service.SaveHealerAnalysis(compReq.EncounterName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	fmt.Printf("\nSuccessfully generated healer analysis for: [%s]\n", compReq.EncounterName)
}
