package firestore

import (
	"fmt"

	"google.golang.org/api/iterator"
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

func (c Client) SaveHealerComposition(encounterName string, healerDetails HealerDetails) error {

	_, _, err := c.Client.Collection("Reports").Doc(encounterName).Collection("HealerComp").Add(*c.Ctx, healerDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Report %s written to collection\n", healerDetails.Report)

	return nil
}

func (c Client) DoesReportExists(reportCode, encounterName string, fightID int) bool {
	reports := c.Client.Collection("Reports").Doc(encounterName).Collection("HealerComp")
	query := reports.Where("Report", "==", reportCode).Where("FightID", "==", fightID)
	result := query.Documents(*c.Ctx)

	if _, next := result.Next(); next == iterator.Done {
		return false
	}

	return true
}

func (c Client) ReturnAllHealerComps(encounterName string) ([]HealerDetails, error) {
	var healerDetails []HealerDetails
	iter := c.Client.Collection("Reports").Doc(encounterName).Collection("HealerComp").Documents(*c.Ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []HealerDetails{}, fmt.Errorf("error iterating through documents for %s: %v", encounterName, err)
		}

		var healerComp HealerDetails
		if err := doc.DataTo(&healerComp); err != nil {
			return []HealerDetails{}, fmt.Errorf("error converting healer comp from database to struct for %s: %v", encounterName, err)
		}

		healerDetails = append(healerDetails, healerComp)
	}

	return healerDetails, nil
}
