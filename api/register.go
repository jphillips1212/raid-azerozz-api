package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// RozTools declares the tooling used to create analysis, database and logstore.
type Api interface {
	// Generate DB
	GenerateEncounterData(w http.ResponseWriter, r *http.Request)

	// Return DB
	GetHealerFrequency(w http.ResponseWriter, r *http.Request)
}

func Register(api Api) {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		// Generate and save data and analysis to DB from WarcraftLogs
		router.Post("/generate/encounter", api.GenerateEncounterData)

		// Return data from DB to consumer
		router.Get("/encounter/{encounterName}/healer-frequency", api.GetHealerFrequency)
	})

	fmt.Println("Registering router...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
	}
}
