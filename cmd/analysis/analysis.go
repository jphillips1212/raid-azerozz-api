package main

import (
	"fmt"
	"time"

	"github.com/jphillips1212/roztools-api/analysis"
)

func main() {
	start := time.Now()

	analysis.AnalyseHealerComp("Stone Legion Generals")

	elapsed := time.Since(start)
	fmt.Printf("Total time taken for analysis: %s\n", elapsed)
}
