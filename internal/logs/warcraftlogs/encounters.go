package warcraftlogs

import (
	"context"

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
			} `graphql:"fightRankings(page: $page)" scalar:"true"`
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
func (c Client) GetReportsForEncounter(fightId, page int) (*QueryFightByID, error) {
	var query QueryFightByID
	variables := map[string]any{
		"id":   graphql.Int(fightId),
		"page": graphql.Int(page),
	}

	err := c.Client.Query(context.Background(), &query, variables)
	if err != nil {
		return &QueryFightByID{}, err
	}

	return &query, nil
}
