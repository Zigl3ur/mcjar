package flags

type ServerType string

const (
	Vanilla  ServerType = "vanilla"
	Paper    ServerType = "paper"
	Folia    ServerType = "folia"
	Velocity ServerType = "velocity"
	Purpur   ServerType = "purpur"
	Fabric   ServerType = "fabric"
	Neoforge ServerType = "neoforge"
	Forge    ServerType = "forge"
)

var ValidServerType = []string{
	Vanilla.String(),
	Paper.String(),
	Folia.String(),
	Velocity.String(),
	Purpur.String(),
	Fabric.String(),
	Neoforge.String(),
	Forge.String(),
}

func (s ServerType) String() string {
	return string(s)
}
