package main

import (
	"fmt"
	"os"

	"github.com/Zigl3ur/mc-jar-fetcher/handlers"
	"github.com/Zigl3ur/mc-jar-fetcher/handlers/flags"
	"github.com/Zigl3ur/mc-jar-fetcher/utils"
	"github.com/spf13/pflag"
)

func main() {
	flagsValues := flags.Init()
	flagsValues.Register()

	defaultFlags := []string{
		"version",
		"type",
		"dest",
	}

	defaultFlagsValues := []string{}

	for _, flag := range defaultFlags {
		if !pflag.Lookup(flag).Changed {
			defaultFlagsValues = append(defaultFlagsValues, fmt.Sprintf("- %s: %s", flag, pflag.Lookup(flag).Value))
		}
	}

	if len(defaultFlagsValues) > 0 {
		fmt.Println("Using default values:")
		for _, value := range defaultFlagsValues {
			fmt.Println(value)
		}
	}

	switch flagsValues.Type {
	case flags.Vanilla:
		if url, err := handlers.GetUrlVanilla(flagsValues.Version); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if err := utils.WriteToFs(url, flagsValues.Path); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid type, valid ones are [vanilla, paper, spigot, mohist, forge, fabric]")
	}
}
