package list

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Zigl3ur/mcli/internal/cli/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/paper"
	"github.com/Zigl3ur/mcli/internal/handlers/purpur"
	"github.com/Zigl3ur/mcli/internal/handlers/vanilla"
	"github.com/spf13/cobra"
)

const invalidServerType string = "Invalid server type, valid ones are [vanilla, paper, spigot, purpur, forge, fabric]"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List version and builds for specified server type",
		Run:   execute,
	}

	cmd.Flags().StringP("type", "t", flags.Vanilla.String(), "The server type to get version / builds list")
	cmd.Flags().StringP("version", "v", "1.21", "The server version to get the list of builds")

	return cmd
}

func execute(cmd *cobra.Command, args []string) {
	serverType := cmd.Flag("type").Value.String()
	version := cmd.Flag("version").Value.String()

	switch serverType {
	case flags.Vanilla.String():
		vlist, err := vanilla.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range vlist.Versions {
			fmt.Printf("- %s\n", v.Id)
		}
	case flags.Paper.String():
		vlist, err := paper.GetVersionsList()
		if err != nil {
			log.Fatal(err)
		}

		if !cmd.Flag("version").Changed {
			for _, v := range vlist.Versions {
				fmt.Printf("- %s:\n", v.Version.Id)
				fmt.Println("  - builds:")
				for _, b := range v.Builds {
					fmt.Printf("\t- %d\n", b)
				}
			}
		} else {
			found := false
			for _, v := range vlist.Versions {
				if v.Version.Id == version {
					fmt.Printf("- %s:\n", version)
					fmt.Println("  - builds:")
					for _, b := range v.Builds {
						fmt.Printf("\t- %d\n", b)
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
		slices.Reverse(vlist)

		if !cmd.Flag("version").Changed {
			for _, v := range vlist {
				fmt.Printf("- %s:\n", v)
				fmt.Println("  - builds:")
				builds, _ := purpur.GetBuildList(v)
				slices.Reverse(builds)
				for _, b := range builds {
					fmt.Printf("\t- %s\n", b)
				}
			}
		} else {
			if slices.Contains(vlist, version) {
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
		}
	// case "fabric":
	// vlist, err := fabric.GetVersionsList()
	// if err != nil {
	// log.Fatal(err)
	// }
	//
	// if !pflag.Lookup("version").Changed {
	// for _, v := range vlist.Versions {
	// fmt.Printf("- %s:\n", v.Version)
	// }
	// } else {
	// found := false
	// for _, v := range vlist.Versions {
	// if v.Version == f.version {
	// fmt.Printf("- %s\n", v.Version)
	// found = true
	// break
	// }
	// }
	// if !found {
	// log.Fatal("fabric doesn't support this version")
	// }
	// }

	default:
		log.Fatal(invalidServerType)
	}
	os.Exit(0)
}
