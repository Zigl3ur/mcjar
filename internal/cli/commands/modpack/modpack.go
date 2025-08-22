package modpack

import (
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/get"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/info"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/search"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modpacks",
		Short: "Search and download modpacks",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.Flags().SortFlags = false
	cmd.AddCommand(search.NewCommand("modpack"))
	cmd.AddCommand(info.NewCommand("mod"))
	cmd.AddCommand(get.NewCommand("modpack"))

	return cmd
}
