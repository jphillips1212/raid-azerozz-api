package pkg

import (
	"fmt"

	wl "github.com/jphillips1212/roztools-api/warcraftlogs"
)

// GetHealerComposition returns the actors with role `healer` for the encounter
func GetHealerComposition(encounterID int) {
	client := wl.Client{
		Client: wl.New(),
	}

	reports, err := client.GetReportsForEncounter(2407)
	if err != nil {
		fmt.Printf("error getting kills for encounter %d: [%v]\n", encounterID, err)
	}

	for _, report := range reports {
		comp, err := client.GetEncounterComposition(report.Code, report.FightID, report.StartTime, report.EndTime)
		if err != nil {
			fmt.Printf("error retrieving composition for fight %d on report %s: [%v]\n", report.FightID, report.Code, err)
		}

		fmt.Println("_________________________________________________")
		fmt.Printf("report %s for encounter %d\n", report.Code, report.FightID)

		for _, actor := range comp {
			if actor.Role == "healer" {
				fmt.Printf("%s as %s %s\n", actor.Name, actor.Spec, actor.Class)
			}
		}

		fmt.Println("_________________________________________________")
	}
}
