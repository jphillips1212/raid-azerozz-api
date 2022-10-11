package main

import (
	"fmt"
	"time"

	"github.com/jphillips1212/roztools-api/pkg"
)

func main() {
	start := time.Now()
	// Hardcoded to Sire Denathrius
	pkg.GenerateHealerCompositions(2417, false)

	elapsed := time.Since(start)
	fmt.Printf("Total time to run: %s\n", elapsed)
}
