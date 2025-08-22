package datapack

import (
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/get"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/info"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/search"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datapacks",
		Short: "Search and download datapacks",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.Flags().SortFlags = false
	cmd.AddCommand(search.NewCommand("plugin"))
	cmd.AddCommand(info.NewCommand("mod"))
	cmd.AddCommand(get.NewCommand("plugin"))

	return cmd
}
