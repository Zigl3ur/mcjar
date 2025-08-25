package get

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Dowload plugin / mod / modpack and datapack from modrinth",
		Long:  "Download a plugin / mod / modpack and datapack from modrinth",
		RunE:  execute,
	}

	cmd.Flags().StringP("destination", "d", ".", "the folder where to put the downloaded jar")

	cmd.Flags().SortFlags = false

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {
	return nil
}
