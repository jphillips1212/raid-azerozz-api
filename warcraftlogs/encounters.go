package warcraftlogs

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

type QueryFightByID struct {
	WorldData struct {
		Encounter struct {
			Name          graphql.String
			Zone          zone
			FightRankings struct {
				Page     int
				Count    int
				Rankings []rankings
			} `scalar:"true"`
		} `graphql:"encounter(id: $id)"`
	}
}

type zone struct {
	Name graphql.String
}

type rankings struct {
	Report report `json:"report"`
}

type report struct {
	Code    string `json:"code"`
	FightID int    `json:"fightID"`
}

type KillDetails struct {
	Code      string
	FightID   int
	StartTime int
	EndTime   int
}

// GetReportsForEncounter returns a slice of killdetails for the provided encounter
func (c Client) GenerateReportsForEncounter(fightId int, encounterReports chan<- KillDetails, n chan<- int) {
	var query QueryFightByID
	variables := map[string]any{
		"id": graphql.Int(fightId),
	}

	err := c.Client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Printf("Error querying warcraftLogs for encounter: [%d]\n", fightId)
	}

	// Hardcode to wait on ten reports
	fmt.Printf("Waitgroup count is %d\n", query.WorldData.Encounter.FightRankings.Count)
	n <- 10

	// Hardcoded to only loop over first ten reports to not hit rate limit
	//for _, encounter := range query.WorldData.Encounter.FightRankings.Rankings {
	for i := 0; i < 10; i++ {
		go func(i int, ch chan<- KillDetails) {
			encounter := query.WorldData.Encounter.FightRankings.Rankings[i] // Remove when putting back for range loop

			fightTimes, err := c.GetFightTimes(encounter.Report.Code, encounter.Report.FightID)
			if err != nil {
				fmt.Printf("error returning fight times for report [%s] for fight ID [%d]: [%v]\n", encounter.Report.Code, fightId, err)
				ch <- KillDetails{}
				return
			}

			ch <- KillDetails{
				Code:      encounter.Report.Code,
				FightID:   encounter.Report.FightID,
				StartTime: fightTimes.StartTime,
				EndTime:   fightTimes.EndTime,
			}

		}(i, encounterReports)
	}
}
