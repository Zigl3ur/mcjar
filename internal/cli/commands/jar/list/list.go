package list

import (
	"fmt"
	"slices"

	"github.com/Zigl3ur/mcjar/internal/cli/flags"
	"github.com/Zigl3ur/mcjar/internal/handlers/fabric"
	"github.com/Zigl3ur/mcjar/internal/handlers/forge"
	"github.com/Zigl3ur/mcjar/internal/handlers/neoforge"
	"github.com/Zigl3ur/mcjar/internal/handlers/paper"
	"github.com/Zigl3ur/mcjar/internal/handlers/purpur"
	"github.com/Zigl3ur/mcjar/internal/handlers/vanilla"
	"github.com/Zigl3ur/mcjar/internal/utils/loader"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List version and builds for specified server type",
		Long:    "List available versions and builds for specified server type",
		PreRunE: validate,
		RunE:    execute,
	}

	cmd.Flags().StringP("type", "t", "", "The server type to get version / builds list")
	cmd.Flags().StringP("version", "v", "1.21", "The server version to get the list of builds (if any)")
	cmd.Flags().BoolP("snapshots", "s", false, "List snapshots versions (if any)")

	cmd.Flags().SortFlags = false

	return cmd
}

func validate(cmd *cobra.Command, args []string) error {
	serverType, _ := cmd.Flags().GetString("type")

	if cmd.Flag("type").Changed && !slices.Contains(flags.ValidServerType, serverType) {
		return fmt.Errorf("invalid type, valid ones are %s", flags.ValidServerType)
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	serverType, _ := cmd.Flags().GetString("type")
	version, _ := cmd.Flags().GetString("version")
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
