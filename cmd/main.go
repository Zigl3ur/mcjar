package main

import "github.com/Zigl3ur/mc-jar-fetcher/handlers/flags"

func main() {
	flagsValues := flags.Init()
	flagsValues.Validate()
	flagsValues.Execute()
}
