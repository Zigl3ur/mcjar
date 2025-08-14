package flags

type ServerType string

const (
	Vanilla  ServerType = "vanilla"
	Paper    ServerType = "paper"
	Purpur   ServerType = "purpur"
	Fabric   ServerType = "fabric"
	Neoforge ServerType = "neoforge"
	Forge    ServerType = "forge"
	Spigot   ServerType = "spigot"
)

func (s ServerType) String() string {
	return string(s)
}
