package info

import (
	"fmt"

	"github.com/Zigl3ur/mcli/internal/handlers/modrinth"
	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
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

	cmd.Flags().SortFlags = false

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {
	slug := args[0]

	loader.Start(fmt.Sprintf("Getting info for \"%s\"", slug))

	result, err := modrinth.Info(slug)
	if err != nil {
		loader.Stop()
		return err
	}

	updatedAtFormated := utils.Iso8601Format(result.UpdatedAt)
	createdAtFormated := utils.Iso8601Format(result.CreatedAt)

	loader.Stop()

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
`, result.Title, result.Description, result.Downloads, updatedAtFormated, createdAtFormated, result.Categories, result.ServerSide, result.ClientSide)

	fmt.Println("  - Versions:")
	for loader, versions := range result.LoadersVersions {
		fmt.Printf("        %s: %s\n", loader, versions)
	}

	return nil
}
