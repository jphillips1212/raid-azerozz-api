package service

type Service interface {
	// Analysis
	AnalyseAndSaveBoss(encounterID int)

	// Logscrape
	ScrapeAndSaveBoss(encounterID int, persist bool)
}
