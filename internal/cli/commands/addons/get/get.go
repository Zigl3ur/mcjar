package get

import (
	"fmt"
	"slices"

	"github.com/Zigl3ur/mcli/internal/cli/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/modrinth"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Dowload plugin / mod / modpack and datapack from modrinth",
		Long:    "Download a plugin / mod / modpack and datapack from modrinth",
		Args:    cobra.ExactArgs(1),
		PreRunE: validate,
		RunE:    execute,
	}

	cmd.Flags().StringP("version", "v", "", "the game version (required)")
	cmd.Flags().StringP("loader", "l", "", "the loader to be compatible with (required)")
	cmd.Flags().StringP("destination", "d", "", "the folder where to put the downloaded jar")

	cmd.MarkFlagsRequiredTogether("loader", "version")
	//nolint:errcheck
	cmd.MarkFlagDirname("destination")

	return cmd
}

func validate(cmd *cobra.Command, args []string) error {
	mcLoader, _ := cmd.Flags().GetString("loader")

	if cmd.Flag("loader").Changed && !slices.Contains(flags.ValidLoaders, mcLoader) {
		return fmt.Errorf("invalid loader provided (given: %s) valid ones are %s", mcLoader, flags.ValidLoaders)
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	slug := args[0]

	version, _ := cmd.Flags().GetString("version")
	mcLoader, _ := cmd.Flags().GetString("loader")
	dir, _ := cmd.Flags().GetString("destination")
	isVerbose, _ := cmd.Flags().GetBool("verbose")

	filePath, err := modrinth.Download(slug, version, mcLoader, dir)
	if err != nil {
		return err
	}

	// not a mod pack no need to extract
	if filePath == "" {
		return nil
	}

	if err = modrinth.MrPackHandler(filePath, dir, isVerbose); err != nil {
		return err
	}

	return nil
}
