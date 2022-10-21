package service

type Service interface {
	// Analysis
	SaveHealerAnalysis(encounterName string)

	// Logscrape
	SaveHealerLogs(encounterName string, persist bool) error
}
