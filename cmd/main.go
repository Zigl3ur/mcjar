package main

import (
	"fmt"
	"os"

	"github.com/Zigl3ur/mc-jar-fetcher/handlers/flags"
	"github.com/Zigl3ur/mc-jar-fetcher/handlers/vanilla"
	"github.com/spf13/pflag"
)

func main() {
	flagsValues := flags.Init()
	flagsValues.Register()

	fmt.Fprintln(os.Stdout, "Using values:")
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Value.String() == "0" {
			return
		}
		fmt.Fprintf(os.Stdout, "- %s: %s\n", f.Name, f.Value)
	})

	switch flagsValues.Type {
	case flags.Vanilla:
		if err := vanilla.VanillaHandler(flagsValues.Version, flagsValues.Path); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid type, valid ones are [vanilla, paper, spigot, mohist, forge, fabric]")
	}
}
