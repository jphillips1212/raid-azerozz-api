package firestore

import (
	"fmt"

	model "github.com/jphillips1212/roztools-api/model"
	"google.golang.org/api/iterator"
)

func (c Client) SaveHealerComposition(encounterName string, healerDetails model.HealerDetails) error {

	_, _, err := c.Client.Collection("Reports").Doc(encounterName).Collection("HealerComp").Add(*c.Ctx, healerDetails)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) GetIfReportExists(reportCode, encounterName string, fightID int) bool {
	reports := c.Client.Collection("Reports").Doc(encounterName).Collection("HealerComp")
	query := reports.Where("Report", "==", reportCode).Where("FightID", "==", fightID)
	result := query.Documents(*c.Ctx)

	if _, next := result.Next(); next == iterator.Done {
		return false
	}

	return true
}

func (c Client) GetAllHealerCompositions(encounterName string) ([]model.HealerDetails, error) {
	var healerDetails []model.HealerDetails
	iter := c.Client.Collection("Reports").Doc(encounterName).Collection("HealerComp").Documents(*c.Ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.HealerDetails{}, fmt.Errorf("error iterating through documents for %s: %v", encounterName, err)
		}

		var healerComp model.HealerDetails
		if err := doc.DataTo(&healerComp); err != nil {
			return []model.HealerDetails{}, fmt.Errorf("error converting healer comp from database to struct for %s: %v", encounterName, err)
		}

		healerDetails = append(healerDetails, healerComp)
	}

	return healerDetails, nil
}
