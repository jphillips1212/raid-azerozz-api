package firestore

import (
	"fmt"

	model "github.com/jphillips1212/roztools-api/model"
	"google.golang.org/api/iterator"
)

// SaveEncounterData saves the data for this particular encounter to the database
func (c Client) SaveEncounterData(encounterName string, encounterData model.EncounterData) error {

	_, _, err := c.Client.Collection("Reports").Doc(encounterName).Collection("EncounterData").Add(*c.Ctx, model.EncounterData{
		StartTime:    encounterData.StartTime,
		EndTime:      encounterData.EndTime,
		Code:         encounterData.Code,
		FightID:      encounterData.FightID,
		GuildID:      encounterData.GuildID,
		GuildName:    encounterData.GuildName,
		GuildFaction: encounterData.GuildFaction,
		Healers:      encounterData.Healers,
		Tanks:        encounterData.Tanks,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetIfReportExists returns if the reportCode already exists in the database
func (c Client) GetIfReportExists(reportCode, encounterName string, fightID int) bool {
	reports := c.Client.Collection("Reports").Doc(encounterName).Collection("EncounterData")
	query := reports.Where("Code", "==", reportCode).Where("FightID", "==", fightID)
	result := query.Documents(*c.Ctx)

	if _, next := result.Next(); next == iterator.Done {
		return false
	}

	return true
}

// GetIfMaxReportsForEncounter returns if the maximum number of reports for this encounter has already been
// queried on WarcraftLogs and added to the database
// Currently set to 50 reports on 50 pages (2500) set by WarcraftLogs
func (c Client) GetIfMaxReportsForEncounter(encounterID int) bool {
	return false
}

// GetAllEncounters returns all of the encounter reports for the provided boss
func (c Client) GetAllEncounters(encounterName string) ([]model.EncounterData, error) {
	var encountersData []model.EncounterData
	iter := c.Client.Collection("Reports").Doc(encounterName).Collection("EncounterData").Documents(*c.Ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.EncounterData{}, fmt.Errorf("error iterating through documents for %s: %v", encounterName, err)
		}

		var encounter model.EncounterData
		if err := doc.DataTo(&encounter); err != nil {
			return []model.EncounterData{}, fmt.Errorf("error converting healer comp from database to struct for %s: %v", encounterName, err)
		}

		encountersData = append(encountersData, encounter)
	}

	return encountersData, nil
}
