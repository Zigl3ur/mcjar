package flags

import (
	"fmt"
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

	// list
	pflag.StringVarP(&flagsVar.list, "list", "l", "", "list available versions for the specified server type")

	// version
	pflag.StringVarP(&flagsVar.version, "version", "v", "1.21", "the server version")

	// type
	pflag.StringVarP((*string)(&flagsVar.serverType), "type", "t", "vanilla", "the server type")

	// build
	pflag.StringVarP(&flagsVar.build, "build", "b", "", "the build version")

	// output
	pflag.StringVarP(&flagsVar.path, "dest", "d", "server.jar", "the destination for the downloaded file")

	pflag.CommandLine.SortFlags = false

	// custom help msg
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -h, --help\t  show this help message\n")
		pflag.VisitAll(func(f *pflag.Flag) {
			fmt.Fprintf(os.Stderr, "  -%s, --%s\t  %s\n", f.Shorthand, f.Name, f.Usage)
		})
		os.Exit(0)
	}

	return flagsVar
}

func (f *flags) Validate() {
	pflag.Parse()

	validServerType := []string{"vanilla", "forge", "mohist", "paper", "fabric", "spigot"}

	if !slices.Contains(validServerType, f.serverType) {
		fmt.Fprintln(os.Stderr, invalidServerType)
		os.Exit(1)
	}

}

func (f *flags) Execute() {
	if f.list != "" {
		switch strings.ToLower(f.list) {
		case "vanilla":
			vlist, err := vanilla.GetVersionsList()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			for _, v := range vlist.Versions {
				fmt.Fprintf(os.Stdout, "- %s\n", v.Id)
			}
			os.Exit(0)
		case "paper":
			vlist, err := paper.GetVersionsList()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			for _, v := range vlist {
				fmt.Fprintf(os.Stdout, "- %s\n", v)
			}

			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, invalidServerType)
			os.Exit(1)
		}
	}

	fmt.Fprintln(os.Stdout, "Using Values:")
	fmt.Fprintf(os.Stdout, "- type: %s\n", f.serverType)
	fmt.Fprintf(os.Stdout, "- version: %s\n", f.version)
	if f.build != "" {
		fmt.Fprintf(os.Stdout, "- build: %s\n", f.build)
	}
	fmt.Fprintf(os.Stdout, "- dest: %s\n", f.path)

	switch f.serverType {
	case "vanilla":
		if err := vanilla.Handler(f.version, f.path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "paper":
		if err := paper.Handler(f.version, f.path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, invalidServerType)
		os.Exit(0)
	}
}
