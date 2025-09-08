package addons

import (
	"github.com/Zigl3ur/mcjar/internal/cli/commands/addons/get"
	"github.com/Zigl3ur/mcjar/internal/cli/commands/addons/info"
	"github.com/Zigl3ur/mcjar/internal/cli/commands/addons/search"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addons",
		Short: "Query plugins / mods / modpacks and datapacks from modrinth",
		Long:  "Search, get Info and Download plugins / mods / modpacks and datapacks from modrinth",
		Run: func(cmd *cobra.Command, args []string) {
			//nolint:errcheck
			cmd.Usage()
		},
	}

	cmd.Flags().SortFlags = false
	cmd.AddCommand(search.NewCommand(), info.NewCommand(), get.NewCommand())

	return cmd
}
