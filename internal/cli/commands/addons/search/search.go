package search

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcjar/internal/cli/flags"
	"github.com/Zigl3ur/mcjar/internal/handlers/modrinth"
	"github.com/Zigl3ur/mcjar/internal/utils"
	"github.com/Zigl3ur/mcjar/internal/utils/loader"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [query]",
		Short:   "Search for plugins / mods / modpacks and datapacks from modrinth",
		Long:    "Search for plugins / mods / modpacks and datapacks from modrinth",
		Args:    cobra.ExactArgs(1),
		PreRunE: validate,
		RunE:    execute,
	}

	cmd.Flags().StringP("type", "t", "", "addons type to search for")
	cmd.Flags().Int("limit", 10, "number of results to display")
	cmd.Flags().StringP("index", "i", "relevance", "sort search result")
	cmd.Flags().StringArrayP("versions", "v", []string{}, "versions that results items must support")
	cmd.Flags().StringP("loader", "l", "", "the minecraft loader")
	cmd.Flags().BoolP("slug", "s", false, "show the slug of each item, it is used for info / get commands")

	cmd.Flags().SortFlags = false

	return cmd
}

func validate(cmd *cobra.Command, args []string) error {
	mcLoader, _ := cmd.Flags().GetString("loader")
	addonsType, _ := cmd.Flags().GetString("type")
	index, _ := cmd.Flags().GetString("index")

	if cmd.Flags().Changed("index") && !slices.Contains(flags.ValidIndex, index) {
		return fmt.Errorf("invalid index type provided (given: %s) valid ones are %s", index, flags.ValidIndex)
	}

	if cmd.Flags().Changed("type") && !slices.Contains(flags.ValidAddons, addonsType) {
		return fmt.Errorf("invalid addons type provided (given: %s) valid ones are %s", addonsType, flags.ValidAddons)
	}

	if cmd.Flags().Changed("loader") && !slices.Contains(flags.ValidLoaders, mcLoader) {
		return fmt.Errorf("invalid loader provided (given: %s) valid ones are %s", mcLoader, flags.ValidLoaders)
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	query := args[0]

	addonsType, _ := cmd.Flags().GetString("type")
	slug, _ := cmd.Flags().GetBool("slug")
	limit, _ := cmd.Flags().GetInt("limit")
	index, _ := cmd.Flags().GetString("index")
	versions, _ := cmd.Flags().GetStringArray("versions")
	mcLoader, _ := cmd.Flags().GetString("loader")

	loader.Start(fmt.Sprintf("Searching for \"%s\"", query))

	facets := utils.FacetsBuilder(versions, mcLoader, addonsType)

	results, err := modrinth.Search(query, index, facets, limit)
	loader.Stop()
	if err != nil {
		return err
	}

	var filters []string
	filters = append(filters, fmt.Sprintf("limit: %d", limit))
	filters = append(filters, fmt.Sprintf("index: %s", index))

	if len(versions) > 0 {
		filters = append(filters, fmt.Sprintf("versions: %v", versions))
	}
	if addonsType != "" {
		filters = append(filters, fmt.Sprintf("type: %v", addonsType))
	}
	if mcLoader != "" {
		filters = append(filters, fmt.Sprintf("loader: %s", mcLoader))
	}

	filtersStr := fmt.Sprintf("(%s)", strings.Join(filters, ", "))

	if results.TotalHits == 0 {
		fmt.Printf("No results found for \"%s\" %s\n", query, filtersStr)
	} else {
		fmt.Printf("Found %d results for \"%s\" %s\n", results.TotalHits, query, filtersStr)
		for _, r := range results.Results {
			if slug {
				fmt.Printf("  - [%s] %s, %s\n", r.Slug, r.Title, r.Description)
			} else {
				fmt.Printf("  - %s, %s\n", r.Title, r.Description)
			}
		}
	}
	return nil
}
