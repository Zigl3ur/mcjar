package flags

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type ServerType string

const (
	Vanilla ServerType = "vanilla"
	Forge   ServerType = "forge"
	Mohist  ServerType = "mohist"
	Paper   ServerType = "paper"
	Fabric  ServerType = "fabric"
	Spigot  ServerType = "spigot"
)

type flags struct {
	Version string     // minecraft version
	Type    ServerType // like forge, mohist, paper, etc
	Build   int        // build version like for fabric and forge
	Path    string     // the name of the file outputted from the download
}

func Init() *flags {
	flagsVar := &flags{}

	// version
	pflag.StringVarP(&flagsVar.Version, "version", "v", "1.21.1", "the server version")

	// type
	pflag.StringVarP((*string)(&flagsVar.Type), "type", "t", string(Vanilla), "the server type")

	// build
	// TODO: not sure all build are gonna be int
	pflag.IntVarP(&flagsVar.Build, "build", "b", 0, "the build version")

	// output
	pflag.StringVarP(&flagsVar.Path, "dest", "d", "server.jar", "the destination for the downloaded file")

	pflag.CommandLine.SortFlags = false

	// custom help msg
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -h, --help\t  show this help message\n")
		pflag.VisitAll(func(f *pflag.Flag) {
			fmt.Fprintf(os.Stderr, "  -%s, --%s\t  %s (default: %s)\n", f.Shorthand, f.Name, f.Usage, f.DefValue)
		})
		os.Exit(0)
	}

	return flagsVar
}

func (f *flags) Register() {
	pflag.Parse()
}
