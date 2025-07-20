package flags

import "flag"

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
	Output  string     // the name of the file outputted from the download
}

func Init() *flags {
	flagsVar := &flags{}

	// version
	flag.StringVar(&flagsVar.Version, "v", "1.21.8", "the server version")
	flag.StringVar(&flagsVar.Version, "version", "1.21.8", "the server version")

	// type
	flag.StringVar((*string)(&flagsVar.Type), "t", string(Vanilla), "the server type")
	flag.StringVar((*string)(&flagsVar.Type), "type", string(Vanilla), "the server type")

	// build
	// TODO: not sure all build are gonna be int
	flag.IntVar(&flagsVar.Build, "b", 0, "the build version")
	flag.IntVar(&flagsVar.Build, "build", 0, "the build version")

	// output
	flag.StringVar(&flagsVar.Output, "o", "server.jar", "the name of the outputted file")
	flag.StringVar(&flagsVar.Output, "output", "server.jar", "the name of the outputted file")

	return flagsVar
}

func (f *flags) Register() {
	flag.Parse()
}
