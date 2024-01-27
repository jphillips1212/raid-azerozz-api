package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func (handler Handler) GetHealerFrequency(w http.ResponseWriter, r *http.Request) {
	encounterName := chi.URLParam(r, "encounterName")
	encounterName = strings.Replace(encounterName, "_", " ", -1)
	if encounterName == "" {
		http.Error(w, "No encounter name provided", 404)
	}

	healerFrequencies, err := handler.Database.GetHealerFrequencies(encounterName)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(healerFrequencies); err != nil {
		http.Error(w, err.Error(), 500)
	}

	fmt.Printf("\nSuccessfully got healer frequency for: [%s]\n", encounterName)
}
