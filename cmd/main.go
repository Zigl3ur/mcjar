package main

import (
	"log"

	"github.com/Zigl3ur/mc-jar-fetcher/handlers/flags"
)

func main() {
	// disable timestamp
	log.SetFlags(0)

	// flags logic
	flagsValues := flags.Init()
	flagsValues.Validate()
	flagsValues.Execute()
}
