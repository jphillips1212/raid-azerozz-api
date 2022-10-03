package pkg

import (
	"fmt"
	"sync"

	wl "github.com/jphillips1212/roztools-api/warcraftlogs"
)

type HealerDetails struct {
	FightID   int
	Report    string
	StartTime int
	EndTime   int
	Healers   []Healer
}

type Healer struct {
	Name  string
	Class string
	Spec  string
}

type KillAnalysis struct {
	Error error

	Code      string
	FightID   int
	StartTime int
	EndTime   int
}

// GenerateHealerCompositions generates a list of healer compositions for the provided encounter id
// var encounterID - The ID of the encounter that healer compositions should be generated for
// var persist - Whether the reports should be continued to be run (persisted) after an already existing report
// is returned from the API
//
// Persisting the reports will generate a lot of requests and should only be used when generating reports for a fight for the first time
func GenerateHealerCompositions(encounterID int, persist bool) {
	client := wl.Client{
		Client: wl.New(),
	}

	query := *client.GenerateReportsForEncounter(encounterID)
	dataChan := make(chan KillAnalysis)
	wg := sync.WaitGroup{}

	// Hardcoded to only loop over first ten reports to not hit rate limit
	//for _, encounter := range query.WorldData.Encounter.FightRankings.Rankings {
	for i := 0; i < 10; i++ {
		encounter := query.WorldData.Encounter.FightRankings.Rankings[i]

		// Check for if this fightID and report has already been analysed
		if (encounter.Report.Code == "27PJCZVx6pb3qmDQ" && encounter.Report.FightID == 33) ||
			(encounter.Report.Code == "XjJx2MTqaVCNwYbt" && encounter.Report.FightID == 52) { // Hardcode to mimic report being found
			fmt.Printf("Report %s for fight %d has already been analysed, skipping analysis for this fight\n", encounter.Report.Code, encounter.Report.FightID)
			if !persist {
				fmt.Printf("Persist is not set, abandoning rest of reports for boss %d\n", encounterID)
				break
			}
		} else {
			wg.Add(1)
			go generateReportForEncounter(&client, encounter, dataChan)
			go listenForEncounterReport(&client, dataChan, &wg)
		}
	}

	wg.Wait()
	close(dataChan)
}

func generateReportForEncounter(client *wl.Client, encounter wl.Rankings, dataChan chan<- KillAnalysis) {
	// Get Fight Times for the encounter for use later
	fightTimes, err := client.GetFightTimes(encounter.Report.Code, encounter.Report.FightID)
	if err != nil {
		dataChan <- KillAnalysis{
			Error:   err,
			Code:    encounter.Report.Code,
			FightID: encounter.Report.FightID,
		}
		return
	}

	dataChan <- KillAnalysis{
		Code:      encounter.Report.Code,
		FightID:   encounter.Report.FightID,
		StartTime: fightTimes.StartTime,
		EndTime:   fightTimes.EndTime,
	}
}

func listenForEncounterReport(client *wl.Client, ch chan KillAnalysis, wg *sync.WaitGroup) {
	report := <-ch
	if report.Error != nil {
		fmt.Printf("error returning fight times for report [%s] for fight ID [%d]: [%v]\n", report.Code, report.FightID, report.Error)
		wg.Done()
		return
	}

	comp, err := client.GetEncounterComposition(report.Code, report.FightID, report.StartTime, report.EndTime)
	if err != nil {
		fmt.Printf("error returning healer composition for report %s for fight id %d: [%v] - abandoning generating healer comp for this fight\n", report.Code, report.FightID, err)
		wg.Done()
		return
	}

	healerDetails := HealerDetails{
		FightID:   report.FightID,
		Report:    report.Code,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
	}

	for _, actor := range comp {
		if actor.Role == "healer" {
			healerDetails.Healers = append(healerDetails.Healers, Healer{
				Name:  actor.Name,
				Class: actor.Class,
				Spec:  actor.Spec,
			})
		}
	}

	fmt.Println(healerDetails)
	wg.Done()
}
