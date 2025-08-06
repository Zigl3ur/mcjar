package flags

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Zigl3ur/mc-jar-fetcher/handlers/paper"
	"github.com/Zigl3ur/mc-jar-fetcher/handlers/vanilla"
	"github.com/spf13/pflag"
)

const invalidServerType string = "Invalid server type, valid ones are [vanilla, paper, spigot, mohist, forge, fabric]"

type flags struct {
	list       string // list version for specified server type
	version    string // minecraft version
	serverType string // like forge, mohist, paper, etc
	build      string // build version like for fabric and forge
	path       string // the name of the file outputted from the download
}

func Init() *flags {
	flagsVar := &flags{}

	// list available versions / builds for a server type
	pflag.StringVarP(&flagsVar.list, "list", "l", "", "list available versions for the specified server type")

	// game version
	pflag.StringVarP(&flagsVar.version, "version", "v", "1.21", "the server version")

	// server type
	pflag.StringVarP((*string)(&flagsVar.serverType), "type", "t", "vanilla", "the server type")

	// build
	pflag.StringVarP(&flagsVar.build, "build", "b", "", "the build version")

	// file destination
	pflag.StringVarP(&flagsVar.path, "dest", "d", "server.jar", "the destination for the downloaded file")

	pflag.CommandLine.SortFlags = false

	// custom help msg
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -h, --help\t  show this help message\n")
		pflag.VisitAll(func(f *pflag.Flag) {
			fmt.Fprintf(os.Stderr, "  -%s, --%s\t  %s\n", f.Shorthand, f.Name, f.Usage)
		})
	}

	return flagsVar
}

// validate data for provided flags without making any api calls
func (f *flags) Validate() {
	pflag.Parse()

	validServerType := []string{"vanilla", "forge", "mohist", "paper", "fabric", "spigot", "purpur"}

	if pflag.Lookup("type").Changed && !slices.Contains(validServerType, f.serverType) {
		log.Fatal(invalidServerType)
	}

	if pflag.Lookup("list").Changed && !slices.Contains(validServerType, f.list) {
		log.Fatal(invalidServerType)
	}

}

// execute functions relative to flags data
func (f *flags) Execute() {
	if pflag.Lookup("list").Changed {
		switch strings.ToLower(f.list) {
		case "vanilla":
			vlist, err := vanilla.GetVersionsList()
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range vlist.Versions {
				fmt.Printf("- %s\n", v.Id)
			}
		case "paper":
			vlist, err := paper.GetVersionsList()
			if err != nil {
				log.Fatal(err)
			}

			if !pflag.Lookup("version").Changed {
				for _, v := range vlist.Versions {
					fmt.Printf("- %s:\n", v.Version.Id)
					fmt.Println("  - builds:")
					for _, b := range v.Builds {
						fmt.Printf("\t- %d\n", b)
					}
				}
			} else {
				found := false
				for _, v := range vlist.Versions {
					if v.Version.Id == f.version {
						fmt.Printf("- %s:\n", f.version)
						fmt.Println("  - builds:")
						for _, b := range v.Builds {
							fmt.Printf("\t- %d\n", b)
						}
						found = true
						break
					}
				}
				if !found {
					log.Fatal("paper doesn't support this version")
				}
			}

		default:
			log.Fatal(invalidServerType)
		}
		os.Exit(0)
	}

	fmt.Println("Using Values:")
	fmt.Printf("- type: %s\n", f.serverType)
	fmt.Printf("- version: %s\n", f.version)
	if pflag.Lookup("build").Changed {
		fmt.Printf("- build: %s\n", f.build)
	}
	fmt.Printf("- dest: %s\n", f.path)

	switch f.serverType {
	case "vanilla":
		if err := vanilla.Handler(f.version, f.path); err != nil {
			log.Fatal(err)
		}
	case "paper":
		if err := paper.Handler(f.version, f.build, f.path); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(invalidServerType)
	}
}
