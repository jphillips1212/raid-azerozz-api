package firestore

import (
	"fmt"
	"sync"
)

type HealerFrequency struct {
	Frequency int
	HealerKey string
	Healers   []string
}

func (c Client) SaveHealerAnalysis(encounterName, healerKey string, healerFrequency HealerFrequency, wg *sync.WaitGroup) {

	_, err := c.Client.Collection("Analysis").Doc(encounterName).Collection("HealerFrequency").Doc(healerKey).Set(*c.Ctx, healerFrequency)
	if err != nil {
		fmt.Printf("\nerror saving analysis for encounter %s and healer key %s\n", encounterName, healerKey)
	}

	wg.Done()
}
