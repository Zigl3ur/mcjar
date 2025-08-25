package mod

import (
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/get"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/info"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/search"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mods",
		Short: "Search and download mods from modrinth",
		Run: func(cmd *cobra.Command, args []string) {
			//nolint:errcheck
			cmd.Usage()
		},
	}

	cmd.Flags().SortFlags = false
	cmd.AddCommand(search.NewCommand("mod"))
	cmd.AddCommand(info.NewCommand("mod"))
	cmd.AddCommand(get.NewCommand("mod"))

	return cmd
}
