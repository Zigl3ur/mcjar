package info

import (
	"fmt"

	"github.com/Zigl3ur/mcli/internal/handlers/modrinth"
	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [slug]",
		Short: "Get Info about a plugin / mod / modpack and datapack from modrinth",
		Long:  "Get detailed information about a specific plugin / mod / modpack and datapack from modrinth",
		Args:  cobra.ExactArgs(1),
		RunE:  execute,
	}

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {
	slug := args[0]

	result, err := modrinth.Info(slug)
	if err != nil {
		return err
	}

	updatedAtFormated, err := utils.Iso8601Format(result.UpdatedAt)
	if err != nil {
		updatedAtFormated = result.UpdatedAt
	}
	createdAtFormated, err := utils.Iso8601Format(result.CreatedAt)
	if err != nil {
		createdAtFormated = result.CreatedAt
	}

	fmt.Printf(`- %s
  %s
  - Downloads: %d
  - Updated: %s
  - Created: %s
  - Categories:
	%s
  - Compatibility:
	Server: %s
	Client: %s
  - Loader:
	%s
  - Game Versions:
	%s
`, result.Title, result.Description, result.Downloads, updatedAtFormated, createdAtFormated, result.Categories, result.ServerSide, result.ClientSide, result.Loaders, result.GameVersions)

	return nil
}
