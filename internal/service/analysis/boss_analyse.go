package analysis

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/jphillips1212/roztools-api/model"
)

// AnalyseAndSaveBoss looks through the database for the encounterID and generates analysis for the boss, saving it to another collection
func (analysis Analysis) AnalyseAndSaveBoss(encounterID int) {
	encountersData, err := analysis.Database.GetAllEncounters(EncounterIDs[encounterID])
	if err != nil {
		fmt.Printf("\nThere was an error retrieving the logs for [%s] for analysis\n", EncounterIDs[encounterID])
	}

	var totalAllianceKills, totalHordeKills, totalRecordedKills int
	for _, encounter := range encountersData {
		if encounter.GuildFaction == 0 {
			totalAllianceKills += 1
		} else if encounter.GuildFaction == 1 {
			totalHordeKills += 1
		}

		totalRecordedKills += 1
	}

	fmt.Printf("\n Total Alliance Kills: [%d]", totalAllianceKills)
	fmt.Printf("\n Total Horde Kills: [%d]", totalHordeKills)
	fmt.Printf("\n Total Recorded Kills: [%d]", totalRecordedKills)

	fmt.Printf("\n Calculating Healer Frequency\n")
	healerFrequency := findFrequencyOfRole("healer", encountersData)
	fmt.Println(healerFrequency[0])
	tankFrequency := findFrequencyOfRole("tank", encountersData)
	fmt.Println(tankFrequency[0])
	guildFrequency := findFrequencyOfGuilds(encountersData)
	fmt.Printf("\n Total Unique Guilds: [%d]", len(guildFrequency))
}

func findFrequencyOfGuilds(encounterData []model.EncounterData) []model.GuildFrequency {
	frequency := make(map[int]int)

	// Calculate frequency of guilds by looping over guild_ids
	for _, encounter := range encounterData {
		guildID := encounter.GuildID
		if frequency[guildID] == 0 {
			frequency[guildID] = 1
		} else {
			frequency[guildID]++
		}
	}

	var guildFrequencies = []model.GuildFrequency{}

	// Convert frequencies from map[int]int to guild frequency struct
	for guildID, count := range frequency {
		guildFrequency := model.GuildFrequency{
			Frequency: count,
			GuildID:   guildID,
		}

		guildFrequencies = append(guildFrequencies, guildFrequency)
	}

	// Sort struct to be sorted by frequency
	sort.Slice(guildFrequencies, func(i, j int) bool {
		return guildFrequencies[i].Frequency > guildFrequencies[j].Frequency
	})

	return guildFrequencies
}

func findFrequencyOfRole(role string, encountersData []model.EncounterData) []model.RoleFrequency {
	frequency := make(map[string]int)

	// Calculate frequency by generating keys for each composition
	for _, encounter := range encountersData {
		var composition []model.Player
		if role == "healer" {
			composition = encounter.Healers
		} else if role == "tank" {
			composition = encounter.Tanks
		}

		key := generateRoleKey(composition)
		if frequency[key] == 0 {
			frequency[key] = 1
		} else {
			frequency[key]++
		}
	}

	var roleFrequencies []model.RoleFrequency

	// Convert frequency with compKey into struct
	for compKey, count := range frequency {
		comp := strings.Split(compKey, ":")

		compFrequency := model.RoleFrequency{
			Frequency:      count,
			CompositionKey: compKey,
			Composition:    comp,
		}

		roleFrequencies = append(roleFrequencies, compFrequency)
	}

	// Sort struct to be sorted by frequency
	sort.Slice(roleFrequencies, func(i, j int) bool {
		return roleFrequencies[i].Frequency > roleFrequencies[j].Frequency
	})

	return roleFrequencies
}

// Generates a string that acts as a key for that specific composition of players in that role (tanks/healers/dps)
func generateRoleKey(composition []model.Player) string {
	// Generate string slice of all players for this role
	var compString []string
	for _, h := range composition {
		compString = append(compString, fmt.Sprintf("%s %s", h.Spec, h.Class))
	}

	// Sort healers alphabetically
	sort.Slice(compString, func(i, j int) bool {
		return compString[i] < compString[j]
	})

	// Convert slice into string to be used as key
	var compKey string
	for _, s := range compString {
		compKey = fmt.Sprintf("%s:%s", compKey, s)
	}

	//Remove first character ":"
	_, i := utf8.DecodeLastRuneInString(compKey)
	return compKey[i:]
}
