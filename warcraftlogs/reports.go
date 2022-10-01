package warcraftlogs

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

type QueryReportForEncounterComposition struct {
	ReportData struct {
		Report struct {
			Table struct {
				Data data
			} `graphql:"table(fightIDs: [$fightId], startTime: $startTime, endTime: $endTime)" scalar:"true"`
		} `graphql:"report(code: $reportCode)"`
	}
}

type data struct {
	Composition []composition `json:"composition"`
}

type composition struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Type string `json:"type"`
	Spec []spec `json:"specs"`
}

type spec struct {
	Spec string `json:"spec"`
	Role string `json:"role"`
}

type EncounterActor struct {
	Name    string
	ActorID int
	Class   string
	Spec    string
	Role    string
}

// GetEncounterComposition returns all actors associated for the provided encounter
func (c Client) GetEncounterComposition(reportCode string, fightId int, startTime int, endTime int) ([]EncounterActor, error) {
	var query QueryReportForEncounterComposition
	variables := map[string]any{
		"reportCode": graphql.String(reportCode),
		"fightId":    graphql.Int(fightId),
		"startTime":  graphql.Float(startTime),
		"endTime":    graphql.Float(endTime),
	}

	err := c.Client.Query(context.Background(), &query, variables)
	if err != nil {
		return []EncounterActor{}, err
	}

	var actors []EncounterActor

	for _, actor := range query.ReportData.Report.Table.Data.Composition {
		actors = append(actors, EncounterActor{
			Name:    actor.Name,
			ActorID: actor.ID,
			Class:   actor.Type,
			Spec:    actor.Spec[0].Spec,
			Role:    actor.Spec[0].Role,
		})
	}

	return actors, nil
}

type QueryReportForFightTime struct {
	ReportData struct {
		Report struct {
			Fights []struct {
				StartTime graphql.Float
				EndTime   graphql.Float
			} `graphql:"fights(fightIDs: [$fightId])"`
		} `graphql:"report(code: $reportCode)"`
	}
}

type FightTimes struct {
	StartTime int
	EndTime   int
}

// GetFightTimes returns the startTime and endTime for the provided encounter within a provided report
func (c Client) GetFightTimes(reportCode string, fightId int) (FightTimes, error) {
	var query QueryReportForFightTime
	variables := map[string]any{
		"reportCode": graphql.String(reportCode),
		"fightId":    graphql.Int(fightId),
	}

	err := c.Client.Query(context.Background(), &query, variables)
	if err != nil {
		return FightTimes{}, err
	}

	// Assume first Fights struct contains kill as we use unique identifier of fightIDs
	return FightTimes{
		StartTime: int(query.ReportData.Report.Fights[0].StartTime),
		EndTime:   int(query.ReportData.Report.Fights[0].EndTime),
	}, nil
}
