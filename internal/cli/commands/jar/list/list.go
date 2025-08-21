package list

import (
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/Zigl3ur/mcli/internal/cli/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/fabric"
	"github.com/Zigl3ur/mcli/internal/handlers/forge"
	"github.com/Zigl3ur/mcli/internal/handlers/neoforge"
	"github.com/Zigl3ur/mcli/internal/handlers/paper"
	"github.com/Zigl3ur/mcli/internal/handlers/purpur"
	"github.com/Zigl3ur/mcli/internal/handlers/vanilla"
	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List version and builds for specified server type",
		Run:   execute,
	}

	cmd.Flags().StringP("type", "t", "", "The server type to get version / builds list")
	cmd.Flags().StringP("version", "v", "1.21", "The server version to get the list of builds (if available)")
	cmd.Flags().BoolP("snapshots", "s", false, "List snapshots versions (if any)")

	return cmd
}

func execute(cmd *cobra.Command, args []string) {
	serverType := cmd.Flag("type").Value.String()
	version := cmd.Flag("version").Value.String()
	snapshots, _ := cmd.Flags().GetBool("snapshots")
	versionChanged := cmd.Flag("version").Changed

	loader.Start(fmt.Sprintf("fetching %s versions", serverType))

	switch serverType {
	case flags.Vanilla.String():
		vanilla.ListHandler(snapshots)
	case flags.Paper.String():
		paper.ListHandler(version, versionChanged, snapshots)
	case flags.Purpur.String():
		purpur.ListHandler(version, versionChanged, snapshots)
	case flags.Fabric.String():
		fabric.ListHandler(version, snapshots)
	case flags.Neoforge.String():
		neoforge.ListHandler(version, versionChanged, snapshots)
	case flags.Forge.String():
		vlist, err := forge.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}
		loader.Stop()

		if !cmd.Flag("version").Changed {
			sortedVersion := utils.SortMcVersions(slices.Collect(maps.Keys(vlist)))
			for _, v := range sortedVersion {
				fmt.Printf("- %s\n", v)
			}
		} else if vlist[version] != nil {
			fmt.Printf("- %s\n", version)
			fmt.Println("  - builds:")
			for _, b := range vlist[version] {
				fmt.Printf("   - %s\n", b)
			}
		} else {
			log.Fatal("forge doesn't support this version")
		}

	default:
		loader.Stop()
		cmd.Usage()
	}
}
