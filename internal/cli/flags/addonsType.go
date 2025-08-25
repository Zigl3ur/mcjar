package flags

type AddOns string

const (
	Plugin   AddOns = "plugin"
	DataPack AddOns = "datapack"
	Mod      AddOns = "mod"
	Modpack  AddOns = "modpack"
)

var ValidAddons = []string{
	Plugin.String(),
	DataPack.String(),
	Mod.String(),
	Modpack.String(),
}

func (s AddOns) String() string {
	return string(s)
}
