package app

import (
	"fmt"

	flags "github.com/Zigl3ur/mc-jar-fetcher/app/handlers"
)

func Run() {
	flagsValues := flags.Init()
	flagsValues.Register()

	fmt.Println(flagsValues)
}
