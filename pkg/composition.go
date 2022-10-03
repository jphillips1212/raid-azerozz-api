package pkg

import (
	"fmt"
	"sync"

	wl "github.com/jphillips1212/roztools-api/warcraftlogs"
)

// GetHealerComposition returns the actors with role `healer` for the encounter
func GetHealerComposition(encounterID int) {
	client := wl.Client{
		Client: wl.New(),
	}

	encounterReports := make(chan wl.KillDetails)
	numEncountersChan := make(chan int)

	go client.GenerateReportsForEncounter(2407, encounterReports, numEncountersChan)

	// Wait for GenerateReportsForEncounter to respond with how many reports have been returned
	numEncounters := <-numEncountersChan
	wg := sync.WaitGroup{}
	wg.Add(numEncounters)

	for i := 1; i <= numEncounters; i++ {
		go func() {
			report := <-encounterReports
			comp, _ := client.GetEncounterComposition(report.Code, report.FightID, report.StartTime, report.EndTime)
			fmt.Println("_________________________________________________")
			fmt.Printf("report %s for encounter %d\n", report.Code, report.FightID)

			for _, actor := range comp {
				if actor.Role == "healer" {
					fmt.Printf("%s as %s %s\n", actor.Name, actor.Spec, actor.Class)
				}
			}

			fmt.Println("_________________________________________________")

			wg.Done()
		}()
	}
	wg.Wait()

	close(encounterReports)
}
