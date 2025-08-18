package list

import (
	"fmt"
	"log"
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
	cmd.Flags().StringP("version", "v", "1.21", "The server version to get the list of builds")

	return cmd
}

func execute(cmd *cobra.Command, args []string) {
	serverType := cmd.Flag("type").Value.String()
	version := cmd.Flag("version").Value.String()

	loader.Start(fmt.Sprintf("fetching %s versions", serverType))

	switch serverType {
	case flags.Vanilla.String():
		vlist, err := vanilla.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}
		loader.Stop()

		for _, v := range vlist.Versions {
			fmt.Printf("- %s\n", v.Id)
		}
	case flags.Paper.String():
		vlist, err := paper.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}
		loader.Stop()

		if !cmd.Flag("version").Changed {
			for _, v := range vlist.Versions {
				fmt.Printf("- %s\n", v.Version.Id)
			}
		} else {
			found := false
			for _, v := range vlist.Versions {
				if v.Version.Id == version {
					fmt.Printf("- %s\n", version)
					fmt.Println("  - builds:")
					for _, b := range v.Builds {
						fmt.Printf("    - %d\n", b)
					}
					found = true
					break
				}
			}
			if !found {
				log.Fatal("paper doesn't support this version")
			}
		}

	case flags.Purpur.String():
		vlist, err := purpur.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}

		loader.Stop()

		if !cmd.Flag("version").Changed {
			for _, v := range vlist {
				fmt.Printf("- %s\n", v)
			}
		} else if slices.Contains(vlist, version) {
			builds, _ := purpur.GetBuildList(version)
			slices.Reverse(builds)
			fmt.Printf("- %s:\n", version)
			fmt.Println("  - builds:")
			for _, b := range builds {
				fmt.Printf("\t- %s\n", b)
			}
		} else {
			log.Fatal("purpur doesn't support this version")
		}
	case flags.Fabric.String():
		vlist, err := fabric.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}
		loader.Stop()

		for _, v := range vlist.Versions {
			fmt.Printf("- %s\n", v.Version)
		}

	case flags.Neoforge.String():
		vlist, err := neoforge.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}
		loader.Stop()

		if !cmd.Flag("version").Changed {
			for k := range vlist {
				fmt.Printf("- %s\n", k)
			}
		} else if vlist[version] != nil {
			fmt.Printf("- %s\n", version)
			fmt.Println("  - neoforge versions:")
			for _, b := range vlist[version] {
				fmt.Printf("    - %s\n", b)
			}
		} else {
			log.Fatal("neoforge doesn't support this version")
		}
	case flags.Forge.String():
		vlist, err := forge.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}
		loader.Stop()

		if !cmd.Flag("version").Changed {
			for k := range vlist {
				fmt.Printf("- %s\n", k)
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
		log.Fatal(utils.InvalidServerType)
	}
}
