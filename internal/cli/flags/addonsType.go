package flags

type AddOns string

const (
	Plugin   AddOns = "plugin"
	DataPack AddOns = "datapack"
	Mod      AddOns = "mod"
	Modpack  AddOns = "modpack"
)

func (s AddOns) String() string {
	return string(s)
}
