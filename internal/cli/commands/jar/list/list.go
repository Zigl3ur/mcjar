package list

import (
	"fmt"

	"github.com/Zigl3ur/mcli/internal/cli/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/fabric"
	"github.com/Zigl3ur/mcli/internal/handlers/forge"
	"github.com/Zigl3ur/mcli/internal/handlers/neoforge"
	"github.com/Zigl3ur/mcli/internal/handlers/paper"
	"github.com/Zigl3ur/mcli/internal/handlers/purpur"
	"github.com/Zigl3ur/mcli/internal/handlers/vanilla"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List version and builds for specified server type",
		Long:  "List available versions and builds for specified server type",
		RunE:  execute,
	}

	cmd.Flags().StringP("type", "t", "", "The server type to get version / builds list")
	cmd.Flags().StringP("version", "v", "1.21", "The server version to get the list of builds (if available)")
	cmd.Flags().BoolP("snapshots", "s", false, "List snapshots versions (if any)")

	cmd.Flags().SortFlags = false

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {
	serverType := cmd.Flag("type").Value.String()
	version := cmd.Flag("version").Value.String()
	snapshots, _ := cmd.Flags().GetBool("snapshots")
	versionChanged := cmd.Flag("version").Changed

	loader.Start(fmt.Sprintf("fetching %s versions", serverType))

	switch serverType {
	case flags.Vanilla.String():
		return vanilla.ListHandler(snapshots)
	case flags.Paper.String():
		return paper.ListHandler(flags.Paper.String(), version, versionChanged, snapshots)
	case flags.Folia.String():
		return paper.ListHandler(flags.Folia.String(), version, versionChanged, snapshots)
	case flags.Velocity.String():
		return paper.ListHandler(flags.Velocity.String(), version, versionChanged, snapshots)
	case flags.Purpur.String():
		return purpur.ListHandler(version, versionChanged, snapshots)
	case flags.Fabric.String():
		return fabric.ListHandler(version, snapshots)
	case flags.Neoforge.String():
		return neoforge.ListHandler(version, versionChanged, snapshots)
	case flags.Forge.String():
		return forge.ListHandler(version, versionChanged, snapshots)
	default:
		loader.Stop()
		//nolint:errcheck
		cmd.Usage()
		return nil
	}
}
