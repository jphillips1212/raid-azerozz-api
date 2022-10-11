package analysis

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	firestore "github.com/jphillips1212/roztools-api/firestore"
)

func AnalyseHealerComp(encounter string) {
	fs := firestore.New()

	healerComps, err := fs.ReturnAllHealerComps(encounter)
	if err != nil {
		fmt.Printf("Error returning healer comp [%v]\n", err)
		return
	}

	healerFrequencies := findFrequencyOfHealers(healerComps)

	wg := sync.WaitGroup{}
	for _, healerFrequency := range healerFrequencies {
		wg.Add(1)
		go fs.SaveHealerAnalysis(encounter, healerFrequency.HealerKey, healerFrequency, &wg)
	}

	wg.Wait()

	fmt.Printf("\nHealer frequencies have been saved for %s, total of %d encounters\n", encounter, len(healerComps))
}

func findFrequencyOfHealers(healerComps []firestore.HealerDetails) []firestore.HealerFrequency {
	frequency := make(map[string]int)

	// Calculate frequency using healerKey
	for _, healers := range healerComps {
		key := generateHealerKey(healers.Healers)
		if frequency[key] == 0 {
			frequency[key] = 1
		} else {
			frequency[key]++
		}
	}

	var healerFrequencies []firestore.HealerFrequency

	// Convert frequency with healerKey into struct
	for healerKey, count := range frequency {
		healers := strings.Split(healerKey, ":")

		healerFrequency := firestore.HealerFrequency{
			Frequency: count,
			HealerKey: healerKey,
			Healers:   healers,
		}

		healerFrequencies = append(healerFrequencies, healerFrequency)
	}

	// Sort struct to be sorted by frequency
	sort.Slice(healerFrequencies, func(i, j int) bool {
		return healerFrequencies[i].Frequency > healerFrequencies[j].Frequency
	})

	return healerFrequencies
}

// Generates a string that acts as a key for that specific composition of healers
func generateHealerKey(healerComp []firestore.Healer) string {
	// Generate string slice of all healers
	var healers []string
	for _, h := range healerComp {
		healers = append(healers, fmt.Sprintf("%s %s", h.Spec, h.Class))
	}

	// Sort healers alphabetically
	sort.Slice(healers, func(i, j int) bool {
		return healers[i] < healers[j]
	})

	// Convert slice into string to be used as key
	var healerKey string
	for _, s := range healers {
		healerKey = fmt.Sprintf("%s:%s", healerKey, s)
	}

	//Remove first character ":"
	_, i := utf8.DecodeLastRuneInString(healerKey)
	return healerKey[i:]
}
