package analysis

import (
	"fmt"
	"sync"

	wl "github.com/jphillips1212/roztools-api/internal/logs/warcraftlogs"
	model "github.com/jphillips1212/roztools-api/model"
)

type killAnalysis struct {
	Error     error
	StartTime int
	EndTime   int
	Code      string
	FightID   int
}

type healerAnalysis struct {
	Error error

	HealerDetails model.HealerDetails
}

// GenerateHealerCompositions generates a list of healer compositions for the provided encounter id
// var encounterID - The ID of the encounter that healer compositions should be generated for
// var persist - Whether the reports should be continued to be run (persisted) after an already existing report
// var n - The number of logs to analyse (only works if persist is set to true other will cancel)
// is returned from the API
//
// Persisting the reports will generate a lot of requests and should only be used when generating reports for a fight for the first time
func (analysis Analysis) SaveHealerLogs(encounterName string, persist bool) error {

	encounterID := EncounterIDs[encounterName]
	reportChan := make(chan killAnalysis)
	healerChan := make(chan healerAnalysis)
	wg := sync.WaitGroup{}

	// n tracks the number of reports that have been analysed, defaulted to stopping after 100
	var n int
	// Loop through up to 20 pages of reports
records:
	for page := 1; page <= 20; page++ {
		query, err := analysis.Logs.GetReportsForEncounter(encounterID, page)
		if err != nil {
			fmt.Printf("\nerror querying warcraft logs, abandoning further queries [%v]\n", err)
			break
		}
		encounterName := string(query.WorldData.Encounter.Name)
		fmt.Printf("\n_____Starting analysis for encounter %s on page %d_____\n", encounterName, page)

		for _, encounter := range query.WorldData.Encounter.FightRankings.Rankings {
			if n >= 100 {
				break records
			}

			// Check for if this fightID and report has already been analysed
			if analysis.Database.GetIfReportExists(encounter.Report.Code, encounterName, encounter.Report.FightID) {
				if !persist {
					fmt.Printf("Persist is not set, abandoning rest of reports for boss %d\n", encounterID)
					break records
				}
				break
			} else {
				n++
				wg.Add(3)
				go analysis.generateReportForEncounter(encounter, reportChan, &wg)
				go analysis.generateHealerComposition(reportChan, healerChan, &wg)
				go analysis.saveHealerComposition(encounterName, healerChan, &wg)
				fmt.Printf("\nAnalysing report %s: Report number %d of 100\n", encounter.Report.Code, n)
			}
		}
	}

	wg.Wait()
	close(reportChan)
	close(healerChan)

	return nil
}

func (analysis Analysis) generateReportForEncounter(encounter wl.Rankings, reportChan chan<- killAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get Fight Times for the encounter for use later
	fightTimes, err := analysis.Logs.GetFightTimes(encounter.Report.Code, encounter.Report.FightID)
	if err != nil {
		reportChan <- killAnalysis{
			Error:   err,
			Code:    encounter.Report.Code,
			FightID: encounter.Report.FightID,
		}
		return
	}

	reportChan <- killAnalysis{
		Code:      encounter.Report.Code,
		FightID:   encounter.Report.FightID,
		StartTime: fightTimes.StartTime,
		EndTime:   fightTimes.EndTime,
	}
}

func (analysis Analysis) generateHealerComposition(reportChan <-chan killAnalysis, healerChan chan<- healerAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	report := <-reportChan
	if report.Error != nil {
		healerChan <- healerAnalysis{
			Error: fmt.Errorf("error returning fight times for report [%s] for fight ID [%d]: [%v]",
				report.Code, report.FightID, report.Error),
			HealerDetails: model.HealerDetails{
				Report:  report.Code,
				FightID: report.FightID,
			},
		}
		return
	}

	comp, err := analysis.Logs.GetEncounterComposition(report.Code, report.FightID, report.StartTime, report.EndTime)
	if err != nil {
		healerChan <- healerAnalysis{
			Error: fmt.Errorf("error returning healer composition for report %s for fight id %d: [%v] - abandoning generating healer comp for this fight",
				report.Code, report.FightID, err),
			HealerDetails: model.HealerDetails{
				Report:  report.Code,
				FightID: report.FightID,
			},
		}
		return
	}

	healers := []model.Healer{}
	for _, actor := range comp {
		if actor.Role == "healer" {
			healers = append(healers, model.Healer{
				Name:  actor.Name,
				Class: actor.Class,
				Spec:  actor.Spec,
			})
		}
	}

	healerChan <- healerAnalysis{
		HealerDetails: model.HealerDetails{
			Report:    report.Code,
			FightID:   report.FightID,
			StartTime: report.StartTime,
			EndTime:   report.EndTime,
			Healers:   healers,
		},
	}
}

func (analysis Analysis) saveHealerComposition(encounterName string, ch chan healerAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	healerComp := <-ch

	if healerComp.Error != nil {
		fmt.Printf("error generating healer composition [%v]\n", healerComp.Error)
		return
	}

	err := analysis.Database.SaveHealerComposition(encounterName, healerComp.HealerDetails)
	if err != nil {
		fmt.Printf("error saving healing composition to database [%v]\n", err)
		return
	}
}
