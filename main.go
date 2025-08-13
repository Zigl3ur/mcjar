package main

import (
	"log"

	"github.com/Zigl3ur/mcli/cmd"
)

func main() {
	// disable date and time
	log.SetFlags(0)
	cmd.Execute()
}
