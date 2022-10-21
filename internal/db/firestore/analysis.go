package firestore

import (
	"fmt"
	"sync"

	model "github.com/jphillips1212/roztools-api/model"
	"google.golang.org/api/iterator"
)

func (c Client) SaveHealerFrequencies(encounterName, healerKey string, healerFrequency model.HealerFrequency, wg *sync.WaitGroup) {

	_, err := c.Client.Collection("Analysis").Doc(encounterName).Collection("HealerFrequency").Doc(healerKey).Set(*c.Ctx, healerFrequency)
	if err != nil {
		fmt.Printf("\nerror saving analysis for encounter %s and healer key %s\n", encounterName, healerKey)
	}

	wg.Done()
}

func (c Client) GetHealerFrequencies(encounterName string) ([]model.HealerFrequency, error) {
	var healerFrequencies []model.HealerFrequency
	iter := c.Client.Collection("Analysis").Doc(encounterName).Collection("HealerFrequency").Documents(*c.Ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []model.HealerFrequency{}, fmt.Errorf("error iterating through analysis documents for %s: %v", encounterName, err)
		}

		var healerFreq model.HealerFrequency
		if err := doc.DataTo(&healerFreq); err != nil {
			return []model.HealerFrequency{}, fmt.Errorf("error converting analysis healer comp from database to struct for %s: %v", encounterName, err)
		}

		healerFrequencies = append(healerFrequencies, healerFreq)
	}

	return healerFrequencies, nil
}
