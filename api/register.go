package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// RozTools declares the tooling used to create analysis, database and logstore.
type Api interface {
	// Analysis
	GetHealerFrequency(w http.ResponseWriter, r *http.Request)

	// Generate
	GenerateHealerComposition(w http.ResponseWriter, r *http.Request)
}

func Register(api Api) {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		router.Get("/encounter/{encounterName}/healer-frequency", api.GetHealerFrequency)

		router.Post("/generate/healers", api.GenerateHealerComposition)
	})

	fmt.Println("Registering router...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
	}
}
