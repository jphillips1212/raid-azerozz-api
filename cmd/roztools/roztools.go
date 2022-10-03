package main

import (
	"fmt"
	"time"

	"github.com/jphillips1212/roztools-api/pkg"
)

func main() {
	start := time.Now()
	// Hardcoded to Sire Denathrius
	pkg.GetHealerComposition(2407)

	elapsed := time.Since(start)
	fmt.Printf("Total time to run: %s", elapsed)
}
