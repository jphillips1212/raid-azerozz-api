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
			Zone          Zone
			FightRankings struct {
				Page     int
				Count    int
				Rankings []Rankings
			} `scalar:"true"`
		} `graphql:"encounter(id: $id)"`
	}
}

type Zone struct {
	Name graphql.String
}

type Rankings struct {
	Report Report `json:"report"`
}

type Report struct {
	Code    string `json:"code"`
	FightID int    `json:"fightID"`
}

// GetReportsForEncounter returns a slice of killdetails for the provided encounter
func (c Client) GenerateReportsForEncounter(fightId int) *QueryFightByID {
	var query QueryFightByID
	variables := map[string]any{
		"id": graphql.Int(fightId),
	}

	err := c.Client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Printf("Error querying warcraftLogs for encounter: [%d]\n", fightId)
	}

	return &query
}
