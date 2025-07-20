package app

import (
	"fmt"
	"os"

	"github.com/Zigl3ur/mc-jar-fetcher/app/handlers"
	"github.com/Zigl3ur/mc-jar-fetcher/app/handlers/flags"
	"github.com/Zigl3ur/mc-jar-fetcher/app/utils"
	"github.com/spf13/pflag"
)

func Run() {
	flagsValues := flags.Init()
	flagsValues.Register()

	defaultValues := make(map[string]string)

	// TODO: prettier way
	if !pflag.Lookup("version").Changed {
		defaultValues["version"] = flagsValues.Version
	}
	if !pflag.Lookup("type").Changed {
		defaultValues["type"] = string(flagsValues.Type)
	}
	if !pflag.Lookup("dest").Changed {
		defaultValues["path"] = flagsValues.Path
	}

	if len(defaultValues) > 0 {
		fmt.Println("Using default values:")
		for key, value := range defaultValues {
			fmt.Printf(" - %s: %s\n", key, value)
		}
	}

	switch flagsValues.Type {
	case flags.Vanilla:
		if url, err := handlers.GetUrlVanilla(flagsValues.Version); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else if err := utils.WriteToFs(url, flagsValues.Path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	}
}
