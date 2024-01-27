package analysis

import (
	"fmt"
	"sync"

	wl "github.com/jphillips1212/roztools-api/internal/logs/warcraftlogs"
	"github.com/jphillips1212/roztools-api/model"
)

// ScrapeAndSaveBoss scrapes all the data for the provided boss encounterID from warcraftLogs and saves it into
// the DB to be used for analysis
// encounterID corresponds to the encounterID of the raid boss
// persist represents if the scraping should continue to run after it hits a report it's saved already before
func (analysis Analysis) ScrapeAndSaveBoss(encounterID int, persist bool) {

	encounterDataChan := make(chan model.EncounterDataChannel)
	wg := sync.WaitGroup{}

	// n tracks the number of reports that have been analysed, defaulted to stopping after 100
	var n int
	// Loop through up to 50 pages of reports
records:
	for page := 1; page <= 50; page++ {
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
				wg.Add(2)
				go analysis.getEncounterData(encounter, encounterDataChan, &wg)
				go analysis.saveEncounterData(encounterName, encounterDataChan, &wg)
				fmt.Printf("\nAnalysing report %s: Report number %d of 100\n", encounter.Report.Code, n)
			}
		}
	}

	wg.Wait()
	close(encounterDataChan)
}

func (analysis Analysis) getEncounterData(encounter wl.Rankings, encounterDataChan chan<- model.EncounterDataChannel, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get Fight Times for the encounter for use later
	fightTimes, err := analysis.Logs.GetFightTimes(encounter.Report.Code, encounter.Report.FightID)
	if err != nil {
		encounterDataChan <- model.EncounterDataChannel{
			Error: err,
			EncounterData: model.EncounterData{
				Code:    encounter.Report.Code,
				FightID: encounter.Report.FightID,
			},
		}
		return
	}

	comp, err := analysis.Logs.GetEncounterComposition(encounter.Report.Code, encounter.Report.FightID, fightTimes.StartTime, fightTimes.EndTime)

	healers, tanks := []model.Player{}, []model.Player{}
	for _, actor := range comp {
		if actor.Role == "healer" {
			healers = append(healers, model.Player{
				Name:  actor.Name,
				Class: actor.Class,
				Spec:  actor.Spec,
			})
		}
		if actor.Role == "tank" {
			tanks = append(tanks, model.Player{
				Name:  actor.Name,
				Class: actor.Class,
				Spec:  actor.Spec,
			})
		}
	}

	encounterDataChan <- model.EncounterDataChannel{
		EncounterData: model.EncounterData{
			StartTime:    fightTimes.StartTime,
			EndTime:      fightTimes.EndTime,
			Code:         encounter.Report.Code,
			FightID:      encounter.Report.FightID,
			GuildID:      encounter.Guild.ID,
			GuildName:    encounter.Guild.Name,
			GuildFaction: encounter.Guild.Faction,
			Healers:      healers,
			Tanks:        tanks,
		},
	}
}

func (analysis Analysis) saveEncounterData(encounterName string, encounterDataChan <-chan model.EncounterDataChannel, wg *sync.WaitGroup) {
	defer wg.Done()
	encounterDataChannel := <-encounterDataChan

	if encounterDataChannel.Error != nil {
		fmt.Printf("error generating data for encounter [%v]\n", encounterDataChannel.Error)
		return
	}

	err := analysis.Database.SaveEncounterData(encounterName, encounterDataChannel.EncounterData)
	if err != nil {
		fmt.Printf("error saving encounter data to database [%v]\n", err)
		return
	}
}
