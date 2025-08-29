package get

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Dowload plugin / mod / modpack and datapack from modrinth",
		Long:  "Download a plugin / mod / modpack and datapack from modrinth",
		Args:  cobra.ExactArgs(1),
		RunE:  execute,
	}

	cmd.Flags().StringP("version", "v", "", "the game version")
	cmd.Flags().StringP("loader", "l", "", "the loader to be compatible with")
	cmd.Flags().StringP("destination", "d", ".", "the folder where to put the downloaded jar")

	cmd.Flags().SortFlags = false

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {
	// cmd .

	// slug := args[0]

	// if err := modrinth.Download(slug, ""); err != nil {
	// return err
	// }

	return nil
}
