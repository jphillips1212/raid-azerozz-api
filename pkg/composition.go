package pkg

import (
	"fmt"
	"sync"

	fs "github.com/jphillips1212/roztools-api/firestore"
	wl "github.com/jphillips1212/roztools-api/warcraftlogs"
)

type killAnalysis struct {
	Error     error
	StartTime int
	EndTime   int

	Report wl.Report
}

type healerAnalysis struct {
	Error error

	HealerDetails fs.HealerDetails
}

// GenerateHealerCompositions generates a list of healer compositions for the provided encounter id
// var encounterID - The ID of the encounter that healer compositions should be generated for
// var persist - Whether the reports should be continued to be run (persisted) after an already existing report
// is returned from the API
//
// Persisting the reports will generate a lot of requests and should only be used when generating reports for a fight for the first time
func GenerateHealerCompositions(encounterID int, persist bool) {
	wl := wl.New()
	fs := fs.New()

	reportChan := make(chan killAnalysis)
	healerChan := make(chan healerAnalysis)
	wg := sync.WaitGroup{}

	// Loop through five pages of reports
	for i := 0; i <= 5; i++ {
		query, err := wl.GenerateReportsForEncounter(encounterID, i)
		if err != nil {
			fmt.Printf("\nError querying warcraft logs, abandoning further queries\n %v", err)
			break
		}
		encounterName := string(query.WorldData.Encounter.Name)
		fmt.Printf("\n_____Starting analysis for encounter %s on page %d_____\n", encounterName, i)

		for _, encounter := range query.WorldData.Encounter.FightRankings.Rankings {

			// Check for if this fightID and report has already been analysed
			if fs.DoesReportExists(encounter.Report.Code, encounterName, encounter.Report.FightID) {
				fmt.Printf("Report %s for fight %d has already been analysed, skipping analysis for this fight\n", encounter.Report.Code, encounter.Report.FightID)
				if !persist {
					fmt.Printf("Persist is not set, abandoning rest of reports for boss %d\n", encounterID)
					break
				}
			} else {
				wg.Add(3)
				go generateReportForEncounter(wl, encounter, reportChan, &wg)
				go generateHealerComposition(wl, reportChan, healerChan, &wg)
				go saveHealerComposition(fs, encounterName, healerChan, &wg)
			}
		}
	}

	wg.Wait()
	close(reportChan)
	close(healerChan)
}

func generateReportForEncounter(client *wl.Client, encounter wl.Rankings, reportChan chan<- killAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get Fight Times for the encounter for use later
	fightTimes, err := client.GetFightTimes(encounter.Report.Code, encounter.Report.FightID)
	if err != nil {
		reportChan <- killAnalysis{
			Error: err,
			Report: wl.Report{
				Code:    encounter.Report.Code,
				FightID: encounter.Report.FightID,
			},
		}
		return
	}

	reportChan <- killAnalysis{
		Report: wl.Report{
			Code:    encounter.Report.Code,
			FightID: encounter.Report.FightID,
		},
		StartTime: fightTimes.StartTime,
		EndTime:   fightTimes.EndTime,
	}
}

func generateHealerComposition(client *wl.Client, reportChan <-chan killAnalysis, healerChan chan<- healerAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	report := <-reportChan
	if report.Error != nil {
		healerChan <- healerAnalysis{
			Error: fmt.Errorf("error returning fight times for report [%s] for fight ID [%d]: [%v]",
				report.Report.Code, report.Report.FightID, report.Error),
			HealerDetails: fs.HealerDetails{
				Report:  report.Report.Code,
				FightID: report.Report.FightID,
			},
		}
		return
	}

	comp, err := client.GetEncounterComposition(report.Report.Code, report.Report.FightID, report.StartTime, report.EndTime)
	if err != nil {
		healerChan <- healerAnalysis{
			Error: fmt.Errorf("error returning healer composition for report %s for fight id %d: [%v] - abandoning generating healer comp for this fight",
				report.Report.Code, report.Report.FightID, err),
			HealerDetails: fs.HealerDetails{
				Report:  report.Report.Code,
				FightID: report.Report.FightID,
			},
		}
		return
	}

	healers := []fs.Healer{}
	for _, actor := range comp {
		if actor.Role == "healer" {
			healers = append(healers, fs.Healer{
				Name:  actor.Name,
				Class: actor.Class,
				Spec:  actor.Spec,
			})
		}
	}

	healerChan <- healerAnalysis{
		HealerDetails: fs.HealerDetails{
			Report:    report.Report.Code,
			FightID:   report.Report.FightID,
			StartTime: report.StartTime,
			EndTime:   report.EndTime,
			Healers:   healers,
		},
	}

	fmt.Printf("Finished analysis for report %s\n", report.Report.Code)
}

func saveHealerComposition(client *fs.Client, encounterName string, ch chan healerAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	healerComp := <-ch

	if healerComp.Error != nil {
		fmt.Printf("error generating healer composition [%v]\n", healerComp.Error)
		return
	}

	err := client.SaveHealerComposition(encounterName, healerComp.HealerDetails)
	if err != nil {
		fmt.Printf("error saving healing composition to database [%v]\n", err)
		return
	}
}
